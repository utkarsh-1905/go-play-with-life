package main

import (
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/utkarsh-1905/conways-game/game"
	"github.com/utkarsh-1905/conways-game/graphics"
)

// Window size
const (
	width  = 1000
	height = 1000

	//defining shaders

	vertexShader = `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShader = `
		#version 410
		out vec4 frag_colour;
		void main(){
			frag_colour = vec4(255,255,255,1);
		}
	` + "\x00"
)

// \x00 is used to terminate the string and it is a requirement without which the shader will not compile

var square = []float32{
	-0.5, 0.5, 0,
	-0.5, -0.5, 0,
	0.5, -0.5, 0,
	-0.5, 0.5, 0,
	0.5, 0.5, 0,
	0.5, -0.5, 0,
}

func init() {
	runtime.LockOSThread() //locking the is thread on which the glfw context is create since the context will work on same thread only
}

func main() {

	window := graphics.InitGlfw(width, height)
	defer glfw.Terminate()
	program := graphics.InitOpenGL(vertexShader, fragmentShader)
	// vao := graphics.MakeVAO(square)
	game := game.InitGame(100, 3) //changing first param changes board size

	for !window.ShouldClose() {
		Draw(game.Matrix, window, program)
	}
}

func Draw(cells [][]*game.Cell, window *glfw.Window, prog *uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) //clearing the window of anything that was previously drawn
	gl.UseProgram(*prog)                                // uses the program memory we created

	// gl.BindVertexArray(vao) //binding the vertex array object to the vao
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square))/3)

	// for x := range cells {
	// 	for _, c := range cells[x] {
	// 		c.Draw()
	// 	}
	// }
	cells[5][7].Draw()
	cells[6][1].Draw()
	glfw.PollEvents()    // to handle keyboard or mouse inputs - not needed
	window.SwapBuffers() // like traditional graphic drivers, it first draws everything on a blank canvas and swaps it with current window display everytime
}
