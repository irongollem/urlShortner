package main

import (
	"database/sql"
	"github.com/irongollem/urlShortner.git/internal/controllers"
	"github.com/irongollem/urlShortner.git/internal/db"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	slite, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer slite.Close()

	if err := db.CreateTable(slite); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			controllers.ShowIndex(w, r)
		} else {
			controllers.Proxy(slite)(w, r)
		}
	})
	http.HandleFunc("/shorten", controllers.ShowShorten(slite))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
