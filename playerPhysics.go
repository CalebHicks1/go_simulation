package main

import (
	"math"

	"github.com/faiface/pixel/pixelgl"
)

func controlPlayer(win pixelgl.Window, player *Player, dt float64) {
	// if win.Pressed(pixelgl.KeyD) {
	// 	// fmt.Print("New Force\n")
	// 	player.xVel = 100

	// }

	switch {
	case win.Pressed(pixelgl.KeyD):
		player.xAcc = 200.0
	case win.Pressed(pixelgl.KeyA):
		player.xAcc = -200.0
	default:
		player.xAcc = 0.0
	}
	if win.Pressed(pixelgl.KeyW) {
		player.yVel = 500
	}

	player.xVel = player.xVel*(1-dt*6) + player.xAcc*(dt*6)
	player.xPos += player.xVel * dt
}

func simulatePlayer(player *Player, dt float64) {

	xForce := 0.0
	yForce := 0.0
	for _, force := range forces {

		xComp := 0.0
		yComp := 0.0
		if force.source == "gravity" {
			if gravEnabled {

				yComp = (force.Strength * player.mass * 100) * -1
			}
		} else {
			distance := math.Max(math.Sqrt(math.Pow(force.xPos-player.xPos, 2)+math.Pow(force.yPos-player.yPos, 2))/200, 1)

			forceAfterDistance := force.Strength / math.Pow(distance, 2)

			angle := math.Atan2(force.xPos-player.xPos, force.yPos-player.yPos)
			// xComp := forceAfterDistance * math.Sin(angle)
			// yComp := forceAfterDistance * math.Cos(angle)
			xComp = forceAfterDistance * math.Sin(angle)
			yComp = forceAfterDistance * math.Cos(angle)

			// fmt.Printf("x dist: %f\ny dist: %f\n angle: %f\n xComp: %f\n\n", force.xPos-atom.xPos, force.yPos-atom.yPos, angle, xComp)

			// fmt.Print(force.source)
			// fmt.Printf("%f %f %f, %f\n", math.Sin(angle), math.Cos(angle), xComp, yComp)
		}
		// fmt.Printf("%f\n", xComp)

		xForce += xComp
		yForce += yComp
		// fmt.Printf("force: %f, %f\n", xForce, yForce)
	}

	// calculate the Atom's x and y acceleration and velocity
	xAcc := xForce / player.mass
	player.xPos = player.xPos + (player.xVel * dt) + (0.5 * xAcc * (dt * dt))
	player.xVel = player.xVel + (xAcc * dt)

	yAcc := yForce / player.mass
	player.yPos = player.yPos + (player.yVel * dt) + (0.5 * yAcc * (dt * dt))
	player.yVel = player.yVel + (yAcc * dt)

	// calculate the atom's new grid location

	// gridXPos := int(math.Min(math.Max((math.Floor(atom.xPos/AtomWidth)), 0), windowWidth/AtomWidth))
	// gridYPos := int(math.Min(math.Max((math.Floor(atom.yPos/AtomWidth)), 0), windowHeight/AtomWidth))

	// check for window boundary collisions at floor and ceiling.
	if player.yPos-float64(player.width) <= 0 && player.yVel < 0 /*|| gridYPos >= windowHeight/AtomWidth*/ {

		// set the Atom's location to its current grid location, to move it back a little
		player.yVel = -0.2 * player.yVel
		player.yPos = float64(player.width)
	}

	// check for window boundary collisions at walls.
	if player.xPos <= 0 && player.xVel < 0 /*|| (gridXPos >= windowWidth/AtomWidth && atom.xVel > 0)*/ {

		// set the Atom's location to its current grid location, to move it back a little
		player.xVel = -player.xVel
		player.xPos -= 1
	}

	// fmt.Printf("New grid location: (%d, %d)\n", gridXPos, gridYPos)

	// check if the new grid location is occupied

	// if grid[gridXPos][gridYPos] != nil && grid[gridXPos][gridYPos] != atom {
	// 	// fmt.Print("collide")
	// 	handleCollision(atom, grid[gridXPos][gridYPos])
	// 	if atom.currGridXPos != gridXPos {
	// 		if grid[gridXPos][atom.currGridYPos] == nil {
	// 			grid[atom.currGridXPos][atom.currGridYPos] = nil
	// 			atom.currGridXPos = gridXPos
	// 			grid[atom.currGridXPos][atom.currGridYPos] = atom
	// 			atom.xPos = float64(gridXPos * AtomWidth)
	// 		} else {

	// 			atom.xPos = float64(atom.currGridXPos * AtomWidth)
	// 		}
	// 	}
	// 	if atom.currGridYPos != gridYPos {
	// 		if grid[atom.currGridXPos][gridYPos] == nil {
	// 			grid[atom.currGridXPos][atom.currGridYPos] = nil
	// 			atom.currGridYPos = gridYPos
	// 			grid[atom.currGridXPos][atom.currGridYPos] = atom
	// 			atom.yPos = float64(gridYPos * AtomWidth)
	// 		} else {

	// 			atom.yPos = float64(atom.currGridYPos * AtomWidth)
	// 		}
	// 	}
	// } else {
	// 	grid[atom.currGridXPos][atom.currGridYPos] = nil
	// 	grid[gridXPos][gridYPos] = atom
	// 	atom.currGridXPos = gridXPos
	// 	atom.currGridYPos = gridYPos
	// }

}
