package conversion

import (
	"encoding/json"
	"fmt"
	"github.com/almerlucke/metabible/pkg/conversion/patterns"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

type Patterns []any

type FieldInfo map[string]Patterns

func (pats Patterns) isInteger() bool {
	for _, pat := range pats {
		f := pat.(float64)
		if f != float64(int64(f)) {
			// if f has any decimals we know that at least one value contains a real float
			return false
		}
	}

	return true
}

func (pats Patterns) investigate() {
	elem := pats[0]
	t := reflect.TypeOf(elem).Kind()
	if t == reflect.Slice {
		// get first element of any slice
		subElem := elem.([]any)[0]
		// check kind of real value of the any object
		subElemKind := reflect.ValueOf(subElem).Kind()
		if subElemKind == reflect.String {
			log.Println("string slice")
			if patterns.IsObjectReference(subElem.(string)) {
				log.Println("reference to objects")
			}
		}
	} else if t == reflect.Float64 {
		if pats.isInteger() {
			log.Println("is integer")
		} else {
			log.Println("is float")
		}
	} else if t == reflect.String {
		possibleValues := map[string]bool{}
		for _, elem := range pats {
			possibleValues[elem.(string)] = true
		}
		if len(possibleValues) < len(pats) && len(possibleValues) < 10 {
			values := []string{}
			for k, _ := range possibleValues {
				values = append(values, k)
			}
			log.Printf("possible values: %v\n", values)
		} else {
			log.Println("is string")
		}
	} else if t == reflect.Bool {
		log.Println("is bool")
	}
}

func showPatterns() {
	fileNames := []string{"books", "chapters", "easton", "events", "people", "peopleGroups", "periods", "places", "verses"}

	for _, fileName := range fileNames {
		f, err := os.Open(fmt.Sprintf("theographic-bible-metadata/json/%s.json", fileName))
		if err != nil {
			log.Fatalln(err)
		}

		defer f.Close()

		jsonData, err := io.ReadAll(f)
		if err != nil {
			log.Fatalln(err)
		}

		var entities []*Entity

		err = json.Unmarshal(jsonData, &entities)
		if err != nil {
			log.Fatalln(err)
		}

		fieldInfo := FieldInfo{}

		for _, entity := range entities {
			for k, v := range entity.Fields {
				patterns, ok := fieldInfo[k]
				if !ok {
					patterns = Patterns{}
				}
				patterns = append(patterns, v)
				fieldInfo[k] = patterns
			}
		}

		keys := []string{}
		for k, _ := range fieldInfo {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		log.Printf("Filename: %s", fileName)

		for _, k := range keys {
			log.Printf("key: %v", k)
			fieldInfo[k].investigate()
		}

		log.Printf("\n\n")
	}
}
