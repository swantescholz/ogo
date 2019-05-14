//header file for fragment and vertex shader
#define PI 3.14159265359                        //pi
#define degToRad(x) ((x) * 0.0174532925199)   //transforms a degree value to radian
#define radToDeg(x) ((x) * 57.295779513082)   //transforms a radian value to degree

varying vec3 vFragPosition; //the position of the fragment
varying vec3 vFragNormal;   //the normal of the fragment

uniform mat3 uNormalMatrix;
uniform mat4 uModelMatrix;
uniform mat4 uViewMatrix;
uniform mat4 uProjectionMatrix;
uniform mat4 uModelViewMatrix;
uniform mat4 uViewProjectionMatrix;
uniform mat4 uModelViewProjectionMatrix;

uniform gl_MaterialParameters uFrontMaterial;
uniform gl_MaterialParameters uBackMaterial;

struct S_LightSource {
	vec4 ambient;
	vec4 diffuse;
	vec4 specular;
	vec3 position;
	vec3 halfVector; //computed in shader
	vec3 spotDirection;
	float spotExponent;
	float spotCutoff;
	float spotCosCutoff;
	float constantAttenuation;
	float linearAttenuation;
	float quadraticAttenuation;
	bool enabled;
	bool positional;
};
uniform S_LightSource uLightSource[gl_MaxLights];
uniform vec4 uLightModelAmbient; //global ambient color
uniform bool uLightModelLocalViewer;
uniform bool uLightModelTwoSide;

struct S_Camera {
	vec3 pos;
	vec3 dir;
};
uniform S_Camera uCamera;

uniform float uFarClipplane; //far clipplane distance to camera
uniform float uElapsedTime; //time since last frame
uniform float uTimeSinceInit; //time since program start
uniform sampler1D uGradient; //a simple 1D gradient texture
uniform sampler2D uTexture; //a simple 2D texture

const mat2 cmat2id = mat2(1,0,0,1);
const mat3 cmat3id = mat3(1,0,0,0,1,0,0,0,1);
const mat4 cmat4id = mat4(1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1);
const vec3 cvec3zero = vec3(0,0,0);
const vec3 cvec3one = vec3(1,1,1);
const vec3 cvec3x = vec3(1,0,0);
const vec3 cvec3y = vec3(0,1,0);
const vec3 cvec3z = vec3(0,0,1);
const vec3 cvec3nx = vec3(-1,0,0);
const vec3 cvec3ny = vec3(0,-1,0);
const vec3 cvec3nz = vec3(0,0,-1);
const vec2 cvec2zero = vec2(0,0);
const vec2 cvec2one = vec2(1,1);
const vec2 cvec2x = vec2(1,0);
const vec2 cvec2y = vec2(0,1);
const vec2 cvec2nx = vec2(-1,0);
const vec2 cvec2ny = vec2(0,-1);

//some useful functions:
float lengthSq(const vec3 v) {return dot(v,v);}
float distanceSq(const vec3 a, const vec3 b) {return dot(b-a,b-a);}

//matrix generation functions:
mat4 mat4transpose(const mat4 m) {
	return mat4(
		m[0][0],m[1][0],m[2][0],m[3][0],
		m[0][1],m[1][1],m[2][1],m[3][1],
		m[0][2],m[1][2],m[2][2],m[3][2],
		m[0][3],m[1][3],m[2][3],m[3][3]);
}
float mat4determinant(const mat4 m) {
	return
		m[0][0]*(m[1][1]*m[2][2]-m[1][2]*m[2][1])-
		m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0])+
		m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0]);
}
mat4 mat4invert(const mat4 m) {
	mat4 r;
	float invDet = mat4determinant(m);
	if(invDet == 0.0) return cmat4id;
	r[0][0] =  invDet * (m[1][1] * m[2][2] - m[1][2] * m[2][1]);
	r[0][1] = -invDet * (m[0][1] * m[2][2] - m[0][2] * m[2][1]);
	r[0][2] =  invDet * (m[0][1] * m[1][2] - m[0][2] * m[1][1]);
	r[0][3] =  0.0f;
	r[1][0] = -invDet * (m[1][0] * m[2][2] - m[1][2] * m[2][0]);
	r[1][1] =  invDet * (m[0][0] * m[2][2] - m[0][2] * m[2][0]);
	r[1][2] = -invDet * (m[0][0] * m[1][2] - m[0][2] * m[1][0]);
	r[1][3] =  0.0f;
	r[2][0] =  invDet * (m[1][0] * m[2][1] - m[1][1] * m[2][0]);
	r[2][1] = -invDet * (m[0][0] * m[2][1] - m[0][1] * m[2][0]);
	r[2][2] =  invDet * (m[0][0] * m[1][1] - m[0][1] * m[1][0]);
	r[2][3] =  0.0f;
	r[3][0] = -(m[3][0] * r[0][0] + m[3][1] * r[1][0] + m[3][2] * r[2][0]);
	r[3][1] = -(m[3][0] * r[0][1] + m[3][1] * r[1][1] + m[3][2] * r[2][1]);
	r[3][2] = -(m[3][0] * r[0][2] + m[3][1] * r[1][1] + m[3][2] * r[2][2]);
	r[3][3] =  1.0f;
	return r;
}
vec3 mat4transformNormal(const vec3 v, const mat4 m) {
	float fLength = length(v);
	if(fLength == 0.0f) return v;
	mat4 mTransform = mat4transpose(mat4invert(m));
	return normalize(vec3(
		v.x * m[0][0] + v.y * m[1][0] + v.z * m[2][0],
		v.x * m[0][1] + v.y * m[1][1] + v.z * m[2][1],
		v.x * m[0][2] + v.y * m[1][2] + v.z * m[2][2])) * fLength;
}
mat4 mat4scaling(const vec3 v) {
	return mat4(v.x, 0,0,0,0, v.y, 0,0,0,0, v.z, 0,0,0,0,1);
}
mat4 mat4translation(const vec3 v) {
	return mat4(0,0,0,0, 0,0,0,0, 0,0,0,0, v.x,v.y,v.z,1);
}
mat4 mat4rotationX(const float x) {
	float c = cos(x), s = sin(x);
	return mat4(
	1,0,0,0,
	0,c,s,0,
	0,-s,c,0,
	0,0,0,1);
}
mat4 mat4rotationY(const float y) {
	float c = cos(y), s = sin(y);
	return mat4(
	c,0,-s,0,
	0,1,0,0,
	s,0,c,0,
	0,0,0,1);
}
mat4 mat4rotationZ(const float z) {
	float c = cos(z), s = sin(z);
	return mat4(
	c,s,0,0,
	-s,c,0,0,
	0,0,1,0,
	0,0,0,1);
}
mat4 mat4rotation(const float x, const float y, const float z) {
	return mat4rotationZ(z)*mat4rotationX(x)*mat4rotationY(y);
}
mat4 mat4rotation(const vec3 v) {
	return mat4rotation(v.x,v.y,v.z);
}
mat4 mat4rotationAxis(const vec3 v, const float f) {
	//calculate sine und cosine
	float fSin = sin(-f);
	float fCos = cos(-f);
	float fOneMinusCos = 1.0f - fCos;
	
	//normalise axis vector
	vec3 vAxis = normalize(v);
	
	return mat4(
		(vAxis.x * vAxis.x) * fOneMinusCos + fCos,
		(vAxis.x * vAxis.y) * fOneMinusCos - (vAxis.z * fSin),
		(vAxis.x * vAxis.z) * fOneMinusCos + (vAxis.y * fSin),
		0.0,
		(vAxis.y * vAxis.x) * fOneMinusCos + (vAxis.z * fSin),
		(vAxis.y * vAxis.y) * fOneMinusCos + fCos,
		(vAxis.y * vAxis.z) * fOneMinusCos - (vAxis.x * fSin),
		0.0,
		(vAxis.z * vAxis.x) * fOneMinusCos - (vAxis.y * fSin),
		(vAxis.z * vAxis.y) * fOneMinusCos + (vAxis.x * fSin),
		(vAxis.z * vAxis.z) * fOneMinusCos + fCos,
		0.0,
		0.0, 0.0, 0.0, 1.0);
}

