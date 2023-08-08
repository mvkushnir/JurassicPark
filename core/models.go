package core

import (
	"golang.org/x/exp/slices"
)

type Species string

const (
	Tyrannosaurus Species = "Tyrannosaurus"
	Velociraptor          = "Velociraptor"
	Spinosaurus           = "Spinosaurus"
	Megalosaurus          = "Megalosaurus"
	Brachiosaurus         = "Brachiosaurus"
	Stegosaurus           = "Stegosaurus"
	Ankylosaurus          = "Ankylosaurus"
	Triceratops           = "Triceratops"
)

var Carnivores = []Species{Tyrannosaurus, Velociraptor, Spinosaurus, Megalosaurus}

type Dinosaur struct {
	Name    string  `form:"name" json:"name" xml:"name"  binding:"required"`
	Species Species `form:"species" json:"species" xml:"species"  binding:"required"`
	Cage    int     `form:"cage" json:"cage" xml:"cage"  binding:"required"`
}

func (d Dinosaur) isCarnivore() bool {
	return slices.Contains(Carnivores, d.Species)
}

type PowerStatus string

const (
	ACTIVE PowerStatus = "ACTIVE"
	DOWN               = "DOWN"
)

type Cage struct {
	Number      int         `form:"number" json:"number" xml:"number"  binding:"required"`
	PowerStatus PowerStatus `form:"powerStatus" json:"powerStatus" xml:"powerStatus"  binding:"required"`
	Capacity    int         `form:"capacity" json:"capacity" xml:"capacity"  binding:"required"`
}
