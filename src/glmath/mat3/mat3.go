package mat3

import "math"
import "fmt"
import . "glmath/vec3"
import . "glmath/qtrnn"

type Mat3 struct {
	M [3 * 3]float64
}

func NewMat3() *Mat3 {
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
	n := NewMat3()
	for i := 0; i < 9; i++ {
		n.M[i] = m.M[i]
	}
	return n
}

func (m *Mat3) SetIdentity() {
	m.M = [3 * 3]float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}

func (m *Mat3) Ptr() *float64 {
	return (*float64)(&(m.M[0]))
}

func (m *Mat3) RotateLocal(angle float64, axis *Vec3) {
	*m = *m.rotate(angle, axis)
}

func (m *Mat3) RotateGlobal(angle float64, axis *Vec3) {
	axis = m.Mulv(axis)
	*m = *m.rotate(angle, axis)
}

func (m *Mat3) GetM() [3 * 3]float64 {
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
	x := math.Atan((-m.M[5]) / m.M[8])
	y := math.Asin(m.M[2])
	z := math.Atan((-m.M[1]) / m.M[0])

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

	m.M[0] = math.Cos(r.Y) * math.Cos(r.Z)
	m.M[1] = -math.Cos(r.Y) * math.Sin(r.Z)
	m.M[2] = math.Sin(r.Y)

	m.M[3] = math.Sin(r.X)*math.Sin(r.Y)*math.Cos(r.Z) +
		math.Cos(r.X)*math.Sin(r.Z)
	m.M[4] = -math.Sin(r.X)*math.Sin(r.Y)*math.Sin(r.Z) +
		math.Cos(r.X)*math.Cos(r.Z)
	m.M[5] = -math.Sin(r.X) * math.Cos(r.Y)

	m.M[6] = -math.Cos(r.X)*math.Sin(r.Y)*math.Cos(r.Z) +
		math.Sin(r.X)*math.Sin(r.Z)
	m.M[7] = math.Cos(r.X)*math.Sin(r.Y)*math.Sin(r.Z) +
		math.Sin(r.X)*math.Cos(r.Z)
	m.M[8] = math.Cos(r.X) * math.Cos(r.Y)
}

// Set M rotation to quateronion 'q'
func (m *Mat3) SetQuaternion(q *Qtrnn) {
	xx, xy, xz, xw := q.X*q.X, q.X*q.Y, q.X*q.Z, q.X*q.W
	yy, yz, yw := q.Y*q.Y, q.Y*q.Z, q.Y*q.W
	zz, zw := q.Z*q.Z, q.Z*q.W

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
	result := NewMat3()

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
func (m *Mat3) rotate(angle float64, axis *Vec3) *Mat3 {
	sinAngle := math.Sin(angle * math.Pi / 180)
	cosAngle := math.Cos(angle * math.Pi / 180)
	oneMinusCos := float64(1) - cosAngle

	axis = axis.Unit()

	xx := axis.X * axis.X
	yy := axis.Y * axis.Y
	zz := axis.Z * axis.Z
	xy := axis.X * axis.Y
	yz := axis.Y * axis.Z
	zx := axis.Z * axis.X
	xs := axis.X * sinAngle
	ys := axis.Y * sinAngle
	zs := axis.Z * sinAngle

	rotation := NewMat3()

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
