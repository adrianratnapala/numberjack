package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"github.com/ratnapala/numberjack"
	"github.com/ratnapala/numberjack/savage"
)

func handler(w http.ResponseWriter, r *http.Request) {
	pfad := r.URL.Path[1:]

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, `<img src="/savage/%s.svg" alt="WTF"></img>` + "\n", pfad)
	fmt.Fprintf(w, "</html>\n")
}

func svgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/svg+xml")
	savage.ThingDoc(w, thing.ExamplePath)
}

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Fatal(r)
		}
	}()
	savage.ThingDoc(os.Stdout, thing.ExamplePath)

	http.HandleFunc("/", handler)
	http.HandleFunc("/savage/", svgHandler)
	http.ListenAndServe(":8080", nil)
}
