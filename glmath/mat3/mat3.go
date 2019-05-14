package mat3

import (
	"math"
	"fmt"
	. "ogo/glmath/vec3"
	. "ogo/glmath/qtrnn"
	. "ogo/common"
)

type Mat3 struct {
	M [3 * 3]Double
}

func New() *Mat3 {
	m := new(Mat3)
	m.SetIdentity()
	return m
}

func (m *Mat3) Fv32() []float32 {
	fv32 := make([]float32, 9)
	for i := 0; i < 9; i++ {
		fv32[i] = float32(m.M[i])
	}
	return fv32
}

func (m *Mat3) Copy() *Mat3 {
	n := New()
	for i := 0; i < 9; i++ {
		n.M[i] = m.M[i]
	}
	return n
}

func (m *Mat3) SetIdentity() {
	m.M = [3 * 3]Double{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}
func Identity() *Mat3 {
	m := New()
	m.M = [9]Double{
		1,0,0,
		0,1,0,
		0,0,1}
	return m
}
// transpose the matrix
func (m *Mat3) Transpose() *Mat3 {	
	return &Mat3{[9]Double{
		m.M[0], m.M[3], m.M[6],
		m.M[1], m.M[4], m.M[7],
		m.M[2], m.M[5], m.M[8]}}
}

func (m *Mat3) Ptr() *Double {
	return (*Double)(&(m.M[0]))
}

func (m *Mat3) RotateLocal(angle Double, axis *Vec3) {
	*m = *m.rotate(angle, axis)
}

func (m *Mat3) RotateGlobal(angle Double, axis *Vec3) {
	axis = m.Mulv(axis)
	*m = *m.rotate(angle, axis)
}

func (m *Mat3) GetM() [3 * 3]Double {
	return m.M
}

func (m *Mat3) Right() *Vec3 {
	return V3(m.M[0], m.M[1], m.M[2])
}

func (m *Mat3) Up() *Vec3 {
	return V3(m.M[3], m.M[4], m.M[5])
}

func (m *Mat3) Forward() *Vec3 {
	return V3(m.M[6], m.M[7], m.M[8])
}

func (m *Mat3) SetRightUpForward(right, up, forward *Vec3) {
	m.M[0] = right.X
	m.M[1] = right.Y
	m.M[2] = right.Z
	m.M[3] = up.X
	m.M[4] = up.Y
	m.M[5] = up.Z
	m.M[6] = forward.X
	m.M[7] = forward.Y
	m.M[8] = forward.Z
}

// Get M rotation as Euler angles in degrees
func (m *Mat3) GetEuler() *Vec3 {
	x := (-m.M[5] / m.M[8]).Atan()
	y := (m.M[2]).Asin()
	z := (-m.M[1] / m.M[0]).Atan()

	// Convert to Degrees
	x *= 180 / math.Pi
	y *= 180 / math.Pi
	z *= 180 / math.Pi

	return V3(x, y, z)
}

// Set M rotation to Euler angles in degrees
func (m *Mat3) SetEuler(r *Vec3) {
	// Convert to Radians
	r.X *= math.Pi / 180
	r.Y *= math.Pi / 180
	r.Z *= math.Pi / 180

	m.M[0] = r.Y.Cos() * r.Z.Cos()
	m.M[1] = -r.Y.Cos() * r.Z.Sin()
	m.M[2] = r.Y.Sin()

	m.M[3] = r.X.Sin()*r.Y.Sin()*r.Z.Cos() +
		r.X.Cos()*r.Z.Sin()
	m.M[4] = -r.X.Sin()*r.Y.Sin()*r.Z.Sin() +
		r.X.Cos()*r.Z.Cos()
	m.M[5] = -r.X.Sin() * r.Y.Cos()

	m.M[6] = -r.X.Cos()*r.Y.Sin()*r.Z.Cos() +
		r.X.Sin()*r.Z.Sin()
	m.M[7] = r.X.Cos()*r.Y.Sin()*r.Z.Sin() +
		r.X.Sin()*r.Z.Cos()
	m.M[8] = r.X.Cos() * r.Y.Cos()
}

// Set M rotation to quateronion 'q'
func (m *Mat3) SetQuaternion(q *Qtrnn) {
	var xx, xy, xz, xw Double = q.X*q.X, q.X*q.Y, q.X*q.Z, q.X*q.W
	var yy, yz, yw Double = q.Y*q.Y, q.Y*q.Z, q.Y*q.W
	var zz, zw Double = q.Z*q.Z, q.Z*q.W

	m.M[0] = 1.0 - 2.0*(yy+zz)
	m.M[1] = 2.0 * (xy - zw)
	m.M[2] = 2.0 * (xz + yw)
	m.M[3] = 2.0 * (xy + zw)
	m.M[4] = 1.0 - 2.0*(xx+zz)
	m.M[5] = 2.0 * (yz - xw)
	m.M[6] = 2.0 * (xz - yw)
	m.M[7] = 2.0 * (yz + xw)
	m.M[8] = 1.0 - 2.0*(xx+yy)
}

// Multiply 'm' by 'o' and return result
func (m *Mat3) Mul(o *Mat3) *Mat3 {
	result := New()

	for row := 0; row < 3; row++ {
		ca := 3 * row
		cb := ca + 1
		cc := ca + 2

		result.M[ca] =
			m.M[ca]*o.M[0] +
				m.M[cb]*o.M[3] +
				m.M[cc]*o.M[6]

		result.M[cb] =
			m.M[ca]*o.M[1] +
				m.M[cb]*o.M[4] +
				m.M[cc]*o.M[7]

		result.M[cc] =
			m.M[ca]*o.M[2] +
				m.M[cb]*o.M[5] +
				m.M[cc]*o.M[8]
	}

	return result
}

// Multiply 'm' by 'v' and return result
func (m *Mat3) Mulv(v *Vec3) *Vec3 {
	return V3(
		v.X*m.M[0]+v.Y*m.M[1]+v.Z*m.M[2],
		v.X*m.M[3]+v.Y*m.M[4]+v.Z*m.M[5],
		v.X*m.M[6]+v.Y*m.M[7]+v.Z*m.M[8])
}

// Unexported rotate used by RotateLocal & RotateGlobal
func (m *Mat3) rotate(angle Double, axis *Vec3) *Mat3 {
	sinAngle := (angle * math.Pi / 180).Sin()
	cosAngle := (angle * math.Pi / 180).Cos()
	oneMinusCos := Double(1) - cosAngle

	axis.Normalize()

	xx := axis.X * axis.X
	yy := axis.Y * axis.Y
	zz := axis.Z * axis.Z
	xy := axis.X * axis.Y
	yz := axis.Y * axis.Z
	zx := axis.Z * axis.X
	xs := axis.X * sinAngle
	ys := axis.Y * sinAngle
	zs := axis.Z * sinAngle

	rotation := New()

	rotation.M[0] = (oneMinusCos * xx) + cosAngle
	rotation.M[1] = (oneMinusCos * xy) - zs
	rotation.M[2] = (oneMinusCos * zx) + ys

	rotation.M[3] = (oneMinusCos * xy) + zs
	rotation.M[4] = (oneMinusCos * yy) + cosAngle
	rotation.M[5] = (oneMinusCos * yz) - xs

	rotation.M[6] = (oneMinusCos * zx) - ys
	rotation.M[7] = (oneMinusCos * yz) + xs
	rotation.M[8] = (oneMinusCos * zz) + cosAngle

	return rotation.Mul(m)
}

func (m *Mat3) String() string {
	s := "["
	for i := 0; i < 9; i++ {
		s += fmt.Sprintf("%v", m.M[i])
		if i < 8 {
			s += ", "
		}
	}
	s += "]"
	return s
}
