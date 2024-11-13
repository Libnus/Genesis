package main

import (
    "fmt"
	"os"
    "os/signal"
    "syscall"
	"flag"
	"log"
	// "time"
	"sync"
	"math/rand"
	"image/color"
	// "strconv"
	"net/http"

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
var MITE_ID = 0

var sigChan = make(chan os.Signal, 1)

var species = make(map[color.Color]int)
var organisms []*Mite
var gridOccupy = [][]*Mite{} // every grid position

// keep replays

// TODO ^^ replace with *Mite instead of bool so we can perform more complex calculations about the organisms occupying particular grid locations

// used to sync occupancy grid
var gridMu = [][]*sync.Mutex{}
var gridLockMu sync.Mutex
var speciesMu sync.Mutex

// func addSpecies(mite *Mite){
// 	speciesMu.Lock()
// 	if _, exists := species[mite.Color]; exists{
// 		species[mite.Color]++
// 	} else{
// 		species[mite.Color] = 1
// 	}
// 	speciesMu.Unlock()
// }

func clearGrid(){
	gridOccupy = [][]*Mite{}
	gridMu = [][]*sync.Mutex{}
}

func createOccupancyGrid(rows int, cols int){
	clearGrid()

	for i := 0; i < rows; i++ {
		gridOccupy = append(gridOccupy, make([]*Mite, cols))
		gridMu = append(gridMu, make([]*sync.Mutex, cols))
		for j := 0; j < cols; j++ {
			// each position is not occupied initially
			gridOccupy[i][j] = nil
			gridMu[i][j] = &sync.Mutex{}
		}
	}
}

var rpy_flag *bool

func main() {

    // Notify the channel on SIGINT and SIGTERM signals
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var draw = flag.Bool("draw", true, "Enable or disable drawing")
	var seed = flag.Int64("seed", 1125, "help message for flag n")
	var loadReplay = flag.Bool("load_replay", false, "load a replay from organisms.rpy")

	rpy_flag = flag.Bool("replay", false, "replay events")
	var maxStep = flag.Int("max_step", 0, "maximum number of steps to simulate")
	var brain = flag.Bool("brain", false, "fetches a brain from a replay")
	flag.Parse()

	fmt.Println(*brain)

	if *brain{
		fmt.Println("getting a brain")
		loadReplayFromFile("organisms.rpy")
		getBrains()
		return // just fetch brain and return
	}

	if *loadReplay{
		fmt.Print("loading replay...")
		*draw = true
		initReplay()
		*rpy_flag = false
		loadReplayFromFile("organisms.rpy")

	}

	// error checking
	 // if len(os.Arg) > 1 {
     //    seed, err := strconv.ParseInt(os.Args[1], 10, 64)
     //    if err != nil {
     //        fmt.Println("Invalid seed:", os.Args[1])
     //        return
     //    }
    rand.Seed(*seed)
    // } else {
    //     rand.Seed(time.Now().UnixNano())
    // }

	game := &Game{drawEnabled: *draw, rpy: *rpy_flag, isReplay: *loadReplay, maxStep: *maxStep}


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
	if *rpy_flag {
		initReplay()
	}


	if !*loadReplay{
		for i := 0; i < 1000; i++{
			mite := createRandomMite(100, i)
			// generateMiteName(mite, nil)
			organisms = append(organisms, mite)
			addSpecies(mite)

			if *rpy_flag {
				//replay.Events[0] = append(replay.Events[0], createEvent(Birth, nil, mite))
				addEvent(Birth, nil, mite)
			}

			treeInsert("root", getName(mite), colorToHex(mite.Color), treeData)
		}
	}


	// genome1 := organisms[0].Genome
	// genome2 := genome1

	// similarity := compareGenomes(genome1, genome2)

	// for i := range 100 {
	// 	genome2 = mutateGenome(genome2)
	// 	similarity := genomeSimilarity(genome1, genome2)
	// 	fmt.Printf("%d   Similarity: %.2f%%\n", i, similarity)
	// }


	// // fmt.Printf("Genome similarity: %.2f%%\n", similarity)
	// fmt.Println(len(genome1), len(genome2))
	// return

	// graphing / logging
	// http.HandleFunc("/", httpserver)
	// http.ListenAndServe(":8081", nil)
	// go func() {
    //     http.HandleFunc("/", httpserver)
    //     http.ListenAndServe(":8081", nil)
    // }()

	go func(){
		http.HandleFunc("/data", dataHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go signalHandler(*rpy_flag)


	if *draw {
        ebiten.SetMaxTPS(60) // Sync TPS with FPS if drawing
    }

	// else {
    //     ebiten.SetMaxTPS(0) // Run at maximum TPS if drawing is disabled
    // }

	ebiten.SetWindowSize(ROWS*5, COLS*5)
	ebiten.SetWindowTitle("Gen 0")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
