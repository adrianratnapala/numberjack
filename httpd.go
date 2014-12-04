package main

import (
	"io"
	"fmt"
	"net/http"
)

var noble_tecumseh = `<svg  width="200" height="200" xmlns="http://www.w3.org/2000/svg">
	<rect
		width="150" height="150"
		style="fill:blue;stroke:pink;stroke-width:5;fill-opacity:0.1;stroke-opacity:0.9"
	/>
</svg>
`


func handler(w http.ResponseWriter, r *http.Request) {
	pfad := r.URL.Path[1:]

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, `<img src="/savage/%s.svg" alt="WTF"></img>` + "\n", pfad)
	fmt.Fprintf(w, "</html>\n")
}

func svgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/svg+xml")
	io.WriteString(w, noble_tecumseh)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/savage/", svgHandler)
	http.ListenAndServe(":8080", nil)
}
