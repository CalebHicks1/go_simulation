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
	// mass := 0.0
	rb.mass = 0.0
	for _, atom := range rb.atoms {

		xCOM += atom.xPos * atom.atomType.mass
		yCOM += atom.yPos * atom.atomType.mass
		rb.mass += atom.atomType.mass
	}
	rb.xPos = xCOM / math.Max(rb.mass, 1)
	rb.yPos = yCOM / math.Max(rb.mass, 1)

}

func rotateAndRenderRigidBody(rb *RigidBody, win *pixelgl.Window, imd *imdraw.IMDraw, dt float64) {

	// for every atom
	// if rb.rotated {
	// fmt.Printf("\rtorque: %f", rb.torque)

	tempXPos := rb.xPos
	// rb.xPos += 1
	// rb.xPos = win.MousePosition().X
	// rb.deltaX = rb.xPos - tempXPos

	tempYPos := rb.yPos
	// rb.yPos += 1
	// rb.yPos = win.MousePosition().Y
	// rb.deltaY = rb.yPos - tempYPos

	// angular acceleration
	aAcc := rb.torque / rb.momentOfInertia

	rb.rotation = (rb.angularVelocity * dt) + (0.5 * aAcc * (dt * dt))

	rb.angularVelocity = rb.angularVelocity + (aAcc * dt)

	// linear movement
	xAcc := rb.xForce / rb.mass
	rb.xPos = rb.xPos + (rb.xVel * dt) + (0.5 * xAcc * (dt * dt))
	rb.xVel = rb.xVel + (xAcc * dt)

	yAcc := rb.yForce / rb.mass
	rb.yPos = rb.yPos + (rb.yVel * dt) + (0.5 * yAcc * (dt * dt))
	rb.yVel = rb.yVel + (yAcc * dt)

	rb.deltaX = rb.xPos - tempXPos
	rb.deltaY = rb.yPos - tempYPos

	// reset some values
	rb.torque = 0.0
	// rb.rotated = false
	// rb.rotation = 0.0
	// rb.deltaX = 0.0
	// rb.deltaY = 0.0
	rb.xForce = 0
	rb.yForce = 0
	// calculateCenterOfMass(rb)

	// center of mass
	rb.xPos = rb.xCOM / math.Max(rb.mass, 1)
	rb.yPos = rb.yCOM / math.Max(rb.mass, 1)
	rb.xCOM = 0
	rb.yCOM = 0
	fmt.Printf("\r%d, %f, %f", rb.numAtoms, rb.mass, rb.prevMass)
	if rb.prevMass != rb.mass {
		if rb.rounds <= 2 {
			fmt.Print("first")
			rb.rounds += 1
			rb.prevMass = rb.mass
			// rb.mass = 0.0
		} else {

			for _, atom := range rb.atoms {
				atom.rigidBody = nil
			}
			for index, rigidBody := range rigidBodies {
				if rb == rigidBody {
					rigidBodies[index] = rigidBodies[len(rigidBodies)-1]
					rigidBodies = rigidBodies[:len(rigidBodies)-1]
				}
			}
		}

	} else {

		// rb.prevMass = rb.mass
		// rb.mass = 0.0
	}
	rb.prevMass = rb.mass
	rb.mass = 0.0
	rb.numAtoms = 0
	// rb.prevMass = rb.mass
	// rb.mass = 0.0

	// rb.yPos = win.MousePosition().Y
}

// }

// take all pixels in RigidBodyAtoms list and try to combine them into a rigid body.
func buildRigidBodies(win *pixelgl.Window, imd *imdraw.IMDraw) {

	for _, atom := range RigidBodyAtoms {

		// if the atom was not added to a rigid body, add a new one
		if atom.atomType.name == "rigidBody" {
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
					0.0,
					0.0,
					0.0,
					0.0,
					0.0,
					0.0,
					0.0,
					0.0,
					0.0,
					0.0,
					0,
					0,
				}
				rb.atoms = depthFirstSearch(atom, &rb)
				// atom.rigidBody = &rb
				rigidBodies = append(rigidBodies, &rb)
				calculateCenterOfMass(&rb)
				calculateMomentOfInertia(&rb)
				fmt.Printf("moment of inertia: %f\n", rb.momentOfInertia)
			}
		}
	}
	RigidBodyAtoms = []*Atom{}
}
