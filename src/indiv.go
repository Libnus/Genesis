package main

import (
	"math/rand"
	"image/color"
	// "fmt"
	// "time"
)

type Mite struct {
	nnet *NeuralNetwork
	genome []string
	id int
	X, Y int
	color color.Color
}

// NOTE:: iterative
func createRandomMite(numGenes int, id int) *Mite {
	genome := createGenome(numGenes, int64(rand.Int()))
	nnet, _ := processGenome(genome)

	var x, y = 0, 0
	for{
		if gridOccupy[x][y] {
			x = rand.Intn(ROWS)
			y = rand.Intn(COLS)
		} else{ break }
	}

	gridOccupy[x][y] = true

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
	var x, y = rand.Intn(ROWS), rand.Intn(COLS)
	for{
		if gridOccupy[x][y] {
			x = rand.Intn(ROWS)
			y = rand.Intn(COLS)
		} else{ break }
	}

	gridOccupy[x][y] = true
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

// constructors ^^^


// NOTE:: does not lock new position on occupancy grid.
// 		  assumed to be done already due to collision checking
//
// 		  NEVER EVER call this function unless it is certain the grid location is not occupied
// func moveMite(indiv *Mite, oldPos [2]int, newPos [2]int){

// 	gridMu[ oldPos[0] ][ oldPos[1] ].Lock()
// 	gridOccupy[ oldPos[0] ][ oldPos[1] ] = false
// 	gridMu[ oldPos[0] ][ oldPos[1] ].Unlock()
// 	// ^ could be bad, relinquishing control of the grid location before the mite has been secured on the new square
// 	// but again we assume the new pos is already locked and ready to go

// 	indiv.X = newPos[0]
// 	indiv.Y = newPos[1]

// 	gridOccupy[ newPos[0] ][ newPos[1] ] = true
// }

// move mite to x and y if possible
func moveMite(indiv *Mite, x int, y int) {
	// check collisions
	if x < 0 || x >= ROWS { return }
	if y < 0 || y >= COLS { return }

	// check for other mites in the new square

	// startTime := time.Now()
	gridMu[x][y].Lock()
	// endTime := time.Now()
	// duration += endTime.Sub(startTime)

	if gridOccupy[x][y] {
		gridMu[x][y].Unlock()
		return
	}

	// now move
	// startTime = time.Now()
	gridMu[indiv.X][indiv.Y].Lock()
	// endTime = time.Now()
	// duration += endTime.Sub(startTime)

	gridOccupy[indiv.X][indiv.Y] = false
	gridMu[indiv.X][indiv.Y].Unlock()

	gridOccupy[x][y] = true
	indiv.X, indiv.Y = x, y

	gridMu[x][y].Unlock()

}
