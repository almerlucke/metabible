#!/usr/bin/env python3
"""Hand-made (human-quality) translations of the Easton texts, merged over the machine translation.

    python3 cmd/easton/hand.py prepare DIR --top 20          # batches of the most-read entries not done yet
    python3 cmd/easton/hand.py prepare DIR --short 300       # all entries of at most 300 characters
    python3 cmd/easton/hand.py merge DIR                      # check DIR/out_*.json and add them to data/easton/nl.hand.json

A batch is a list of {id, name, text}; every verse link in the text is replaced by a placeholder
(⟦1⟧, ⟦2⟧…) that the translator must keep. The translator writes out_NN.json: {id: translated text}.
Merging checks the placeholders and puts the links back, with Dutch book abbreviations.
"""
import json, os, re, sys

lang = 'nl'
LINK = re.compile(r'\[([^\]]*)\]\(([^)]*)\)')
sys.path.insert(0, os.path.dirname(__file__))
from translate_labels import localise_label  # noqa: E402

def load_entities():
    out = {}
    for f in ('people', 'places'):
        for x in json.load(open(f'data/json/{f}.json')):
            if (x.get('dictText') or '').strip():
                out[x['id']] = x
    return out

def hand_path(): return f'data/easton/{lang}.hand.json'
def done(): return json.load(open(hand_path())) if os.path.exists(hand_path()) else {}

def placeholders(text):
    links = []
    def sub(m):
        links.append(m.group(0)); return f'⟦{len(links)}⟧'
    return LINK.sub(sub, text.strip()), links

def prepare(dir_, top=0, short=0, per_batch=10):
    ents, have = load_entities(), done()
    rows = [(int(x.get('verseCount') or 0), id_) for id_, x in ents.items() if id_ not in have]
    if short:
        picked = [id_ for _, id_ in sorted(rows, reverse=True) if len(placeholders(ents[id_]['dictText'])[0]) <= short]
    else:
        picked = [id_ for _, id_ in sorted(rows, reverse=True)[:top]]
    os.makedirs(dir_, exist_ok=True)
    names = json.load(open(f'data/names/{lang}.json'))
    people = {x['id']: x['displayTitle'] for x in json.load(open('data/json/people.json'))}
    places = {x['id']: x['displayTitle'] for x in json.load(open('data/json/places.json'))}
    gloss = {}
    for kind, en in (('people', people), ('places', places)):
        for id_, nl in names[kind].items():
            gloss.setdefault(en[id_], nl)
    for k in range(0, len(picked), per_batch):
        batch, texts = [], []
        for id_ in picked[k:k + per_batch]:
            t, links = placeholders(ents[id_]['dictText'])
            batch.append({'id': id_, 'name': ents[id_]['displayTitle'], 'text': t, 'links': len(links)})
            texts.append(t)
        json.dump(batch, open(f'{dir_}/batch_{k // per_batch:02d}.json', 'w'), ensure_ascii=False, indent=1)
    words = set(re.findall(r"[A-Z][A-Za-z\-']+", ' '.join(b['text'] for b in (json.load(open(f'{dir_}/{f}')) for f in os.listdir(dir_) if f.startswith('batch_')) for b in b)))
    names_ = {w for w in words}
    with open(f'{dir_}/glossary.tsv', 'w') as g:
        for en, nl in sorted(gloss.items()):
            if '(' not in en and en in names_ and en != nl:
                g.write(f'{en}\t{nl}\n')
    print(len(picked), 'entries in', (len(picked) + per_batch - 1) // per_batch, 'batches ->', dir_)

def merge(dir_):
    ents, have = load_entities(), done()
    added, problems = 0, []
    for f in sorted(os.listdir(dir_)):
        if not f.startswith('out_'): continue
        batch = {b['id']: b for b in json.load(open(f'{dir_}/batch_{f[4:]}'))}
        for id_, text in json.load(open(f'{dir_}/{f}')).items():
            if id_ not in batch:
                problems.append((id_, 'not in the batch')); continue
            src, links = placeholders(ents[id_]['dictText'])
            found = re.findall(r'⟦(\d+)⟧', text)
            if sorted(map(int, found)) != list(range(1, len(links) + 1)):
                problems.append((id_, f'placeholders {sorted(map(int, found))} for {len(links)} links')); continue
            if len(src) > 120 and not 0.6 < len(text) / len(src) < 1.7:
                problems.append((id_, f'length ratio {len(text) / len(src):.2f}')); continue
            if re.search(r'\b(the|and|which|his|their)\b', text):
                problems.append((id_, 'English words left')); continue
            out = re.sub(r'⟦(\d+)⟧', lambda m: localise_label(links[int(m.group(1)) - 1]), text.strip())
            have[id_] = out; added += 1
    json.dump(have, open(hand_path(), 'w'), ensure_ascii=False, indent=0, sort_keys=True)
    print(added, 'added;', len(have), 'in total;', len(problems), 'problems')
    for p in problems: print('  ', p)

if __name__ == '__main__':
    cmd, dir_ = sys.argv[1], sys.argv[2]
    args = dict(zip(sys.argv[3::2], sys.argv[4::2]))
    if cmd == 'prepare': prepare(dir_, int(args.get('--top', 0)), int(args.get('--short', 0)), int(args.get('--per', 10)))
    else: merge(dir_)
