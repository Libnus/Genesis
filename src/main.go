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
	_ "net/http/pprof"

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
var miteIdMu sync.Mutex

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
var saveLastStep *bool

func main() {

	go func() {
		// Start the HTTP server for pprof at :6060
		http.ListenAndServe(":6060", nil)
	}()
    // Notify the channel on SIGINT and SIGTERM signals
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var draw = flag.Bool("draw", true, "Enable or disable drawing")
	var seed = flag.Int64("seed", 1125, "help message for flag n")
	var loadReplay = flag.Bool("load_replay", false, "load a replay from organisms.rpy")
	var loadFromReplay = flag.Bool("load_from_replay", false, "start a sim from the last step of a replay")
	saveLastStep = flag.Bool("save_last", false, "saves only the last step of sim. Still saves event data")

	rpy_flag = flag.Bool("replay", false, "replay events")
	var maxStep = flag.Int("max_step", 0, "maximum number of steps to simulate")
	var brain = flag.Bool("brain", false, "fetches a brain from a replay")
	flag.Parse()


	if *loadFromReplay{
		// *loadReplay = true
		*rpy_flag = true
	}

	if *saveLastStep{
		*rpy_flag = true
		*loadReplay = false
	}

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

	game := &Game{drawEnabled: *draw, saveLast: *saveLastStep, rpy: *rpy_flag, isReplay: *loadReplay, maxStep: *maxStep}


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

	// mite := createRandomMite(100, 0)

	// return


	if !*loadReplay && !*loadFromReplay{
		for i := 0; i < 1000; i++{
			mite := createRandomMite(100, i)
			// generateMiteName(mite, nil)
			organisms = append(organisms, mite)

			treeNode := treeInsert("root", getName(mite), colorToHex(mite.Color))
			addSpecies(mite, treeNode)

			if *rpy_flag {
				addEvent(Birth, nil, mite)
			}

		}
	} else if *loadFromReplay{
		fmt.Println("laoding")
		loadReplayFromFile("organisms.rpy")
		//gridOccupy = replay.replayGrid[0]
		step := 0
		if len(replay.ReplayGrid) > 1{ // construct event data
			for step, _ = range replay.ReplayGrid{
				fmt.Println("event logging wokring", step)
				for _, event := range replay.Events[step] {
					organism := replay.Organisms[event.Source]

					if event.Type == Birth{
						// CURR_POP++
						//
						target := replay.Organisms[event.Target]

						var treeNode *TreeNode
						if organism != nil{
							treeNode = treeInsert(getName(organism), getName(target), colorToHex(target.Color))
						} else{
							treeNode = treeInsert("root", getName(target), colorToHex(target.Color))
						}
						addSpecies(target, treeNode)
						MITE_ID++
					} else if event.Type == Kill {
						addKill()
					} else if event.Type == Death {
						// CURR_POP--
						removeSpecies(organism)
					}
				}
			}
		}
		CURR_STEP =  step

		for i, _ := range replay.ReplayGrid[len(replay.ReplayGrid)-1]{
			for j, pos := range replay.ReplayGrid[len(replay.ReplayGrid)-1][i]{
				if pos == -1 {
					gridOccupy[i][j] = nil
					continue
				}
				mite := replay.Organisms[pos]
				fmt.Println(mite.Dead, mite.Birth, CURR_STEP, getName(mite))
				return
				// mite := createMite(replay.Organisms[pos].Genome)
				// mite.Id = pos

				// // NOTE quick fix from bug identified in indiv.CellDivide
				// if _, ok := speciesData[getName(mite)]; !ok {
				// 	continue
				// }

				organisms = append(organisms, mite)
				gridOccupy[i][j] = mite


				if len(replay.ReplayGrid) == 1{
					addSpecies(replay.Organisms[pos], treeData)
				}
			}
		}
		fmt.Println("starting from replay step:", CURR_STEP)
	}

	// ============== TESTING
	// mite1 := organisms[0]
	// mite2 := mite1
	// var mite2 *Mite

	// genome1 := organisms[0].Genome
	// genome2 := genome1

	// similarity := 0.0
	// similarity = genomeSimilarity(mite1, mite2) // comparing same mite
	// fmt.Printf("%d   Similarity: %.2f%%\n", -1, similarity*100)
	// // return

	// for i := range 100 {
	// 	mite2 = createMite(mutateGenome(mite2.Genome))
	// 	similarity = genomeSimilarity(mite1, mite2)
	// 	fmt.Printf("%d   Similarity: %.2f%%\n", i, similarity*100)
	// }


	// fmt.Printf("Genome similarity: %.2f%%\n", similarity*100)
	// return

	// graphing / logging
	// http.HandleFunc("/", httpserver)
	// http.ListenAndServe(":8081", nil)
	// go func() {
    //     http.HandleFunc("/", httpserver)
    //     http.ListenAndServe(":8081", nil)
    // }()

	go func(){
		http.HandleFunc("/brain", brainHandler)
		http.HandleFunc("/data", dataHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	go signalHandler(*rpy_flag, *saveLastStep)


	if *draw {
        ebiten.SetMaxTPS(60) // Sync TPS with FPS if drawing
    } else {
        ebiten.SetMaxTPS(10000) // Run at maximum TPS if drawing is disabled
    }

	ebiten.SetWindowSize(ROWS*5, COLS*5)
	ebiten.SetWindowTitle("Gen 0")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
