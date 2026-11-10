#!/usr/bin/env python3
"""Writes data/artwork/nl.json: Dutch titles for the engravings, keyed by the cleaned English title
(what atlas.cleanArtTitle makes of the Commons file name). Hand-translated; file names that are only
a catalogue number get a description of what the picture shows.
Run from the repository root: python3 cmd/names/artwork_nl.py"""
import json, re

enrich = json.load(open('data/enrichment.json'))

def clean(t):  # the same rules as cleanArtTitle in pkg/atlas/gallery.go
    t = t.replace('_', ' ')
    t = re.sub(r'\s*\(\d{5,}\)', '', t)
    t = re.sub(r'(?i)^(by\s+)?gustave\s+dor[eé]\s*(\([^)]*\))?\s*[,;:\-–—]*\s*', '', t)
    t = re.sub(r'(?i)\s*[-–—]\s*google art project$|\s*(retouched|retusche|cropped)\s*$', '', t)
    t = t.strip().strip('-–—,; ')
    return (t[:1].upper() + t[1:]) if t else 'Engraving'

NL = {
 '1-concubine-dore': 'De bijvrouw van de Leviet', '2 Kings 2 - 11': '2 Koningen 2: Elia opgenomen in de hemel',
 'Abraham Journeying into the Land of Canaan': 'Abraham trekt naar het land Kanaän', 'Abraham and the Three Angels': 'Abraham en de drie engelen',
 'Abraham,God and two angels': 'Abraham, God en twee engelen', 'Achan Stoned to Death': 'Achan wordt gestenigd', 'Achanstoned (cropped)': 'Achan wordt gestenigd',
 'Adam and Eve Driven out of Eden': 'Adam en Eva uit Eden verdreven', 'Adam and Eve expelled from Paradise': 'Adam en Eva uit het paradijs verdreven', 'Adamandevesin': 'De zondeval van Adam en Eva',
 "Angel decimates Sennacherib's army": 'Een engel slaat het leger van Sanherib', 'Ark of the Covenant': 'De ark van het verbond', 'Athaliah': 'Athalia', 'Athaliah1': 'Athalia',
 'Baruch': 'Baruch', "Baruch Writing Jeremiah's Prophecies": 'Baruch schrijft de profetieën van Jeremia op', 'BaruchWritingJeremiah sProphecies': 'Baruch schrijft de profetieën van Jeremia op',
 'Bibel death of absalom': 'De dood van Absalom', "Bull's head": 'Stierenkop', 'By Dore, Gustave; La Sainte Trinite': 'De Heilige Drie-eenheid', 'Dore, Gustave; La Sainte Trinite': 'De Heilige Drie-eenheid',
 'Cain and Abel Offering their Sacrifices': 'Kaïn en Abel brengen hun offers', 'Cain kills Abel': 'Kaïn doodt Abel', 'Christ Presented to the PeopleDore': 'Christus aan het volk getoond',
 'Christ on the Cross': 'Christus aan het kruis', "Christ's Entry into Jerusalem": 'Christus trekt Jeruzalem binnen', 'Christian Martyrs, 1869-1871, NGA 73471': 'De christelijke martelaren',
 'Colorized King Solomon in Old Age': 'Koning Salomo op hoge leeftijd (ingekleurd)', 'Combat between Soldiers of Ish-bosheth and David V-1': 'Strijd tussen de krijgslieden van Isboseth en David',
 'Confusion of Tongues': 'De spraakverwarring', 'Crucifixion of Jesus': 'De kruisiging van Jezus', 'Crucifixion-dore': 'De kruisiging', 'Crucifixiondarkness': 'Duisternis bij de kruisiging',
 "DESTRUCTION OF SENNACHERIB'S HOST.": 'De ondergang van het leger van Sanherib', 'David Attacks the Ammonites. Vorleser': 'David valt de Ammonieten aan', 'David-goliath28': 'David en Goliath',
 'Death of Abimelech': 'De dood van Abimelech', 'Death of Absalom': 'De dood van Absalom', 'Death of Agag, La Sainte Bible, 1866': 'De dood van Agag', 'Deluge': 'De zondvloed', 'Deluge gustave dore': 'De zondvloed',
 'DescendantsOfCain': 'De nakomelingen van Kaïn', 'Destruction of Leviathan': 'De vernietiging van de Leviathan', 'Destruction of the Army of the Amorites': 'De Amorieten verslagen',
 'Dispute of Jesus and the Pharisees': 'Jezus in discussie met de Farizeeën', 'Dore - Flight into Egypt': 'De vlucht naar Egypte', 'Dore 084': 'Bijbelillustratie van Gustave Doré',
 'Dore Death of Korah, Dathan and Abiram': 'De dood van Korach, Datan en Abiram', 'Dore DeathOfAbimelech 01': 'De dood van Abimelech', 'Dore DeathOfAbimelech 02': 'De dood van Abimelech',
 'Dore Moses Striking the Rock in Horeb': 'Mozes slaat op de rots in Horeb', 'Dore Return of the Spies From the Land of Promise': 'De verspieders keren terug uit het beloofde land',
 'Dore SaulAndTheWitchOfEndor': 'Saul en de waarzegster van En-Dor', 'Dore Solomon Proverbs': 'De spreuken van Salomo', 'Dore adultress': 'De overspelige vrouw', 'Dore jonah': 'Jona',
 'Dore jonah whale': 'Jona en de grote vis', 'Dore joshua crossing': 'Jozua trekt over de Jordaan', 'Dore joshua sun': 'Jozua gebiedt de zon stil te staan', 'Dore-SamsonSlayingTheLion': 'Simson doodt de leeuw',
 'DoreArrivalatCalvary': 'Aankomst op Golgotha', 'DoreDeathofAnanias': 'De dood van Ananias', 'DoreHagar': 'Hagar', 'DoreJesusHealingtheLunatic': 'Jezus geneest de maanziekige', 'DoreJesusPrayingintheGarden': 'Jezus bidt in de hof',
 'DoreJesusSeaGalilee': 'Jezus bij de zee van Galilea', 'DoreJohntheBaptistPreachingintheWilderness': 'Johannes de Doper preekt in de woestijn', 'DoreStPaulPreachingtotheThessalonians': 'Paulus preekt tot de Thessalonicenzen',
 'DoreTheDeadchrist': 'De dode Christus', 'Doré, Le Christ condamné': 'Christus ter dood veroordeeld', 'DoréMaryMagdaleneRepentant': 'Maria Magdalena in berouw', 'Dove Sent Forth from the Ark': 'De duif wordt uit de ark gezonden',
 'Dragon page column separator': 'Draak (sierstuk)', 'Eagle': 'Adelaar', 'El sueño de la esposa de Pilato': 'De droom van de vrouw van Pilatus', 'Elephant page column separator': 'Olifant (sierstuk)',
 'Eliezer and Rebekah': 'Eliëzer en Rebekka', 'Elijah taken up into heaven': 'Elia wordt opgenomen in de hemel', 'Expulsion of Ishmael and His Mother': 'Ismaël en zijn moeder weggestuurd',
 'Ezekiel Prophesying — Doré Bible illustration (Cassell, c.1870)': 'Ezechiël profeteert', 'Ezekiel’s Vision of the Valley of Dry Bones': 'Ezechiëls visioen van het dal vol doodsbeenderen',
 'Family drowning in flood - Dore': 'Een gezin verdrinkt in de vloed', 'Flight of Lot': 'De vlucht van Lot', 'Flood destroying the world': 'De vloed vernietigt de wereld', 'Formation of Eve': 'De schepping van Eva',
 'Gusta Dore - The Ascension': 'De hemelvaart', 'Gustave dure holy bible deluge edited': 'De zondvloed', 'Gustavo-Dore-Bibbia-Giaele-e-Sisara': 'Jaël en Sisera', 'Gustavo-Dore-Bibbia-un-angelo-appare-a-Giosue': 'Een engel verschijnt aan Jozua',
 'HealingGustaveDore': 'Jezus geneest', 'Houri on Camels': 'Vrouwen op kamelen', 'Isaac Blessing Jacob': 'Izak zegent Jakob', 'Jacob Wrestling with the Angel (cropped)': 'Jakob worstelt met de engel', 'Jacob-angel': 'Jakob en de engel',
 'Jephtha': 'Jefta', "Jephtha's doughter": 'De dochter van Jefta', 'Jesus and the disciples at the Last Supper': 'Jezus en de discipelen bij het laatste avondmaal', 'Jesus asks the Samaritan woman for a draft from the well': 'Jezus vraagt de Samaritaanse vrouw om water uit de put',
 'Jesus converses with the learned ones in the Temple': 'Jezus in gesprek met de schriftgeleerden in de tempel', 'Jesus healing the sick': 'Jezus geneest de zieken', 'Jesus is nailed to the cross': 'Jezus wordt aan het kruis genageld',
 'Jesus is scourged': 'Jezus wordt gegeseld', 'Jesus raises the daughter of Jairus from the dead': 'Jezus wekt de dochter van Jaïrus op', 'Jesus rides into Jerusalem on a donkey on Palm Sunday': 'Jezus rijdt op een ezel Jeruzalem binnen',
 'Jesus suffers agony in the garden of Gethseman': 'Jezus worstelt in de hof van Gethsemane', 'Jesus talks of the tribute money': 'Jezus spreekt over de belastingpenning', 'Jesus talks with Mary and Martha in their house': 'Jezus bij Maria en Martha thuis',
 'Jesus walks on the sea': 'Jezus wandelt over het water', 'JesusCalmingtheTempestDore': 'Jezus stilt de storm', 'JesusCuresamute': 'Jezus geneest een stomme', 'JesusFeedingMultitude': 'Jezus spijst de menigte', 'JesusPharisees': 'Jezus en de Farizeeën',
 'John the Baptist baptizes Jesus': 'Johannes de Doper doopt Jezus', 'Jonah and the whale': 'Jona en de grote vis', 'Joshua Commanding the Sun to Stand Still': 'Jozua gebiedt de zon stil te staan', 'Joshua Crossing Jordan': 'Jozua trekt over de Jordaan',
 'Judas betrays Jesus with a kiss': 'Judas verraadt Jezus met een kus', 'Judgement of Solomon': 'Het oordeel van Salomo', 'Judith and holofernes': 'Judit en Holofernes', 'Jésus guérit les malades, par Gustave Doré': 'Jezus geneest de zieken',
 'King Solomon in Old Age higher-contrast version (square cropped closeup on face)': 'Koning Salomo op hoge leeftijd (close-up)', "L'Ascension": 'De hemelvaart',
 'La Vision de Daniel': 'Daniëls visioen', 'La pêche miraculeuse de Gustave Doré': 'De wonderbare visvangst', 'Lazarus and the Rich Man': 'Lazarus en de rijke man', 'Le Christ cloué sur le croix': 'Christus aan het kruis genageld',
 'Le Prophète Daniel': 'De profeet Daniël', 'Le Prophète Ezechiel': 'De profeet Ezechiël', 'Le prophète Jérémie à la dictée, par Gustave Doré': 'De profeet Jeremia dicteert', "Le retour de l'enfant prodigue, par Gustave Doré": 'De terugkeer van de verloren zoon',
 'Lucifer3': 'Lucifer', 'Marriage at Cana engraving by Gustave Doré': 'De bruiloft te Kana', 'Menorah Candlestick': 'De menora (kandelaar)', 'Michal Gustave Doré': 'Michal', 'Michal lets David escape': 'Michal laat David ontsnappen',
 'Morte Agag': 'De dood van Agag', 'Moses dore': 'Mozes', 'Moses radiant': 'Mozes met stralend gelaat', 'Noah Cursing Canaan': 'Noach vervloekt Kanaän', 'Noah curses Ham': 'Noach vervloekt Cham', 'Noah-Curses-Ham': 'Noach vervloekt Cham',
 'Outrage at Gibeah': 'De gruweldaad te Gibea', 'Paul Addresses the Crowd After His Arrest': 'Paulus spreekt de menigte toe na zijn gevangenneming', 'Paul Addresses the Crowd After His Arrest by Gustave Doré': 'Paulus spreekt de menigte toe na zijn gevangenneming',
 'People suffering from the plague in Jerusalem. Wood engravin Wellcome V0010576': 'Mensen lijden aan de pest in Jeruzalem', 'Peter denies that he is one of Jesus’ disciples': 'Petrus ontkent dat hij een discipel van Jezus is',
 "Peter'sEscapefromPrison": 'Petrus ontsnapt uit de gevangenis', 'Prophet Amos': 'De profeet Amos', 'Rahab dore': 'Rachab', 'Raquel com cântaro': 'Rachel met een kruik', 'Rules are made to be broken II': 'Bijbelillustratie van Gustave Doré',
 'Samson tue un lion': 'Simson doodt een leeuw', 'Samson-in-DagonTemple': 'Simson in de tempel van Dagon', 'Snake page column separator': 'Slang (sierstuk)', 'Sodoma': 'Sodom', 'Sodoma - Doré': 'Sodom', 'Sodoma Doré2': 'De vlucht uit Sodom',
 'St Paul in prison': 'Paulus in de gevangenis', 'St Peter Denying Christ': 'Petrus verloochent Christus', "St. John's vision of death": 'Het visioen van Johannes over de dood', 'StPeterattheHouseofCornelius': 'Petrus in het huis van Cornelius',
 'Study for "Jacob\'s Dream" - Walters 371319': 'Studie voor "Jakobs droom"', 'Study for "The Judas Kiss" - Walters 371387': 'Studie voor "De kus van Judas"', 'The Annunciation': 'De aankondiging aan Maria',
 'The Bible (1865) - Zechariah 6-5': 'Zacharia 6: de vier wagens', 'The Birth of Jesus': 'De geboorte van Jezus', 'The Christian Martyrs by G.Dore (1871)': 'De christelijke martelaren',
 'The Dead Christ — Doré Bible illustration (Cassell, c.1870)': 'De dode Christus', 'The Death of Jezebel': 'De dood van Izebel', 'The Death of Samson': 'De dood van Simson', 'The Death of Stephen': 'De dood van Stefanus',
 'The Death of Stephen by Gustave Doré': 'De dood van Stefanus', 'The Holy Bible - Judith XIV - Judith showing the head of Holofernes - original': 'Judit toont het hoofd van Holofernes', 'The Holy Bible - Plate CXLI, The Judas Kiss': 'De kus van Judas',
 'The Holy Bible - Plate I, The Deluge (cropped)': 'De zondvloed', 'The House of Caiaphas': 'Het huis van Kajafas', 'The Levite and the concubine': 'De Leviet en zijn bijvrouw', 'The Pharisee and the publican': 'De Farizeeër en de tollenaar',
 'The Prophet Jeremiah (retusche)': 'De profeet Jeremia', 'The Prophet Jeremiah': 'De profeet Jeremia', 'The Transfiguration': 'De verheerlijking op de berg', 'The Vision of The Valley of The Dry Bones': 'Het visioen van het dal vol doodsbeenderen',
 'The beggar Lazarus begs at the gate of the rich man': 'De bedelaar Lazarus bij de poort van de rijke man', "The dream of Pilate's wife by Alphonse François": 'De droom van de vrouw van Pilatus',
 'The prodigal son decides to return to his father': 'De verloren zoon besluit naar zijn vader terug te keren', 'The walls of Jericho falling down — Doré Bible illustration (Cassell, c.1870)': 'De muren van Jericho storten in',
 'TheCrownofThornsDore': 'De doornenkroon', 'Vindication of Susanna': 'De onschuld van Susanna bewezen', 'Winged Bull': 'Gevleugelde stier', "Witch of Endor. Dore 1866": 'De waarzegster van En-Dor', 'World Destroyed by Water': 'De wereld door water vernietigd',
 'World Destroyed by Water — Doré Bible illustration (Cassell, c.1870)': 'De wereld door water vernietigd', "Zechariah's Vision of Four Chariots": 'Zacharia’s visioen van de vier wagens',
 'Ιησούς και Σαμαρείτιδα': 'Jezus en de Samaritaanse vrouw', 'Ο ιησούς και η μοιχαλίδα': 'Jezus en de overspelige vrouw', 'Гюстав Доре. Смерть Самсона': 'De dood van Simson',
}
# series: Limoges windows, Paris Ascension, Russian magazine pages, MET plates
def series(t):
    if t.startswith('Limoges'):
        return 'Glasraam in de kerk Saint-Pierre-du-Queyroix, Limoges (' + t.split()[-1] + ')'
    if t.startswith("L'Ascension"):
        n = re.search(r'- (\d+)$', t)
        return 'De hemelvaart' + (f' ({n.group(1)})' if n else '')
    if t.startswith('La Sainte Bible (vol. 2)'):
        return 'Illustratie uit La Sainte Bible, deel 2'
    if t.startswith('Живописное обозрение'):
        return 'Pagina uit het Russische tijdschrift Zhivopisnoe Obozrenie'
    if t.startswith('Рицпа'):
        return 'Rizpa verjaagt de vogels'
    return None

out, missing = {}, []
for a in enrich['art']:
    c = clean(a['title'])
    nl = NL.get(c) or series(c)
    if nl: out[c] = nl
    else: missing.append(c)
print(len(out), 'titles;', len(missing), 'missing:', missing)
json.dump(out, open('data/artwork/nl.json', 'w'), ensure_ascii=False, indent=0, sort_keys=True)
open('data/artwork/nl.json', 'a').write('\n')
