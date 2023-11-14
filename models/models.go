package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
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
2023/10/31 07:104

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
2023/10/31 07:11:041:

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
	//cleanMap["easton"] = createCleanObjectSlice(meta.Easton, meta, EastonAnyFromEntity)
	//cleanMap["events"] = createCleanObjectSlice(meta.Events, meta, EventAnyFromEntity)

	for k, v := range cleanMap {
		err := writeCleanedJson(v, k)
		if err != nil {
			return err
		}
	}

	return nil
}

type FieldType int

const (
	UnsupportedType FieldType = iota
	IntType
	FloatType
	BoolType
	StringType
	StringsType
)

func FieldTypeFromValue(v reflect.Value) FieldType {
	switch v.Type().Kind() {
	case reflect.Bool:
		return BoolType
	case reflect.Int:
		return IntType
	case reflect.Float64:
		return FloatType
	case reflect.String:
		return StringType
	case reflect.Slice:
		return StringsType
	default:
		return UnsupportedType
	}
}

func StringsFromValue(v reflect.Value) []string {
	var strs []string

	for _, elem := range v.Interface().([]interface{}) {
		if elem != nil {
			strs = append(strs, elem.(string))
		}
	}

	return strs
}

func setFieldFromEntityField(srcFieldValue reflect.Value, dstFieldValue reflect.Value) {
	srcFieldType := FieldTypeFromValue(srcFieldValue)
	dstFieldType := FieldTypeFromValue(dstFieldValue)

	switch dstFieldType {
	case IntType:
		dstFieldValue.SetInt(int64(srcFieldValue.Float()))
	case StringsType:
		dstFieldValue.Set(reflect.ValueOf(StringsFromValue(srcFieldValue)))
	case StringType:
		if srcFieldType == StringsType {
			dstFieldValue.Set(reflect.ValueOf(StringsFromValue(srcFieldValue)[0]))
		} else {
			dstFieldValue.Set(srcFieldValue)
		}
	default:
		dstFieldValue.Set(srcFieldValue)
	}
}

func setEmptyField(dstFieldValue reflect.Value) {
	dstFieldType := FieldTypeFromValue(dstFieldValue)

	switch dstFieldType {
	case StringsType:
		dstFieldValue.Set(reflect.ValueOf([]string{}))
	default:
		break
	}
}

func ConvertFromEntity(src *Entity, dst any) {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer {
		return
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}

	vt := v.Type()

	zeroValue := reflect.Value{}

	idField := v.FieldByName("ID")
	if idField == zeroValue {
		return
	}

	idField.SetString(src.ID)

	for i := 0; i < v.NumField(); i++ {
		entityFieldName := "__skip__"

		// Default use json tag
		if jsonFieldName, ok := vt.Field(i).Tag.Lookup("json"); ok {
			jsonFieldName = strings.Split(jsonFieldName, ",")[0]
			if jsonFieldName != "" {
				entityFieldName = jsonFieldName
			}
		}

		// Maybe use src tag
		if srcFieldName, ok := vt.Field(i).Tag.Lookup("src"); ok {
			srcFieldName := strings.Split(srcFieldName, ",")[0]
			if srcFieldName != "" {
				entityFieldName = srcFieldName
			}
		}

		if entityFieldName == "__skip__" {
			continue
		}

		fieldsMap := reflect.ValueOf(src.Fields)
		fieldValue := fieldsMap.MapIndex(reflect.ValueOf(entityFieldName))
		if fieldValue == zeroValue {
			setEmptyField(v.Field(i))
			continue
		}

		setFieldFromEntityField(fieldValue.Elem(), v.Field(i))
	}
}

type Book struct {
	ID           string   `json:"id" db:"id" src:"__skip__"`
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

	ConvertFromEntity(entity, book)

	var writers []string

	// Convert personLookup id's to entity id's as all other ref id's are entity id's
	for _, writeID := range book.Writers {
		persons := strings.Split(writeID, ",")
		for _, person := range persons {
			for _, entity := range metadata.People {
				personLookup := entity.Str("personLookup")
				if person == personLookup && !slices.Contains(writers, entity.ID) {
					writers = append(writers, entity.ID)
					break
				}
			}
		}
	}

	book.Writers = writers

	return book
}

func BookAnyFromEntity(entity *Entity, metadata *Metadata) any {
	return BookFromEntity(entity, metadata)
}

type Chapter struct {
	ID          string   `json:"id" db:"id" src:"__skip__"`
	Book        string   `json:"book" db:"book"` // ref TODO make single item instead of array (consistency)
	ChapterNum  int      `json:"chapterNum" db:"chapter_num"`
	OsisRef     string   `json:"osisRef" db:"osis_ref"`
	PeopleCount int      `json:"peopleCount" db:"people_count"`
	PlacesCount int      `json:"placesCount" db:"places_count"`
	Slug        string   `json:"slug" db:"slug"`
	Verses      []string `json:"verses"`               // refs
	Writers     []string `json:"writers" src:"writer"` // refs
	WriterCount int      `json:"writerCount" db:"writer_count" src:"writer count"`
}

func ChapterFromEntity(entity *Entity) *Chapter {
	chapter := &Chapter{}

	ConvertFromEntity(entity, chapter)

	return chapter
}

func ChapterAnyFromEntity(entity *Entity, _ *Metadata) any {
	return ChapterFromEntity(entity)
}

type Easton struct {
	ID           string   `json:"id" db:"id"`
	DefID        string   `json:"defId" db:"def_id"`
	DictLookup   string   `json:"dictLookup" db:"dict_lookup"`
	DictText     string   `json:"dictText" db:"dict_text"`
	HasList      bool     `json:"hasList" db:"has_list"`
	Index        int      `json:"index" db:"index"`
	ItemNum      int      `json:"itemNum" db:"item_num"`
	MatchSlugs   string   `json:"matchSlugs" db:"match_slugs"`
	MatchType    string   `json:"matchType" db:"match_type"` // [unmatched person place multi]
	PersonLookup []string `json:"personLookup"`              // refs
	PlaceLookup  []string `json:"placeLookup"`               // refs
	TermID       string   `json:"termId" db:"term_id"`
	TermLabel    string   `json:"termLabel" db:"term_label"`
}

func EastonFromEntity(entity *Entity) *Easton {
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

func EastonAnyFromEntity(entity *Entity, _ *Metadata) any {
	return EastonFromEntity(entity)
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

func EventFromEntity(entity *Entity) *Event {
	event := &Event{}

	event.ID = entity.ID
	event.LagType = entity.Str("Lag Type")
	event.SortKey = entity.Float("Sort Key")
	event.Duration = entity.Str("duration")
	event.Groups = entity.Strs("groups")
	event.Lag = entity.Str("lag")
	event.Locations = entity.Strs("locations")
	event.Notes = entity.Str("notes")
	event.PartOf = entity.Strs("partOf")
	event.Participants = entity.Strs("participants")
	event.PeopleFromVerses = entity.Strs("people (from verses)")
	event.PlacesFromVerses = entity.Strs("places (from verses)")
	event.Predecessor = entity.Strs("predecessor")
	event.RangeFlag = entity.Bool("rangeFlag")
	event.StartDate = entity.Str("startDate")
	event.Title = entity.Str("title")
	event.VerseSort = entity.Str("verseSort")
	event.Verses = entity.Strs("verses")

	return event
}

func EventAnyFromEntity(entity *Entity, _ *Metadata) any {
	return EventFromEntity(entity)
}

type People struct { // TODO
	ID                     string   `json:"id" db:"id"`
	AlphaGroup             string   `json:"alphaGroup" db:"alpha_group" src:"fields.alpha"`
	AlsoCalled             string   `json:"alsoCalled" db:"also_called"`
	Ambiguous              bool     `json:"ambiguous" db:"ambiguous"`
	BirthPlace             []string `json:"birthPlace"`
	BirthYear              []string `json:"birthYear"`
	ChaptersWritten        []string `json:"chaptersWritten"`
	Children               []string `json:"children"`
	DeathPlace             []string `json:"deathPlace"`
	DeathYear              []string `json:"deathYear"`
	DictText               []string `json:"dictText" db:"dict_text"` // TODO: convert db from and to string slice
	DictionaryLink         string   `json:"dictionaryLink" db:"dictionary_link"`
	DictionaryText         string   `json:"dictionaryText" db:"dictionary_text"`
	Disambiguation         string   `json:"disambiguation" db:"disambiguation"`
	DisplayTitle           string   `json:"displayTitle" db:"display_title"`
	Eastons                []string `json:"eastons"`
	EastonsCount           int      `json:"eastonsCount" db:"eastons_count"`
	EventGroups            []string `json:"eventGroups" db:"event_groups"` // TODO: convert db from and to string slice
	Events                 []string `json:"events"`
	Father                 string   `json:"father"` // TODO: Single reference to father people
	Gender                 string   `json:"gender" db:"gender"`
	HalfSiblingsSameFather []string `json:"halfSiblingsSameFather"`
	HalfSiblingsSameMother []string `json:"halfSiblingsSameMother"`
	IsProperName           bool     `json:"isProperName" db:"is_proper_name"`
	MaxYear                int      `json:"maxYear" db:"max_year"`
	MemberOf               []string `json:"memberOf"`
	MinYear                int      `json:"minYear" db:"min_year"`
	Mother                 string   `json:"mother"` // TODO: Single reference to mother people
	Name                   string   `json:"name" db:"name"`
	Partners               []string `json:"partners"`
	PersonID               int      `json:"personId" db:"person_id"`
	PersonLookup           string   `json:"personLookup" db:"person_lookup"`
	Siblings               []string `json:"siblings"`
	Slug                   string   `json:"slug" db:"slug"`
	Status                 string   `json:"status" db:"status"`
	Surname                string   `json:"surname" db:"surname"`
	Timeline               []string `json:"timeline"`
	VerseCount             int      `json:"verseCount" db:"verse_count"`
	Verses                 []string `json:"verses"`
}

func PeopleFromEntity(entity *Entity) *People {
	people := &People{}

	people.ID = entity.ID

	people.AlphaGroup = entity.Str("alphaGroup")
	people.AlsoCalled = entity.Str("alsoCalled")
	people.Ambiguous = entity.Bool("ambiguous")
	people.BirthPlace = entity.Strs("birthPlace")
	people.BirthYear = entity.Strs("birthYear")
	people.ChaptersWritten = entity.Strs("chaptersWritten")
	people.Children = entity.Strs("children")
	people.DeathPlace = entity.Strs("deathPlace")
	people.DeathYear = entity.Strs("deathYear")
	people.DictText = entity.Strs("dictText")
	people.DictionaryLink = entity.Str("dictionaryLink")
	people.DictionaryText = entity.Str("dictionaryText")
	people.Disambiguation = entity.Str("Disambiguation (temp)")
	/*
		2023/10/31 07:11:04 key: Disambiguation (temp)
		2023/10/31 07:11:04 is string
		2023/10/31 07:11:04 key: Easton's Count
		2023/10/31 07:11:04 is integer
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

	*/

	return people
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
