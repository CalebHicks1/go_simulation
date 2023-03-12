package main

import (
	"fmt"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type Force struct {
	source     string
	xComponent float64
	yComponent float64
	Strength   float64
	xPos       float64
	yPos       float64
}

// defines the physical properties of a type of atom
type AtomType struct {
	name                string
	collisionElasticity float64 // elasticity when colliding with other atoms
	borderElasticity    float64 // elasticity when colliding with walls
	color               pixel.RGBA
	extraWidth          int // number of pixels of width to add when rendering atom
	mass                float64
	friction            float64
}

type Atom struct {
	atomType     AtomType
	xPos         float64 // the actual x location of the atom
	yPos         float64 // the actual y location of the atom
	xVel         float64
	yVel         float64
	currGridXPos int // the x location on the pixel grid of the atom
	currGridYPos int // the y location on the pixel grid of the atom
}

// constant definitions
const AtomWidth = 4
const Gravity = 20
const windowWidth = 1200 // must be a multiple of AtomWidth
const windowHeight = 600
const mouseGrav = 2000

// variable definitions
var (
	gravEnabled = true
	frames      = 0
	second      = time.Tick(time.Second)
	color       = pixel.RGB(0.3, 0.3, 1).Mul(pixel.Alpha(0.7))
	WindowColor = pixel.RGB(0.9, 0.3, 0.5)
	atoms       []*Atom
	forces      []Force
	grid        [windowWidth/AtomWidth + 1][windowHeight/AtomWidth + 1]*Atom
	timeElapsed float64
	currType    = 0
)

var AtomTypes = []AtomType{
	AtomType{
		"water",
		0.97,
		0.2,
		pixel.RGB(0.3, 0.3, 1).Mul(pixel.Alpha(0.5)),
		4,
		0.5,
		0.5,
	},
	AtomType{
		"stone",
		0.5,
		0.2,
		pixel.RGB(0.2, 0.2, 0.2),
		0,
		5,
		0.5,
	},
}

func drawGrid() *imdraw.IMDraw {
	grid := imdraw.New(nil)
	grid.Color = WindowColor.Mul(pixel.RGB(0.8, 0.8, 0.8))
	for x := 0.0; x <= windowWidth; x += AtomWidth {
		grid.Push(pixel.V(x, 0))
		grid.Push(pixel.V(x, windowHeight))
		grid.Line(1)
		grid.Push(pixel.V(0, x))
		grid.Push(pixel.V(windowWidth, x))
		grid.Line(1)
	}
	return grid
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Particles!",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	imd := imdraw.New(nil)

	// win.Clear(colornames.Skyblue)

	// create forces
	forces = append(forces, Force{"gravity", 0, 0, Gravity, windowWidth / 2, -100000})

	last := time.Now()

	// grid := drawGrid()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		imd.Clear() // clear the window

		// keypress
		handleKeyPress(*win)

		// simulate Atoms
		tempGrid := grid
		for x := 0; x < windowWidth/AtomWidth; x++ {
			for y := 0; y < windowHeight/AtomWidth; y++ {
				if tempGrid[x][y] != nil {
					renderAtom(*tempGrid[x][y], imd)
					updatePostion(tempGrid[x][y], dt)
				}
			}
		}

		// for _, force := range forces {
		// 	renderForce(force, imd)
		// }

		// draw to screen
		win.Clear(WindowColor)
		imd.Draw(win)
		// grid.Draw(win)
		win.Update()
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
