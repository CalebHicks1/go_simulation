package main

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const AtomWidth = 4
const AtomMass = 0.5
const Gravity = 1000
const Friction = 0.5
const Elasticity = 0.2
const particleCollision = 0.97
const windowWidth = 1200 // must be a multiple of AtomWidth
const windowHeight = 600
const mouseGrav = 2000
const extraAtomWidth = 2

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
	AtomTypes   = []string{"water", "stone"}
)

type Atom struct {
	color        pixel.RGBA
	mass         float64
	xPos         float64 // the actual x location of the atom
	yPos         float64 // the actual y location of the atom
	xVel         float64
	yVel         float64
	currGridXPos int // the x location on the pixel grid of the atom
	currGridYPos int // the y location on the pixel grid of the atom
}

type Force struct {
	source     string
	xComponent float64
	yComponent float64
	Strength   float64
	xPos       float64
	yPos       float64
}

func handleCollision(atom1 *Atom, atom2 *Atom) {

	// m1*v1 + m2*v2 = m1*v3 +m2*v4

	tempYVel := particleCollision * atom1.yVel
	atom1.yVel = particleCollision * atom2.yVel
	atom2.yVel = tempYVel

	// tempYVel := Elasticity * atom1.yVel
	// atom1.yVel = Elasticity * atom2.yVel
	// atom2.yVel = tempYVel

	atom1.yPos = float64(atom1.currGridYPos * AtomWidth)
	atom2.yPos = float64(atom2.currGridYPos * AtomWidth)

	tempXVel := particleCollision * atom1.xVel
	atom1.xVel = particleCollision * atom2.xVel
	atom2.xVel = tempXVel

	// tempXVel := Elasticity * atom1.xVel
	// atom1.xVel = Elasticity * atom2.xVel
	// atom2.xVel = tempXVel

	atom1.xPos = float64(atom1.currGridXPos * AtomWidth)
	atom2.xPos = float64(atom2.currGridXPos * AtomWidth)

}

// set the atom's position and velocity
func updatePostion(atom *Atom, dt float64) {

	xForce := 0.0
	yForce := 0.0
	for _, force := range forces {

		// calculate direction to force
		// xDirection := 0.0
		// yDirection := 0.0
		// if force.xPos-atom.xPos != 0 {

		// 	xDirection = (force.xPos - atom.xPos) / math.Abs(force.xPos-atom.xPos)
		// }
		// if force.yPos-atom.yPos != 0 {

		// 	yDirection = (force.yPos - atom.yPos) / math.Abs(force.yPos-atom.yPos)
		// }

		// calculate angle towards force

		xComp := 0.0
		yComp := 0.0
		if gravEnabled && force.source == "gravity" {
			yComp = force.Strength * -1
		} else {
			distance := math.Max(math.Sqrt(math.Pow(force.xPos-atom.xPos, 2)+math.Pow(force.yPos-atom.yPos, 2))/200, 1)

			forceAfterDistance := force.Strength / math.Pow(distance, 2)

			angle := math.Atan2(force.xPos-atom.xPos, force.yPos-atom.yPos)
			// xComp := forceAfterDistance * math.Sin(angle)
			// yComp := forceAfterDistance * math.Cos(angle)
			xComp = forceAfterDistance * math.Sin(angle)
			yComp = forceAfterDistance * math.Cos(angle)
		}
		// fmt.Printf("%f\n", xComp)
		// fmt.Printf("%f %f %f, %f\n", math.Sin(angle), math.Cos(angle), xComp, yComp)

		xForce += xComp
		yForce += yComp
		// fmt.Printf("force: %f, %f\n", xForce, yForce)
	}

	// check for window boundary collisions at floor and ceiling.
	if (atom.yPos <= 0 && atom.yVel < 0) || atom.yPos >= windowHeight {
		// p = mv
		// fmt.Print("floor\n")

		// set the Atom's location to its current grid location, to move it back a little
		atom.yVel = Elasticity * -atom.yVel
		atom.yPos = float64(atom.currGridYPos * AtomWidth)
	}

	// check for window boundary collisions at walls.
	if (atom.xPos <= 0 && atom.xVel < 0) || (atom.xPos >= windowWidth && atom.xVel > 0) {

		// set the Atom's location to its current grid location, to move it back a little
		atom.xVel = Elasticity * -atom.xVel
		atom.xPos = float64(atom.currGridXPos * AtomWidth)
	}

	// calculate the Atom's x and y acceleration and velocity
	xAcc := xForce / atom.mass
	atom.xPos = atom.xPos + (atom.xVel * dt) + (0.5 * xAcc * (dt * dt))
	atom.xVel = atom.xVel + (xAcc * dt)

	yAcc := yForce / atom.mass
	atom.yPos = atom.yPos + (atom.yVel * dt) + (0.5 * yAcc * (dt * dt))
	atom.yVel = atom.yVel + (yAcc * dt)

	// calculate the atom's new grid location

	gridXPos := int(math.Min(math.Max((math.Floor(atom.xPos/AtomWidth)), 0), windowWidth/AtomWidth))
	gridYPos := int(math.Min(math.Max((math.Floor(atom.yPos/AtomWidth)), 0), windowHeight/AtomWidth))

	// fmt.Printf("New grid location: (%d, %d)\n", gridXPos, gridYPos)

	// check if the new grid location is occupied

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

// draw the atom to the screen
func renderAtom(atom Atom, imd *imdraw.IMDraw) {
	imd.Color = atom.color

	renderXPos := atom.currGridXPos * AtomWidth
	renderYPos := atom.currGridYPos * AtomWidth
	// renderXPos := int(math.Floor(atom.xPos/AtomWidth) * AtomWidth)
	// renderYPos := int(math.Floor(atom.yPos/AtomWidth) * AtomWidth)

	// imd.Push(pixel.V(float64(atom.xPos), float64(atom.yPos)))
	// imd.Push(pixel.V(float64(atom.xPos)+AtomWidth, float64(atom.yPos)+AtomWidth))
	imd.Push(pixel.V(float64(renderXPos)-extraAtomWidth, float64(renderYPos)-extraAtomWidth))
	imd.Push(pixel.V(float64(renderXPos)+AtomWidth+extraAtomWidth, float64(renderYPos)+AtomWidth+extraAtomWidth))
	imd.Rectangle(0)
}

func renderForce(force Force, imd *imdraw.IMDraw) {
	imd.Color = pixel.RGB(1, 0, 0)

	renderXPos := int(math.Floor(force.xPos/AtomWidth) * AtomWidth)
	renderYPos := int(math.Floor(force.yPos/AtomWidth) * AtomWidth)

	// imd.Push(pixel.V(float64(atom.xPos), float64(atom.yPos)))
	// imd.Push(pixel.V(float64(atom.xPos)+AtomWidth, float64(atom.yPos)+AtomWidth))
	imd.Push(pixel.V(float64(renderXPos), float64(renderYPos)-1000))
	imd.Push(pixel.V(float64(renderXPos), float64(renderYPos)+1000000))
	imd.Line(1)
	imd.Push(pixel.V(float64(renderXPos)-1000, float64(renderYPos)))
	imd.Push(pixel.V(float64(renderXPos)+1000, float64(renderYPos)))
	// imd.Rectangle(0)
	imd.Line(1)
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

		// add new atom to screen
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			// newColor := pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
			// newColor := pixel.RGB(0.7, 1, 0.7)
			// newColor := pixel.RGB(0.2, 0.2, 1).Mul(pixel.Alpha(0.2))
			for x := -10.0; x < 10; x++ {
				for y := -10.0; y < 10; y++ {

					// fmt.Printf("MousePosition: (%f, %f)\n", win.MousePosition().X+x*AtomWidth, win.MousePosition().Y+y*AtomWidth)
					// fmt.Printf("GridPosition: (%d, %d)\n", int(math.Floor(win.MousePosition().X/AtomWidth)), int(math.Floor(win.MousePosition().Y/AtomWidth)))

					// create a new atom at the location of the mouse.
					// calculate the grid location of the atom

					gridX := int((math.Floor(win.MousePosition().X+x*AtomWidth) / AtomWidth))
					gridY := int((math.Floor(win.MousePosition().Y+y*AtomWidth) / AtomWidth))

					if gridX >= 0 && gridX < windowWidth/AtomWidth && gridY >= 0 && gridY < windowHeight/AtomWidth {
						if grid[gridX][gridY] == nil {
							newAtom := Atom{
								color,
								AtomMass,
								win.MousePosition().X + x*AtomWidth,
								win.MousePosition().Y + y*AtomWidth,
								0,
								0,
								gridX,
								gridY,
							}
							atoms = append(atoms, &newAtom)
							grid[gridX][gridY] = &newAtom

						}
					}

				}
			}
		}

		if win.JustPressed(pixelgl.MouseButtonRight) {
			// fmt.Print("New Force\n")
			newForce := Force{
				"mouse",
				0,
				0,
				mouseGrav,
				win.MousePosition().X,
				win.MousePosition().Y,
			}
			forces = append(forces, newForce)
		}

		if win.JustPressed(pixelgl.KeyG) {
			gravEnabled = !gravEnabled
		}

		if win.Pressed(pixelgl.KeyV) {
			// fmt.Print("New Force\n")
			newForce := Force{
				"mouse",
				0,
				0,
				mouseGrav,
				win.MousePosition().X,
				win.MousePosition().Y,
			}
			forces = []Force{}
			forces = append(forces, newForce)
			forces = append(forces, Force{"gravity", 0, 0, Gravity, windowWidth / 2, -100000})
		}

		if win.JustReleased(pixelgl.KeyV) {
			forces = []Force{}
			forces = append(forces, Force{"gravity", 0, 0, Gravity, windowWidth / 2, -100000})

		}

		if win.JustPressed(pixelgl.KeyC) {
			fmt.Print("Cleared Force\n")

			forces = []Force{}
			forces = append(forces, Force{"gravity", 0, 0, Gravity, windowWidth / 2, -100000})
		}

		// simulate every atom
		// for _, atom := range atoms {
		// 	renderAtom(*atom, imd)
		// 	updatePostion(atom, dt)
		// }

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
