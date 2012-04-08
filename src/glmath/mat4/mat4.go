package mat4

import "math"
import "fmt"
import . "glmath/vec3"
import . "glmath/qtrnn"

type Mat4 struct {
	M [4 * 4]float64
}

var Identity = &Mat4{[16]float64{1,0,0,0, 0,1,0,0, 0,0,1,0, 0,0,0,1}}

func New() *Mat4 {
	m := new(Mat4)
	m.SetIdentity()
	return m
}

func (m *Mat4) Fv32() []float32 {
	fv32 := make([]float32, 16)
	for i := 0; i < 16; i++ {
		fv32[i] = float32(m.M[i])
	}
	return fv32
}

func (m *Mat4) Copy() *Mat4 {
	n := New()
	for i := 0; i < 4*4; i += 1 {
		n.M[i] = m.M[i]
	}
	return n
}

func (m *Mat4) Set(n *Mat4) {
	for i := 0; i < 16; i ++ {
		m.M[i] = n.M[i]
	}
}

func (m *Mat4) SetIdentity() {
	m.M = [4 * 4]float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func (m *Mat4) Ptr() *float64 {
	return (*float64)(&(m.M[0]))
}

func Translation(pos *Vec3) *Mat4 {
	m := New()
	m.M = [16]float64{
		1,0,0,0,
		0,1,0,0,
		0,0,1,0,
		pos.X,pos.Y,pos.Z,1}
	return m
}

func Scaling(scale *Vec3) *Mat4 {
	m := New()
	m.M = [16]float64{
		scale.X,0,0,0,
		0,scale.Y,0,0,
		0,0,scale.Z,0,
		0,0,0,1}
	return m
}

func RotationAxis(axis *Vec3, f float64) *Mat4 {
	m := New()
	fSin := math.Sin(-f)
	fCos := math.Cos(-f)
	fOneMinusCos := 1.0 - fCos
	vAxis := axis.Unit()
	m.M = [16]float64{(vAxis.X * vAxis.X) * fOneMinusCos + fCos,
		              (vAxis.X * vAxis.Y) * fOneMinusCos - (vAxis.Z * fSin),
		              (vAxis.X * vAxis.Z) * fOneMinusCos + (vAxis.Y * fSin),
		              0.0,
		              (vAxis.Y * vAxis.X) * fOneMinusCos + (vAxis.Z * fSin),
		              (vAxis.Y * vAxis.Y) * fOneMinusCos + fCos,
		              (vAxis.Y * vAxis.Z) * fOneMinusCos - (vAxis.X * fSin),
		              0.0,
		              (vAxis.Z * vAxis.X) * fOneMinusCos - (vAxis.Y * fSin),
		              (vAxis.Z * vAxis.Y) * fOneMinusCos + (vAxis.X * fSin),
		              (vAxis.Z * vAxis.Z) * fOneMinusCos + fCos,
		              0.0,
		              0.0, 0.0, 0.0, 1.0}
	return m
}

func Projection(fovx, ratio, znear, zfar float64) *Mat4 {
	m := New()
	dxInv := 1.0/(2.0 * znear * math.Tan(fovx*0.5));
	dyInv := dxInv*ratio;
	dzInv := 1.0/(zfar - znear);
	X := 2.0 * znear;
	C := -(znear+zfar)*dzInv;
	D := -znear*zfar*dzInv;
	m.M = [16]float64{X*dxInv,0,0,0,  0,X*dyInv,0,0,  0,0,C,-1, 0,0,D,0}
	return m
}

func Camera(pos, dir, up *Vec3) *Mat4 {
	mtrans := Translation(pos.Muls(-1))
	vZAxis := dir.Unit().Muls(-1)
	vXAxis := up.Cross(vZAxis).Unit()
	vYAxis := vZAxis.Cross(vXAxis).Unit()
	mrot := New()
	mrot.M = [16]float64{
		vXAxis.X, vYAxis.X, vZAxis.X, 0.0,
		vXAxis.Y, vYAxis.Y, vZAxis.Y, 0.0,
		vXAxis.Z, vYAxis.Z, vZAxis.Z, 0.0,
		0.0,      0.0,      0.0,      1.0};
	return mtrans.Mul(mrot)
}

// returns the determinant
func (m *Mat4) Determinant() float64 {
	//determinant of the left 3*3 matrix
	return m.M[0] * (m.M[5] * m.M[9] - m.M[6] * m.M[9]) -
		m.M[1] * (m.M[4] * m.M[9] - m.M[6] * m.M[8]) +
		m.M[2] * (m.M[4] * m.M[8] - m.M[5] * m.M[8])
}
/*
func (m *Mat4) Invert() *Mat4 {
	//calculate inversed value of the determinant
	fInvDet := m.Determinant()
	if fInvDet == 0.0f {return New()}
	fInvDet = 1.0f / fInvDet;
	
	mResult := New()
	mResult.m11 =  fInvDet * (m.m22 * m.m33 - m.m23 * m.m32);
	mResult.m12 = -fInvDet * (m.m12 * m.m33 - m.m13 * m.m32);
	mResult.m13 =  fInvDet * (m.m12 * m.m23 - m.m13 * m.m22);
	mResult.m14 =  0.0f;
	mResult.m21 = -fInvDet * (m.m21 * m.m33 - m.m23 * m.m31);
	mResult.m22 =  fInvDet * (m.m11 * m.m33 - m.m13 * m.m31);
	mResult.m23 = -fInvDet * (m.m11 * m.m23 - m.m13 * m.m21);
	mResult.m24 =  0.0f;
	mResult.m31 =  fInvDet * (m.m21 * m.m32 - m.m22 * m.m31);
	mResult.m32 = -fInvDet * (m.m11 * m.m32 - m.m12 * m.m31);
	mResult.m33 =  fInvDet * (m.m11 * m.m22 - m.m12 * m.m21);
	mResult.m34 =  0.0f;
	mResult.m41 = -(m.m41 * mResult.m11 + m.m42 * mResult.m21 + m.m43 * mResult.m31);
	mResult.m42 = -(m.m41 * mResult.m12 + m.m42 * mResult.m22 + m.m43 * mResult.m32);
	mResult.m43 = -(m.m41 * mResult.m13 + m.m42 * mResult.m23 + m.m43 * mResult.m33);
	mResult.m44 =  1.0f;
	return mResult;
}//*/

// transpose the matrix
func (m *Mat4) Transpose() *Mat4 {	
	return &Mat4{[16]float64{
		m.M[0], m.M[4], m.M[8 ], m.M[12],
		m.M[1], m.M[5], m.M[9 ], m.M[13],
		m.M[2], m.M[6], m.M[10], m.M[14],
		m.M[3], m.M[7], m.M[11], m.M[15]}}
}

//#########################

func (m *Mat4) RotateLocal(angle float64, axis *Vec3) {
	*m = *m.rotate(angle, axis)
}

func (m *Mat4) RotateGlobal(angle float64, axis *Vec3) {
	axis = m.Mulv(axis)
	*m = *m.rotate(angle, axis)
}

func (m *Mat4) Right() *Vec3 {
	return V3(m.M[0], m.M[1], m.M[2])
}

func (m *Mat4) Up() *Vec3 {
	return V3(m.M[4], m.M[5], m.M[6])
}

func (m *Mat4) Forward() *Vec3 {
	return V3(m.M[8], m.M[9], m.M[10])
}

func (m *Mat4) SetRightUpForward(right, up, forward *Vec3) {
	m.SetIdentity()
	m.M[0] = right.X
	m.M[1] = right.Y
	m.M[2] = right.Z
	m.M[4] = up.X
	m.M[5] = up.Y
	m.M[6] = up.Z
	m.M[8] = forward.X
	m.M[9] = forward.Y
	m.M[10] = forward.Z
}

// Get M rotation as Euler angles in degrees
func (m *Mat4) GetEuler() *Vec3 {
	x := math.Atan((-m.M[6]) / m.M[10])
	y := math.Asin(m.M[2])
	z := math.Atan((-m.M[1]) / m.M[0])

	// Convert to Degrees
	x *= 180 / math.Pi
	y *= 180 / math.Pi
	z *= 180 / math.Pi

	return V3(x, y, z)
}

// Set M rotation to Euler angles in degrees
func (m *Mat4) SetEuler(r *Vec3) {
	// Convert to Radians
	r.X *= math.Pi / 180
	r.Y *= math.Pi / 180
	r.Z *= math.Pi / 180
	
	m.SetIdentity()
	m.M[0] = math.Cos(r.Y) * math.Cos(r.Z)
	m.M[1] = -math.Cos(r.Y) * math.Sin(r.Z)
	m.M[2] = math.Sin(r.Y)

	m.M[4] = math.Sin(r.X)*math.Sin(r.Y)*math.Cos(r.Z) +
		math.Cos(r.X)*math.Sin(r.Z)
	m.M[5] = -math.Sin(r.X)*math.Sin(r.Y)*math.Sin(r.Z) +
		math.Cos(r.X)*math.Cos(r.Z)
	m.M[6] = -math.Sin(r.X) * math.Cos(r.Y)

	m.M[8] = -math.Cos(r.X)*math.Sin(r.Y)*math.Cos(r.Z) +
		math.Sin(r.X)*math.Sin(r.Z)
	m.M[9] = math.Cos(r.X)*math.Sin(r.Y)*math.Sin(r.Z) +
		math.Sin(r.X)*math.Cos(r.Z)
	m.M[10] = math.Cos(r.X) * math.Cos(r.Y)
}

// Set M rotation to quateronion 'q'
func (m *Mat4) SetQuaternion(q *Qtrnn) {
	xx, xy, xz, xw := q.X*q.X, q.X*q.Y, q.X*q.Z, q.X*q.W
	yy, yz, yw := q.Y*q.Y, q.Y*q.Z, q.Y*q.W
	zz, zw := q.Z*q.Z, q.Z*q.W
	
	m.SetIdentity()
	m.M[0] = 1.0 - 2.0*(yy+zz)
	m.M[1] = 2.0 * (xy - zw)
	m.M[2] = 2.0 * (xz + yw)
	m.M[4] = 2.0 * (xy + zw)
	m.M[5] = 1.0 - 2.0*(xx+zz)
	m.M[6] = 2.0 * (yz - xw)
	m.M[8] = 2.0 * (xz - yw)
	m.M[9] = 2.0 * (yz + xw)
	m.M[10] = 1.0 - 2.0*(xx+yy)
}

// Multiply 'm' by 'o' and return result
func (m *Mat4) Mul(o *Mat4) *Mat4 {
	result := New()
	for row := 0; row < 4; row++ {
		da := 4 * row
		db := da + 1
		dc := da + 2
		dd := da + 3

		result.M[da] =
			m.M[da]*o.M[0] +
			m.M[db]*o.M[4] +
			m.M[dc]*o.M[8] +
			m.M[dd]*o.M[12]
		result.M[db] =
			m.M[da]*o.M[1] +
			m.M[db]*o.M[5] +
			m.M[dc]*o.M[9] +
			m.M[dd]*o.M[13]
		result.M[dc] =
			m.M[da]*o.M[2] +
			m.M[db]*o.M[6] +
			m.M[dc]*o.M[10] +
			m.M[dd]*o.M[14]
		result.M[dd] =
			m.M[da]*o.M[3] +
			m.M[db]*o.M[7] +
			m.M[dc]*o.M[11] +
			m.M[dd]*o.M[15]
	}
	return result
}

// Multiply 'm' by 'v' and return result
func (m *Mat4) Mulv(v *Vec3) *Vec3 {
	return V3(
		v.X*m.M[0]+v.Y*m.M[1]+v.Z*m.M[2]+1.0*m.M[3],
		v.X*m.M[4]+v.Y*m.M[5]+v.Z*m.M[6]+1.0*m.M[7],
		v.X*m.M[8]+v.Y*m.M[9]+v.Z*m.M[10]+1.0*m.M[11])
}

// Unexported rotate used by RotateLocal & RotateGlobal
func (m *Mat4) rotate(angle float64, axis *Vec3) *Mat4 {
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

	rotation := New()

	rotation.M[0] = (oneMinusCos * xx) + cosAngle
	rotation.M[1] = (oneMinusCos * xy) - zs
	rotation.M[2] = (oneMinusCos * zx) + ys

	rotation.M[4] = (oneMinusCos * xy) + zs
	rotation.M[5] = (oneMinusCos * yy) + cosAngle
	rotation.M[6] = (oneMinusCos * yz) - xs

	rotation.M[8] = (oneMinusCos * zx) - ys
	rotation.M[9] = (oneMinusCos * yz) + xs
	rotation.M[10] = (oneMinusCos * zz) + cosAngle

	return rotation.Mul(m)
}

func (m *Mat4) String() string {
	s := "["
	for i := 0; i < 16; i++ {
		s += fmt.Sprintf("%v", m.M[i])
		if i < 15 {
			s += ", "
		}
	}
	s += "]"
	return s
}

/*
func ProjectionMatrix(fovx, ratio, znear, zfar float64) *Mat4 {
	m := New()
	const piOver360 = 0.0087266462599716478846184538424431
	ymax := znear * math.Tan(float64(fovx * piOver360))
	ymin := -ymax
	//xmax := ymax * ratio
	xmin := ymin * ratio
		
	width := ymax - xmin
	height := ymax - ymin
	
	depth := zfar - znear
	q := -(zfar + znear) / depth
	qn := -2 * (zfar * znear) / depth
	
	w := 2 * znear / width
	w = w / ratio
	h := 2 * znear / height

	m.M[0]  = w
	m.M[1]  = 0
	m.M[2]  = 0
	m.M[3]  = 0
	
	m.M[4]  = 0
	m.M[5]  = h
	m.M[6]  = 0
	m.M[7]  = 0
	
	m.M[8]  = 0
	m.M[9]  = 0
	m.M[10] = q
	m.M[11] = -1
	
	m.M[12] = 0
	m.M[13] = 0
	m.M[14] = qn
	m.M[15] = 0
	return m
}//*/



