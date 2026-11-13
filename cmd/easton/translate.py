#!/usr/bin/env python3
"""Machine-translates the Easton dictionary text of people and places into data/easton/<lang>.json.

    python3 cmd/easton/translate.py nl            # all entries
    python3 cmd/easton/translate.py nl 40         # only the first 40, for a quick look

Offline: CTranslate2 with an open English->Dutch model (Argos Translate's package, converted
from the OPUS-MT models). Set MT_HOME to the folder holding lib/ (pip install --target lib ctranslate2
sentencepiece) and model/ (the unzipped .argosmodel). The links to verses stay as they are; only the
text between them is translated. The result is machine translation and the site says so.
"""
import json, os, re, sys, time

home = os.environ.get('MT_HOME', os.path.expanduser('~/.cache/metabible-mt'))
sys.path.insert(0, home + '/lib')
import ctranslate2, sentencepiece as spm
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

lang = sys.argv[1] if len(sys.argv) > 1 else 'nl'
limit = int(sys.argv[2]) if len(sys.argv) > 2 else 0
model = [d for d in os.listdir(home + '/model') if d.startswith('translate-en_' + lang)][0]
M = f'{home}/model/{model}/'
sp = spm.SentencePieceProcessor(model_file=M + 'sentencepiece.model')
tr = ctranslate2.Translator(M + 'model', device='cpu', inter_threads=4, intra_threads=4, compute_type='int8')

LINK = re.compile(r'\[[^\]]*\]\([^)]*\)')
ABBR = set("""Heb Gr Lat Syr Arab Chald Egypt St Mt Ch ver vers comp cf viz etc ib lit No sing pl ex Ps Is Dr Rev Gen Ex Lev Num Deut Josh Judg
Sam Kings Chr Neh Esth Prov Eccl Isa Jer Ezek Dan Hos Obad Mic Nah Hab Zeph Hag Zech Mal Matt Rom Cor Gal Eph Phil Col Thess Tim Philem Jas
Pet Jude Ap Jos B.C A.D i.e e.g v vv""".split())
CAND = re.compile(r'[.!?]\s+(?=[A-Z])')

def split_sentences(s):
    """Sentences, without cutting after an abbreviation (Heb. 'abel) or an initial."""
    out, last = [], 0
    for m in CAND.finditer(s):
        word = re.search(r'([A-Za-z.]+)$', s[last:m.start()])
        if s[m.start()] == '.' and word and (word.group(1) in ABBR or (len(word.group(1)) == 1 and word.group(1).isupper())):
            continue
        out.append(s[last:m.start() + 1])
        last = m.end()
    out.append(s[last:])
    return [x for x in out if x.strip()]

from translate_labels import localise_label  # noqa: E402

def core(seg):
    """(leading, text, trailing): punctuation and spaces around the part worth translating."""
    m = re.match(r'^([\s(\[;:,.\-–—)\]]*)(.*?)([\s(\[;:,\-–—]*)$', seg, re.S)
    return m.group(1), m.group(2), m.group(3)

def translate_batch(texts, beam=4):
    toks = [sp.encode(t, out_type=str) for t in texts]
    res = tr.translate_batch(toks, beam_size=beam, max_batch_size=24, batch_type='examples',
                             max_decoding_length=[int(len(t) * 2 + 12) for t in toks] and 512, repetition_penalty=1.1)
    out = []
    for t, r in zip(texts, res):
        o = ''.join(r.hypotheses[0]).replace('\u2581', ' ').strip()
        # a runaway or empty result is worse than the English
        if not o or len(o) >= 3 * len(t) + 40:
            o = t
        else:
            # the model finishes every piece like a sentence; the source may go on after a link
            if o[-1] in '.!?' and t[-1] not in '.!?':
                o = o[:-1]
            if t[0].islower() and o[0].isupper() and not o[:2].isupper():
                o = o[0].lower() + o[1:]
        out.append(o)
    return out

data = []
for f in ('people', 'places'):
    for x in json.load(open(f'data/json/{f}.json')):
        if (x.get('dictText') or '').strip():
            data.append((x['id'], x['dictText']))
if limit:
    data = data[:limit]

# 1. cut every entry into pieces: ('raw', link or punctuation) or ('tx', sentence)
plan, sentences = {}, []
for id_, text in data:
    paras = []
    for para in text.strip().split('\n\n'):
        parts, last = [], 0
        for m in list(LINK.finditer(para)) + [None]:
            seg = para[last:m.start()] if m else para[last:]
            lead, body, trail = core(seg)
            parts.append(('raw', lead))
            if re.search(r'[A-Za-z]{2}', body):
                for s in split_sentences(body):
                    parts.append(('tx', len(sentences)))
                    sentences.append(s.strip())
                    parts.append(('raw', ' '))
                parts.pop()  # the space after the last sentence
            else:
                parts.append(('raw', body))
            parts.append(('raw', trail))
            if m:
                parts.append(('raw', localise_label(m.group(0))))
                last = m.end()
        paras.append(parts)
    plan[id_] = paras

# 2. translate, longest first so batches are even
order = sorted(range(len(sentences)), key=lambda i: len(sentences[i]))
done = [None] * len(sentences)
t0 = time.time()
B = 240
for k in range(0, len(order), B):
    idx = order[k:k + B]
    for i, o in zip(idx, translate_batch([sentences[i] for i in idx])):
        done[i] = o
    if (k // B) % 10 == 0:
        print(f'{k + len(idx)}/{len(order)} sentences, {time.time() - t0:.0f}s', file=sys.stderr, flush=True)

# 3. put the entries back together
out = {}
for id_, paras in plan.items():
    out[id_] = '\n\n'.join(''.join(v if kind == 'raw' else done[v] for kind, v in parts) for parts in paras)
os.makedirs('data/easton', exist_ok=True)
path = f'data/easton/{lang}.json' if not limit else f'data/easton/{lang}.sample.json'
json.dump(out, open(path, 'w'), ensure_ascii=False, indent=0, sort_keys=True)
print(len(out), 'entries,', len(sentences), 'sentences,', f'{time.time() - t0:.0f}s ->', path, file=sys.stderr)
