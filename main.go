package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/utkarsh-1905/conways-game/game"
	"github.com/utkarsh-1905/conways-game/graphics"
)

// constants
const (
	width  = 1000
	height = 1000
)

func GetShaders() (string, string) {
	vertexByte, err := os.ReadFile("shaders/vertex.glsl")
	if err != nil {
		panic(err)
	}
	vertex := string(vertexByte) + "\x00" // \x00 is used to terminate the string and it is a requirement without which the shader will not compile
	fragmentByte, err := os.ReadFile("shaders/fragment.glsl")
	if err != nil {
		panic(err)
	}
	fragment := string(fragmentByte) + "\x00"

	return vertex, fragment
}

func init() {
	runtime.LockOSThread() //locking the is thread on which the glfw context is create since the context will work on same thread only
}

func main() {

	window := graphics.InitGlfw(width, height)
	defer glfw.Terminate()

	program := graphics.InitOpenGL(GetShaders())

	//init flags
	var matSize int
	var prob float64
	var fps int

	flag.IntVar(&matSize, "mat", 100, "Enter the size of Matrix.")
	flag.Float64Var(&prob, "prob", 0.08, "Enter the spawning probability.")
	flag.IntVar(&fps, "fps", 5, "Enter desired fps.")
	flag.Parse()

	//fot play/pause with space key
	shouldPlay := true
	window.SetKeyCallback(func() glfw.KeyCallback {
		return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeySpace && action == glfw.Press {
				shouldPlay = !shouldPlay
			}
		}
	}())

	//for closing the window with shift+w
	window.SetKeyCallback(func() glfw.KeyCallback {
		return func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
			if key == glfw.KeyW && action == glfw.Press && mods == glfw.ModShift {
				window.SetShouldClose(true)
			}
		}
	}())

	game := game.InitGame(matSize, prob) //changing first param changes board size

	for !window.ShouldClose() {
		fmt.Println(shouldPlay)
		Play(game, window, program, fps, shouldPlay)
	}
}

func Play(game *game.Game, window *glfw.Window, program *uint32, fps int, shouldPlay bool) {
	t := time.Now()
	// fmt.Println("Generation ", game.Iterations)

	Draw(game.Matrix, window, program)
	if shouldPlay {
		game.UpdateGame()
		game.Iterations++
	}

	time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
}

func Draw(cells [][]*game.Cell, window *glfw.Window, prog *uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) //clearing the window of anything that was previously drawn
	gl.UseProgram(*prog)                                // uses the program memory we created

	for x := range cells {
		for _, c := range cells[x] {
			c.Draw()
		}
	}

	glfw.PollEvents()    // to handle keyboard or mouse inputs
	window.SwapBuffers() // like traditional graphic drivers, it first draws everything on a blank canvas and swaps it with current window display everytime
}

//Todo - generation number, benchmarking of functions, click to spawn, show matrix grid too
// Feature - add water in matrix to affect spawn
