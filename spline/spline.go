package spline

import (
	"github.com/banthar/gl"
	"fmt"	
	"math"
	"sort"
	"ogo/glmath"
	"ogo/glmath/mat4"
	"ogo/glmath/plane"
	. "ogo/glmath/vec2"
	. "ogo/glmath/vec3"
	"ogo/vbo"
	"ogo/program"
	"ogo/globals"
	. "ogo/common"
)

const Dtime = 0.00001

type Node struct {
	P *Vec3 //position
	V *Vec3 //tangent = velocity
	T Double
}

func NewNode(p, v *Vec3, t Double) *Node {
	return &Node{P: p.Copy(), V: v.Copy(), T: t}
}

type Spline struct {
	Nodes  []*Node
	vbo    *vbo.Vbo
	t1, t2 Double
}

func New() *Spline {
	s := &Spline{Nodes: nil, vbo: nil}
	s.t1 = Double(math.Inf(1))
	s.t2 = Double(math.Inf(-1))
	return s
}

func (s *Spline) Destroy() {
	s.Nodes = nil
	s.vbo.Destroy()
	s.vbo = nil
}

func (s *Spline) String() string {
	return fmt.Sprintf("Spline(len: %v)", len(s.Nodes))
}

func (s *Spline) StartTime() Double {
	return s.t1
}
func (s *Spline) EndTime() Double {
	return s.t2
}

func (s *Spline) AddNode(p, v *Vec3, t Double) {
	s.Nodes = append(s.Nodes, NewNode(p, v, t))
	if t > s.t2 {
		s.t2 = t
	}
	if t < s.t1 {
		s.t1 = t
	}
}
func (s *Spline) AddNode2(p, v *Vec3) {
	var t Double = 0.0
	l := len(s.Nodes)
	if l > 0 {
		t = s.Nodes[l-1].T + Dtime
	}
	s.AddNode(p, v, t)
}
func (s *Spline) AddNode3(p *Vec3) {
	s.AddNode2(p, V3(0, 1, 0))
}
func (s *Spline) Sort() {
	sort.Sort(s)
}
func (s *Spline) Len() int {return len(s.Nodes)}
func (s *Spline) Less(a,b int) bool {
	return s.Nodes[a].T < s.Nodes[b].T
}
func (s *Spline) Swap(a,b int) {
	s.Nodes[a], s.Nodes[b] = s.Nodes[b], s.Nodes[a]
}


//spline nodes have to be sorted
func (s *Spline) Normalize() {
	l := len(s.Nodes)
	if l < 2 {
		panic("too few spline nodes for normalization")
	}
	t1 := s.Nodes[0].T
	dt := s.Nodes[l-1].T - t1
	finv := 1.0 / dt
	for i := 0; i < l; i++ {
		s.Nodes[i].T -= t1
		s.Nodes[i].T *= finv
	}
	s.t1, s.t2 = 0.0, 1.0
}
func (s *Spline) MakeConsistentTimeSpacing(excludeEnds bool, t1, t2 Double) {
	l := len(s.Nodes)
	if l < 2 {
		panic("too few spline nodes for consistent time spacing")
	}
	if t1 > t2 {
		t1, t2 = t2, t1
	}
	dt := t2 - t1
	part := dt / Double(l-1)
	time := t1
	if excludeEnds {
		part = dt / Double(l-3)
		time -= part
	}
	for i := 0; i < l; i++ {
		s.Nodes[i].T = time
		time += part
	}
	s.t1, s.t2 = t1, t2
}

//auto calculation of tangents
func (s *Spline) ApplyCatmull(k Double) {
	n := s.Nodes
	l := len(n)
	if l < 2 {
		panic("too few spline nodes for catmull-rom")
	}
	n[0].V = n[1].P.Sub(n[0].P).Muls(2.0 * k)
	n[l-1].V = n[l-1].P.Sub(n[l-2].P).Muls(2.0 * k)
	for i := 1; i < l-1; i++ {
		n[i].V = n[i+1].P.Sub(n[i-1].P).Muls(k)
	}
}

func (s *Spline) NodeBefore(t Double) *Node {
	n := s.Nodes
	l := len(n)
	if l == 0 {
		panic("too few spline nodes for node before")
	}
	if l == 1 {
		return n[0]
	}
	for i := 1; i < l-1; i++ {
		if n[i].T > t {
			return n[i-1]
		}
	}
	return n[l-2]
}
func (s *Spline) NodeAfter(t Double) *Node {
	n := s.Nodes
	l := len(n)
	if l == 0 {
		panic("too few spline nodes for node after")
	}
	if l == 1 {
		return n[0]
	}
	for i := 1; i < l; i++ {
		if n[i].T > t {
			return n[i]
		}
	}
	return n[l-1]
}

//returns the interpolated position value
func (s *Spline) At(time Double) *Vec3 {
	a := s.NodeBefore(time)
	b := s.NodeAfter(time)
	dt := b.T - a.T
	t := (time - a.T) / dt
	return glmath.InterpolateHermiteVec3(a.P, a.V, b.P, b.V, t)
}

//returns the interpolated value of the first derivative
func (s *Spline) AtD1(time Double) *Vec3 {
	a := s.NodeBefore(time)
	b := s.NodeAfter(time)
	dt := b.T - a.T
	t := (time - a.T) / dt
	return glmath.InterpolateHermiteD1Vec3(a.P, a.V, b.P, b.V, t)
}

func (s *Spline) CreateVbo(nsteps, ncorners int, radius Double) {
	s.CreateVbo_(nsteps, ncorners, radius, false, V2(0,0),V2(0,0))
}
func (s *Spline) CreateVboTex(nsteps, ncorners int, radius Double,
	tex1, tex2 *Vec2) {
	s.CreateVbo_(nsteps, ncorners, radius, true, tex1, tex2)
}
func (s *Spline) CreateVbo_(nsteps, ncorners int, radius Double,
	usetex bool, tex1, tex2 *Vec2) {
	if ncorners < 3 {
		panic("too few corners")
	}
	nump := (ncorners+1) * (nsteps + 1)
	numn := nump
	numt := nump
	p := make([]*Vec3, 0, nump)
	n := make([]*Vec3, 0, numn)
	var t, tcircle []*Vec2
	if usetex {
		t = make([]*Vec2, 0, numt)
		tcircle = make([]*Vec2, 0, ncorners)
	}
	fixedaxis := V3(0.00001, 1.0, -.00002).Unit()
	angle := -2.0 * math.Pi / Double(ncorners)
	lastDir0, lastPos := V3(1,0,0), V3(0,0,0)
	processCircle := func(pos, normal *Vec3, step int) {
		dir0 := normal.Perp(fixedaxis)
		//correct axis swift: shortest distance between this circle's dir0 and the last circle's
		if step > 0 {
			dir0 = plane.ByPointAndNormal(pos, normal).NearestPoint(lastPos.Add(lastDir0)).Sub(pos).Unit()
		}
		lastDir0 = dir0
		lastPos = pos
		texy := (tex2.Y-tex1.Y)*Double(step)/Double(nsteps) + tex1.Y
		for i := 0; i < ncorners+1; i++ {
			mrot := mat4.RotationAxis(normal, Double(i)*angle)
			dir := mrot.Mulv(dir0)
			p = append(p, pos.Add(dir.Muls(radius)))
			n = append(n, dir)
			if usetex {
				texx := (tex2.X-tex1.X)*Double(i)/Double(ncorners) + tex1.X
				t = append(t, V2(texx, texy))
			}
		}
	}
	t1, t2 := s.StartTime(), s.EndTime()
	dt := t2 - t1
	part := dt / Double(nsteps)
	time := t1
	for i := 0; i < nsteps+1; i++ {
		pos := s.At(time)
		vel := s.AtD1(time)
		processCircle(pos, vel.Unit(), i)
		time += part
	}
	//front/back texCoords:
	if usetex {
		clockhand := V3(0, 1, 0)
		for i := 0; i < ncorners; i++ {
			mrot := mat4.RotationAxis(V3(0, 0, -1), Double(i)*angle)
			dir := mrot.Mulv(clockhand)
			tcircle = append(tcircle, V2(dir.X, dir.Y).Add(V2(1, 1)).Muls(0.5))
		}
	}
	//filling vbo:
	s.vbo = vbo.New()
	v := s.vbo
	numi := (2*3*nsteps * (ncorners+1)) + 2*(ncorners-2)
	v.P = make([]*Vec3, 0, numi)
	v.N = make([]*Vec3, 0, numi)
	if usetex {
		v.T = make([]*Vec2, 0, numi)
	}
	//front + back
	fillVboWithCircle := func(index int, ccw bool) {
		normal := s.AtD1(t1).Muls(-1.0) //front face normal
		if !ccw {
			normal = s.AtD1(t2).Muls(1.0)  //back face normal
		}
		doIndex := func(i1, i2 int) {
			v.P = append(v.P, p[index],p[index+i1],p[index+i2])
			v.N = append(v.N, normal, normal, normal)
			if usetex {
				v.T = append(v.T, tcircle[0], tcircle[i1], tcircle[i2])
			}
		}
		for i := 1; i < ncorners-1; i++ {
			if ccw {
				doIndex(i+1, i)
			} else {
				doIndex(i, i+1)
			}
		}
	}
	fillVboWithCircle(0, true)
	fillVboWithCircle(len(p)-(ncorners+1), false)
	
	appendVertexTriangle := func(a, b, c int) {
		v.P = append(v.P, p[a], p[b], p[c])
		v.N = append(v.N, n[a], n[b], n[c])
		if usetex {
			v.T = append(v.T, t[a], t[b], t[c])
		}
	}
	fillVboWithCylinder := func(index int) {
		for i := 0; i < ncorners-1; i++ {
			appendVertexTriangle(index+i+0, index+i+1, index+i+0+ncorners+1)
			appendVertexTriangle(index+i+1+ncorners+1, index+i+0+ncorners+1, index+i+1)
		}
		//last, closing rectangle
		appendVertexTriangle(index+ncorners-1, index+ncorners+0, index+2*ncorners+0)
		appendVertexTriangle(index+2*ncorners+1, index+2*ncorners+0, index+ncorners+0)
	}
	for i := 0; i < nsteps; i++ {
		fillVboWithCylinder(i * (ncorners+1))
	}
	s.vbo.Create(gl.STATIC_DRAW)
}

func (s *Spline) Render() {
	if globals.UseShader {
		program.SetModelMatrix(mat4.Identity())
	}
	s.vbo.Enable()
	s.vbo.Draw(gl.TRIANGLES)
	s.vbo.Disable()
}



