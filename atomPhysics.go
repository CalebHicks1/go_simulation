package main

import (
	"math"
	"math/rand"
)

func updatePostion(atom *Atom, dt float64) {

	switch atom.atomType.name {
	case "water":
		simulateWater(atom, dt)
	// case "stone":
	// 	simulateStone(atom, dt)
	default:
		break
	}
	simulateAtom(atom, dt)
}

func handleCollision(atom1 *Atom, atom2 *Atom) {

	// v1f = (((m1-m2)/(m1+m2))*v1i) + (((2*m2) / (m1+m2))*v2i)
	// v2f = (((2*m1)/(m1+m2)) * v1i) + (((m2-m1) / (m1+m2))*v2i)

	// from https://www2.tntech.edu/leap/murdock/books/v1chap7.pdf

	// TODO: try to separate x and y, don't want a y collision to stop x movement

	// y

	tempYVel := atom1.atomType.collisionElasticity*(((2*atom1.atomType.mass)/(atom1.atomType.mass+atom2.atomType.mass))*atom1.yVel) + (((atom2.atomType.mass - atom1.atomType.mass) / (atom1.atomType.mass + atom2.atomType.mass)) * atom2.yVel)
	atom1.yVel = atom2.atomType.collisionElasticity*(((atom1.atomType.mass-atom2.atomType.mass)/(atom1.atomType.mass+atom2.atomType.mass))*atom1.yVel) + (((2 * atom2.atomType.mass) / (atom1.atomType.mass + atom2.atomType.mass)) * atom2.yVel)
	atom2.yVel = tempYVel

	atom1.yPos = float64(atom1.currGridYPos * AtomWidth)
	// atom2.yPos = float64(atom2.currGridYPos * AtomWidth)

	// x

	tempXVel := atom1.atomType.collisionElasticity*(((2*atom1.atomType.mass)/(atom1.atomType.mass+atom2.atomType.mass))*atom1.xVel) + (((atom2.atomType.mass - atom1.atomType.mass) / (atom1.atomType.mass + atom2.atomType.mass)) * atom2.xVel)
	atom1.xVel = atom2.atomType.collisionElasticity*(((atom1.atomType.mass-atom2.atomType.mass)/(atom1.atomType.mass+atom2.atomType.mass))*atom1.xVel) + (((2 * atom2.atomType.mass) / (atom1.atomType.mass + atom2.atomType.mass)) * atom2.xVel)
	atom2.xVel = tempXVel

	atom1.xPos = float64(atom1.currGridXPos * AtomWidth)
	// atom2.xPos = float64(atom2.currGridXPos * AtomWidth)
	// atom1.xPos = float64((atom1.currGridXPos * AtomWidth) + AtomWidth/2)
	// atom2.xPos = float64((atom2.currGridXPos * AtomWidth) + AtomWidth/2)

}

// set the atom's position and velocity
func simulateAtom(atom *Atom, dt float64) {

	xForce := 0.0
	yForce := 0.0
	for _, force := range forces {

		xComp := 0.0
		yComp := 0.0
		if gravEnabled && force.source == "gravity" {
			yComp = (force.Strength * atom.atomType.mass * 100) * -1
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

	// calculate the Atom's x and y acceleration and velocity
	xAcc := xForce / atom.atomType.mass
	atom.xPos = atom.xPos + (atom.xVel * dt) + (0.5 * xAcc * (dt * dt))
	atom.xVel = atom.xVel + (xAcc * dt)

	yAcc := yForce / atom.atomType.mass
	atom.yPos = atom.yPos + (atom.yVel * dt) + (0.5 * yAcc * (dt * dt))
	atom.yVel = atom.yVel + (yAcc * dt)

	// calculate the atom's new grid location

	gridXPos := int(math.Min(math.Max((math.Floor(atom.xPos/AtomWidth)), 0), windowWidth/AtomWidth))
	gridYPos := int(math.Min(math.Max((math.Floor(atom.yPos/AtomWidth)), 0), windowHeight/AtomWidth))

	// check for window boundary collisions at floor and ceiling.
	if (atom.yPos <= 0 && atom.yVel < 0) || gridYPos >= windowHeight/AtomWidth {

		// set the Atom's location to its current grid location, to move it back a little
		atom.yVel = atom.atomType.borderElasticity * -atom.yVel
		atom.yPos = float64(atom.currGridYPos * AtomWidth)
	}

	// check for window boundary collisions at walls.
	if (atom.xPos <= 0 && atom.xVel < 0) || (gridXPos >= windowWidth/AtomWidth && atom.xVel > 0) {

		// set the Atom's location to its current grid location, to move it back a little
		atom.xVel = atom.atomType.borderElasticity * -atom.xVel
		atom.xPos = float64(atom.currGridXPos * AtomWidth)
	}

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

// set the atom's position and velocity
func simulateWater(atom *Atom, dt float64) {
	if math.Abs(atom.xVel) < 50 {

		if rand.Float64() > 0.5 {

			atom.xVel = 100
		} else {
			atom.xVel = -100
		}
	}
}
