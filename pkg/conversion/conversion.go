package conversion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/almerlucke/metabible/pkg/models"
)

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

// https://wiki.crosswire.org/OSIS_Bibles

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

func Conversion() {
	metadata, err := LoadMetadata("data/theographic-bible-metadata/json")
	if err != nil {
		log.Fatalln(err)
	}

	err = metadata.CleanJson()
	if err != nil {
		log.Fatalln(err)
	}
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

	f, err := os.Create(fmt.Sprintf("data/json/%s.json", fileName))
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

	cleanMap["books"] = createCleanObjectSlice(meta.Books, meta, BookFromEntity)
	cleanMap["chapters"] = createCleanObjectSlice(meta.Chapters, meta, ChapterFromEntity)
	cleanMap["easton"] = createCleanObjectSlice(meta.Easton, meta, EastonFromEntity)
	cleanMap["events"] = createCleanObjectSlice(meta.Events, meta, EventFromEntity)
	cleanMap["people"] = createCleanObjectSlice(meta.People, meta, PersonFromEntity)
	cleanMap["peopleGroups"] = createCleanObjectSlice(meta.PeopleGroups, meta, PeopleGroupFromEntity)
	cleanMap["periods"] = createCleanObjectSlice(meta.Periods, meta, PeriodFromEntity)
	cleanMap["places"] = createCleanObjectSlice(meta.Places, meta, PlaceFromEntity)
	cleanMap["verses"] = createCleanObjectSlice(meta.Verses, meta, VerseFromEntity)

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
			srcFieldName = strings.Split(srcFieldName, ",")[0]
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

type BookPosition struct {
	Book    string
	Chapter string
	Verse   string
}

func BookPositionFromString(s string) *BookPosition {
	var (
		book    string
		chapter string
		verse   string
	)
	parts := strings.Split(s, ".")
	if len(parts) > 0 {
		book = parts[0]
	}
	if len(parts) > 1 {
		chapter = parts[1]
	}
	if len(parts) > 2 {
		verse = parts[2]
	}
	return &BookPosition{
		Book:    book,
		Chapter: chapter,
		Verse:   verse,
	}
}

func (pos *BookPosition) IsBook() bool {
	return pos.Book != "" && pos.Chapter == "" && pos.Verse == ""
}

func (pos *BookPosition) IsChapter() bool {
	return pos.Book != "" && pos.Chapter != "" && pos.Verse == ""
}

func (pos *BookPosition) IsVerse() bool {
	return pos.Book != "" && pos.Chapter != "" && pos.Verse != ""
}

func (pos *BookPosition) String() string {
	var str string
	if pos.Book != "" {
		str += pos.Book
	}
	if pos.Chapter != "" {
		str += "." + pos.Chapter
	}
	if pos.Verse != "" {
		str += "." + pos.Verse
	}
	return str
}

func replaceReferences(text string, reg *regexp.Regexp, group int, replaceFunc func(string) string) string {
	textParts := reg.Split(text, -1)
	replacedText := textParts[0]
	indices := reg.FindAllSubmatchIndex([]byte(text), -1)
	for i, index := range indices {
		replacedText += replaceFunc(text[index[group*2]:index[group*2+1]]) + textParts[i+1]
	}
	return replacedText
}

func replaceBookPositionReferences(text string, meta *Metadata) string {
	reg := regexp.MustCompile(`/[0-9]?[a-z]+#(.+?\.[0-9]+(\.[0-9]+)?)`)

	return replaceReferences(text, reg, 1, func(str string) string {
		pos := BookPositionFromString(str)
		ref := pos.String()

		if pos.IsBook() {
			for _, book := range meta.Books {
				if book.Str("osisName") == ref {
					ref = fmt.Sprintf("/book/%s", book.ID)
					break
				}
			}
		} else if pos.IsChapter() {
			for _, chapter := range meta.Chapters {
				if chapter.Str("osisRef") == ref {
					ref = fmt.Sprintf("/chapter/%s", chapter.ID)
					break
				}
			}
		} else if pos.IsVerse() {
			for _, verse := range meta.Verses {
				if verse.Str("osisRef") == ref {
					ref = fmt.Sprintf("/verse/%s", verse.ID)
					break
				}
			}
		} else {
			log.Printf("unknown ref: %v", str)
		}

		return ref
	})
}

func BookFromEntity(entity *Entity, metadata *Metadata) any {
	book := &models.Book{}

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

func ChapterFromEntity(entity *Entity, _ *Metadata) any {
	chapter := &models.Chapter{}

	ConvertFromEntity(entity, chapter)

	return chapter
}

func EastonFromEntity(entity *Entity, _ *Metadata) any {
	easton := &models.Easton{}

	ConvertFromEntity(entity, easton)

	easton.HasList = entity.Str("has_list") == "True"

	return easton
}

func EventFromEntity(entity *Entity, _ *Metadata) any {
	event := &models.Event{}

	ConvertFromEntity(entity, event)

	return event
}

func PersonFromEntity(entity *Entity, meta *Metadata) any {
	person := &models.Person{}

	ConvertFromEntity(entity, person)

	var dictText string

	dictTextSlice := entity.Strs("dictText")
	if len(dictTextSlice) > 0 {
		dictText = dictTextSlice[0]
	}

	person.DictText = replaceBookPositionReferences(dictText, meta)

	return person
}

func PeopleGroupFromEntity(entity *Entity, _ *Metadata) any {
	group := &models.PeopleGroup{}

	ConvertFromEntity(entity, group)

	return group
}

func PeriodFromEntity(entity *Entity, _ *Metadata) any {
	period := &models.Period{}

	ConvertFromEntity(entity, period)

	return period
}

func PlaceFromEntity(entity *Entity, meta *Metadata) any {
	place := &models.Place{}

	ConvertFromEntity(entity, place)

	var dictText string

	dictTextSlice := entity.Strs("dictText")
	if len(dictTextSlice) > 0 {
		dictText = dictTextSlice[0]
	}

	place.DictText = replaceBookPositionReferences(dictText, meta)

	return place
}

func VerseFromEntity(entity *Entity, meta *Metadata) any {
	verse := &models.Verse{}

	ConvertFromEntity(entity, verse)

	verse.VerseNum, _ = strconv.Atoi(entity.Str("verseNum"))

	replaceFunc := func(ref string) string {
		parts := strings.Split(ref, "/")
		if parts[1] == "place" {
			for _, place := range meta.Places {
				if place.Str("slug") == parts[2] {
					ref = fmt.Sprintf("/place/%s", place.ID)
					break
				}
			}
		} else if parts[1] == "person" {
			for _, person := range meta.People {
				if person.Str("slug") == parts[2] {
					ref = fmt.Sprintf("/person/%s", person.ID)
					break
				}
			}
		}

		return ref
	}

	cleanReg := regexp.MustCompile(`http:/+`)
	mdText := cleanReg.ReplaceAllString(entity.Str("mdText"), "/")
	richText := cleanReg.ReplaceAllString(entity.Str("richText"), "/")

	reg := regexp.MustCompile(`\[?/.+?/.+?_[0-9]+`)
	verse.MDText = replaceReferences(mdText, reg, 0, replaceFunc)
	verse.RichText = replaceReferences(richText, reg, 0, replaceFunc)

	return verse
}
