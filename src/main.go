package main

import (
    // "fmt"
	// "os"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)


// func hexTo32BitBinary(hexString string) (string, error) {
//     // Parse the hexadecimal string to an integer
//     hexInt, err := strconv.ParseInt(hexString, 16, 64)
//     if err != nil {
//         return "", err
//     }

//     // Format the integer as a 32-bit binary string and left-pad with zeros
//     binaryString := fmt.Sprintf("%032b", hexInt)

//     return binaryString, nil
// }

/*

IDEA: specify connection frequency for an organism which changes how frequently it attempts
to make new connections and delete unused ones (for allowing the neural network to learn)

    also specify a frequency to change weights (how oftento change the weight) and specify a value
    to change the weight by that amount.

    These are all values in the genome of the organism which can mutate

*/

/*
// 71526A73
// 0 1110001 0 1010010 0110101001110011
//    113       82       27251
   01110001010100100110101001110011

   // 64DBD2E9
   // 0 1100100 1 1011011 1101001011101001
   // 01100100110110111101001011101001

   648FD2E9
   0 1100100 1 0001111 1101001011101001
   01100100100011111101001011101001

   8FA0D2E9
   1 0001111 1 0100000 1101001011101001
   10001111101000001101001011101001

   8F8FD2F9
   1 0001111 1 0001111 1101001011111001
   10001111100011111101001011111001

	A00076F9
   1 0100000 0 0000000 0111011011111001
   10100000000000000111011011111001


   Source type: input
   Source ID: Px
   Dest type: internal neuron
   Dest ID: 0
   Weight: -1.39807128906

*/

var ROWS = 128
var COLS = 128

var organisms []*Mite
var gridOccupy = [][]bool{} // every grid position
// TODO ^^ replace with *Mite instead of bool so we can perform more complex calculations about the organisms occupying particular grid locations

// used to sync occupancy grid
var gridMu = [][]*sync.Mutex{}

func clearGrid(){
	gridOccupy = [][]bool{}
	gridMu = [][]*sync.Mutex{}
}

func createOccupancyGrid(rows int, cols int){
	clearGrid()

	for i := 0; i < rows; i++ {
		gridOccupy = append(gridOccupy, make([]bool, cols))
		gridMu = append(gridMu, make([]*sync.Mutex, cols))
		for j := 0; j < cols; j++ {
			// each position is not occupied initially
			gridOccupy[i][j] = false
			gridMu[i][j] = &sync.Mutex{}
		}
	}
}

func main() {

	createOccupancyGrid(ROWS, COLS)

    // fmt.Println("Processing genome test")
    // // genome := []string{"648FD2E9", "8FA0D2E9", "8F8FD2F9", "A00076F9"}
    // genome := createGenome(1000, 1125)
    // nnet, _ := processGenome(genome)

	// if extractEdgeList(nnet) != nil{
	// 	os.Exit(1)
	// }

    // printNeural(nnet)

	// _ = calcNeuralPotential(nnet)
	// fmt.Println("______")

	for i := 0; i < 1000; i++{
		organisms = append(organisms, createRandomMite(1000, i))
	}

	ebiten.SetMaxTPS(60)

	ebiten.SetWindowSize(ROWS*5, COLS*5)
	ebiten.SetWindowTitle("Gen 0")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
