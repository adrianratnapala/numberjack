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

type xmlWriter struct {
	w io.Writer
	lineStarted, longAttr, tagOpen bool
	path []string
}

func (x* xmlWriter) fmt(f string, args ...interface{}) {
	if _, err := fmt.Fprintf(x.w, f, args...); err != nil {
		panic(err)
	}
	x.lineStarted = true
}

// writes a formatted string, but indents unless lineStarted
func (x* xmlWriter) iFmt(f string, args ...interface{}) {
	if !x.lineStarted {
		w := x.w
		for k := 0; k < len(x.path); k++ {
			if _, err := io.WriteString(w, "    "); err != nil {
				panic(err)
			}
		}
	}
	x.fmt(f, args...)
}

// Like iFmt except it writes only a single space if not right after NewLine.
func (x* xmlWriter) wsFmt(f string, args ...interface{}) {
	if !x.lineStarted {
		x.iFmt(f, args...)
		return
	}

	if _, err := io.WriteString(x.w, " "); err != nil {
		panic(err)
	}
	x.fmt(f, args...)
}


// Emit "<tag" and be ready to accept attributes, or TagDone.
func (x* xmlWriter) tagStart(tag string) {
	x.iFmt("<%s", tag)
	x.path = append(x.path, tag)
	x.tagOpen = true
}

// Writes `key="value"`
func (x* xmlWriter) attr(key string, value interface{}) {
	x.wsFmt(`%s="%v"`, key, value)
}

// Adds a newline and keeps tracking information used by the indenter
func (x* xmlWriter) newLine() {
	x.longAttr = x.tagOpen
	x.lineStarted = false
	if _, err := io.WriteString(x.w, "\n"); err != nil {
		panic(err);
	}
}

func (x* xmlWriter) pop() string {
	tags := x.path
	n := len(tags) - 1
	tag := tags[n]
	x.path = tags[:n]
	return tag
}

// Ends the current tag with `>`, unless pop, then end it with `/>`.
func (x* xmlWriter) _tagEnd(txt string)  {
	if x.longAttr {
		x.newLine()
	}

	x.tagOpen = false
	x.iFmt(txt)
	x.newLine()
}

// Ends the current tag with `>`.
func (x* xmlWriter) tagEnd() {
	x._tagEnd(">")
}

// Closes the currently element with an end tag or `/>`.
func (x* xmlWriter) tagPop() {
	if x.tagOpen {
		x._tagEnd("/>")
		x.pop()
		return
	}
	x.iFmt("</%s>", x.pop())
	x.newLine()
}

// Emit an entire <tag ../> element, `gen()` generates the attributes and body.
func (x* xmlWriter) element(tag string, gen func()) {
	x.tagStart(tag)
	gen()
	x.tagPop()
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
	defer func() {
		r := recover()
		if r != nil {
			log.Fatal(r)
		}
	}()

	x := &xmlWriter{w : os.Stdout}
	x.element("die", func() {
		x.attr("xmlns", "http://www.w3.org/2000/svg")
		x.newLine()
		x.attr("width", 1000)
		x.attr("height", 1000)
		x.tagEnd()

		x.element("quck", func() {
			x.attr("id", 14)
		})

		x.element("longer", func() {
			x.newLine()
			x.attr("id", "tweedle-dee-is-dum")
		})
	})
}

func Main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/savage/", svgHandler)
	http.ListenAndServe(":8080", nil)
}
