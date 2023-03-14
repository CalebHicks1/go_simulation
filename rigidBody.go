package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func depthFirstSearch(atom *Atom, rb *RigidBody) []*Atom {

	// fmt.Printf("Current Atom Location: (%d, %d)\n", atom.currGridXPos, atom.currGridYPos)
	// check up, down, left, right

	// returnString := fmt.Sprintf("(%d, %d)", atom.currGridXPos, atom.currGridYPos)
	atom.rigidBody = rb
	atom.color = pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
	currX := atom.currGridXPos
	currY := atom.currGridYPos

	rigidBodyAtoms := []*Atom{atom}

	// check up
	if grid[currX][currY+1] != nil && grid[currX][currY+1].rigidBody == nil && grid[currX][currY+1].atomType.name == "rigidBody" {
		rigidBodyAtoms = append(rigidBodyAtoms, depthFirstSearch(grid[currX][currY+1], rb)...)
	}

	// check right
	if grid[currX+1][currY] != nil && grid[currX+1][currY].rigidBody == nil && grid[currX+1][currY].atomType.name == "rigidBody" {
		rigidBodyAtoms = append(rigidBodyAtoms, depthFirstSearch(grid[currX+1][currY], rb)...)
	}

	// check down
	if grid[currX][currY-1] != nil && grid[currX][currY-1].rigidBody == nil && grid[currX][currY-1].atomType.name == "rigidBody" {
		rigidBodyAtoms = append(rigidBodyAtoms, depthFirstSearch(grid[currX][currY-1], rb)...)
	}

	// check left
	if grid[currX-1][currY] != nil && grid[currX-1][currY].rigidBody == nil && grid[currX-1][currY].atomType.name == "rigidBody" {
		rigidBodyAtoms = append(rigidBodyAtoms, depthFirstSearch(grid[currX-1][currY], rb)...)
	}
	return rigidBodyAtoms

}

func calculateMomentOfInertia(rb *RigidBody) {
	moi := 0.0
	for _, atom := range rb.atoms {

		xDist := atom.xPos - rb.xPos
		yDist := atom.yPos - rb.yPos

		// todo: redundant calculation?
		distance := math.Sqrt(math.Pow(xDist, 2) + math.Pow(yDist, 2))
		moi += atom.atomType.mass * math.Pow(distance, 2)
	}
	rb.momentOfInertia = moi
}

func calculateCenterOfMass(rb *RigidBody) {
	xCOM := 0.0
	yCOM := 0.0
	mass := 0.0
	for _, atom := range rb.atoms {

		xCOM += atom.xPos * atom.atomType.mass
		yCOM += atom.yPos * atom.atomType.mass
		mass += atom.atomType.mass
	}
	rb.xPos = xCOM / math.Max(mass, 1)
	rb.yPos = yCOM / math.Max(mass, 1)

}

func rotateAndRenderRigidBody(rb *RigidBody, imd *imdraw.IMDraw, dt float64) {

	// for every atom
	// if rb.rotated {
	// fmt.Printf("\rtorque: %f", rb.torque)

	// angular acceleration
	aAcc := rb.torque / rb.momentOfInertia

	rb.rotation = (rb.angularVelocity * dt) + (0.5 * aAcc * (dt * dt))

	rb.angularVelocity = rb.angularVelocity + (aAcc * dt)
	// atom.yPos = atom.yPos + (atom.yVel * dt) + (0.5 * yAcc * (dt * dt))

	// atom.yVel = atom.yVel + (yAcc * dt)
	//theta = wt + 1/2 \alpha t^{2}

	for _, atom := range rb.atoms {

		newXPos := ((atom.xPos - rb.xPos) * math.Cos(rb.rotation)) - ((atom.yPos - rb.yPos) * math.Sin(rb.rotation)) + rb.xPos
		newYPos := ((atom.xPos - rb.xPos) * math.Sin(rb.rotation)) + ((atom.yPos - rb.yPos) * math.Cos(rb.rotation)) + rb.yPos

		atom.xPos = newXPos
		atom.yPos = newYPos

		gridXPos := int(math.Floor(atom.xPos / AtomWidth))
		gridYPos := int(math.Floor(atom.yPos / AtomWidth))

		// renderXPos := 0.0
		// renderYPos := 0.0

		// if grid[gridXPos][gridYPos] == nil {

		grid[atom.currGridXPos][atom.currGridYPos] = nil
		atom.currGridXPos = gridXPos
		atom.currGridYPos = gridYPos
		grid[atom.currGridXPos][atom.currGridYPos] = atom

		// }

		// renderXPos = float64(atom.currGridXPos * AtomWidth)
		// renderYPos = float64(atom.currGridYPos * AtomWidth)
		// imd.Color = atom.atomType.color
		// imd.Push(pixel.V(float64(renderXPos)-float64(atom.atomType.extraWidth), float64(renderYPos)-float64(atom.atomType.extraWidth)))
		// imd.Push(pixel.V(float64(renderXPos)+AtomWidth+float64(atom.atomType.extraWidth), float64(renderYPos)+AtomWidth+float64(atom.atomType.extraWidth)))
		// imd.Rectangle(0)

		// fmt.Printf("atom.xPos: %f, atom.yPos: %f, (%d, %d)\n", atom.xPos, atom.yPos, atom.currGridXPos, atom.currGridYPos)
	}
	rb.torque = 0.0
	rb.rotated = false
	rb.rotation = 0.0
}

// }

// take all pixels in RigidBodyAtoms list and try to combine them into a rigid body.
func buildRigidBodies(win *pixelgl.Window, imd *imdraw.IMDraw) {

	for _, atom := range RigidBodyAtoms {

		// if the atom was not added to a rigid body, add a new one
		if atom.rigidBody == nil {
			rb := RigidBody{
				0.0,
				0.0,
				0.0,
				nil,
				pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64()),
				true,
				0.0,
				0.0,
				0.0,
			}
			rb.atoms = depthFirstSearch(atom, &rb)
			// atom.rigidBody = &rb
			rigidBodies = append(rigidBodies, &rb)
			calculateCenterOfMass(&rb)
			calculateMomentOfInertia(&rb)
			fmt.Printf("moment of inertia: %f\n", rb.momentOfInertia)
		}
	}
	RigidBodyAtoms = []*Atom{}
}
