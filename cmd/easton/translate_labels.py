"""Dutch book abbreviations for the verse links in the Easton texts (shared by translate.py and hand.py)."""
import re

BOOKS = {  # Easton's abbreviation -> the usual Dutch one
 'Gen.':'Gen.','Ex.':'Ex.','Lev.':'Lev.','Num.':'Num.','Deut.':'Deut.','Josh.':'Joz.','Judg.':'Richt.','Ruth':'Ruth','Sam.':'Sam.','Kings':'Kon.','Chr.':'Kron.','Ezra':'Ezra','Neh.':'Neh.',
 'Esther':'Esth.','Esth.':'Esth.','Job':'Job','Ps.':'Ps.','Prov.':'Spr.','Eccl.':'Pred.','Cant.':'Hoogl.','Song':'Hoogl.','Isa.':'Jes.','Jer.':'Jer.','Lam.':'Klaagl.','Ezek.':'Ezech.','Dan.':'Dan.',
 'Hos.':'Hos.','Joel':'Joël','Amos':'Amos','Obad.':'Obad.','Jonah':'Jona','Mic.':'Micha','Nah.':'Nah.','Hab.':'Hab.','Zeph.':'Zef.','Hag.':'Hag.','Zech.':'Zach.','Mal.':'Mal.','Matt.':'Matt.',
 'Mark':'Marc.','Luke':'Luc.','John':'Joh.','Acts':'Hand.','Rom.':'Rom.','Cor.':'Kor.','Gal.':'Gal.','Eph.':'Ef.','Phil.':'Fil.','Col.':'Kol.','Thess.':'Thess.','Tim.':'Tim.','Titus':'Tit.',
 'Philem.':'Filem.','Heb.':'Hebr.','James':'Jak.','Pet.':'Petr.','Jude':'Jud.','Rev.':'Openb.'}

def localise_label(link):
    m = re.match(r'^\[(\d?\s?)([A-Za-z]+\.?)(\s+[\d].*?)\]\((.*)\)$', link, re.S)
    if m and m.group(2) in BOOKS:
        return f'[{m.group(1)}{BOOKS[m.group(2)]}{m.group(3)}]({m.group(4)})'
    return link
