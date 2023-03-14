package main

import (
	"image"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// draw the atom to the screen
func renderAtom(atom Atom, imd *imdraw.IMDraw) {

	renderXPos := 0.0
	renderYPos := 0.0
	if atom.rigidBody == nil {
		imd.Color = atom.atomType.color
		renderXPos = float64(atom.currGridXPos * AtomWidth)
		renderYPos = float64(atom.currGridYPos * AtomWidth)
	} else {
		imd.Color = atom.rigidBody.color
		renderXPos = math.Floor(atom.xPos/AtomWidth) * AtomWidth
		renderYPos = math.Floor(atom.yPos/AtomWidth) * AtomWidth
		// renderXPos = atom.xPos
		// renderYPos = atom.yPos
	}
	imd.Push(pixel.V(float64(renderXPos)-float64(atom.atomType.extraWidth), float64(renderYPos)-float64(atom.atomType.extraWidth)))
	imd.Push(pixel.V(float64(renderXPos)+AtomWidth+float64(atom.atomType.extraWidth), float64(renderYPos)+AtomWidth+float64(atom.atomType.extraWidth)))
	imd.Rectangle(0)
}

func renderRigidBody(rb *RigidBody, imd *imdraw.IMDraw) {
	// y := windowHeight / 2
	imd.Color = pixel.RGB(1, 0, 0)
	imd.Push(pixel.V(rb.xPos, rb.yPos))
	imd.Circle(10, 0)
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

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func renderPlayer(imd *imdraw.IMDraw, player Player) {
	imd.Color = player.color
	imd.Push(pixel.V(player.xPos, player.yPos-float64(player.width)))
	imd.Push(pixel.V(player.xPos+float64(player.width), player.yPos+float64(player.width)))
	// imd.Circle(float64(player.width), 0)
	imd.Rectangle(0)
}
