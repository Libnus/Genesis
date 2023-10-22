package main

import (
	"math/rand"
	"image/color"
	"fmt"
	// "os"
	// "time"
)

type Mite struct {
	nnet *NeuralNetwork
	genome []string
	id int
	X, Y int
	birth int
	nutrition float64
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

	// red := uint8(rand.Intn(256))
	// green := uint8(rand.Intn(256))
	// blue := uint8(rand.Intn(256))
	// alpha := uint8(rand.Intn(256))

	red, green, blue := getIndivColor(genome)
	alpha := 255

	// fmt.Println("\nMite")
	// for _, gene := range genome {
	// 	fmt.Printf("%s ", gene)
	// }
	// fmt.Println()
	// os.Exit(1)

	return &Mite{
		nnet: nnet,
		genome: genome,
		id: id,
		X: x,
		Y: y,
		birth: 0,
		nutrition: 0.2,
		color: color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)},
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
	// red := uint8(rand.Intn(256))
	// green := uint8(rand.Intn(256))
	// blue := uint8(rand.Intn(256))
	// alpha := uint8(rand.Intn(256))

	red, green, blue := getIndivColor(genome)
	alpha := 255

	newMite := &Mite{
		nnet: nnet,
		genome: genome,
		id: rand.Intn(10000), // TODO fix lol
		X: 0,
		Y: 0,
		nutrition: 0.2,
		birth: CURR_STEP,
		color: color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)},
	}

	randomizePos(newMite)
	return newMite
}

func cellDivide(mite *Mite) *Mite {

	// get a position NOTE: will not spawn if all neighbor grids are occupied

	// a set of all neighbors
	neighbors := [][2]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}
	//neighbours = [(-1, -1), (0, -1), (1, -1), (-1, 0), (1, 0), (-1, 1), (0, 1), (1, 1)]

	var newX int
	var newY int

	for{
		// randomly choose unchecked neighbor
		rn := rand.Intn(len(neighbors))
		var n [2]int = neighbors[ rn ]

		newX = mite.X + n[0]
		newY = mite.Y + n[1]


		if (newX >= 0 && newX < 128) && (newY >= 0 && newY < 128) {
			gridMu[ newX ][ newY ].Lock()
			if !gridOccupy[ newX ][ newY ] {
				break
			}
			gridMu[ newX ][ newY ].Unlock()
		}

		// if occupied then remove the neighbor and try the next
		neighbors[rn] = neighbors[len(neighbors)-1]
		neighbors = neighbors[:len(neighbors)-1]

		// out of neighbors. no baby mite, sorry
		if len(neighbors) == 0 { return nil }
	}

	gridOccupy[ newX ][ newY ] = true

	genome := mutateGenome(mite.genome)
	nnet, _ := processGenome(genome)

	red, green, blue := getIndivColor(genome)
	alpha := 255

	newMite := &Mite{
		nnet: nnet,
		genome: genome,
		id: rand.Intn(10000), // TODO fix lol
		X: newX,
		Y: newY,
		birth: CURR_STEP,
		nutrition: 0.2,
		color: color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)},
	}
	gridMu[ newX ][ newY ].Unlock()

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
	if x == indiv.X && y == indiv.Y { return }

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
