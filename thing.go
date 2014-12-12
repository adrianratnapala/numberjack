package thing

type Thing interface {
	AsPath()(*Path, bool)
}

type Path struct {
	vertices [][2]float64
}

var ExamplePath = &Path {
	vertices : [][2]float64 {
		{10, 10},
		{10, 90},
		{90, 10},
	},
}

// FIX: why not just use golang type assertions?
func (p *Path) AsPath()(*Path, bool) {
	return p, p != nil
}


func (p *Path) Coords2()([][2] float64) {
	return p.vertices
}
