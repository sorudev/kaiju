{"Name":"sprite_transparent","Vertex":"content\\renderer\\src\\sprite.vert","VertexFlags":"","Fragment":"content\\renderer\\src\\sprite.frag","FragmentFlags":"-DOIT","Geometry":"","GeometryFlags":"","TessellationControl":"","TessellationControlFlags":"","TessellationEvaluation":"","TessellationEvaluationFlags":"","LayoutGroups":[{"Type":"Vertex","Layouts":[{"Location":0,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec3","Name":"Position","Source":"in","Fields":null},{"Location":1,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec3","Name":"Normal","Source":"in","Fields":null},{"Location":2,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"Tangent","Source":"in","Fields":null},{"Location":3,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec2","Name":"UV0","Source":"in","Fields":null},{"Location":4,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"Color","Source":"in","Fields":null},{"Location":5,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"ivec4","Name":"JointIds","Source":"in","Fields":null},{"Location":6,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"JointWeights","Source":"in","Fields":null},{"Location":7,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec3","Name":"MorphTarget","Source":"in","Fields":null},{"Location":-1,"Binding":0,"Set":0,"InputAttachment":-1,"Type":"UniformBufferObject","Name":"","Source":"uniform","Fields":[{"Type":"mat4","Name":"view"},{"Type":"mat4","Name":"projection"},{"Type":"mat4","Name":"uiView"},{"Type":"mat4","Name":"uiProjection"},{"Type":"vec3","Name":"cameraPosition"},{"Type":"vec3","Name":"uiCameraPosition"},{"Type":"vec2","Name":"screenSize"},{"Type":"float","Name":"time"}]},{"Location":8,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"mat4","Name":"model","Source":"in","Fields":null},{"Location":12,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"uvs","Source":"in","Fields":null},{"Location":13,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"fgColor","Source":"in","Fields":null},{"Location":0,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"fragColor","Source":"out","Fields":null},{"Location":4,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec2","Name":"fragTexCoord","Source":"out","Fields":null}]},{"Type":"Fragment","Layouts":[{"Location":0,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"fragColor","Source":"in","Fields":null},{"Location":4,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec2","Name":"fragTexCoord","Source":"in","Fields":null},{"Location":-1,"Binding":1,"Set":-1,"InputAttachment":-1,"Type":"sampler2D","Name":"texSampler","Source":"uniform","Fields":null},{"Location":0,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"vec4","Name":"outColor","Source":"out","Fields":null},{"Location":1,"Binding":-1,"Set":-1,"InputAttachment":-1,"Type":"float","Name":"reveal","Source":"out","Fields":null}]}]}