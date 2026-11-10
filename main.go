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

	errs, warns := 0, 0
	for _, i := range a.Validate() {
		if i.Severity == "error" {
			errs++
		} else {
			warns++
		}
	}
	log.Printf("data check: %d errors, %d warnings (details at /api/atlas/validate)", errs, warns)

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	controllers.MountAtlas(mux, a)
	mux.HandleFunc("GET /sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Service-Worker-Allowed", "/")
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, "assets/sw.js")
	})

	fmt.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", controllers.Middleware(mux)))
}
