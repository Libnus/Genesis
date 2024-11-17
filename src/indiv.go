package main

import (
	"math/rand"
	"image/color"
	// "fmt"
	// "os"
	"log"
	// "time"
)

type Mite struct {
	Genus string
	Species string
	Nnet *NeuralNetwork
	Genome []string
	didMove bool // did the mite move this step
	Id int
	X, prevX, Y, prevY int
	Birth int
	Nutrition float64
	Dead bool // flag to check if the organism is dead
	Color color.Color
	Sick int // is the creature sick. 0 is not
			 // if > 0 then Sick acts as a timer for how many more steps they are sick
}

func getName(mite *Mite) string{
	return mite.Genus + " " + mite.Species
}

func generateMiteName(mite1, mite2 *Mite){

	if mite2 == nil {
		mite1.Genus = generateGenusName()
		mite1.Species = generateSpeciesName()
		return
	}

	speciesThreshold := 0.95
	genusThreshold := 0.85

	genus := mite1.Genus
	species := mite1.Species

	similarity := genomeSimilarity(mite1, mite2)
	// fmt.Printf("genome similarity %.2f%%\n", similarity*100.0)

	if similarity < genusThreshold{ // make new genus and name
		mite2.Genus = generateGenusName()
		mite2.Species = generateSpeciesName()
	} else if similarity < speciesThreshold{ // just new speices name
		mite2.Genus = genus
		mite2.Species = generateSpeciesName()
	} else{
		mite2.Genus = genus
		mite2.Species = species
	}
}

// NOTE:: iterative
func createRandomMite(numGenes int, id int) *Mite {
	genome := createGenome(numGenes, int64(rand.Int()))
	nnet, _ := processGenome(&genome)

	// fmt.Println("NEW LENGTH", len(genome))

	var x, y = 0, 0
	for{
		if gridOccupy[x][y] != nil {
			x = rand.Intn(ROWS)
			y = rand.Intn(COLS)
		} else{ break }
	}


	// red := uint8(rand.Intn(256))
	// green := uint8(rand.Intn(256))
	// blue := uint8(rand.Intn(256))
	// alpha := uint8(rand.Intn(256))


	// fmt.Println("\nMite")
	// for _, gene := range genome {
	// 	fmt.Printf("%s ", gene)
	// }
	// fmt.Println()
	// os.Exit(1)

	newMite := &Mite{
		Nnet: nnet,
		Genome: genome,
		Id: MITE_ID,
		X: x,
		Y: y,
		Birth: 0,
		Nutrition: 1.0,
		Dead: false,
	}


	speciesMu.Lock()
	MITE_ID++
	speciesMu.Unlock()
	generateMiteName(newMite, nil)

	red, green, blue := getIndivColor(getName(newMite))

	alpha := 255
	newMite.Color = color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)}

	gridOccupy[x][y] = newMite
	return newMite
}

func randomizePos(mite *Mite) {
	var x, y = rand.Intn(ROWS), rand.Intn(COLS)
	for{
		if gridOccupy[x][y] != nil{
			x = rand.Intn(ROWS)
			y = rand.Intn(COLS)
		} else{ break }
	}

	gridOccupy[x][y] = mite
	mite.X, mite.Y = x, y
}

func createMite(genome []string) *Mite {
	nnet, _ := processGenome(&genome)

	// TODO move to get random color function
	// red := uint8(rand.Intn(256))
	// green := uint8(rand.Intn(256))
	// blue := uint8(rand.Intn(256))
	// alpha := uint8(rand.Intn(256))


	newMite := &Mite{
		Nnet: nnet,
		Genome: genome,
		Id: MITE_ID, // TODO fix lol
		X: 0,
		Y: 0,
		Nutrition: 1.0,
		Birth: CURR_STEP,
		Dead: false,
	}

	MITE_ID++

	generateMiteName(newMite, nil)
	red, green, blue := getIndivColor(getName(newMite))
	alpha := 255

	newMite.Color = color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)}

	randomizePos(newMite)
	return newMite
}
func createMiteEmpty(genome []string, x, y int) *Mite {
	nnet, _ := processGenome(&genome)

	// TODO move to get random color function
	// red := uint8(rand.Intn(256))
	// green := uint8(rand.Intn(256))
	// blue := uint8(rand.Intn(256))
	// alpha := uint8(rand.Intn(256))


	newMite := &Mite{
		Nnet: nnet,
		Genome: genome,
		Id: MITE_ID, // TODO fix lol
		X: x,
		Y: y,
		Nutrition: 1.0,
		Birth: CURR_STEP,
		Dead: false,
	}

	miteIdMu.Lock()
	MITE_ID++
	miteIdMu.Unlock()

	// generateMiteName(newMite, nil)
	red, green, blue := getIndivColor(getName(newMite))
	alpha := 255

	newMite.Color = color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)}

	randomizePos(newMite)
	return newMite
}

func cellDivide(mite *Mite) *Mite {
	if mite.Dead { return nil } // don't divide if dead
	// mites no longer birth if dead

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
			if gridOccupy[ newX ][ newY ] == nil {
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


	genome := mutateGenome(mite.Genome)
	nnet, _ := processGenome(&genome)


	// newMite := createMiteEmpty(genome, newX, newY)

	newMite := &Mite{
		Nnet: nnet,
		Genome: genome,
		Id: MITE_ID, // TODO fix lol
		X: newX,
		Y: newY,
		Birth: CURR_STEP,
		Nutrition: 1.0,
		Dead: false,
		// Color: color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)},
	}

	miteIdMu.Lock()
	MITE_ID++
	miteIdMu.Unlock()


	speciesMu.Lock()
	firstOrganism := speciesData[getName(mite)].original // get first organism in species
		if firstOrganism == nil {
		return nil } // organism died while dividing
	speciesMu.Unlock()
	if speciesData[getName(mite)].num == 0{
		log.Fatal("fatal error in cell divide, species extinct but dividng")
		// fmt.Println(" error in cell divide, species extinct but dividng", firstOrganism)
		// data := speciesData[getName(mite)]
		// data.num++
		// speciesData[getName(mite)] = data
		// speciesMu.Unlock()
		//
	// 	mite.Dead = true // just kill it ;) problem solved ;) TODO fix whatever causes this bug
	// 	// 														// creatures are null when loading from replay
	// 	gridMu[ newX ][ newY ].Unlock()
	// 	return nil
	}
	generateMiteName(firstOrganism, newMite)
	red, green, blue := getIndivColor(getName(newMite))
	alpha := 255
	newMite.Color = color.RGBA{uint8(red),uint8(green),uint8(blue),uint8(alpha)}

	mite.Nutrition -= 0.5

	gridOccupy[ newX ][ newY ] = newMite

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

	if gridOccupy[x][y] != nil{
		gridMu[x][y].Unlock()
		return
	}

	// now move
	// startTime = time.Now()
	gridMu[indiv.X][indiv.Y].Lock()

	// endTime = time.Now()
	// duration += endTime.Sub(startTime)

	gridOccupy[indiv.X][indiv.Y] = nil
	gridMu[indiv.X][indiv.Y].Unlock()

	gridOccupy[x][y] = indiv
	indiv.prevX = indiv.X
	indiv.prevY = indiv.Y
	indiv.X, indiv.Y = x, y

	gridMu[x][y].Unlock()
}
