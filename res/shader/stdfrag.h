//#include std.h

//LIGHT + MATERIALS

void pointLight(in int i,
		inout vec4 ambient,
		inout vec4 diffuse,
		inout vec4 specular) {
	// vector from surface to light position
	vec3 L = uLightSource[i].position - vFragPosition;
	float d = length(L);     // distance between surface and light position
	L = normalize(L);
	vec3 R = normalize(-reflect(L, vFragNormal)); //perfect reflection vector
	vec3 E = normalize(uCamera.pos-vFragPosition); //to-eye vector
	float attenuation = 1.0 / (
		uLightSource[i].constantAttenuation +
		uLightSource[i].linearAttenuation * d +
		uLightSource[i].quadraticAttenuation * d * d);
	//direction of maximum highlights:
	//const vec3 halfVector = normalize(L + uCameraPosition-vFragPosition);
	float nDotL = max(0.0, dot(vFragNormal, L));
	float rDotE = max(0.0, dot(R,E));
	float pf; //power factor
	if (nDotL == 0.0) pf = 0.0;
	else pf = pow(rDotE, uFrontMaterial.shininess);
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
	for (int i = 0; i < gl_MaxLights; i++) {
		if (uLightSource[i].enabled) {
			pointLight(i, amb, diff, spec);
		}
	}
	color = uLightModelAmbient + uFrontMaterial.emission +
		amb  * uFrontMaterial.ambient * 0.2 +
		diff * uFrontMaterial.diffuse * 0.5 +
		spec * uFrontMaterial.specular * 1.0;
	color = clamp(color, vec4(0,0,0,0), vec4(1,1,1,1));
	color.a = 1.0;
	return color;
}
