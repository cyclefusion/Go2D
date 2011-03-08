package go2d

import (
	"fmt"
	"sdl"
)

//Instance
var g_game *Game
var g_running bool

//Basic game type

type Game struct {
	initFun   func()
	updateFun func()
	drawFun   func()

	mouseupFun   func(int16, int16)
	mousedownFun func(int16, int16)
	mousemoveFun func(int16, int16)
	keydownFun   func(int)
	textinputFun func(uint8)

	title         string
	width, height int
	d3d           bool

	window   *sdl.Window
	renderer *sdl.Renderer
}

//Create a new Game instance
func NewGame(_title string) (game *Game) {
	game = &Game{}

	game.title = _title

	g_game = game
	return game
}

func (game *Game) SetInitFun(_init func()) {
	game.initFun = _init
}

func (game *Game) SetUpdateFun(_update func()) {
	game.updateFun = _update
}

func (game *Game) SetDrawFun(_draw func()) {
	game.drawFun = _draw
}

func (game *Game) SetMouseUpFun(_mouseup func(int16, int16)) {
	game.mouseupFun = _mouseup
}

func (game *Game) SetMouseDownFun(_mousedown func(int16, int16)) {
	game.mousedownFun = _mousedown
}

func (game *Game) SetMouseMoveFun(_mousemove func(int16, int16)) {
	game.mousemoveFun = _mousemove
}

func (game *Game) SetKeyDownFun(_keydown func(int)) {
	game.keydownFun = _keydown
}

func (game *Game) SetTextInputFun(_textinput func(uint8)) {
	game.textinputFun = _textinput
}

//Internal initalization (executes when game starts)
func (game *Game) initialize() {
	//Create the window
	var err string
	game.window, err = sdl.CreateWindow(game.title, game.width, game.height)
	if err != "" {
		panic(fmt.Sprintf("Go2D Error: Creating window: %s", err))
	}

	//Create the renderer

	//Find our available renderers
	openglIndex := 0
	d3dIndex := -1
	numRenderers := sdl.GetNumRenderDrivers()
	for i := 0; i < numRenderers; i++ {
		rendererName := sdl.GetRenderDriverName(i)
		if rendererName == "opengl" {
			openglIndex = i
		} else if rendererName == "direct3d" {
			d3dIndex = i
		}
	}

	//Default renderer is OpenGL
	rendererIndex := openglIndex

	//If we want to use Direct3D and we found it, use it
	if game.d3d && d3dIndex != -1 {
		rendererIndex = d3dIndex
	}

	game.renderer, err = sdl.CreateRenderer(game.window, rendererIndex)
	if err != "" {
		panic(fmt.Sprintf("Go2D Error: Creating renderer: %s", err))
	}

	//initialize font rendering
	sdl.InitTTF()

	//Call the user-defined init function
	if game.initFun != nil {
		game.initFun()
	}

	g_running = true
}

//Internal update function
func (game *Game) update() {
	//Call user-defined update function
	if game.updateFun != nil {
		game.updateFun()
	}
}

//Internal draw function
func (game *Game) draw() {
	//Clear the screen
	sdl.RenderClear(game.renderer)

	//Call user-defined draw function
	if game.drawFun != nil {
		game.drawFun()
	}

	//Render everything
	sdl.RenderPresent(game.renderer)
}

//Set window dimensions
func (game *Game) SetDimensions(_width, _height int) {
	game.width = _width
	game.height = _height
}

//Try to use D3D or not
func (game *Game) SetD3D(_d3d bool) {
	game.d3d = _d3d
}

//Game loop
func (game *Game) Run() {
	defer game.Exit()

	if game.initFun == nil {
		fmt.Println("Go2D Warning: No init function set!")
	}

	if game.updateFun == nil {
		fmt.Println("Go2D Warning: No update function set!")
	}

	if game.drawFun == nil {
		fmt.Println("Go2D Warning: No draw function set!")
	}

	//Initialize the game
	game.initialize()

	for g_running {
		//Check for events and handle them
		for {
			event, present := sdl.PollEvent()
			if present {
				EventHandler(event)
			} else {
				break
			}
		}

		//Update
		game.update()

		//Draw
		game.draw()

		//Give the CPU some time to do other stuff
		sdl.Delay(1)
	}
}

//Release all resources
func (game *Game) Exit() {
	freeResources()

	//Destroy the renderer
	game.renderer.Release()

	//Destroy the window
	sdl.DestroyWindow(game.window)

	//Quit SDL_ttf
	sdl.QuitTTF()

	//Quit SDL
	sdl.Quit()
}