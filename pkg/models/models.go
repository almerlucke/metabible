package models

import (
	"encoding/json"
	"io"
	"os"
)

type Book struct {
	ID           string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	BookDiv      string   `json:"bookDiv" bson:"book_div" db:"book_div"`
	BookName     string   `json:"bookName" bson:"book_name" db:"book_name"`
	BookOrder    int      `json:"bookOrder" bson:"book_order" db:"book_order"`
	ChapterCount int      `json:"chapterCount" bson:"chapter_count" db:"chapter_count"`
	Chapters     []string `json:"chapters" bson:"chapters"` // refs
	OsisName     string   `json:"osisName" bson:"osis_name" db:"osis_name"`
	PeopleCount  int      `json:"peopleCount" bson:"people_count" db:"people_count"`
	PlaceCount   int      `json:"placeCount" bson:"place_count" db:"place_count"`
	PlaceWritten []string `json:"placeWritten" bson:"place_written"` // refs
	ShortName    string   `json:"shortName" bson:"short_name" db:"short_name"`
	Slug         string   `json:"slug" bson:"slug" db:"slug"`
	Testament    string   `json:"testament" bson:"testament" db:"testament"` // [Old Testament, New Testament]
	VerseCount   int      `json:"verseCount" bson:"verse_count" db:"verse_count"`
	Verses       []string `json:"verses" bson:"verses"`            // refs
	Writers      []string `json:"writers" bson:"writers"`          // refs
	YearWritten  []string `json:"yearWritten" bson:"year_written"` // refs
}

type Chapter struct {
	ID          string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	Book        string   `json:"book" bson:"book" db:"book"` // ref
	ChapterNum  int      `json:"chapterNum" bson:"chapter_num" db:"chapter_num"`
	OsisRef     string   `json:"osisRef" bson:"osis_ref" db:"osis_ref"`
	PeopleCount int      `json:"peopleCount" bson:"people_count" db:"people_count"`
	PlacesCount int      `json:"placesCount" bson:"places_count" db:"places_count"`
	Slug        string   `json:"slug" bson:"slug" db:"slug"`
	Verses      []string `json:"verses" bson:"verses"`                // refs
	Writers     []string `json:"writers" bson:"writers" src:"writer"` // refs
	WriterCount int      `json:"writerCount" bson:"writer_count" db:"writer_count" src:"writer count"`
}

type Easton struct {
	ID           string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	DefID        string   `json:"defId" bson:"def_id" db:"def_id" src:"def_id"`
	DictLookup   string   `json:"dictLookup" bson:"dict_lookup" db:"dict_lookup"`
	DictText     string   `json:"dictText" bson:"dict_text" db:"dict_text"`
	HasList      bool     `json:"hasList" bson:"has_list" db:"has_list" src:"__skip__"`
	Index        int      `json:"index" bson:"index" db:"index"`
	ItemNum      int      `json:"itemNum" bson:"item_num" db:"item_num"`
	MatchSlugs   string   `json:"matchSlugs" bson:"match_slugs" db:"match_slugs"`
	MatchType    string   `json:"matchType" bson:"match_type" db:"match_type"` // [unmatched person place multi]
	PersonLookup []string `json:"personLookup" bson:"person_lookup"`           // refs
	PlaceLookup  []string `json:"placeLookup" bson:"place_lookup"`             // refs
	TermID       string   `json:"termId" bson:"term_id" db:"term_id" src:"termID"`
	TermLabel    string   `json:"termLabel" bson:"term_label" db:"term_label"`
}

type Event struct {
	ID               string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	Duration         string   `json:"duration" bson:"duration" db:"duration"`
	Groups           []string `json:"groups" bson:"groups"` // refs
	Lag              string   `json:"lag" bson:"lag" db:"lag"`
	LagType          string   `json:"lagType" bson:"lag_type" db:"lag_type" src:"Lag Type"` // [FS SS]
	Locations        []string `json:"locations" bson:"locations"`                           // refs
	Notes            string   `json:"notes" bson:"notes" db:"notes"`
	PartOf           []string `json:"partOf" bson:"part_of"`                                                 // refs
	Participants     []string `json:"participants" bson:"participants"`                                      // refs
	PeopleFromVerses []string `json:"peopleFromVerses" bson:"people_from_verses" src:"people (from verses)"` // refs
	PlacesFromVerses []string `json:"placesFromVerses" bson:"places_from_verses" src:"places (from verses)"` // refs
	Predecessor      []string `json:"predecessor" bson:"predecessor"`                                        // refs
	RangeFlag        bool     `json:"rangeFlag" bson:"range_flag" db:"range_flag"`
	SortKey          float64  `json:"sortKey" bson:"sort_key" db:"sort_key" src:"Sort Key"` // key: Sort Key
	StartDate        string   `json:"startDate" bson:"start_date" db:"start_date"`
	Title            string   `json:"title" bson:"title" db:"title"`
	VerseSort        string   `json:"verseSort" bson:"verse_sort" db:"verse_sort"`
	Verses           []string `json:"verses" bson:"verses"` // refs
}

type Person struct {
	ID                     string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	AlphaGroup             string   `json:"alphaGroup" bson:"alpha_group" db:"alpha_group"`
	AlsoCalled             string   `json:"alsoCalled" bson:"also_called" db:"also_called"`
	Ambiguous              bool     `json:"ambiguous" bson:"ambiguous" db:"ambiguous"`
	BirthPlace             []string `json:"birthPlace" bson:"birth_place"`
	BirthYear              []string `json:"birthYear" bson:"birth_year"`
	ChaptersWritten        []string `json:"chaptersWritten" bson:"chapters_written"`
	Children               []string `json:"children" bson:"children"`
	DeathPlace             []string `json:"deathPlace" bson:"death_place"`
	DeathYear              []string `json:"deathYear" bson:"death_year"`
	DictText               string   `json:"dictText" bson:"dict_text" db:"dict_text" src:"__skip__"`
	DictionaryLink         string   `json:"dictionaryLink" bson:"dictionary_link" db:"dictionary_link"`
	DictionaryText         string   `json:"dictionaryText" bson:"dictionary_text" db:"dictionary_text"`
	Disambiguation         string   `json:"disambiguation" bson:"disambiguation" db:"disambiguation" src:"Disambiguation (temp)"`
	DisplayTitle           string   `json:"displayTitle" bson:"display_title" db:"display_title"`
	Eastons                []string `json:"eastons" bson:"eastons"`
	EastonsCount           int      `json:"eastonsCount" bson:"eastons_count" db:"eastons_count" src:"Easton's Count"`
	EventGroups            []string `json:"eventGroups" bson:"event_groups" db:"event_groups"` // TODO: convert db from and to string slice
	Events                 []string `json:"events" bson:"events"`
	Father                 string   `json:"father" bson:"father"`
	Gender                 string   `json:"gender" bson:"gender" db:"gender"`
	HalfSiblingsSameFather []string `json:"halfSiblingsSameFather" bson:"half_siblings_same_father"`
	HalfSiblingsSameMother []string `json:"halfSiblingsSameMother" bson:"half_siblings_same_mother"`
	IsProperName           bool     `json:"isProperName" bson:"is_proper_name" db:"is_proper_name"`
	MaxYear                int      `json:"maxYear" bson:"max_year" db:"max_year"`
	MemberOf               []string `json:"memberOf" bson:"member_of"`
	MinYear                int      `json:"minYear" bson:"min_year" db:"min_year"`
	Mother                 string   `json:"mother" bson:"mother"`
	Name                   string   `json:"name" bson:"name" db:"name"`
	Partners               []string `json:"partners" bson:"partners"`
	PersonID               int      `json:"personId" bson:"person_id" db:"person_id" src:"personID"`
	PersonLookup           string   `json:"personLookup" bson:"person_lookup" db:"person_lookup"`
	Siblings               []string `json:"siblings" bson:"siblings"`
	Slug                   string   `json:"slug" bson:"slug" db:"slug"`
	Status                 string   `json:"status" bson:"status" db:"status"`
	Surname                string   `json:"surname" bson:"surname" db:"surname"`
	Timeline               []string `json:"timeline" bson:"timeline"`
	VerseCount             int      `json:"verseCount" bson:"verse_count" db:"verse_count"`
	Verses                 []string `json:"verses" bson:"verses"`
}

type PeopleGroup struct {
	ID        string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	Events    []string `json:"events" bson:"events"`
	EventsDev []string `json:"eventsDev" bson:"events_dev" src:"events_dev"`
	GroupName string   `json:"groupName" bson:"group_name" db:"group_name"`
	Members   []string `json:"members" bson:"members"`
	PartOf    []string `json:"partOf" bson:"part_of"`
	Verses    []string `json:"verses" bson:"verses"`
}

type Period struct {
	ID            string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	BcAD          string   `json:"bcAd" bson:"bc_ad" db:"bc_ad" src:"BC-AD"`
	BooksWritten  []string `json:"booksWritten" bson:"books_written"`
	Events        []string `json:"events" bson:"events"`
	FormattedYear string   `json:"formattedYear" bson:"formatted_year" db:"formatted_year"`
	ISOYear       int      `json:"isoYear" bson:"iso_year" db:"iso_year"`
	PeopleBorn    []string `json:"peopleBorn" bson:"people_born"`
	PeopleDied    []string `json:"peopleDied" bson:"people_died"`
	YearNum       string   `json:"yearNum" bson:"year_num" db:"year_num"`
}

type Place struct {
	ID               string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	Aliases          string   `json:"aliases" bson:"aliases" db:"aliases"`
	AlphaGroup       string   `json:"alphaGroup" bson:"alpha_group" db:"alpha_group"`
	Ambiguous        bool     `json:"ambiguous" bson:"ambiguous" db:"ambiguous"`
	BooksWritten     []string `json:"booksWritten" bson:"books_written"`
	Comment          string   `json:"comment" bson:"comment" db:"comment"`
	DictText         string   `json:"dictText" bson:"dict_text" db:"dict_text" src:"__skip__"`
	DictionaryLink   string   `json:"dictionaryLink" bson:"dictionary_link" db:"dictionary_link"`
	DictionaryText   string   `json:"dictionaryText" bson:"dictionary_text" db:"dictionary_text"`
	DisplayTitle     string   `json:"displayTitle" bson:"display_title" db:"display_title"`
	DuplicateOf      []string `json:"duplicateOf" bson:"duplicate_of" src:"duplicate_of"`
	Eastons          []string `json:"eastons" bson:"eastons"`
	ESVName          string   `json:"esvName" bson:"esv_name" db:"esv_name"`
	EventGroups      []string `json:"eventGroups" bson:"event_groups" db:"event_groups"` // TODO: convert slice to single string for db
	EventsHere       []string `json:"eventsHere" bson:"events_here"`
	FeatureSubType   string   `json:"featureSubType" bson:"feature_sub_type" db:"feature_sub_type"` // [Tower River Gate Country Spring Well Garden Cave]
	FeatureType      string   `json:"featureType" bson:"feature_type" db:"feature_type"`            // [Path Water Region City Landmark Mountain Island Valley]
	HasBeenHere      string   `json:"hasBeenHere"  bson:"has_been_here" db:"has_been_here"`
	KJVName          string   `json:"kjvName" bson:"kjv_name" db:"kjv_name"`
	Latitude         string   `json:"latitude" bson:"latitude" db:"latitude"`
	Longitude        string   `json:"longitude" bson:"longitude" db:"longitude"`
	OpenBibleLat     string   `json:"openBibleLat" bson:"open_bible_lat" db:"open_bible_lat"`
	OpenBibleLong    string   `json:"openBibleLong" bson:"open_bible_long" db:"open_bible_long"`
	PeopleBorn       []string `json:"peopleBorn" bson:"people_born"`
	PeopleDied       []string `json:"peopleDied" bson:"people_died"`
	PlaceID          int      `json:"placeId" bson:"place_id" db:"place_id" src:"placeID"`
	PlaceLookup      string   `json:"placeLookup" bson:"place_lookup" db:"place_lookup"`
	Precision        string   `json:"precision" bson:"precision" db:"precision"` //  [Rough Unlocated Related-Surrounding Related-Within Center Precise]
	RecogitoComments string   `json:"recogitoComments" bson:"recogito_comments" db:"recogito_comments"`
	RecogitoLabel    string   `json:"recogitoLabel" bson:"recogito_label" db:"recogito_label"`
	RecogitoLat      string   `json:"recogitoLat" bson:"recogito_lat" db:"recogito_lat"`
	RecogitoLon      string   `json:"recogitoLon" bson:"recogito_lon" db:"recogito_lon"`
	RecogitoStatus   string   `json:"recogitoStatus" bson:"recogito_status" db:"recogito_status"` //[UNVERIFIED VERIFIED NOT_IDENTIFIABLE]
	RecogitoType     string   `json:"recogitoType" bson:"recogito_type" db:"recogito_type"`
	RecogitoUID      string   `json:"recogitoUid" bson:"recogito_uid" db:"recogito_uid"`
	RecogitoURI      string   `json:"recogitoUri" bson:"recogito_uri" db:"recogito_uri"`
	RootID           []string `json:"rootId" bson:"root_id" src:"rootID"`
	Slug             string   `json:"slug" bson:"slug" db:"slug"`
	Status           string   `json:"status" bson:"status" db:"status"` // [wip publish]
	Timeline         []string `json:"timeline" bson:"timeline"`
	VerseCount       int      `json:"verseCount" bson:"verse_count" db:"verse_count"`
	Verses           []string `json:"verses" bson:"verses"`
}

func (p *Place) HasLocation() bool {
	return p.Latitude != "" && p.Longitude != ""
}

type Verse struct {
	ID              string   `json:"id" bson:"_id" db:"id" src:"__skip__"`
	Book            string   `json:"book" bson:"book" db:"book"`
	Chapter         string   `json:"chapter" bson:"chapter" db:"chapter"`
	EventsDescribed []string `json:"eventsDescribed" bson:"events_described"`
	MDText          string   `json:"mdText" bson:"md_text" db:"md_text" src:"__skip__"`
	OsisRef         string   `json:"osisRef" bson:"osis_ref" db:"osis_ref"`
	People          []string `json:"people" bson:"people"`
	PeopleCount     int      `json:"peopleCount" bson:"people_count" db:"people_count"`
	PeopleGroups    []string `json:"peopleGroups" bson:"people_groups"`
	Places          []string `json:"places" bson:"places"`
	PlacesCount     int      `json:"placesCount" bson:"places_count" db:"places_count"`
	RichText        string   `json:"richText" bson:"rich_text" db:"rich_text" src:"__skip__"`
	Status          string   `json:"status" bson:"status" db:"status"` //  [publish wip]
	Timeline        []string `json:"timeline" bson:"timeline"`
	VerseID         string   `json:"verseId" bson:"verse_id" db:"verse_id" src:"verseID"`
	VerseNum        int      `json:"verseNum" bson:"verse_num" db:"verse_num" src:"__skip__"`
	VerseText       string   `json:"verseText" bson:"verse_text" db:"verse_text"`
	YearNum         int      `json:"yearNum" bson:"year_num" db:"year_num"`
}

func (v *Verse) ShortExcerpt() string {
	excerptLen := 64

	if len(v.VerseText) > excerptLen {
		return v.VerseText[:excerptLen] + "..."
	}

	return v.VerseText
}

func Load(file string, v any) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func LoadDocs[T any](file string) ([]T, error) {
	var docs []T

	err := Load(file, &docs)
	if err != nil {
		return nil, err
	}

	return docs, err
}
