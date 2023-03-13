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
	imd.Color = atom.atomType.color

	renderXPos := atom.currGridXPos * AtomWidth
	renderYPos := atom.currGridYPos * AtomWidth
	// renderXPos := int(math.Floor(atom.xPos/AtomWidth) * AtomWidth)
	// renderYPos := int(math.Floor(atom.yPos/AtomWidth) * AtomWidth)

	// imd.Push(pixel.V(float64(atom.xPos), float64(atom.yPos)))
	// imd.Push(pixel.V(float64(atom.xPos)+AtomWidth, float64(atom.yPos)+AtomWidth))
	imd.Push(pixel.V(float64(renderXPos)-float64(atom.atomType.extraWidth), float64(renderYPos)-float64(atom.atomType.extraWidth)))
	imd.Push(pixel.V(float64(renderXPos)+AtomWidth+float64(atom.atomType.extraWidth), float64(renderYPos)+AtomWidth+float64(atom.atomType.extraWidth)))
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
	imd.Push(pixel.V(player.xPos, player.yPos))
	imd.Circle(float64(player.width), 0)
}
