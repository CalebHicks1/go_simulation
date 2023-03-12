package main

import (
	"fmt"
	"math"

	"github.com/faiface/pixel/pixelgl"
)

func handleKeyPress(win pixelgl.Window) {

	// add new atom to screen
	if win.Pressed(pixelgl.MouseButtonLeft) {
		for x := -10.0; x < 10; x++ {
			for y := -10.0; y < 10; y++ {

				// create a new atom at the location of the mouse.
				// calculate the grid location of the atom

				gridX := int((math.Floor(win.MousePosition().X+x*AtomWidth) / AtomWidth))
				gridY := int((math.Floor(win.MousePosition().Y+y*AtomWidth) / AtomWidth))

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

	if win.JustPressed(pixelgl.KeyT) {

		currType = (currType + 1) % len(AtomTypes)
		fmt.Printf("Atom Type: %s\n", AtomTypes[currType].name)

	}
}
