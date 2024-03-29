package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func handleKeyPress(win pixelgl.Window, imd *imdraw.IMDraw) {

	width := math.Floor(cursorSize)

	// add new atom to screen
	if win.Pressed(pixelgl.MouseButtonLeft) {
		for x := -width; x < width; x++ {
			for y := -width; y < width; y++ {

				// create a new atom at the location of the mouse.
				// calculate the grid location of the atom

				gridX := int((math.Floor(win.MousePosition().X) / AtomWidth) + x)
				gridY := int((math.Floor(win.MousePosition().Y) / AtomWidth) + y)

				if gridX >= 0 && gridX < windowWidth/AtomWidth && gridY >= 0 && gridY < windowHeight/AtomWidth {
					if grid[gridX][gridY] == nil {
						newAtom := Atom{
							AtomTypes[currType],
							win.MousePosition().X + x*AtomWidth,
							win.MousePosition().Y + y*AtomWidth,
							0,
							0,
							gridX,
							gridY,
							"particle",
							nil,
							pixel.RGB(0, 0, 0),
						}
						atoms = append(atoms, &newAtom)
						grid[gridX][gridY] = &newAtom

					}
				}

			}
		}
	}

	// add new atom to screen
	if win.Pressed(pixelgl.KeyE) {

		width := math.Floor(cursorSize)

		imd.Color = pixel.RGB(1, 0, 0)
		imd.Push(pixel.V(win.MousePosition().X-width*AtomWidth, win.MousePosition().Y-width*AtomWidth))
		imd.Push(pixel.V(win.MousePosition().X+width*AtomWidth, win.MousePosition().Y+width*AtomWidth))
		imd.Rectangle(1)

		for x := -width; x < width; x++ {
			for y := -width; y < width; y++ {

				// create a new atom at the location of the mouse.
				// calculate the grid location of the atom

				gridX := int((math.Floor(win.MousePosition().X) / AtomWidth) + x)
				gridY := int((math.Floor(win.MousePosition().Y) / AtomWidth) + y)

				if gridX >= 0 && gridX < windowWidth/AtomWidth && gridY >= 0 && gridY < windowHeight/AtomWidth {
					if grid[gridX][gridY] != nil {

						grid[gridX][gridY] = nil

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

	if win.JustPressed(pixelgl.KeyR) {
		buildRigidBodies(&win, imd)
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

	if win.JustPressed(pixelgl.KeyT) {

		currType = (currType + 1) % len(AtomTypes)
		fmt.Printf("Atom Type: %s\n", AtomTypes[currType].name)
		typeText.Clear()
		typeText.Color = color.Black
		fmt.Fprintf(typeText, "Atom Type: ")
		typeText.Color = AtomTypes[currType].color
		fmt.Fprintf(typeText, "%s", AtomTypes[currType].name)
	}
}
