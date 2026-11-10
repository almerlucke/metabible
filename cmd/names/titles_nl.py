#!/usr/bin/env python3
"""Writes data/titles/nl.json: Dutch titles of the events in data/json/events.json.
Birth/Lifetime/Death titles are made from the name table (go run ./cmd/names); the rest were translated by hand below.
Run from the repository root: python3 cmd/names/titles_nl.py"""
import json, re

events = json.load(open('data/json/events.json'))
people = json.load(open('data/json/people.json'))
table = json.load(open('data/names/nl.json'))['people']
glue = json.load(open('data/names/nl.glue.json'))['names']

# English name -> Dutch, from the person with most verses of that name
best = {}
for p in people:
    n = p['displayTitle']; vc = int(p['verseCount'] or 0)
    if p['id'] in table and (n not in best or vc > best[n][0]):
        best[n] = (vc, table[p['id']])
fix = {'Enoch': 'Henoch', 'Cainan': 'Kenan', 'Mahlaleel': 'Mahalaleël', 'Jehoahash': 'Joahaz', 'Jeroboam II': 'Jerobeam II', 'Abijam': 'Abiam'}  # the table has the other person of that name
def nm(en):
    return fix.get(en) or glue.get(en) or (best[en][1] if en in best else en)

kings = {  # names in "Reign of / Prophecies of / Judgeship of ..."
 'Jair':'Jaïr','Tola':'Tola','Eli':'Eli','Ibzan':'Ibzan','Elon':'Elon','Abdon':'Abdon','Samson':'Simson','Saul':'Saul','David':'David','Solomon':'Salomo','Rehoboam':'Rehabeam',
 'Jeroboam I':'Jerobeam I','Jeroboam II':'Jerobeam II','Abijam':'Abiam','Asa':'Asa','Nadab':'Nadab','Baasha':'Baësa','Elah':'Ela','Zimri':'Zimri','Tibni':'Tibni','Omri':'Omri','Ahab':'Achab',
 'Jehoshaphat':'Josafat','Elijah':'Elia','Ahaziah':'Ahazia','Elisha':'Elisa','Joram':'Joram','Jehoram':'Joram','Ahaziah (Jehoahaz)':'Ahazia (Joahaz)','Jehu':'Jehu','Athaliah':'Athalia',
 'Joash':'Joas','Jehoahaz':'Joahaz','Jehoash':'Joas','Amaziah':'Amazia','Jonah':'Jona','Joel':'Joël','Uzziah':'Uzzia','Amos':'Amos','Hosea':'Hosea','Zachariah':'Zacharia','Shallum':'Sallum',
 'Menahem':'Menahem','Pekahiah':'Pekahia','Pekah':'Pekah','Isaiah':'Jesaja','Jotham':'Jotham','Micah':'Micha','Ahaz':'Achaz','Hoshea':'Hosea','Hezekiah':'Hizkia','Manasseh':'Manasse',
 'Nahum':'Nahum','Amon':'Amon','Josiah':'Josia','Jeremiah':'Jeremia','Zephaniah':'Zefanja','Johoahaz (Shallum)':'Joahaz (Sallum)','Jehoiakim':'Jojakim','Habakkuk':'Habakuk','Daniel':'Daniël',
 'Jeconiah':'Jechonia','Zedekiah':'Zedekia','Ezekiel':'Ezechiël','Obadiah':'Obadja','Haggai':'Haggaï','Zechariah':'Zacharia','Malachi':'Maleachi',
 'Cushan-rishathaim':'Kusan-Rishathaïm','Othniel':'Othniël','Eglon':'Eglon','Ehud':'Ehud','Jabin':'Jabin','Barak And Deborah':'Barak en Debora','Midian':'Midian','Gideon':'Gideon',
 'Abimelech':'Abimelech','Jephthah':'Jefta','Ammon and Philistia':'Ammon en Filistea'}
prefix = [('Reign of ','Regering van '),('Prophecies of ','Profetieën van '),('Judgeship of ','Richterschap van '),('Deliverance by ','Bevrijding door '),
          ('Subjugation by ','Onderwerping door '),('Usurpation by ','Machtsovername door '),('Oppression by ','Onderdrukking door ')]
manual = {
 'Birth of Jacob and Esau':'Geboorte van Jakob en Ezau','Death of Eli, Ark Captured':'Dood van Eli; de ark buitgemaakt',
 'Creation of all things':'Schepping van alle dingen','Creation of Adam and Eve':'Schepping van Adam en Eva','The Fall':'De zondeval','Cain kills Abel':'Kaïn doodt Abel','God curses Cain':'God vervloekt Kaïn',
 'Enoch taken up':'Henoch opgenomen','God decides to destroy every living thing':'God besluit al het leven te vernietigen','God instructs Noah to build the ark':'God draagt Noach op de ark te bouwen',
 'The Great Flood begins':'De grote zondvloed begint','The Great Flood':'De grote zondvloed','Water covers the earth':'Het water bedekt de aarde','The flood subsides':'De vloed zakt','Leaving the Ark':'Het verlaten van de ark',
 'God makes a covenant with Noah':'God sluit een verbond met Noach','Noah gets drunk and Canaan is cursed':'Noach wordt dronken en Kanaän wordt vervloekt','Tower of Babel':'De toren van Babel',
 'Abraham enters Canaan':'Abraham komt in Kanaän','Abrahamic Covenant':'Het verbond met Abraham','Abraham goes to Egypt':'Abraham gaat naar Egypte','Sojourn in Canaan and Egypt':'Verblijf in Kanaän en Egypte',
 'Abraham and Lot Separate':'Abraham en Lot gaan uiteen','Sodom Destroyed':'Sodom verwoest','Isaac Marries Rebekah':'Izak trouwt met Rebekka',"Nimrod's kingdom begins":'Het koninkrijk van Nimrod begint',
 'Rachel dies giving birth to Benjamin':'Rachel sterft bij de geboorte van Benjamin','Joseph sold to Egypt':'Jozef verkocht naar Egypte','Joseph imprisoned':'Jozef gevangengezet','Joseph Promoted in Egypt':'Jozef verheven in Egypte',
 '7 years of plenty':'Zeven jaren van overvloed','7 years of famine':'Zeven jaren van hongersnood',"Joseph's brothers go to Egypt":'Jozefs broers gaan naar Egypte',"Joseph's brothers return to Egypt":'Jozefs broers keren terug naar Egypte',
 "Job's Trial":'Jobs beproeving','Moses Flees Pharaoh':'Mozes vlucht voor Farao','Exodus from Egypt':'Uittocht uit Egypte','Wilderness Wanderings':'Omzwervingen in de woestijn','Ten Commandments Given':'De Tien Geboden gegeven',
 'Tabernacle Built':'De tabernakel gebouwd','Crossing the Jordan into Canaan':'Over de Jordaan Kanaän in','War with Canaanites':'Oorlog met de Kanaänieten','Dividing the land':'Verdeling van het land',
 'Philistine Dominion':'Filistijnse overheersing','Sampson destroys Temple of Dagon':'Simson verwoest de tempel van Dagon','United Kingdom':'Het verenigde koninkrijk','David Kills Goliath':'David doodt Goliath',
 "Construction of Solomon's Temple":'Bouw van de tempel van Salomo','Southern Kingdom (Judah)':'Zuidelijk koninkrijk (Juda)','Northern Kingdom (Israel)':'Noordelijk koninkrijk (Israël)','Elijah Translated':'Elia opgenomen',
 "John's birth predicted":'De geboorte van Johannes wordt aangekondigd','Espousal of Mary':'Verloving van Maria','The Annuciation':'De aankondiging aan Maria','Mary Visits Elizabeth':'Maria bezoekt Elisabet',
 'An Angel Speaks to Joseph in a Dream':'Een engel spreekt Jozef in een droom','Wise Men Visit Herod':'De wijzen bezoeken Herodes','Wise Men Visit Jesus':'De wijzen bezoeken Jezus','Joseph and Mary Flee to Egypt':'Jozef en Maria vluchten naar Egypte',
 'Herod slays male children':'Herodes laat de jongetjes doden','Jesus Circumsized':'Jezus wordt besneden','Joseph and Mary return from Egypt':'Jozef en Maria keren terug uit Egypte','Temple Presentation of Jesus':'Jezus wordt in de tempel opgedragen',
 'Joseph and Mary return to Nazareth':'Jozef en Maria keren terug naar Nazareth','Childhood of Jesus':'Jeugd van Jezus','Jesus as a child in the Temple':'Jezus als kind in de tempel','Jesus Growing Up':'Jezus groeit op',
 'John the Baptist Begins Ministry':'Johannes de Doper begint zijn bediening','John Baptizes Jesus':'Johannes doopt Jezus','Temptations of Jesus':'De verzoekingen van Jezus',"John's Witness":'Het getuigenis van Johannes',
 'First Disciples Gathered':'De eerste discipelen','Water to Wine':'Water wordt wijn','Visiting Capernaum':'Bezoek aan Kapernaüm','First Passover and Temple Cleansing':'Eerste Pascha en reiniging van de tempel','Jesus and Nicodemus':'Jezus en Nicodemus',
 'John the Baptist Witnesses':'Johannes de Doper getuigt','John the Baptist Imprisoned':'Johannes de Doper gevangengezet','Jesus Leaves Judea toward Galilee':'Jezus verlaat Judea richting Galilea','The Samaritan Woman':'De Samaritaanse vrouw',
 'Jesus visits Galilee':'Jezus bezoekt Galilea','Jesus Rejected at Nazareth':'Jezus verworpen in Nazareth','Jesus in Capernaum':'Jezus in Kapernaüm','Fishers of Men':'Vissers van mensen','The Demoniac':'De bezetene',
 "Peter's Mother-in-law and Others Healed":"Petrus' schoonmoeder en anderen genezen","Jesus' 1st Tour of Galilee":"Jezus' eerste reis door Galilea",'Healing the Leper':'Genezing van de melaatse','Healing the Man with Palsy':'Genezing van de verlamde',
 'Jesus Calls Matthew':'Jezus roept Mattheüs','Fasting Question':'De vraag over het vasten','Parable of the Winseskins':'Gelijkenis van de wijnzakken','Bethesda pool healing':'Genezing bij de vijver Bethesda',
 'Corn on the Sabbath':'Aren plukken op de sabbat','Healing a Withered Hand':'Genezing van de verdorde hand','Healing Multitudes':'Genezing van de menigten','Apostles Chosen':'De apostelen gekozen','Sermon on the Mount':'De Bergrede',
 "Healing Centurion's Servant":'Genezing van de knecht van de hoofdman',"Widow's Son Resurrected":'De zoon van de weduwe opgewekt',"Jesus' Message on John the Baptist":"Jezus' boodschap over Johannes de Doper",
 'Woe to Chorazin and Bethsaida':'Wee u, Chorazin en Bethsaïda','Anointed by a Sinner Woman':'Gezalfd door een zondares',"Jesus' 2nd Tour of Galilee":"Jezus' tweede reis door Galilea",
 'Blind and Dumb Demoniac and Following Discourse':'De blinde en stomme bezetene en de toespraak daarna','The Holy Spirit is promised':'De Heilige Geest wordt beloofd','Jesus ascends to Heaven':'Jezus vaart op naar de hemel',
 'Matthias replaces Judas':'Matthias volgt Judas op','The Holy Spirit comes':'De Heilige Geest komt','Peter preaches at Pentecost':'Petrus preekt met Pinksteren','The church grows':'De gemeente groeit','Lame man healed':'De kreupele man genezen',
 "Peter preaches in Solomon's Portico":'Petrus preekt in de zuilengang van Salomo',"Peter and John's Trial":'Het verhoor van Petrus en Johannes','Believers pray for boldness':'De gelovigen bidden om vrijmoedigheid',
 'Believers share all things':'De gelovigen delen alles','Ananias and Sapphira lie to God':'Ananias en Saffira liegen tegen God',"Apostles work signs and wonders in Solomon's Portico":'Apostelen doen tekenen en wonderen in de zuilengang van Salomo',
 'Apostles arrested, put on trial':'Apostelen gearresteerd en verhoord','Gamaliel advises the counsel and Apostles freed':'Gamaliël adviseert de raad en de apostelen worden vrijgelaten','Seven chosen to serve':'Zeven mannen gekozen om te dienen',
 'Stephen is seized':'Stefanus wordt gegrepen','Stephen addresses the council':'Stefanus spreekt de raad toe','Stephen is stoned':'Stefanus wordt gestenigd','Saul persecutes the church':'Saulus vervolgt de gemeente',
 'Philip preaches in Samaria':'Filippus preekt in Samaria','Conversion of Ethiopian Eunuch':'Bekering van de Ethiopische kamerling','Saul is converted':'Saulus wordt bekeerd','Saul proclaims Jesus':'Saulus verkondigt Jezus',
 'Peter heals Aeneas and Dorcas':'Petrus geneest Eneas en Dorkas','Peter meets Cornelius':'Petrus ontmoet Cornelius','Peter preaches to gentiles in Caesarea':'Petrus preekt tot heidenen in Cesarea',
 'Peter defends Gentile salvation':'Petrus verdedigt het heil voor de heidenen','Disciples first called Christians':'Discipelen voor het eerst christenen genoemd','Herod kills James':'Herodes doodt Jakobus','Herod imprisons Peter':'Herodes zet Petrus gevangen',
 'An angel rescues Peter from prison':'Een engel bevrijdt Petrus uit de gevangenis','Herod dies':'Herodes sterft','First missionary journey begins':'Eerste zendingsreis begint','First Missionary Journey':'Eerste zendingsreis',
 'Barnabas and Saul go to Cyprus':'Barnabas en Saulus gaan naar Cyprus','Mission to Antioch in Pisidia':'Zending naar Antiochië in Pisidië','Mission to Iconium':'Zending naar Ikonium','Mission to Lystra and Derbe':'Zending naar Lystra en Derbe',
 'Return to Antioch in Syria':'Terugkeer naar Antiochië in Syrië','Jerusalem Council':'Het concilie van Jeruzalem','Second Missionary Journey':'Tweede zendingsreis','Paul and Barnabas part ways':'Paulus en Barnabas gaan uiteen',
 'Paul Writes Galatians':'Paulus schrijft aan de Galaten','Second Missionary Journey begins':'Tweede zendingsreis begint','Timothy Joins Paul and Silas':'Timotheüs voegt zich bij Paulus en Silas',
 'Mission to Phrygia, Galatia and Asia':'Zending naar Frygië, Galatië en Azië','Call to Macedonia':'De roep naar Macedonië',"Lydia's Conversion":'Bekering van Lydia','Paul casts demon out of soothsayer':'Paulus drijft een boze geest uit een waarzegster',
 'Paul and Silas imprisoned':'Paulus en Silas gevangengezet','Philippian jailer converted':'De gevangenbewaarder van Filippi bekeerd','Paul and Silas released':'Paulus en Silas vrijgelaten','Uproar in Thessalonica':'Oproer in Thessalonica',
 'Mission to Berea':'Zending naar Berea','Paul preaches in Athens':'Paulus preekt in Athene','Mission to Corinth, 1&2 Thess Written':'Zending naar Korinthe; 1 en 2 Thessalonicenzen geschreven',
 'Return from Second Missionary Journey':'Terugkeer van de tweede zendingsreis','Third Missionary Journey Begins':'Derde zendingsreis begint','Third Missionary Journey':'Derde zendingsreis',
 'Mission to Ephesus/1 Cor Written':'Zending naar Efeze; 1 Korintiërs geschreven','Riot in Ephesus':'Oproer in Efeze','Mission to Macedonia and Greece':'Zending naar Macedonië en Griekenland','Eutychus revived':'Eutychus levend gemaakt',
 'Voyage to Miletus':'Reis naar Milete','Paul speaks to Ephesian Elders':'Paulus spreekt de oudsten van Efeze toe','Voyage from Miletus to Jerusalem':'Reis van Milete naar Jeruzalem','Paul visits elders in Jerusalem':'Paulus bezoekt de oudsten in Jeruzalem',
 'Paul arrested in the Temple':'Paulus gearresteerd in de tempel','Paul gives his testimony':'Paulus legt zijn getuigenis af','Paul goes before the Council':'Paulus voor de Raad','The conspiracy to kill Paul':'De samenzwering om Paulus te doden',
 'Paul sent to Felix':'Paulus naar Felix gezonden','Paul appeals to Caesar':'Paulus beroept zich op de keizer','Agrippa and Bernice meet Festus':'Agrippa en Berenice bezoeken Festus','Paul goes before Agrippa':'Paulus voor Agrippa',
 'Journey to Rome begins':'Reis naar Rome begint',"Paul's Journey to Rome":"Paulus' reis naar Rome",'Paul is shipwrecked by a storm':'Paulus lijdt schipbreuk in een storm','Paul stays at Malta':'Paulus verblijft op Malta',
 'Paul arrives at Rome':'Paulus komt aan in Rome',"Paul's First Roman imprisonment":"Paulus' eerste gevangenschap in Rome",
}
pattern = re.compile(r'^(Birth|Lifetime|Death) of (.+)$')
word = {'Birth': 'Geboorte van ', 'Lifetime': 'Leven van ', 'Death': 'Dood van '}
out, missing = {}, []
for e in events:
    t = e['title']
    m = pattern.match(t) if t not in manual else None
    if m:
        out[e['id']] = word[m.group(1)] + nm(m.group(2)); continue
    if t in manual:
        out[e['id']] = manual[t]; continue
    for en, nl in prefix:
        if t.startswith(en) and t[len(en):] in kings:
            out[e['id']] = nl + kings[t[len(en):]]; break
    else:
        missing.append(t)
print(len(out), 'titles;', len(missing), 'missing', missing)
json.dump(out, open('data/titles/nl.json', 'w'), ensure_ascii=False, indent=0, sort_keys=True)
open('data/titles/nl.json', 'a').write('\n')
