// atlasdev serves the atlas on its own port, handy for trying data variations
// (for example a different journeys/corrections/enrichment directory).
//
//	go run ./cmd/atlasdev -journeys /tmp/trial/journeys.json -addr :3001
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/almerlucke/metabible/pkg/atlas"
	"github.com/almerlucke/metabible/pkg/controllers"
)

func main() {
	journeys := flag.String("journeys", "data/journeys.json", "journeys file; corrections.json and enrichment.json are read from the same directory")
	addr := flag.String("addr", ":3001", "listen address")
	flag.Parse()

	a, err := atlas.Load("data/json", *journeys)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	controllers.MountAtlas(mux, a)
	fmt.Println("atlasdev on", *addr)
	log.Fatal(http.ListenAndServe(*addr, controllers.Middleware(mux)))
}
