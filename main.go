package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/almerlucke/metabible/pkg/atlas"
	"github.com/almerlucke/metabible/pkg/controllers"
)

func main() {
	a, err := atlas.Load("data/json", "data/journeys.json")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	controllers.MountAtlas(mux, a)

	fmt.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
