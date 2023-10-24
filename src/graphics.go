package main

import (
	"sync"
	"image/color"
	"math"
	"math/rand"
	"fmt"
	// "os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

var CURR_GEN = 0
var CURR_STEP = 1
var MAX_STEP = 300

// var duration time.Duration

func killWest() {
	var newOrganisms []*Mite

	for _, mite := range organisms{
		if mite.X >= (COLS / 2){
			newOrganisms = append(newOrganisms, mite)
		}
	}

	organisms = newOrganisms
}

func (g *Game) Update() error {
	if CURR_STEP % MAX_STEP == 0 { // new generation
		// killWest() // selection criteria function



		CURR_GEN++
		CURR_STEP++
		ebiten.SetWindowTitle(fmt.Sprintf("Gen %d", CURR_GEN))
		// now reproduce
		// reposition all mites

		// gridOccupy = map[int]bool{}


		// ===== STEP GENERATION
		// createOccupancyGrid(ROWS, COLS)

		// var children []*Mite
		// for{
		// 	if len(children) >= 1000 { break }
		// 	// pick random organism
		// 	parentMite := organisms[rand.Intn(len(organisms))]
		// 	childGenome := mutateGenome(parentMite.genome)

		// 	// randomizePos(parentMite)
		// 	children = append(children, createMite(childGenome))

		// 	// checkGenomeRep(childGenome, parentMite.genome)
		// }

		// organisms = children

		// // check rep
		// // if len(organisms) != 1000 {
		// // 	fmt.Println("organisms != 1000", len(organisms))
		// // 	os.Exit(1)
		// // }

		//CURR_STEP = 0
		// fmt.Println("Done mutating; starting gen", CURR_GEN)
		// // os.Exit(1)


	} else {

		// if(CURR_STEP == 10) { os.Exit(1) }
		// timing
		start := time.Now()

		children := []*Mite{}
		var updateMu sync.Mutex


		groups := 4

		var wgMite sync.WaitGroup
		for i := 0; i < groups; i++ {
			startSlice := int(math.Ceil( float64(len(organisms)) / float64(groups)) )*i
			endSlice := startSlice + int(math.Ceil( float64(len(organisms)) / float64(groups) ))

			wgMite.Add(1)

			go func(startSlice int, endSlice int, children *[]*Mite, updateMu *sync.Mutex) {
				defer wgMite.Done()
				for mite := startSlice; mite < endSlice; mite++{
					if mite >= len(organisms) { break }
					stepOrganism(organisms[mite])

					if organisms[mite].nutrition > 1.0 { organisms[mite].nutrition = 1.0 }

					// check if organism is fucking dead
					if CURR_STEP - organisms[mite].birth < 200 && organisms[mite].nutrition > 0.0 && !organisms[mite].dead {
						updateMu.Lock()
						*children = append(*children, organisms[mite])
						updateMu.Unlock()
					}

					//decide if organism will divide
						// if rand.Float64() < 0.005 {
					if organisms[mite].nutrition >= 0.7 && rand.Float64() < 0.01{
						updateMu.Lock()

						newMite := cellDivide(organisms[mite])
						if newMite != nil { *children = append(*children, newMite) }

						updateMu.Unlock()
						organisms[mite].nutrition = 0.1
					}
				}
			}(startSlice, endSlice, &children, &updateMu)
		}
		wgMite.Wait() // wait for all organisms to finish before finishing update and drawing

		organisms = children

		CURR_STEP++
		end := time.Now()
		fmt.Println(end.Sub(start), len(organisms))
		// fmt.Println(duration)
		// duration = 0
	}

	gridOccupy = [][]*Mite{}
	createOccupancyGrid(ROWS, COLS)
	for _, mite := range organisms {
		gridOccupy[mite.X][mite.Y] = mite
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.Fill(color.RGBA{16, 42, 67, 255})
	screen.Fill(color.RGBA{34, 34, 34, 255})

	// if(CURR_STEP >= MAX_STEP) { return }
	for _, mite := range organisms {
		vector.DrawFilledCircle(screen, float32((mite.X*10.0)+5.0), float32((mite.Y*10.0)+5.0), 5, mite.color, false)
	}

	// for i := 0; i < 128; i+=1 {
	// 	for j := 0; j < 128; j+=1 {

	// 		vector.DrawFilledCircle(screen, float32(j*10)+5.0, float32(i*10)+5.0, 5, color.RGBA{255,0,0,255}, false)
	// 	}
	// }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ROWS * 10, COLS * 10
}

// func drawGame() {
// 	ebiten.SetWindowSize(640, 640)
// 	ebiten.SetWindowTitle("Hello, World!")
// 	if err := ebiten.RunGame(&Game{}); err != nil {
// 		log.Fatal(err)
// 	}
// }
