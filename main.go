package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const AtomWidth = 8
const AtomMass = 1
const Gravity = -3000
const Friction = 0.5

var (
	frames      = 0
	second      = time.Tick(time.Second)
	color       = pixel.RGB(1, 0, 0)
	WindowColor = pixel.RGB(0.2, 0.2, 0.2)
	atoms       []*Atom
	forces      []Force
	grid        [512][512]*Atom
	timeElapsed float64
)

type Atom struct {
	color        pixel.RGBA
	mass         float64
	xPos         float64
	yPos         float64
	xVel         float64
	yVel         float64
	currGridXPos int
	currGridYPos int
}

type Force struct {
	xComponent float64
	yComponent float64
}

func handleCollision(atom1 *Atom, atom2 *Atom) {

	// m1*v1 + m2*v2 = m1*v3 +m2*v4

	tempVel := 0.8 * atom1.yVel
	atom1.yVel = 0.8 * atom2.yVel
	atom2.yVel = tempVel

	atom1.yPos = float64(atom1.currGridYPos)
	atom2.yPos = float64(atom2.currGridYPos)

	// var atom1Pos int
	// var sameDirection = false
	// if atom1.currGridYPos > atom2.currGridYPos {
	// 	atom1Pos = 1
	// } else if atom1.currGridYPos < atom2.currGridYPos {
	// 	atom1Pos = 3
	// } else {
	// 	atom1Pos = 2

	// }
	// fmt.Printf("Atom1 position: %d", atom1Pos)
	// atom1.yVel = -atom1.yVel

	// switch atom1Pos {
	// case 1:
	// 	atom1.yVel = -atom1.yVel
	// case 3:
	// 	atom1.yVel = 2 * atom1.yVel
	// }
}

// set the atom's position and velocity
func updatePostion(atom *Atom, dt float64) {

	xForce := 0.0
	yForce := 0.0
	for _, force := range forces {
		xForce += force.xComponent
		yForce += force.yComponent
	}
	if atom.yPos <= 0 { // touching the ground
		// p = mv
		// fmt.Print("floor\n")
		atom.yVel = 0.8 * -atom.yVel
		atom.yPos = float64(atom.currGridYPos)
	}

	xAcc := xForce / atom.mass
	atom.xPos = atom.xPos + (atom.xVel * dt) + (0.5 * xAcc * (dt * dt))
	atom.xVel = atom.xVel + (xAcc * dt)

	yAcc := yForce / atom.mass
	atom.yPos = atom.yPos + (atom.yVel * dt) + (0.5 * yAcc * (dt * dt))
	atom.yVel = atom.yVel + (yAcc * dt)

	gridXPos := int(math.Max((math.Floor(atom.xPos/AtomWidth) * AtomWidth), 0))
	gridYPos := int(math.Max((math.Floor(atom.yPos/AtomWidth) * AtomWidth), 0))

	if grid[gridXPos][gridYPos] != nil && grid[gridXPos][gridYPos] != atom {
		// fmt.Print("collide")
		handleCollision(atom, grid[gridXPos][gridYPos])
	} else {
		grid[atom.currGridXPos][atom.currGridYPos] = nil
		grid[gridXPos][gridYPos] = atom
		atom.currGridXPos = gridXPos
		atom.currGridYPos = gridYPos
	}

}

func drawGrid() *imdraw.IMDraw {
	grid := imdraw.New(nil)
	grid.Color = WindowColor.Mul(pixel.RGB(0.8, 0.8, 0.8))
	for x := 0.0; x <= 512; x += AtomWidth {
		grid.Push(pixel.V(x, 0))
		grid.Push(pixel.V(x, 512))
		grid.Line(1)
		grid.Push(pixel.V(0, x))
		grid.Push(pixel.V(512, x))
		grid.Line(1)
	}
	return grid
}

// draw the atom to the screen
func renderAtom(atom Atom, imd *imdraw.IMDraw) {
	imd.Color = atom.color

	renderXPos := int(math.Floor(atom.xPos/AtomWidth) * AtomWidth)
	renderYPos := int(math.Floor(atom.yPos/AtomWidth) * AtomWidth)

	// imd.Push(pixel.V(float64(atom.xPos), float64(atom.yPos)))
	// imd.Push(pixel.V(float64(atom.xPos)+AtomWidth, float64(atom.yPos)+AtomWidth))
	imd.Push(pixel.V(float64(renderXPos), float64(renderYPos)))
	imd.Push(pixel.V(float64(renderXPos)+AtomWidth, float64(renderYPos)+AtomWidth))
	imd.Rectangle(0)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 512, 512),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	imd := imdraw.New(nil)

	win.Clear(colornames.Skyblue)

	// create forces
	forces = append(forces, Force{0, Gravity})

	last := time.Now()

	grid := drawGrid()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		imd.Clear() // clear the window

		// add new atom to screen
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			newColor := pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
			newAtom := Atom{
				newColor,
				10.0,
				win.MousePosition().X,
				win.MousePosition().Y,
				0,
				0,
				int(math.Floor(win.MousePosition().X/AtomWidth) * AtomWidth),
				int(math.Floor(win.MousePosition().Y/AtomWidth) * AtomWidth),
			}
			atoms = append(atoms, &newAtom)
		}

		// simulate every atom
		for _, atom := range atoms {
			renderAtom(*atom, imd)
			updatePostion(atom, dt)
		}

		// draw to screen
		win.Clear(WindowColor)
		imd.Draw(win)
		grid.Draw(win)
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
