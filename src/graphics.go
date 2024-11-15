package main

import (
	"sync"
	"image/color"
	"math"
	"math/rand"
	// "strconv"
	"fmt"
	"syscall"
	// "os"
	"time"
	// "net/http"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	// "github.com/go-echarts/go-echarts/v2/charts"
	// "github.com/go-echarts/go-echarts/v2/components"
    // "github.com/go-echarts/go-echarts/v2/opts"
    // "github.com/go-echarts/go-echarts/v2/types"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{
	drawEnabled bool
	rpy bool
	isReplay bool
	maxStep int
	mouseClicked   bool
    mouseX, mouseY int
}

var GAME_PAUSED = false 
var previousSpacePressed = false

var CURR_GEN = 0
var CURR_SICK = 0
var PREV_INTERP = 0
var CURR_KILL = 0
var CURR_STEP = 0
var MAX_STEP = 300

//var genData []opts.LineData
// var genData []int
// var killingData []float64
// var killingLineData []opts.LineData
// var genDataX []string
// var events = []opts.MarkLineNameCoordItem{
//     {
//         Name:        fmt.Sprintf("Gen %d", CURR_GEN),
//         Coordinate0: []interface{}{CURR_STEP, 0},              // x = CURR_STEP, y = 0
//         Coordinate1: []interface{}{CURR_STEP, 1000}, // x = CURR_STEP, y = maxY
//     },
// }

// var evolutionTree = []opts.TreeData{
// 	{
// 		Name: "root",
// 		Children: []*opts.TreeData{},
// 	},
// }


// func insertEvolutionTree(parent string, child string, node *opts.TreeData) bool{
// 	// if len(node.Children) == 0{
// 	// 	return false
// 	// }

// 	if node.Name == parent{
// 		newSpecies := opts.TreeData{
// 			Name: child,
// 			Children: []*opts.TreeData{},
// 			ItemStyle: &opts.ItemStyle{
// 				Color: child,
// 			},
// 		}

// 		node.Children = append(node.Children, &newSpecies)
// 	}

// 	for _, next := range node.Children {
// 		if insertEvolutionTree(parent, child, next){
// 			return true
// 		}
// 	}

// 	return false
// }

// func generateEventSeries() []opts.LineData {
//     eventData := make([]opts.LineData, 7)
//     for i := 0; i < 7; i++ {
//         eventData[i] = opts.LineData{Value: 200} // This example adds a horizontal line at value 200
//     }
//     return eventData
// }

// func getMaxVal(data []opts.LineData) int {
//     max := data[0].Value.(int)
//     for _, item := range data {
//         if item.Value.(int) > max {
//             max = item.Value.(int)
//         }
//     }
//     return max
// }

// func getSpeciesData() ([]string, []opts.BarData){
// 	xAxis := make([]string, 0, len(species))
// 	barData := make([]opts.BarData, 0, len(species))

// 	speciesMu.Lock()
// 	for color, count := range species {
// 		if count >= 10{
// 			colorHex := colorToHex(color)
// 			xAxis = append(xAxis, colorHex)
// 			barData = append(barData, opts.BarData{Value: count, ItemStyle: &opts.ItemStyle{Color: colorHex},})
// 		}
// 	}
// 	speciesMu.Unlock()

// 	return xAxis, barData
// }

// func createSpeciesChart() *charts.Bar {
// 	bar := charts.NewBar()

// 	xAxis, data := getSpeciesData()

// 	bar.SetGlobalOptions(
//         charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
//         charts.WithTitleOpts(opts.Title{Title: "Species Count Chart"}),
//     )
// 	bar.SetXAxis(xAxis).AddSeries("Count", data)

// 	return bar
// }

// func createEvolutionTree() *charts.Tree{
// 	tree := charts.NewTree()
// 	tree.SetGlobalOptions(
// 		charts.WithInitializationOpts(opts.Initialization{ChartID: "evolutionTreeChart"}),
// 		charts.WithTitleOpts(opts.Title{Title: "Evolution Tree"}),
// 		//charts.WithTooltipOpts(opts.Tooltip{Show: false}),
// 	)

// 	tree.AddSeries("tree", evolutionTree).
// 		SetSeriesOptions(
// 			charts.WithTreeOpts(
// 				opts.TreeChart{
// 					Layout:           "radial",
// 					Orient:           "TB",
// 					InitialTreeDepth: -1,
// 					Leaves: &opts.TreeLeaves{
// 						Label: &opts.Label{Show: opts.Bool(false), Position: "right", Color: "Black"},
// 					},
// 					Roam: opts.Bool(true),
// 				},
// 			),
// 			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "top", Color: "Black"}),
// 		)
// 	return tree
// }


// func treeHandler(w http.ResponseWriter, r *http.Request) *charts.Tree {
// 	// Build the tree chart

// 	tree := charts.NewTree()
// 	tree.SetGlobalOptions(
// 		charts.WithInitializationOpts(opts.Initialization{}),
// 		charts.WithTitleOpts(opts.Title{Title: "Evolution Tree"}),
// 		//charts.WithTooltipOpts(opts.Tooltip{Show: false}),
// 	)

// 	tree.AddSeries("tree", evolutionTree).
// 		SetSeriesOptions(
// 			charts.WithTreeOpts(
// 				opts.TreeChart{
// 					Layout:           "radial",
// 					Orient:           "TB",
// 					InitialTreeDepth: -1,
// 					Leaves: &opts.TreeLeaves{
// 						Label: &opts.Label{Show: opts.Bool(false), Position: "right", Color: "Black"},
// 					},
// 					Roam: opts.Bool(true),
// 				},
// 			),
// 			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "top", Color: "Black"}),
// 		)

// 	// Serve the chart with embedded JavaScript for search
// 	// w.Header().Set("Content-Type", "text/html")

// 	// Embed the HTML and JavaScript directly in the response
// 	return tree
// }

// func createPopChart() *charts.Line {
// 	line := charts.NewLine()
// 	// set some global options like Title/Legend/ToolTip or anything else
// 	line.SetGlobalOptions(
// 		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
// 		charts.WithTitleOpts(opts.Title{
// 			Title:    "Line example in Westeros theme",
// 			Subtitle: "Line chart rendered by the http server this time",
// 		}))



// 	// thresholdValues := []int{150, 230, 224, 218, 135, 147, 260}
// 	// Put data into instance
// 	line.SetXAxis(genDataX).
// 		AddSeries("Killings", killingLineData).
// 		AddSeries("Number of organisms", genData).
// 		// AddSeries("Event Threshold", convertToLineData(thresholdValues)).

// 		// SetGlobalOptions(
// 		// 	charts.WithYAxisOpts(opts.YAxis{
// 		// 		Type:      "category",
// 		// 		Data:      thresholdValues,			// charts.WithMarkLineNameTypeItemOpts(opts.MarkLineNameTypeItem{
// 		// 	}),
// 		// ).
// 		SetSeriesOptions(
// 			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
// 			charts.WithMarkLineNameCoordItemOpts(events...),
// 				// opts.MarkLineNameCoordItem{
// 				// 	Name:  "Event Line",
// 				// 	Coordinate0: []interface{}{0, 0},
// 				// 	Coordinate1: []interface{}{0,getMaxVal(genData)},// Starting point of the vertical line
//             	// },
// 				// opts.MarkLineNameCoordItem{
// 				// 	Name:  "Event Line 2",
// 				// 	Coordinate0: []interface{}{100, 0},
// 				// 	Coordinate1: []interface{}{100,getMaxVal(genData)},// Starting point of the vertical line
// 				// },
// 			// 	Name: "Line",
// 			// 	Type: "value",
// 			// }),
// 			// charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
// 			// 	Label: &opts.Label{
// 			// 		Show:      opts.Bool(true),
// 			// 		Formatter: "{a}: {b}",
// 			// 	},
// 			// }),
// 		)
// 	return line
// }

// func httpserver(w http.ResponseWriter, h *http.Request) {
// 	page := components.NewPage()
// 	page.SetLayout(components.PageFlexLayout) // arrange charts in a flexible layout
//     page.AddCharts(
//         createPopChart(),
//         createSpeciesChart(),
// 		createEvolutionTree(),
//     )

// 	page.Render(w)
// }


func takeLastN(data []int, n int) []int {
    if len(data) <= n {
        return data // If data has less than or equal to n elements, return it as is
    }
    return data[len(data)-n:] // Slice to get the last n elements
}

// Linear interpolation function
// func linearInterpolate() {
//     // interpolated := make([]float64, len(killingData))
//     // copy(interpolated, killingData) // Copy original data to keep its structure

//     for i := PREV_INTERP; i < len(killingData); i++ {
//         if killingData[i] == -1 { // Detect missing data point (-1)
//             // Find the previous and next valid data points
//             prev, next := i-1, i+1
//             for prev >= 0 && killingData[prev] == -1 {
//                 prev--
//             }
//             for next < len(killingData) && killingData[next] == -1 {
//                 next++
//             }
//             // Interpolate if we found valid neighbors
//             if prev >= 0 && next < len(killingData) {
//                 // killingLineData = append(killingLineData, opts.LineData{Value: int(killingData[prev] + (killingData[next]-killingData[prev])*(float64(i-prev)/float64(next-prev)))})
//             }
//         }
//     }
// }

// func appendKillingData(){

// 	// append current kills
// 	killingData = append(killingData, float64(CURR_KILL))
// 	CURR_KILL = 0
// 	linearInterpolate()
// 	PREV_INTERP = CURR_STEP

// 	// // take last 1000 for moving average and sum
// 	// last1000 := takeLastN(killingData, 100)
// 	// sum := 0
// 	// for _, data := range last1000 {
// 	// 	sum += data
// 	// }

// 	// fmt.Println("yesssssssssssssssssssssssssssssssssssssss", sum)

// 	// // calc new average and append it
// 	// smaKillingData = append(smaKillingData, opts.LineData{Value: sum/100})
// }

func addKill(){
	CURR_KILL++
}

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


var CURR_POP = 0
var AVG_BRAIN_SIZE = 0


func (g *Game) Update() error {
	spacePressed := ebiten.IsKeyPressed(ebiten.KeySpace)
	//
	if spacePressed && !previousSpacePressed {
        if !GAME_PAUSED {
            GAME_PAUSED = true
        } else {
            GAME_PAUSED = false
        }
    }
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
        // Get the current mouse position
        x, y := ebiten.CursorPosition()

        if !g.mouseClicked { // Only register a click if it wasn't already clicked
            // fmt.Printf("Mouse clicked at: (%d, %d)\n", x/10, y/10)
            g.mouseX, g.mouseY = x/10, y/10
            g.mouseClicked = true // Register the click
			openBrainBrowser(g.mouseX, g.mouseY)
        }

    } else {
        g.mouseClicked = false // Reset the click state when mouse is released
    }

	previousSpacePressed = spacePressed
	
    // If paused, do nothing in update
    if GAME_PAUSED {
        return nil
    }
	gridLockMu.Lock()

	if g.rpy{
		saveGrid();
	}

	if CURR_STEP % MAX_STEP == 0 { // new generation
		// killWest() // selection criteria function
		//
		appendKillingData()

		CURR_GEN++
		// CURR_STEP++
		if g.drawEnabled {ebiten.SetWindowTitle(fmt.Sprintf("Gen %d", CURR_GEN))}

		// newEvent := opts.MarkLineNameCoordItem{
		// 	Name:        fmt.Sprintf("Gen %d", CURR_GEN),
		// 	Coordinate0: []interface{}{CURR_STEP, 0},     // x = 50, y = 0
		// 	Coordinate1: []interface{}{CURR_STEP, getMaxVal(genData)},  // x = 50, y = maxY
		// }

		// events = append(events, newEvent)
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


	}
	if g.isReplay {
		// CURR_STEP++
		// fmt.Println(ebiten.CurrentTPS())
		// start := time.Now()

		//fmt.Println(len(replay.Events[CURR_STEP]))
		for _, event := range replay.Events[CURR_STEP] {
			organism := replay.Organisms[event.Source]


			if event.Type == Birth{
				CURR_POP++
				target := replay.Organisms[event.Target]

				if organism != nil{
					treeInsert(getName(organism), getName(target), colorToHex(target.Color), treeData)
				} else{
					treeInsert("root", getName(target), colorToHex(target.Color), treeData)
				}
				addSpecies(target)
			} else if event.Type == Kill {
				addKill()
			} else if event.Type == Death {
				CURR_POP--
				removeSpecies(organism)
			}

		}
		appendGenData(CURR_POP)

		CURR_STEP++
		// end := time.Now()
		// fmt.Println(ebiten.CurrentTPS(), end.Sub(start))

		// fmt.Println(ebiten.CurrentTPS(), len(replay.Events[CURR_STEP]), CURR_POP)
		// fmt.Println(CURR_POP)
	}else {

		if g.maxStep > 0 && CURR_STEP == g.maxStep {
			sigChan <- syscall.SIGINT
		}
		// if(CURR_STEP == 10) { os.Exit(1) }
		// timing
		// start := time.Now()

		children := []*Mite{}
		var updateMu sync.Mutex

		CURR_SICK = 0
		AVG_BRAIN_SIZE = 0


		// killingData = append(killingData, -1)

		if rand.Float64() < 0.05 { // 0.01
			// fmt.Println("the plague.")
			index := rand.Intn(len(organisms))
			mite := organisms[index]
			mite.Sick = rand.Intn(50) + 1
			// CURR_SICK++
		}

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

					AVG_BRAIN_SIZE += len(organisms[mite].Genome)
					stepOrganism(organisms[mite])

					if organisms[mite].Nutrition > 1.0 { organisms[mite].Nutrition = 1.0 }

					updateMu.Lock()

					if organisms[mite].Sick > 0{
						CURR_SICK++
					}

					// less than absolute max age && has nutrition && not murdered
					if organisms[mite].Nutrition > 0.0 && !organisms[mite].Dead {
						*children = append(*children, organisms[mite])
						// pass()
					} else{
						if g.rpy { addEvent(Death, organisms[mite], nil) } // if the organism wasn't murdered then add death event
						removeSpecies(organisms[mite]);
					}
					updateMu.Unlock()

					//decide if organism will divide
						// if rand.Float64() < 0.005 {
					if organisms[mite].Nutrition >= 0.7 && rand.Float64() < 0.01{
						updateMu.Lock()

						newMite := cellDivide(organisms[mite])

						if newMite != nil {
							*children = append(*children, newMite)
							addSpecies(newMite)

							births++

							if g.rpy{
								addEvent(Birth, organisms[mite], newMite)
							}

							if(organisms[mite].Color != newMite.Color){
								treeInsert(getName(organisms[mite]), getName(newMite), colorToHex(newMite.Color), treeData)
							}

						}

						updateMu.Unlock()
						organisms[mite].Nutrition = 0.1
					}
				}
			}(startSlice, endSlice, &children, &updateMu)
		}
		wgMite.Wait() // wait for all organisms to finish before finishing update and drawing


		AVG_BRAIN_SIZE /= len(organisms)
		organisms = children

		// end := time.Now()
		// fmt.Println(end.Sub(start), len(organisms))
		appendGenData(len(organisms))
		appendSickData()

		appendBrainData()

		CURR_STEP++

		// fmt.Println(duration)
		// duration = 0
		//
		gridOccupy = [][]*Mite{}
		createOccupancyGrid(ROWS, COLS)
		for _, mite := range organisms {
			gridOccupy[mite.X][mite.Y] = mite
		}
	}


	gridLockMu.Unlock()
	return nil
}

func (g *Game) UpdateReplay() error{
	CURR_STEP++
	return nil
}

// func (g *Game) Update() error {
// 	if g.isReplay  {
// 		return g.UpdateReplay()
// 	}else{
// 		return g.UpdateGame()
// 	}

// }

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.drawEnabled { return }

	//screen.Fill(color.RGBA{16, 42, 67, 255})
	screen.Fill(color.RGBA{34, 34, 34, 255})

	// if(CURR_STEP >= MAX_STEP) { return }

	if !g.isReplay{
		for _, mite := range organisms {
			vector.DrawFilledCircle(screen, float32((mite.X*10.0)+5.0), float32((mite.Y*10.0)+5.0), 5, mite.Color, false)
		}
	} else{
		start := time.Now()
		grid := replay.ReplayGrid[CURR_STEP]
		for i := range grid{
			for j := range grid[i]{
				if grid[i][j] != -1{
					mite := replay.Organisms[grid[i][j]]
					vector.DrawFilledCircle(screen, float32((i*10.0)+5.0), float32((j*10.0)+5.0), 5, mite.Color, false)
				}
			}
		}
		end := time.Now()
		fmt.Println(end.Sub(start))
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
