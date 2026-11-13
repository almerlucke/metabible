#!/usr/bin/env python3
"""Hand translations of the Wikipedia summaries in data/enrichment.json (CC BY-SA 4.0).

    python3 cmd/easton/wiki_hand.py prepare DIR --per 40     # batches of {id, kind, name, desc, summary} + glossary.tsv
    python3 cmd/easton/wiki_hand.py merge DIR                # check DIR/out_*.json and write data/enrichment.nl.json
    python3 cmd/easton/wiki_hand.py merge DIR --allow-english  # after reading the flagged ones (English names/titles on purpose)

The translator writes out_NN.json: {id: {"summary": ..., "desc": ...}}. The result keeps the English article's
title and link (it is a translation of that article) and replaces summary and desc; photos stay those of
data/enrichment.json. pkg/atlas reads data/enrichment.<lang>.json when the site is shown in that language.
"""
import json, os, re, sys, time

lang = 'nl'

def load():
    return json.load(open('data/enrichment.json'))

def prepare(dir_, per):
    en = load()
    names = {}
    for f, key in (('people', 'people'), ('places', 'places')):
        titles = {x['id']: x['displayTitle'] for x in json.load(open(f'data/json/{f}.json'))}
        for id_, nl in json.load(open(f'data/names/{lang}.json'))[key].items():
            names.setdefault(titles[id_], nl)
    items = []
    for kind in ('places', 'people'):
        titles = {x['id']: x['displayTitle'] for x in json.load(open(f'data/json/{kind}.json'))}
        for id_, w in en[kind].items():
            items.append({'id': id_, 'kind': kind[:-1], 'name': titles.get(id_, w['title']), 'article': w['title'], 'desc': w.get('desc', ''), 'summary': w['summary']})
    os.makedirs(dir_, exist_ok=True)
    for k in range(0, len(items), per):
        json.dump(items[k:k + per], open(f'{dir_}/batch_{k // per:02d}.json', 'w'), ensure_ascii=False, indent=1)
    text = ' '.join(i['name'] + ' ' + i['article'] + ' ' + i['summary'] + ' ' + i['desc'] for i in items)
    words = set(re.findall(r"[A-Z][A-Za-z\-']+", text))
    # the book names of the site, English -> Dutch
    en_loc, nl_loc = json.load(open('pkg/i18n/locales/en.json')), json.load(open(f'pkg/i18n/locales/{lang}.json'))
    with open(f'{dir_}/glossary.tsv', 'w') as g:
        for k, v in sorted(en_loc.items()):
            if k.startswith('book.'):
                g.write(f'{v}\t{nl_loc[k]}\t(Bible book)\n')
        for e, n in sorted(names.items()):
            if '(' not in e and e in words and e != n:
                g.write(f'{e}\t{n}\n')
    print(len(items), 'entries in', (len(items) + per - 1) // per, 'batches ->', dir_)

ENGLISH = re.compile(r'\b(the|and|which|were|his|their|with|from|that|also)\b')

def merge(dir_, allow_english=False):
    en = load()
    spelling = json.load(open(f'data/names/{lang}.glue.json')).get('spelling', {})
    def modern(t): return re.sub(r'[^\W\d_]+', lambda m: spelling.get(m.group(0), m.group(0)), t)
    path = f'data/enrichment.{lang}.json'
    out = json.load(open(path)) if os.path.exists(path) else {'lang': lang, 'places': {}, 'people': {}}
    problems, added = [], 0
    for f in sorted(os.listdir(dir_)):
        if not f.startswith('out_'): continue
        for id_, t in json.load(open(f'{dir_}/{f}')).items():
            kind = 'places' if id_ in en['places'] else 'people' if id_ in en['people'] else None
            if not kind: problems.append((id_, 'unknown id')); continue
            src = en[kind][id_]
            s, d = (t.get('summary') or '').strip(), (t.get('desc') or '').strip()
            if not s: problems.append((id_, 'no summary')); continue
            if src.get('desc') and not d: problems.append((id_, 'no desc')); continue
            r = len(s) / max(1, len(src['summary']))
            if not 0.6 < r < 1.8: problems.append((id_, f'length ratio {r:.2f}')); continue
            if not allow_english and ENGLISH.search(s): problems.append((id_, 'English words left')); continue
            out[kind][id_] = {'title': src['title'], 'summary': modern(s), 'desc': modern(d), 'page': src['page']}
            added += 1
    out['generated'] = time.strftime('%Y-%m-%d')
    json.dump(out, open(path, 'w'), ensure_ascii=False, indent=1, sort_keys=True)
    print(added, 'added;', len(out['places']), 'places and', len(out['people']), 'people in total;', len(problems), 'problems')
    for p in problems: print('  ', p)

if __name__ == '__main__':
    cmd, dir_ = sys.argv[1], sys.argv[2]
    args = dict(zip(sys.argv[3::2], sys.argv[4::2]))
    prepare(dir_, int(args.get('--per', 40))) if cmd == 'prepare' else merge(dir_, '--allow-english' in sys.argv)
