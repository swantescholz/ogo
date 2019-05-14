package nurbs

import (
	"github.com/banthar/gl"
	"fmt"
	"ogo/glmath"
	"ogo/glmath/mat4"
	. "ogo/glmath/vec2"
	. "ogo/glmath/vec3"
	"ogo/vbo"
	"ogo/model"
	"ogo/program"
	"ogo/globals"
	. "ogo/common"
)

type Node struct {
	P *Vec3 //position
	Tu,Tv *Vec3 //tangent u, tangent v
	Tc *Vec2 //tex coord
	U,V Double // [0...1]
}

func NewNode(p, tu, tv *Vec3, tc *Vec2, u,v Double) *Node {
	return &Node{p, tu,tv, tc, u,v}
}
func (n *Node) String() string {
	return fmt.Sprintf("NurbNode(\n\tP: %v\n\tTu: %v\n\tTv: %v\n\tTc: %v\n\tU: %v\n\tV: %v)",
	n.P,n.Tu,n.Tv,n.Tc,n.U,n.V)
}
func (n *Node) InterpolateU(o *Node, t Double) *Node {
	p := glmath.InterpolateHermiteVec3(n.P,n.Tu, o.P,o.Tu, t)
	tu := glmath.InterpolateHermiteD1Vec3(n.P,n.Tu, o.P,o.Tu, t)
	tv := n.Tv.Interpolate(o.Tv, t)
	tc := n.Tc.Interpolate(o.Tc, t)
	u,v := n.U+(o.U-n.U)*t, n.V+(o.V-n.V)*t
	return NewNode(p,tu,tv,tc,u,v)
}
func (n *Node) InterpolateV(o *Node, t Double) *Node {
	p := glmath.InterpolateHermiteVec3(n.P,n.Tv, o.P,o.Tv, t)
	tu := n.Tu.Interpolate(o.Tu, t)
	tv := glmath.InterpolateHermiteD1Vec3(n.P,n.Tv,o.P,o.Tv,t)
	tc := n.Tc.Interpolate(o.Tc, t)
	u,v := n.U+(o.U-n.U)*t, n.V+(o.V-n.V)*t
	return NewNode(p,tu,tv,tc,u,v)
}
func (n *Node) InterpolateBilinear(b,c,d *Node, u,v Double) *Node {
	x := n.InterpolateU(b, u)
	y := c.InterpolateU(d, u)
	node := x.InterpolateV(y, v)
	return node
}
func (n *Node) Normal() *Vec3 {
	return n.Tu.Cross(n.Tv).Unit()
}
func (n *Node) Acceleration(Fg *Vec3) *Vec3 {
	normal := n.Normal()
	return Fg.Perp(normal).Muls(Fg.Length()*(normal.Muls(-1.0).Angle(Fg)).Sin())
}

type Nurbs struct {
	Nodes  [][]*Node
	vbo    *vbo.Vbo
	u1,u2,v1,v2 Double
	nu,nv int
}

func New(nu, nv int) *Nurbs {
	t := &Nurbs{Nodes: nil, vbo: nil}
	t.nu, t.nv = nu,nv
	east := V3(1,0,0)
	north := V3(0,0,-1)
	t.Nodes = make([][]*Node, t.nu)
	for u := range t.Nodes {
		t.Nodes[u] = make([]*Node, t.nv)
		for v := range t.Nodes[u] {
			t.Nodes[u][v] = NewNode(V3(0,0,0), east, north, V2(0,0),0.0,0.0)
		}
	}
	t.SetUVBounds(0.0,1.0,0.0,1.0)
	t.ArrangeGrid(V3(0,0,0),V3(1,0,0),V3(0,0,-1),V3(1,0,-1))
	t.SetTexCoords(V2(0,0),V2(1,0),V2(0,1),V2(1,1))
	return t
}

func (t *Nurbs) Destroy() {
	t.Nodes = nil
	t.vbo.Destroy()
	t.vbo = nil
}

func (t *Nurbs) ForEachNode(f func(n *Node)) {
	for u := 0; u < len(t.Nodes); u++ {
		for v := 0; v < len(t.Nodes[u]); v++ {
			f(t.Nodes[u][v])
		}
	}
}

func (t *Nurbs) SetUVBounds(u1,u2,v1,v2 Double) {
	t.u1, t.u2, t.v1, t.v2 = u1,u2,v1,v2
}
func (t *Nurbs) SetTexCoords(a,b,c,d *Vec2) {
	for u := 0; u < len(t.Nodes); u++ {
		for v := 0; v < len(t.Nodes[u]); v++ {
			fu,fv := t.Nodes[u][v].U, t.Nodes[u][v].V
			t.Nodes[u][v].Tc = glmath.InterpolateBilinearVec2(a,b,c,d, fu,fv)
		}
	}
}
func (t *Nurbs) ArrangeGrid(a,b,c,d *Vec3) {
	for u := 0; u < len(t.Nodes); u++ {
		for v := 0; v < len(t.Nodes[u]); v++ {
			fu := Double(u-1) / Double(t.nu-3)
			fv := Double(v-1) / Double(t.nv-3)
			t.Nodes[u][v].U, t.Nodes[u][v].V = fu,fv
			t.Nodes[u][v].P = glmath.InterpolateBilinearVec3(a,b,c,d, fu,fv)
		}
	}
}

func (t *Nurbs) AutoCalcTangents(k Double) {
	//center
	for u := 1; u < len(t.Nodes)-1; u++ {
		for v := 1; v < len(t.Nodes[u])-1; v++ {
			t.Nodes[u][v].Tu = t.Nodes[u+1][v].P.Sub(t.Nodes[u-1][v].P).Muls(k)
			t.Nodes[u][v].Tv = t.Nodes[u][v+1].P.Sub(t.Nodes[u][v-1].P).Muls(k)
		}
	}
	setTuv := func(u,v, du1,dv1, du2,dv2 int, fac1, fac2 Double) {
		t.Nodes[u][v].Tu = t.Nodes[u+du1][v+dv1].P.Sub(t.Nodes[u][v].P).Muls(fac1)
		t.Nodes[u][v].Tv = t.Nodes[u+du2][v+dv2].P.Sub(t.Nodes[u][v].P).Muls(fac2)
	}
	//left/right side
	for v := 1; v < t.nv-1; v++ {
		setTuv(0     ,v,  1,0, 0,1,  2*k, 2*k)
		setTuv(t.nu-1,v, -1,0, 0,1, -2*k, 2*k)
	}
	//lower/upper side
	for u := 1; u < t.nu-1; u++ {
		setTuv(u,0     , 1,0, 0, 1, 2*k,  2*k)
		setTuv(u,t.nv-1, 1,0, 0,-1, 2*k, -2*k)
	}
	//corners:
	setTuv(     0,      0,  1,0, 0, 1,  2*k,  2*k)
	setTuv(t.nu-1,      0, -1,0, 0, 1, -2*k,  2*k)
	setTuv(     0, t.nv-1,  1,0, 0,-1,  2*k, -2*k)
	setTuv(t.nu-1, t.nv-1, -1,0, 0,-1, -2*k, -2*k)
}

func (t *Nurbs) InterpolateP(u,v Double) *Vec3 {
	return t.Interpolate(u,v).P
}
func (t *Nurbs) InterpolateN(u,v Double) *Vec3 {
	node := t.Interpolate(u,v)
	return node.Normal()
}
func (t *Nurbs) NodesAround(u,v Double) (a,b,c,d *Node) {
	//u = math.Max(math.Min(u, 0.99999),0.0)
	//v = math.Max(math.Min(v, 0.99999),0.0)
	x := int(1+Double(t.nu-3)*u)
	y := int(1+Double(t.nv-3)*v)
	a = t.Nodes[x  ][y  ]
	b = t.Nodes[x+1][y  ]
	c = t.Nodes[x  ][y+1]
	d = t.Nodes[x+1][y+1]
	return
}
func (t *Nurbs) ConvertCoords(u,v Double) (newu,newv Double) {
	newu = (u-t.u1)/(t.u2-t.u1)
	newv = (v-t.v1)/(t.v2-t.v1)
	return
}
//with conversion
func (t *Nurbs) Interpolate(u,v Double) *Node {
	u,v = t.ConvertCoords(u,v)
	return t.Interpolate_(u,v)
}
//without conversion
func (t *Nurbs) Interpolate_(u,v Double) *Node {
	a,b,c,d := t.NodesAround(u,v)
	u1 := (a.U+c.U)*.5
	u2 := (b.U+d.U)*.5
	v1 := (a.V+b.V)*.5
	v2 := (c.V+d.V)*.5
	u = (u-u1)/(u2-u1)
	v = (v-v1)/(v2-v1)
	return a.InterpolateBilinear(b,c,d, u,v)
}

func (t *Nurbs) CreateVbo(nusteps, nvsteps int, usetex bool) {
	//getting vertex data
	nump := (nusteps+1)*(nvsteps+1)
	p := make([]*Vec3, 0, nump)
	n := make([]*Vec3, 0, nump)
	var tc []*Vec2
	if usetex {
		tc = make([]*Vec2, 0, nump)
	}
	ustep := 1.0/Double(nusteps)
	vstep := 1.0/Double(nvsteps)
	for u := 0; u <= nusteps; u++ {
		for v := 0; v <= nvsteps; v++ {
			fu, fv := Double(u)*ustep, Double(v)*vstep
			node := t.Interpolate_(fu,fv)
			pos, normal := node.P.Copy(), node.Normal()
			p = append(p, pos)
			n = append(n, normal)
			if usetex {
				tc = append(tc, node.Tc)
			}
		}
	}
	//filling vbo
	t.vbo = vbo.New()
	v := t.vbo
	numi := (2*3*nusteps*nvsteps)
	v.P = make([]*Vec3, 0, numi)
	v.N = make([]*Vec3, 0, numi)
	if usetex {
		v.T = make([]*Vec2, 0, numi)
	}
	appendVertex := func(index int) {
		v.P = append(v.P, p[index])
		v.N = append(v.N, n[index])
		if usetex {
			v.T = append(v.T, tc[index])
		}
	}
	appendQuad := func(x,y int) {
		a := x*(nvsteps+1) + y
		b,c,d := a+nusteps+1, a+1, a+nusteps+2
		appendVertex(a)
		appendVertex(b)
		appendVertex(d)
		appendVertex(a)
		appendVertex(d)
		appendVertex(c)
	}
	for x := 0; x < nusteps; x++ {
		for y := 0; y < nvsteps; y++ {
			appendQuad(x,y)
		}
	}
	t.vbo.Create(gl.STATIC_DRAW)
}

func (t *Nurbs) Render() {
	t.vbo.Enable()
	if globals.UseShader {
		program.SetModelMatrix(mat4.Identity())
		t.vbo.Draw(gl.TRIANGLES)
	} else {
		gl.PushMatrix()
		gl.MultMatrixf(mat4.Identity().Ptr32())
		t.vbo.Draw(gl.TRIANGLES)
		gl.PopMatrix()
	}
	t.vbo.Disable()
}
func (t *Nurbs) DebugRender() {
	t.Render()
	m := model.Get("cube")
	t.ForEachNode(func(n *Node) {
		m.Mat = mat4.Scaling(V3(1,1,1).Muls(1.1))
		m.Mat = m.Mat.Mul(mat4.Translation(n.P))
		m.Render()
		model.RenderArrow(n.P, n.Tu.Muls(1.0))
		model.RenderArrow(n.P, n.Tv.Muls(1.0))
	})
}

func (t *Nurbs) String() string {
	return fmt.Sprintf("Nurbs(len: %v x %v)", t.nu, t.nv)
}

