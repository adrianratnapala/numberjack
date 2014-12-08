package main

import (
	"io"
	"fmt"
	"net/http"
	"os"
	"log"
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

type svgWriter struct {
	open_tags []string
	w io.Writer
}

func (s *svgWriter) push(tag string) error {
	s.open_tags = append(s.open_tags, tag)
	_, err := fmt.Fprintf(s.w, "<%s>", tag)
	return err
}

func (s *svgWriter) pop() error {
	tags := s.open_tags

	n := len(tags) - 1
	if n < 0  {
		panic(fmt.Errorf("%v: too many pop()s"))
	}

	if _, err := fmt.Fprintf(s.w, "</%s>", tags[n]); err != nil {
		return err
	}

	s.open_tags = tags[:n]

	return nil
}

func (s *svgWriter) End() error {
	tags := s.open_tags
	//log.Printf("End: tags = %v", tags)

	switch n := len(tags); {
	case n <= 0: return nil
	case n > 1: panic(fmt.Errorf(
		"%v: unexpected End(), while %s is open.",
		s, tags[n-1]))
	}

	return s.pop();
}


func newSvgWriter(w io.Writer) (*svgWriter, error) {
	s := &svgWriter { w :w }
	s.push("svg")
	return s, nil
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
	{
		sw, err := newSvgWriter(os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
		sw.End()
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/savage/", svgHandler)
	http.ListenAndServe(":8080", nil)
}
