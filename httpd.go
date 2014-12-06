package main

import (
	"io"
	"fmt"
	"net/http"
)

var savage_style =
`
path {
	stroke: #000000;
	fill-opacity: 0.05;
}
`

type vertex struct
{
	x,y float64
}

type Path struct {
	vertices []vertex
}

// Write pa to w as the inside of an SVG path string (not including quotes).
func writePathData(w io.Writer, pa *Path)  {
	vertices := pa.vertices
	if len(vertices) < 1 {
		return
	}

	pt := vertices[0]
	fmt.Fprintf(w, "M%g %g", pt.x, pt.y)

	for _, pt := range vertices[1:] {
		fmt.Fprintf(w, " L%g %g", pt.x, pt.y)
	}

	io.WriteString(w, " Z")
}

func handler(w http.ResponseWriter, r *http.Request) {
	pfad := r.URL.Path[1:]

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html>\n")
	fmt.Fprintf(w, `<img src="/savage/%s.svg" alt="WTF"></img>` + "\n", pfad)
	fmt.Fprintf(w, "</html>\n")
}

func svgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/svg+xml")
	var corner = Path {
		vertices : []vertex {
			{10, 10},
			{10, 90},
			{90, 10},
		},
	}

	io.WriteString(w, `<svg  width="1000" height="1000" xmlns="http://www.w3.org/2000/svg">` + "\n")
	io.WriteString(w, `	<style tep="text/css"> <![CDATA[`)
        io.WriteString(w, savage_style)
	io.WriteString(w, `	]]> </style>` + "\n")
        io.WriteString(w, `     <path d="`)
	writePathData(w, &corner)
	io.WriteString(w, `"/>` + "\n")
	io.WriteString(w, `</svg>` + "\n")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/savage/", svgHandler)
	http.ListenAndServe(":8080", nil)
}
