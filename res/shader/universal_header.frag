#define M_PI 3.14159265359                        //pi
#define DEG_TO_RAD(x) ((x) * 0.0174532925199)   //transforms a degree value to radian
#define RAD_TO_DEG(x) ((x) * 57.295779513082)   //transforms a radian value to degree

//SWITCHES:
#define MAX_LIGHTS 2
//#define FOG_EXP2
//#define FOG_LINEAR

varying vec3 vFragPosition; //the position of the fragment
varying vec3 vFragNormal; //the normal for the fragment
uniform vec3 uCameraPosition; //the position of the eye
uniform vec3 uCameraDirection; //the direction the eye is looking
uniform float uElapsedTime; //time since last frame
uniform float uTimeSinceInit; //time since program start
uniform int uNumLights; //number of lights in the scene
uniform sampler1D uGradient; //a simple 1D gradient texture
uniform sampler2D uTexture;
uniform vec4 uFogColor;
uniform float uFogDensity;
uniform bool uSeparateSpecular; //shall the specular material be calculated seperately?
const mat2 cmat2identity = mat2(1,0,0,1);
const mat3 cmat3identity = mat3(1,0,0,0,1,0,0,0,1);
const mat4 cmat4identity = mat4(1,0,0,0,0,1,0,0,0,0,1,0,0,0,0,1);
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

struct S_LightSource {
	vec4 ambient;
	vec4 diffuse;
	vec4 specular;
	vec3 position;
	vec3 spotDirection;
	float spotExponent;
	float spotCutoff;
	float constantAttenuation;
	float linearAttenuation;
	float quadraticAttenuation;
};
uniform S_LightSource uLightSource[gl_MaxLights];
uniform vec4 uLightModelAmbient; //global ambient color

//some useful functions:
float lengthSq(const vec3 v) {return dot(v,v);}
float distanceSq(const vec3 a, const vec3 b) {return dot(b-a,b-a);}

float calcFogFactor(float distance, float fogDensity) {
	//fog function:
	float fog;
#ifdef FOG_EXP2
	fog = 1-exp(-fogDensity*fogDensity * distance*distance); //exp2
#else
#ifdef FOG_LINEAR
	fog = distance*fogDensity; //linear
#else
	fog = 1-exp(-fogDensity * distance); //simple exp
#endif
#endif
	return clamp(fog, 0.0, 1.0);
}
float calcGroundFogDistance(float groundFogMin, float groundFogMax,
		vec3 eyePos, vec3 fragPos) {
	float h;
	if ((eyePos.y <= groundFogMin && fragPos.y <= groundFogMin) ||
		(eyePos.y >= groundFogMax && fragPos.y >= groundFogMax)) {
		return 0.0; //no fog
	}
	if ((eyePos.y >= groundFogMax && fragPos.y <= groundFogMin) ||
		(eyePos.y <= groundFogMin && fragPos.y >= groundFogMax) ) {
		return distance(eyePos, fragPos)*(groundFogMax-groundFogMin)/abs(eyePos.y-fragPos.y); //both out
	}
	if (eyePos.y < groundFogMax && eyePos.y > groundFogMin &&
		fragPos.y < groundFogMax && fragPos.y > groundFogMin) {
		return distance(eyePos, fragPos); //both in
	}
	if (eyePos.y < groundFogMax && eyePos.y > groundFogMin) { //just camera in
		h = (fragPos.y > groundFogMax)?(groundFogMax-eyePos.y):(eyePos.y-groundFogMin);
		return distance(eyePos, fragPos)*h/abs(eyePos.y-fragPos.y);
	}
	//just frag in
	h = (eyePos.y > groundFogMax)?(groundFogMax-fragPos.y):(fragPos.y-groundFogMin);
	return distance(eyePos, fragPos)*h/abs(eyePos.y-fragPos.y);
}
float calcSphereFogDistance(vec3 center, float radius, vec3 eyePos, vec3 fragPos) {
	const float radiusSq = radius*radius;
	const vec3 u = fragPos-eyePos;
	const vec3 c = eyePos-center;
	const float udot = dot(u, u);
	const float phalf = dot(c,u)/udot;
	const float disc = phalf*phalf+(radiusSq-dot(c,c))/udot;
	if (disc < 0.0) return 0.0; //no intersection
	if (distanceSq(eyePos, center) < radiusSq) { //camera in
		if (distanceSq(fragPos, center) < radiusSq) {
			return distance(eyePos, fragPos); //both in
		}
		//just camera in
		const float t = -phalf + sqrt(disc);
		return t*length(u);
	}
	if (distanceSq(fragPos, center) < radiusSq) { //just fragment in
		const float t = -phalf - sqrt(disc);
		return (1-t)*length(u);
	}
	//both out:
	if (phalf > 0.0 || phalf < -1.0) return 0.0;
	const float dtime = 2*sqrt(disc);
	return length(u)*dtime;
}
void boxIntersectionTime(float p, float u, float boxPos, float boxHalfH, out float t1, out float t2) {
	const float upper = boxPos + boxHalfH;
	const float lower = boxPos - boxHalfH;
	if (p >= upper && p+u >= upper) {t1 = -2; t2 = -1;}
	if (p <= lower && p+u <= lower) {t1 = -2; t2 = -1;}
	if (abs(u) < 0.001) {t1 = -2; t2 = -1;} //nearly parallel
	if (u < 0) {
		t1 = clamp((upper-p)/u, 0.0, 1.0);
		t2 = clamp((lower-p)/u, 0.0, 1.0);
	}
	else {
		t2 = clamp((upper-p)/u, 0.0, 1.0);
		t1 = clamp((lower-p)/u, 0.0, 1.0);
	}
}
float calcBoxFogDistance(vec3 center, vec3 halfsize, vec3 eyePos, vec3 fragPos) {
	const vec3 u = fragPos-eyePos;
	float tmin[3];
	float tmax[3];
	boxIntersectionTime(eyePos.x, u.x, center.x, halfsize.x, tmin[0], tmax[0]);
	boxIntersectionTime(eyePos.y, u.y, center.y, halfsize.y, tmin[1], tmax[1]);
	boxIntersectionTime(eyePos.z, u.z, center.z, halfsize.z, tmin[2], tmax[2]);
	const float min = max(max(tmin[0],tmin[1]),tmin[2]);
	const float max = min(min(tmax[0],tmax[1]),tmax[2]);
	return (max-min)*length(u);
}

//LIGHT + MATERIALS

void pointLight(in int i,
		inout vec4 ambient,
		inout vec4 diffuse,
		inout vec4 specular) {
	vec3 L = uLightSource[i].position - vFragPosition; // Compute vector from surface to light position
	const float d = length(L);     // Compute distance between surface and light position
	L = normalize(L); // Normalize the vector from surface to light position
	const vec3 R = normalize(-reflect(L, vFragNormal)); //perfect reflection vector
	const vec3 E = normalize(uCameraPosition-vFragPosition); //to-eye vector
	const float attenuation = 1.0 / (
		uLightSource[i].constantAttenuation +
		uLightSource[i].linearAttenuation * d +
		uLightSource[i].quadraticAttenuation * d * d);
	//const vec3 halfVector = normalize(L + uCameraPosition-vFragPosition); //direction of maximum highlights
	const float nDotL = max(0.0, dot(vFragNormal, L));
	const float rDotE = max(0.0, dot(R,E));
	float pf; //power factor
	if (nDotL == 0.0) pf = 0.0;
	else pf = pow(rDotE, gl_FrontMaterial.shininess);
	ambient += attenuation * uLightSource[i].ambient;
	diffuse += attenuation * nDotL * uLightSource[i].diffuse;
	specular += attenuation * pf * uLightSource[i].specular;
	//diffuse = mix(vec4(0,0,1,1),vec4(1,0,0,1), attenuation);
}
 
vec4 calcPerFragmentLighting() {
	vec4 color = vec4(0.0, 0.0, 0.0, 0.0);
	vec4 amb = vec4(0.0);
	vec4 diff = vec4(0.0);
	vec4 spec = vec4(0.0);
	for (int i=0;i<MAX_LIGHTS;i++) {
		pointLight(i, amb, diff, spec);
	}
	color = uLightModelAmbient +
		amb * gl_FrontMaterial.ambient +
		diff * gl_FrontMaterial.diffuse +
		spec * gl_FrontMaterial.specular;
	//if (uSeparateSpecular) gl_FrontSecondaryColor = vec4(spec * gl_FrontMaterial.specular, 1.0);
	//else color += spec * gl_FrontMaterial.specular;
	//color = spec;
	color = clamp(color, vec4(0,0,0,0), vec4(1,1,1,1));
	color.a = 1.0;
	return color;
}
