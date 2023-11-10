package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

/*
2023/10/31 07:11:04 Filename: books
2023/10/31 07:11:04 key: bookDiv
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: bookName
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: bookOrder
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: chapterCount
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: chapters
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: osisName
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: peopleCount
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: placeCount
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: placeWritten
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: shortName
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: slug
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: testament
2023/10/31 07:11:04 possible values: [Old Testament New Testament]
2023/10/31 07:11:04 key: verseCount
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: verses
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: writers
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 key: yearWritten
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04

2023/10/31 07:11:04 Filename: chapters
2023/10/31 07:11:04 key: book
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: chapterNum
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: modified
2023/10/31 07:11:04 possible values: [2020-08-07T13:12:06.000Z 2020-08-07T13:12:13.000Z 2020-08-07T13:12:18.000Z 2019-08-08T01:43:21.000Z 2020-08-07T13:11:53.000Z]
2023/10/31 07:11:04 key: osisRef
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: peopleCount
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: placesCount
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: slug
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: verses
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: writer
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: writer count
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04

2023/10/31 07:11:04 Filename: easton
2023/10/31 07:11:04 key: def_id
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: dictLookup
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: dictText
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: has_list
2023/10/31 07:11:04 possible values: [False True]
2023/10/31 07:11:04 key: index
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: itemNum
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: matchSlugs
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: matchType
2023/10/31 07:11:04 possible values: [unmatched person place multi]
2023/10/31 07:11:04 key: personLookup
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: placeLookup
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: termID
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: termLabel
2023/10/31 07:11:04 is string
2023/10/31 07:11:04

2023/10/31 07:11:04 Filename: events
2023/10/31 07:11:04 key: ID
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: Lag Type
2023/10/31 07:11:04 possible values: [FS SS]
2023/10/31 07:11:04 key: Sort Key
2023/10/31 07:11:04 is float
2023/10/31 07:11:04 key: duration
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: groups
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: lag
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: locations
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: modified
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: notes
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: partOf
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: participants
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: people (from verses)
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: places (from verses)
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: predecessor
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: rangeFlag
2023/10/31 07:11:04 is bool
2023/10/31 07:11:04 key: startDate
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: title
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: verseSort
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: verses
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04

2023/10/31 07:11:04 Filename: people
2023/10/31 07:11:04 key: Disambiguation (temp)
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: Easton's Count
2023/10/31 07:11:04 is integer
2023/10/31 07:11:04 key: alphaGroup
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: alsoCalled
2023/10/31 07:11:04 is string
2023/10/31 07:11:04 key: ambiguous
2023/10/31 07:11:04 is bool
2023/10/31 07:11:04 key: birthPlace
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: birthYear
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: chaptersWritten
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: children
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: deathPlace
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: deathYear
2023/10/31 07:11:04 string slice
2023/10/31 07:11:04 reference to objects
2023/10/31 07:11:04 key: dictText
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 key: dictionaryLink
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: dictionaryText
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: displayTitle
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: eastons
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: eventGroups
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 key: events
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: father
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: gender
2023/10/31 07:11:05 possible values: [Male Female]
2023/10/31 07:11:05 key: halfSiblingsSameFather
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: halfSiblingsSameMother
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: isProperName
2023/10/31 07:11:05 is bool
2023/10/31 07:11:05 key: maxYear
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: memberOf
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: minYear
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: modified
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: mother
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: name
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: partners
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: personID
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: personLookup
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: siblings
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: slug
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: status
2023/10/31 07:11:05 possible values: [publish wip]
2023/10/31 07:11:05 key: surname
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: timeline
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: verseCount
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: verses
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05

2023/10/31 07:11:05 Filename: peopleGroups
2023/10/31 07:11:05 key: events
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: events_dev
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: groupName
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: members
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: modified
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: partOf
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: verses
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05

2023/10/31 07:11:05 Filename: periods
2023/10/31 07:11:05 key: BC-AD
2023/10/31 07:11:05 possible values: [AD BC]
2023/10/31 07:11:05 key: booksWritten
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: events
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: formattedYear
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: isoYear
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: modified
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: peopleBorn
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: peopleDied
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: yearNum
2023/10/31 07:11:05 is string
2023/10/31 07:11:05

2023/10/31 07:11:05 Filename: places
2023/10/31 07:11:05 key: aliases
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: alphaGroup
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: ambiguous
2023/10/31 07:11:05 is bool
2023/10/31 07:11:05 key: booksWritten
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: comment
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: dictText
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 key: dictionaryLink
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: dictionaryText
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: displayTitle
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: duplicate_of
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: eastons
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: esvName
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: eventGroups
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 key: eventsHere
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: featureSubType
2023/10/31 07:11:05 possible values: [Tower River Gate Country Spring Well Garden Cave]
2023/10/31 07:11:05 key: featureType
2023/10/31 07:11:05 possible values: [Path Water Region City Landmark Mountain Island Valley]
2023/10/31 07:11:05 key: hasBeenHere
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: kjvName
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: latitude
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: longitude
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: modified
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: openBibleLat
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: openBibleLong
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: peopleBorn
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: peopleDied
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: placeID
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: placeLookup
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: precision
2023/10/31 07:11:05 possible values: [Rough Unlocated Related-Surrounding Related-Within Center Precise]
2023/10/31 07:11:05 key: recogitoComments
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: recogitoLabel
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: recogitoLat
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: recogitoLon
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: recogitoStatus
2023/10/31 07:11:05 possible values: [UNVERIFIED VERIFIED NOT_IDENTIFIABLE]
2023/10/31 07:11:05 key: recogitoType
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: recogitoUID
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: recogitoUri
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: rootID
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: slug
2023/10/31 07:11:05 is string
2023/10/31 07:11:05 key: status
2023/10/31 07:11:05 possible values: [wip publish]
2023/10/31 07:11:05 key: timeline
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05 key: verseCount
2023/10/31 07:11:05 is integer
2023/10/31 07:11:05 key: verses
2023/10/31 07:11:05 string slice
2023/10/31 07:11:05 reference to objects
2023/10/31 07:11:05

2023/10/31 07:11:06 Filename: verses
2023/10/31 07:11:06 key: book
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: chapter
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: eventsDescribed
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: mdText
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: modified
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: osisRef
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: people
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: peopleCount
2023/10/31 07:11:06 is integer
2023/10/31 07:11:06 key: peopleGroups
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: places
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: placesCount
2023/10/31 07:11:06 is integer
2023/10/31 07:11:06 key: richText
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: status
2023/10/31 07:11:06 possible values: [publish wip]
2023/10/31 07:11:06 key: timeline
2023/10/31 07:11:06 string slice
2023/10/31 07:11:06 reference to objects
2023/10/31 07:11:06 key: verseID
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: verseNum
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: verseText
2023/10/31 07:11:06 is string
2023/10/31 07:11:06 key: yearNum
2023/10/31 07:11:06 is integer
2023/10/31 07:11:06
*/

type Entity struct {
	ID     string         `json:"id"`
	Fields map[string]any `json:"fields"`
}

func (e *Entity) Field(field string) any {
	return e.Fields[field]
}

func (e *Entity) Str(field string) string {
	value := e.Field(field)
	if value != nil {
		return value.(string)
	}

	return ""
}

func (e *Entity) Strs(field string) []string {
	value := e.Field(field)
	if value != nil {
		elems := e.Fields[field].([]any)
		strs := []string{}

		for _, elem := range elems {
			if elem != nil {
				strs = append(strs, elem.(string))
			}
		}

		return strs
	}

	return []string{}
}

func (e *Entity) Float(field string) float64 {
	value := e.Field(field)
	if value != nil {
		return e.Field(field).(float64)
	}

	return 0.0
}

func (e *Entity) Int(field string) int {
	value := e.Field(field)
	if value != nil {
		return int(e.Float(field))
	}

	return 0
}

func (e *Entity) Bool(field string) bool {
	value := e.Field(field)
	if value != nil {
		return e.Field(field).(bool)
	}

	return false
}

type Metadata struct {
	Books        []*Entity
	Chapters     []*Entity
	Easton       []*Entity
	Events       []*Entity
	People       []*Entity
	PeopleGroups []*Entity
	Periods      []*Entity
	Places       []*Entity
	Verses       []*Entity
}

func loadEntityJson(filePath string) ([]*Entity, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	jsonData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var entities []*Entity

	err = json.Unmarshal(jsonData, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func LoadMetadata(jsonDir string) (*Metadata, error) {
	fileNames := []string{"books", "chapters", "easton", "events", "people", "peopleGroups", "periods", "places", "verses"}
	metadata := &Metadata{}

	for _, fileName := range fileNames {
		entities, err := loadEntityJson(fmt.Sprintf("%s/%s.json", jsonDir, fileName))
		if err != nil {
			return nil, err
		}
		switch fileName {
		case "books":
			metadata.Books = entities
		case "chapters":
			metadata.Chapters = entities
		case "easton":
			metadata.Easton = entities
		case "events":
			metadata.Events = entities
		case "people":
			metadata.People = entities
		case "peopleGroups":
			metadata.PeopleGroups = entities
		case "periods":
			metadata.Periods = entities
		case "places":
			metadata.Places = entities
		case "verses":
			metadata.Verses = entities
		}
	}

	return metadata, nil
}

func writeCleanedJson(obj any, fileName string) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("cleaned-json/%s.json", fileName))
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func createCleanObjectSlice(entities []*Entity, meta *Metadata, f func(*Entity, *Metadata) any) []any {
	objs := make([]any, len(entities))
	for index, entity := range entities {
		objs[index] = f(entity, meta)
	}
	return objs
}

func (meta *Metadata) CleanJson() error {
	cleanMap := map[string]any{}

	cleanMap["books"] = createCleanObjectSlice(meta.Books, meta, BookAnyFromEntity)
	cleanMap["chapters"] = createCleanObjectSlice(meta.Chapters, meta, ChapterAnyFromEntity)
	cleanMap["easton"] = createCleanObjectSlice(meta.Easton, meta, EastonAnyFromEntity)

	for k, v := range cleanMap {
		err := writeCleanedJson(v, k)
		if err != nil {
			return err
		}
	}

	return nil
}

type Book struct {
	ID           string   `json:"id" db:"id"`
	BookDiv      string   `json:"bookDiv" db:"book_div"`
	BookName     string   `json:"bookName" db:"book_name"`
	BookOrder    int      `json:"bookOrder" db:"book_order"`
	ChapterCount int      `json:"chapterCount" db:"chapter_count"`
	Chapters     []string `json:"chapters"` // refs
	OsisName     string   `json:"osisName" db:"osis_name"`
	PeopleCount  int      `json:"peopleCount" db:"people_count"`
	PlaceCount   int      `json:"placeCount" db:"place_count"`
	PlaceWritten []string `json:"placeWritten"` // refs
	ShortName    string   `json:"shortName" db:"short_name"`
	Slug         string   `json:"slug" db:"slug"`
	Testament    string   `json:"testament" db:"testament"` // [Old Testament, New Testament]
	VerseCount   int      `json:"verseCount" db:"verse_count"`
	Verses       []string `json:"verses"`      // refs
	Writers      []string `json:"writers"`     // refs
	YearWritten  []string `json:"yearWritten"` // refs
}

func BookFromEntity(entity *Entity, metadata *Metadata) *Book {
	book := &Book{}

	book.ID = entity.ID
	book.BookDiv = entity.Str("bookDiv")
	book.BookName = entity.Str("bookName")
	book.BookOrder = entity.Int("bookOrder")
	book.ChapterCount = entity.Int("chapterCount")
	book.Chapters = entity.Strs("chapters")
	book.OsisName = entity.Str("osisName")
	book.PeopleCount = entity.Int("peopleCount")
	book.PlaceCount = entity.Int("placeCount")
	book.PlaceWritten = entity.Strs("placeWritten")
	book.ShortName = entity.Str("shortName")
	book.Slug = entity.Str("slug")
	book.Testament = entity.Str("testament")
	book.VerseCount = entity.Int("verseCount")
	book.Verses = entity.Strs("verses")

	writers := entity.Strs("writers")
	book.Writers = make([]string, len(writers))

	// Convert personLookup id's to entity id's as all other ref id's are entity id's
	for index, writeID := range writers {
		for _, entity := range metadata.People {
			if writeID == entity.Str("personLookup") {
				book.Writers[index] = entity.ID
			}
		}
	}

	book.YearWritten = entity.Strs("yearWritten")

	return book
}

func BookAnyFromEntity(entity *Entity, metadata *Metadata) any {
	return BookFromEntity(entity, metadata)
}

type Chapter struct {
	ID          string   `json:"id" db:"id"`
	Book        string   `json:"book" db:"book"` // ref TODO make single item instead of array (consistency)
	ChapterNum  int      `json:"chapterNum" db:"chapter_num"`
	OsisRef     string   `json:"osisRef" db:"osis_ref"`
	PeopleCount int      `json:"peopleCount" db:"people_count"`
	PlacesCount int      `json:"placesCount" db:"places_count"`
	Slug        string   `json:"slug" db:"slug"`
	Verses      []string `json:"verses"`  // refs
	Writers     []string `json:"writers"` // refs
	WriterCount int      `json:"writerCount" db:"writercount"`
}

func ChapterFromEntity(entity *Entity, metadata *Metadata) *Chapter {
	chapter := &Chapter{}

	chapter.ID = entity.ID
	chapter.Book = entity.Strs("book")[0]
	chapter.ChapterNum = entity.Int("chapterNum")
	chapter.OsisRef = entity.Str("osisRef")
	chapter.PeopleCount = entity.Int("peopleCount")
	chapter.PlacesCount = entity.Int("placesCount")
	chapter.Slug = entity.Str("slug")
	chapter.Verses = entity.Strs("verses")
	chapter.Writers = entity.Strs("writer")
	chapter.WriterCount = entity.Int("writer count")

	return chapter
}

func ChapterAnyFromEntity(entity *Entity, metadata *Metadata) any {
	return ChapterFromEntity(entity, metadata)
}

type Easton struct {
	ID           string   `json:"id" db:"id"`
	DefID        string   `json:"defId" db:"def_id"`
	DictLookup   string   `json:"dictLookup" db:"dict_lookup"`
	DictText     string   `json:"dictText" db:"dict_text"`
	HasList      bool     `json:"hasList" db:"has_list"` // TODO convert [False True] string to bool
	Index        int      `json:"index" db:"index"`
	ItemNum      int      `json:"itemNum" db:"item_num"`
	MatchSlugs   string   `json:"matchSlugs" db:"match_slugs"`
	MatchType    string   `json:"matchType" db:"match_type"` // [unmatched person place multi]
	PersonLookup []string `json:"personLookup"`              // refs
	PlaceLookup  []string `json:"placeLookup"`               // refs
	TermID       string   `json:"termId" db:"term_id"`
	TermLabel    string   `json:"termLabel" db:"term_label"`
}

func EastonFromEntity(entity *Entity, metadata *Metadata) *Easton {
	easton := &Easton{}

	easton.ID = entity.ID
	easton.DefID = entity.Str("def_id")
	easton.DictLookup = entity.Str("dictLookup")
	easton.DictText = entity.Str("dictText")
	easton.HasList = entity.Str("has_list") == "True"
	easton.Index = entity.Int("index")
	easton.ItemNum = entity.Int("itemNum")
	easton.MatchSlugs = entity.Str("matchSlugs")
	easton.MatchType = entity.Str("matchType")
	easton.PersonLookup = entity.Strs("personLookup")
	easton.PlaceLookup = entity.Strs("placeLookup")
	easton.TermID = entity.Str("termID")
	easton.TermLabel = entity.Str("termLabel")

	return easton
}

func EastonAnyFromEntity(entity *Entity, metadata *Metadata) any {
	return EastonFromEntity(entity, metadata)
}

type Event struct {
	ID               string   `json:"id" db:"id"`
	Duration         string   `json:"duration" db:"duration"`
	Groups           []string `json:"groups"` // refs
	Lag              string   `json:"lag" db:"lag"`
	LagType          string   `json:"lagType" db:"lag_type"` // [FS SS] key: Lag Type
	Locations        []string `json:"locations"`             // refs
	Notes            string   `json:"notes" db:"notes"`
	PartOf           []string `json:"partOf"`           // refs
	Participants     []string `json:"participants"`     // refs
	PeopleFromVerses []string `json:"peopleFromVerses"` // refs key: people (from verses)
	PlacesFromVerses []string `json:"placesFromVerses"` // refs key: places (from verses)
	Predecessor      []string `json:"predecessor"`      // refs
	RangeFlag        bool     `json:"rangeFlag" db:"range_flag"`
	SortKey          float64  `json:"sortKey" db:"sort_key"` // key: Sort Key
	StartDate        string   `json:"startDate" db:"start_date"`
	Title            string   `json:"title" db:"title"`
	VerseSort        string   `json:"verseSort" db:"verse_sort"`
	Verses           []string `json:"verses"` // refs
}

func EventFromEntity(entity *Entity, metadata *Metadata) *Event {
	event := &Event{}

	/*
		023/10/31 07:11:04 Filename: events
		2023/10/31 07:11:04 key: ID
		2023/10/31 07:11:04 is integer
		2023/10/31 07:11:04 key: Lag Type
		2023/10/31 07:11:04 possible values: [FS SS]
		2023/10/31 07:11:04 key: Sort Key
		2023/10/31 07:11:04 is float
		2023/10/31 07:11:04 key: duration
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: groups
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: lag
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: locations
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: modified
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: notes
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: partOf
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: participants
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: people (from verses)
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: places (from verses)
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: predecessor
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
		2023/10/31 07:11:04 key: rangeFlag
		2023/10/31 07:11:04 is bool
		2023/10/31 07:11:04 key: startDate
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: title
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: verseSort
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: verses
		2023/10/31 07:11:04 string slice
		2023/10/31 07:11:04 reference to objects
	*/

	event.ID = entity.ID

	return event
}

func EventAnyFromEntity(entity *Entity, metadata *Metadata) any {
	return EventFromEntity(entity, metadata)
}

type People struct { // TODO
	ID             string `db:"id"`
	PersonLookup   string `db:"person_lookup"`
	PersonID       int    `db:"person_id"`
	Name           string `db:"name"`
	IsProperName   bool   `db:"is_proper_name"`
	Gender         string `db:"gender"`
	DictionaryLink string `db:"dictionary_link"`
	DictionaryText string `db:"dictionary_text"`
	DisplayTitle   string `db:"display_title"`
	VerseCount     int    `db:"verse_count"`
	MinYear        int    `db:"min_year"`
	MaxYear        int    `db:"max_year"`
	AlphaGroup     string `db:"alpha_group"`
	Slug           string `db:"slug"`
	EastonsCount   int    `db:"eastons_count"`
	DictText       string `db:"dict_text"`
	// BirthYear
	// DeathYear
	// MemberOf
	// BirthPlace
	// DeathPlace
	// Verses
	// Siblings
	// Mother
	// Father
	// Children
	// Partners
	// Eastons
	// Timeline

}

type Verse struct {
	ID          string `db:"id"`
	OsisRef     string `db:"osis_ref"`
	VerseNum    int    `db:"verse_num"`
	VerseText   string `db:"verse_text"`
	YearNum     int    `db:"year_num"`
	Status      string `db:"status"`
	MDText      string `db:"md_text"`
	RichText    string `db:"rich_text"`
	VerseID     string `db:"verse_id"`
	PeopleCount int    `db:"people_count"`
	PlacesCount int    `db:"places_count"`

	// Book
	// EventsDescribed
	// People
	// Chapter
	// Timeline

	/*
		2023/10/28 10:05:50 key: richText
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: osisRef
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: places
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: verseText
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: chapter
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: peopleCount
		2023/10/28 10:05:50 is integer

		2023/10/28 10:05:50 key: placesCount
		2023/10/28 10:05:50 is integer

		2023/10/28 10:05:50 key: people
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: yearNum
		2023/10/28 10:05:50 is integer

		2023/10/28 10:05:50 key: mdText
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: verseID
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: timeline
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: eventsDescribed
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: peopleGroups
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: verseNum
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: modified
		2023/10/28 10:05:50 is string

		2023/10/28 10:05:50 key: book
		2023/10/28 10:05:50 string slice
		2023/10/28 10:05:50 reference to objects

		2023/10/28 10:05:50 key: status
		2023/10/28 10:05:50 possible values: [publish wip]
	*/
}
