package models

/*


2023/10/28 10:14:28 Filename: chapters
2023/10/28 10:14:28 key: book
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: chapterNum
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: modified
2023/10/28 10:14:28 possible values: [2020-08-07T13:12:13.000Z 2020-08-07T13:12:18.000Z 2019-08-08T01:43:21.000Z 2020-08-07T13:11:53.000Z 2020-08-07T13:12:06.000Z]
2023/10/28 10:14:28 key: osisRef
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: peopleCount
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: placesCount
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: slug
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: verses
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: writer
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: writer count
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28

2023/10/28 10:14:28 Filename: easton
2023/10/28 10:14:28 key: def_id
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: dictLookup
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: dictText
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: has_list
2023/10/28 10:14:28 possible values: [False True]
2023/10/28 10:14:28 key: index
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: itemNum
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: matchSlugs
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: matchType
2023/10/28 10:14:28 possible values: [unmatched person place multi]
2023/10/28 10:14:28 key: personLookup
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: placeLookup
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: termID
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: termLabel
2023/10/28 10:14:28 is string
2023/10/28 10:14:28

2023/10/28 10:14:28 Filename: events
2023/10/28 10:14:28 key: ID
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: Lag Type
2023/10/28 10:14:28 possible values: [FS SS]
2023/10/28 10:14:28 key: Sort Key
2023/10/28 10:14:28 is float
2023/10/28 10:14:28 key: duration
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: groups
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: lag
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: locations
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: modified
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: notes
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: partOf
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: participants
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: people (from verses)
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: places (from verses)
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: predecessor
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: rangeFlag
2023/10/28 10:14:28 is bool
2023/10/28 10:14:28 key: startDate
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: title
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: verseSort
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: verses
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28

2023/10/28 10:14:28 Filename: people
2023/10/28 10:14:28 key: Disambiguation (temp)
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: Easton's Count
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: alphaGroup
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: alsoCalled
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: ambiguous
2023/10/28 10:14:28 is bool
2023/10/28 10:14:28 key: birthPlace
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: birthYear
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: chaptersWritten
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: children
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: deathPlace
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: deathYear
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: dictText
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 key: dictionaryLink
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: dictionaryText
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: displayTitle
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: eastons
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: eventGroups
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 key: events
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: father
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: gender
2023/10/28 10:14:28 possible values: [Male Female]
2023/10/28 10:14:28 key: halfSiblingsSameFather
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: halfSiblingsSameMother
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: isProperName
2023/10/28 10:14:28 is bool
2023/10/28 10:14:28 key: maxYear
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: memberOf
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: minYear
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: modified
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: mother
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: name
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: partners
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: personID
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: personLookup
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: siblings
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: slug
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: status
2023/10/28 10:14:28 possible values: [publish wip]
2023/10/28 10:14:28 key: surname
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: timeline
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: verseCount
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: verses
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28

2023/10/28 10:14:28 Filename: peopleGroups
2023/10/28 10:14:28 key: events
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: events_dev
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: groupName
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: members
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: modified
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: partOf
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: verses
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28

2023/10/28 10:14:28 Filename: periods
2023/10/28 10:14:28 key: BC-AD
2023/10/28 10:14:28 possible values: [BC AD]
2023/10/28 10:14:28 key: booksWritten
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: events
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: formattedYear
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: isoYear
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: modified
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: peopleBorn
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: peopleDied
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: yearNum
2023/10/28 10:14:28 is string
2023/10/28 10:14:28

2023/10/28 10:14:28 Filename: places
2023/10/28 10:14:28 key: aliases
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: alphaGroup
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: ambiguous
2023/10/28 10:14:28 is bool
2023/10/28 10:14:28 key: booksWritten
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: comment
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: dictText
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 key: dictionaryLink
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: dictionaryText
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: displayTitle
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: duplicate_of
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: eastons
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: esvName
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: eventGroups
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 key: eventsHere
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: featureSubType
2023/10/28 10:14:28 possible values: [Spring Well Garden Cave Tower River Gate Country]
2023/10/28 10:14:28 key: featureType
2023/10/28 10:14:28 possible values: [City Landmark Mountain Island Valley Path Water Region]
2023/10/28 10:14:28 key: hasBeenHere
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: kjvName
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: latitude
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: longitude
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: modified
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: openBibleLat
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: openBibleLong
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: peopleBorn
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: peopleDied
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: placeID
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: placeLookup
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: precision
2023/10/28 10:14:28 possible values: [Rough Unlocated Related-Surrounding Related-Within Center Precise]
2023/10/28 10:14:28 key: recogitoComments
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: recogitoLabel
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: recogitoLat
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: recogitoLon
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: recogitoStatus
2023/10/28 10:14:28 possible values: [UNVERIFIED VERIFIED NOT_IDENTIFIABLE]
2023/10/28 10:14:28 key: recogitoType
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: recogitoUID
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: recogitoUri
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: rootID
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: slug
2023/10/28 10:14:28 is string
2023/10/28 10:14:28 key: status
2023/10/28 10:14:28 possible values: [wip publish]
2023/10/28 10:14:28 key: timeline
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28 key: verseCount
2023/10/28 10:14:28 is integer
2023/10/28 10:14:28 key: verses
2023/10/28 10:14:28 string slice
2023/10/28 10:14:28 reference to objects
2023/10/28 10:14:28

2023/10/28 10:14:29 Filename: verses
2023/10/28 10:14:29 key: book
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: chapter
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: eventsDescribed
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: mdText
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: modified
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: osisRef
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: people
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: peopleCount
2023/10/28 10:14:29 is integer
2023/10/28 10:14:29 key: peopleGroups
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: places
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: placesCount
2023/10/28 10:14:29 is integer
2023/10/28 10:14:29 key: richText
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: status
2023/10/28 10:14:29 possible values: [publish wip]
2023/10/28 10:14:29 key: timeline
2023/10/28 10:14:29 string slice
2023/10/28 10:14:29 reference to objects
2023/10/28 10:14:29 key: verseID
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: verseNum
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: verseText
2023/10/28 10:14:29 is string
2023/10/28 10:14:29 key: yearNum
2023/10/28 10:14:29 is integer
2023/10/28 10:14:29
*/

type Book struct {
	ID        string `json:"id" db:"id"`
	BookDiv   string `db:"book_div"`
	BookName  string `json:"bookName" db:"book_name"`
	BookOrder int    `json:"bookOrder" db:"book_order"`

	VerseCount   int    `json:"verseCount" db:"verse_count"`
	PeopleCount  int    `db:"people_count"`
	PlaceCount   int    `db:"place_count"`
	ChapterCount int    `db:"chapter_count"`
	Testament    string `db:"testament"`
	OsisName     string `db:"osis_name"`

	ShortName string `db:"short_name"`

	Slug string `db:"slug"`

	/*
		2023/10/28 10:14:28 Filename: books
		2023/10/28 10:14:28 key: bookDiv
		2023/10/28 10:14:28 is string

		2023/10/28 10:14:28 key: bookName
		2023/10/28 10:14:28 is string

		2023/10/28 10:14:28 key: bookOrder
		2023/10/28 10:14:28 is integer

		2023/10/28 10:14:28 key: chapterCount
		2023/10/28 10:14:28 is integer

		2023/10/28 10:14:28 key: chapters
		2023/10/28 10:14:28 string slice
		2023/10/28 10:14:28 reference to objects

		2023/10/28 10:14:28 key: osisName
		2023/10/28 10:14:28 is string

		2023/10/28 10:14:28 key: peopleCount
		2023/10/28 10:14:28 is integer

		2023/10/28 10:14:28 key: placeCount
		2023/10/28 10:14:28 is integer

		2023/10/28 10:14:28 key: placeWritten
		2023/10/28 10:14:28 string slice
		2023/10/28 10:14:28 reference to objects

		2023/10/28 10:14:28 key: shortName
		2023/10/28 10:14:28 is string

		2023/10/28 10:14:28 key: slug
		2023/10/28 10:14:28 is string

		2023/10/28 10:14:28 key: testament
		2023/10/28 10:14:28 possible values: [Old Testament New Testament]

		2023/10/28 10:14:28 key: verseCount
		2023/10/28 10:14:28 is integer

		2023/10/28 10:14:28 key: verses
		2023/10/28 10:14:28 string slice
		2023/10/28 10:14:28 reference to objects

		2023/10/28 10:14:28 key: writers
		2023/10/28 10:14:28 string slice

		2023/10/28 10:14:28 key: yearWritten
		2023/10/28 10:14:28 string slice
		2023/10/28 10:14:28 reference to objects
	*/
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

type Chapter struct {
	ID          string `db:"id"`
	OsisRef     string `db:"osis_ref"`
	ChapterNum  int    `db:"chapter_num"`
	Slug        string `db:"slug"`
	PeopleCount int    `db:"people_count"`
	PlacesCount int    `db:"places_count"`
	WriterCount int    `db:"writer_count"`

	/*
		2023/10/28 09:46:48 key: verses
		2023/10/28 09:46:48 string slice
		2023/10/28 09:46:48 reference to objects
		2023/10/28 09:46:48 key: chapterNum
		2023/10/28 09:46:48 is integer
		2023/10/28 09:46:48 key: modified
		2023/10/28 09:46:48 possible values: [2019-08-08T01:43:21.000Z 2020-08-07T13:11:53.000Z 2020-08-07T13:12:06.000Z 2020-08-07T13:12:13.000Z 2020-08-07T13:12:18.000Z]
		2023/10/28 09:46:48 key: book
		2023/10/28 09:46:48 string slice
		2023/10/28 09:46:48 reference to objects
		2023/10/28 09:46:48 key: peopleCount
		2023/10/28 09:46:48 is integer
		2023/10/28 09:46:48 key: placesCount
		2023/10/28 09:46:48 is integer
		2023/10/28 09:46:48 key: writer count
		2023/10/28 09:46:48 is integer
		2023/10/28 09:46:48 key: osisRef
		2023/10/28 09:46:48 is string
		2023/10/28 09:46:48 key: slug
		2023/10/28 09:46:48 is string
		2023/10/28 09:46:48 key: writer
		2023/10/28 09:46:48 string slice
		2023/10/28 09:46:48 reference to objects
	*/
}

/*
EASTON

2023/10/28 09:47:40 key: has_list
2023/10/28 09:47:40 possible values: [False True]
2023/10/28 09:47:40 key: index
2023/10/28 09:47:40 is integer
2023/10/28 09:47:40 key: dictLookup
2023/10/28 09:47:40 is string
2023/10/28 09:47:40 key: termID
2023/10/28 09:47:40 is string
2023/10/28 09:47:40 key: termLabel
2023/10/28 09:47:40 is string
2023/10/28 09:47:40 key: matchType
2023/10/28 09:47:40 possible values: [multi unmatched person place]
2023/10/28 09:47:40 key: personLookup
2023/10/28 09:47:40 string slice
2023/10/28 09:47:40 reference to objects
2023/10/28 09:47:40 key: def_id
2023/10/28 09:47:40 is string
2023/10/28 09:47:40 key: placeLookup
2023/10/28 09:47:40 string slice
2023/10/28 09:47:40 reference to objects
2023/10/28 09:47:40 key: itemNum
2023/10/28 09:47:40 is integer
2023/10/28 09:47:40 key: dictText
2023/10/28 09:47:40 is string
2023/10/28 09:47:40 key: matchSlugs
2023/10/28 09:47:40 is string

*/

type Event struct {
	ID        string  `db:"id"`
	Title     string  `db:"title"`
	StartDate string  `db:"start_date"`
	Duration  string  `db:"duration"`
	Lag       string  `db:"lag"`
	LagType   string  `db:"lag_type"`
	VerseSort string  `db:"verse_sort"`
	SortKey   float64 `db:"sort_key"`
	// Participants
	// Locations
	// Verses
	// Predecessor
	// PartOf
	// People from verses

	/*
		2023/10/28 09:49:58 key: modified
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: rangeFlag
		2023/10/28 09:49:58 is bool

		2023/10/28 09:49:58 key: lag
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: Lag Type
		2023/10/28 09:49:58 possible values: [FS SS]

		2023/10/28 09:49:58 key: verseSort
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: startDate
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: Sort Key
		2023/10/28 09:49:58 is float

		2023/10/28 09:49:58 key: people (from verses)
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: participants
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: places (from verses)
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: locations
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: predecessor
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: verses
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: title
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: notes
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: ID
		2023/10/28 09:49:58 is integer

		2023/10/28 09:49:58 key: duration
		2023/10/28 09:49:58 is string

		2023/10/28 09:49:58 key: partOf
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects

		2023/10/28 09:49:58 key: groups
		2023/10/28 09:49:58 string slice
		2023/10/28 09:49:58 reference to objects
	*/
}

type People struct {
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

	/*
		2023/10/28 09:54:20 key: verses
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: dictionaryLink
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: siblings
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: status
		2023/10/28 09:54:20 possible values: [publish wip]

		2023/10/28 09:54:20 key: father
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: maxYear
		2023/10/28 09:54:20 is integer

		2023/10/28 09:54:20 key: dictionaryText
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: name
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: surname
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: eastons
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: verseCount
		2023/10/28 09:54:20 is integer

		2023/10/28 09:54:20 key: displayTitle
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: ambiguous
		2023/10/28 09:54:20 is bool

		2023/10/28 09:54:20 key: events
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: chaptersWritten
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: memberOf
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: deathYear
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: minYear
		2023/10/28 09:54:20 is integer

		2023/10/28 09:54:20 key: birthYear
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: alsoCalled
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: halfSiblingsSameFather
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: children
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: gender
		2023/10/28 09:54:20 possible values: [Female Male]

		2023/10/28 09:54:20 key: slug
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: partners
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: alphaGroup
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: eventGroups
		2023/10/28 09:54:20 string slice

		2023/10/28 09:54:20 key: personLookup
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: Easton's Count
		2023/10/28 09:54:20 is integer

		2023/10/28 09:54:20 key: personID
		2023/10/28 09:54:20 is integer

		2023/10/28 09:54:20 key: timeline
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: deathPlace
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: birthPlace
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: dictText
		2023/10/28 09:54:20 string slice

		2023/10/28 09:54:20 key: Disambiguation (temp)
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: mother
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects

		2023/10/28 09:54:20 key: modified
		2023/10/28 09:54:20 is string

		2023/10/28 09:54:20 key: isProperName
		2023/10/28 09:54:20 is bool

		2023/10/28 09:54:20 key: halfSiblingsSameMother
		2023/10/28 09:54:20 string slice
		2023/10/28 09:54:20 reference to objects
	*/
}

/*
peopleGroups

2023/10/28 10:00:02 key: modified
2023/10/28 10:00:02 is string

2023/10/28 10:00:02 key: verses
2023/10/28 10:00:02 string slice
2023/10/28 10:00:02 reference to objects

2023/10/28 10:00:02 key: partOf
2023/10/28 10:00:02 string slice
2023/10/28 10:00:02 reference to objects

2023/10/28 10:00:02 key: events
2023/10/28 10:00:02 string slice
2023/10/28 10:00:02 reference to objects

2023/10/28 10:00:02 key: groupName
2023/10/28 10:00:02 is string

2023/10/28 10:00:02 key: members
2023/10/28 10:00:02 string slice
2023/10/28 10:00:02 reference to objects

2023/10/28 10:00:02 key: events_dev
2023/10/28 10:00:02 string slice
2023/10/28 10:00:02 reference to objects


Periods:
2023/10/28 10:01:25 key: events
2023/10/28 10:01:25 string slice
2023/10/28 10:01:25 reference to objects

2023/10/28 10:01:25 key: formattedYear
2023/10/28 10:01:25 is string

2023/10/28 10:01:25 key: modified
2023/10/28 10:01:25 is string

2023/10/28 10:01:25 key: isoYear
2023/10/28 10:01:25 is integer

2023/10/28 10:01:25 key: BC-AD
2023/10/28 10:01:25 possible values: [BC AD]

2023/10/28 10:01:25 key: yearNum
2023/10/28 10:01:25 is string

2023/10/28 10:01:25 key: peopleBorn
2023/10/28 10:01:25 string slice
2023/10/28 10:01:25 reference to objects

2023/10/28 10:01:25 key: peopleDied
2023/10/28 10:01:25 string slice
2023/10/28 10:01:25 reference to objects

2023/10/28 10:01:25 key: booksWritten
2023/10/28 10:01:25 string slice
2023/10/28 10:01:25 reference to objects

Places:
2023/10/28 10:02:06 key: featureType
2023/10/28 10:02:06 possible values: [Water Region City Landmark Mountain Island Valley Path]

2023/10/28 10:02:06 key: longitude
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: precision
2023/10/28 10:02:06 possible values: [Related-Surrounding Related-Within Center Precise Rough Unlocated]

2023/10/28 10:02:06 key: hasBeenHere
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: ambiguous
2023/10/28 10:02:06 is bool

2023/10/28 10:02:06 key: recogitoLabel
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: displayTitle
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: dictionaryLink
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: aliases
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: status
2023/10/28 10:02:06 possible values: [wip publish]

2023/10/28 10:02:06 key: verses
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: eastons
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: placeID
2023/10/28 10:02:06 is integer

2023/10/28 10:02:06 key: recogitoUri
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: recogitoComments
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: modified
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: dictText
2023/10/28 10:02:06 string slice

2023/10/28 10:02:06 key: recogitoType
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: duplicate_of
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: featureSubType
2023/10/28 10:02:06 possible values: [Country Spring Well Garden Cave Tower River Gate]

2023/10/28 10:02:06 key: comment
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: recogitoLat
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: recogitoLon
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: eventsHere
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: slug
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: openBibleLat
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: openBibleLong
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: kjvName
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: timeline
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: eventGroups
2023/10/28 10:02:06 string slice

2023/10/28 10:02:06 key: recogitoStatus
2023/10/28 10:02:06 possible values: [UNVERIFIED VERIFIED NOT_IDENTIFIABLE]

2023/10/28 10:02:06 key: latitude
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: alphaGroup
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: placeLookup
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: booksWritten
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: recogitoUID
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: verseCount
2023/10/28 10:02:06 is integer

2023/10/28 10:02:06 key: esvName
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: rootID
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: peopleDied
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects

2023/10/28 10:02:06 key: dictionaryText
2023/10/28 10:02:06 is string

2023/10/28 10:02:06 key: peopleBorn
2023/10/28 10:02:06 string slice
2023/10/28 10:02:06 reference to objects
*/
