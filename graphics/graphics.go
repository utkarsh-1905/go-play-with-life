package graphics

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// initialize opengl
func InitOpenGL(vertexShader, fragmentShader string) *uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Opengl version", version)

	vertex, err := CompileShader(vertexShader, gl.VERTEX_SHADER) //compiling the vertex shader
	if err != nil {
		panic(err)
	}
	fragment, err := CompileShader(fragmentShader, gl.FRAGMENT_SHADER) //compiling the fragment shader
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertex)
	gl.AttachShader(prog, fragment)
	gl.LinkProgram(prog) // this program gives a reference to store shaders in the future
	return &prog
}

// initialize glfw
func InitGlfw(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	//the init function create a context on a thread, the rest of program should run on same thread, hence thread is locked

	//global properties
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Visual SHA", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent() //making it run on current thread

	return window
}

// generating buffers to store and display data points
func MakeVAO(points []float32) uint32 {
	var vbo uint32 //making the vertex buffer object
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo) //binding the buffer to the array buffer
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)
	//the last parameter is the type of draw, static draw means that the data will not change
	//we provide the size of buffer i.e 4*len(points) since each point is of 4 bytes
	//we provide the pointer to the first element of the array (Ptr function returns gl compatible pointer)

	var vao uint32 //making the vertex array object
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao) //binding the vertex array object to the vao
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

// ALl above is boilerplate code to display shaders on screen
//This function will compile the shaders and return the shader

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType) //creating a shader of type shaderType
	csources, free := gl.Strs(source)     //converting the source code to gl compatible string in C
	//the free function must be called in the end to free the memory

	gl.ShaderSource(shader, 1, csources, nil) //attaching the source code to the shader
	free()

	gl.CompileShader(shader) //compiling the shader

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status) //getting the status of compilation

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength) //getting the length of log

		log := strings.Repeat("\x00", int(logLength+1))          //creating a log string of length logLength
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log)) //getting the log

		return 0, fmt.Errorf("failed to compile %v: %v", source, log) //returning the error
	}

	return shader, nil
}
