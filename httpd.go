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

type Thing interface {
	AsPath()(*Path, bool)
}

type vertex struct
{
	x,y float64
}

type Path struct {
	vertices []vertex
}

// FIX: why not just use golang type assertions?
func (p *Path) AsPath()(*Path, bool) {
	return p, p != nil
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
func (x* xmlWriter) _tagEnd(txt string, noNl bool)  {
	if x.longAttr {
		x.newLine()
	}

	x.tagOpen = false
	x.iFmt(txt)
	if !noNl {
		x.newLine()
        }
}

// Ends the current tag with `>`.
func (x* xmlWriter) tagEnd() {
	x._tagEnd(">", false)
}

// Closes the currently element with an end tag or `/>`.
func (x* xmlWriter) tagPop() {
	if x.tagOpen {
		x._tagEnd("/>", false)
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

// Like tagEnd, except we expect cdata for a body. `gen()` writes the desired
// bytes to it `w`.  If `hard`, then this is a true CDATA element, otherwise
// the text go straight to the output stream
func (x* xmlWriter) cdata(hard bool, gen func(w io.Writer)) {
	x._tagEnd(">", true)
	if hard {
		x.fmt("<![CDATA[")
		gen(x.w)
		x.fmt("]]>")
	} else {
		gen(x.w)
	}
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

func writeSavage(w io.Writer) {
	var corner = Path {
		vertices : []vertex {
			{10, 10},
			{10, 90},
			{90, 10},
		},
	}

	x := &xmlWriter{w : w}
	savageDoc(x, func(){savageThing(x, &corner)})
}

func savageThing(x *xmlWriter, t Thing){
	if path, ok := t.AsPath(); ok {
		x.element("path", func() {
			x.wsFmt(`d="`)
			writePathData(x.w, path)
			x.fmt(`"`)
		})
	}
}

func savageDoc(x *xmlWriter, gen func()) {
	x.element("svg", func() {
		x.attr("xmlns", "http://www.w3.org/2000/svg")
		x.newLine()
		x.attr("width", 1000)
		x.attr("height", 1000)
		x.tagEnd();

		x.element("style", func() {
			x.attr("type", "text/css")
			x.cdata(true, func (w io.Writer) {
				io.WriteString(w, savage_style)
			})
		})

		gen()
	})
}



func svgHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "image/svg+xml")
	writeSavage(w)
}

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Fatal(r)
		}
	}()
	writeSavage(os.Stdout)

	http.HandleFunc("/", handler)
	http.HandleFunc("/savage/", svgHandler)
	http.ListenAndServe(":8080", nil)
}
