//#include stdvert.h

void main(void)
{
	vec4 vertex = uModelMatrix * gl_Vertex;
	vFragNormal = uNormalMatrix * gl_Normal.xyz;
	vFragPosition = vertex.xyz;
	gl_TexCoord[0] = gl_MultiTexCoord0;
	vertex = (uProjectionMatrix * (uViewMatrix * vertex));
	vertex = vertex + vec4(0,0.0,vertex.x*vertex.x,0);
	gl_Position = vertex;
}
