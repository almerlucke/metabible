// atlasdev serves only the atlas pages and APIs, without needing Mongo.
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
	fmt.Println("atlasdev on :3001")
	log.Fatal(http.ListenAndServe(":3001", mux))
}
