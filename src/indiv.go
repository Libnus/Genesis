package main

import (
	"math/rand"
	"image/color"
	// "fmt"
)

type Mite struct {
	nnet *NeuralNetwork
	genome []string
	id int
	X, Y int
	color color.Color
}

func createRandomMite(numGenes int, id int) *Mite {
	genome := createGenome(numGenes, int64(rand.Int()))
	nnet, _ := processGenome(genome)

	var x, y = 20, 64
	for{
		if _, ok := gridOccupy[y*128 + x]; ok {
			x = rand.Intn(128)
			y = rand.Intn(128)
		} else{ break }
	}

	gridOccupy[y*128+x] = true

	red := uint8(rand.Intn(256))
	green := uint8(rand.Intn(256))
	blue := uint8(rand.Intn(256))
	alpha := uint8(rand.Intn(256))

	return &Mite{
		nnet: nnet,
		genome: genome,
		id: id,
		X: x,
		Y: y,
		color: color.RGBA{red,green,blue,alpha},
	}
}
func randomizePos(mite *Mite) {
	var x, y = rand.Intn(128), rand.Intn(128)
	for{
		if _, ok := gridOccupy[y*128 + x]; ok {
			x = rand.Intn(128)
			y = rand.Intn(128)
		} else{ break }
	}

	gridOccupy[y*128+x] = true
	mite.X, mite.Y = x, y
}

func createMite(genome []string) *Mite {
	nnet, _ := processGenome(genome)

	// TODO move to get random color function
	red := uint8(rand.Intn(256))
	green := uint8(rand.Intn(256))
	blue := uint8(rand.Intn(256))
	alpha := uint8(rand.Intn(256))

	newMite := &Mite{
		nnet: nnet,
		genome: genome,
		id: rand.Intn(10000), // TODO fix lol
		X: 0,
		Y: 0,
		color: color.RGBA{red,green,blue,alpha},
	}

	randomizePos(newMite)
	return newMite
}

