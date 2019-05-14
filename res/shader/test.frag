//#include stdfrag.h

void main (void) 
{
	//if (distanceSq(uCamera.pos, vFragPosition.xyz) > uFarClipplane*uFarClipplane) {
	//	discard;
	//}
	vec4 color = vec4(vFragPosition.xyz,1);
	color = vec4(vFragNormal.xyz, 1);
	vec4 tex = texture2D(uTexture, gl_TexCoord[0].st);
	vec4 pfl = calcPerFragmentLighting();
	vec4 fragColor = 0.002*color + 0.95*tex + 0.05*pfl + 0.0*vec4(vFragPosition.y*0.1,0,0,0);
	fragColor = clamp(fragColor, vec4(0,0,0,0), vec4(1,1,1,1));
	fragColor.a = 1.0;
	gl_FragColor = fragColor;
}
