package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/lukevers/glu"
	"github.com/pkg/errors"
)

const (
	playAreaWidth  = defaultWidth * 2
	playAreaHeight = defaultHeight * 2
)

// titleReg is the entity registry to be used during the title screen
var titleReg entityRegistry

// gameplayReg is the entity registry to be used during gameplay
var gameplayReg entityRegistry

func init() {
	// Required by GLFW and probably OpenGL
	runtime.LockOSThread()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "muted" {
		musicMuted = true
	}

	// Initialize GLFW
	err := glfw.Init()
	if err != nil {
		panic("failed to initialize GLFW " + err.Error())
	}
	defer glfw.Terminate()

	// Add a player for every connected joystick
	players := []*spaceShipListener{}
	for _, joy := range connectedJoysticks() {
		ssl := newSpaceShipListener(joy)
		gameplayReg.addEntity(ssl)
		players = append(players, ssl)
	}

	// Create a GLFW window
	glfwWin, err := glfw.CreateWindow(defaultWidth, defaultHeight, "Vertex", nil, nil)
	if err != nil {
		panic(err)
	}
	glfwWin.MakeContextCurrent()
	glfwWin.SetSizeCallback(onResize)
	mainWindow.glfw = glfwWin

	// Initialize OpenGL
	err = initGL()
	if err != nil {
		panic(err)
	}
	err = loadModels()
	if err != nil {
		panic(errors.Wrap(err, "loading models"))
	}
	err = loadAudio()
	if err != nil {
		panic(errors.Wrap(err, "loading music"))
	}

	titleReg.addEntity(newTitle())

	gameplayReg.addEntity(&mainCamera)
	gameplayReg.addEntity(newGrid(playAreaWidth, playAreaHeight))
	gameplayReg.addEntity(newBoundaries(playAreaWidth, playAreaHeight))
	gameplayReg.addEntity(newEnemySpawner())
	gameplayReg.addEntity(newPlayerRespawner(players))

	currentReg = &titleReg

	resetPulse()
	playMusic(titleMusic, 3*time.Second)

	for !glfwWin.ShouldClose() {
		drawFrame(glfwWin)
		tick()
		updatePulse()
		pollJoysticks()

		mainWindow.calcDelta()

		glfw.PollEvents()
	}

}

// tick runs per-frame non-graphical operations.
func tick() {
	currentReg.tick()
}

// drawFrame draws a single frame to the framebuffer.
func drawFrame(window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.PushMatrix()
	gl.Translated(-mainCamera.loc.x, -mainCamera.loc.y, -1500.0*mainCamera.loc.z)
	currentReg.draw()
	gl.PopMatrix()

	window.SwapBuffers()
}

// initGL initializes OpenGL with our configurations.
func initGL() error {
	err := gl.Init()
	if err != nil {
		return err
	}

	gl.MatrixMode(gl.PROJECTION)

	gl.ShadeModel(gl.SMOOTH)

	gl.Enable(gl.DEPTH_TEST)

	gl.Enable(gl.NORMALIZE)
	gl.Enable(gl.LIGHTING)
	gl.Enable(gl.LIGHT0)

	// Set global lighting configuration
	lightModelAmbient := []float32{0.2, 0.2, 0.2, 1.0}
	gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, fPtr(lightModelAmbient))
	gl.LightModeli(gl.LIGHT_MODEL_LOCAL_VIEWER, gl.FALSE)

	// Configure LIGHT0
	var (
		lightAmbient  = []float32{0.0, 0.0, 0.0, 1}
		lightDiffuse  = []float32{1.0, 1.0, 1.0, 1}
		lightSpecular = []float32{1.0, 1.0, 1.0, 1}
		lightPos      = []float32{0.33, 0.33, 0.33, 0}
	)

	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, fPtr(lightAmbient))
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, fPtr(lightDiffuse))
	gl.Lightfv(gl.LIGHT0, gl.SPECULAR, fPtr(lightSpecular))
	gl.Lightfv(gl.LIGHT0, gl.POSITION, fPtr(lightPos))

	// Size up the viewing area and camera
	onResize(nil, defaultWidth, defaultHeight)

	return nil
}

// onResize notifies OpenGL of window resizes.
func onResize(window *glfw.Window, width, height int) {
	mainWindow.width = width
	mainWindow.height = height

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	glu.Perspective(45, float64(width)/float64(height), 1, 9999)
	gl.MatrixMode(gl.MODELVIEW)
}
