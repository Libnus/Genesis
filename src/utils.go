package main

import (
	// "reflect"
	"fmt"
	"math/rand"
	"encoding/csv"
	"image/color"
	"os"
	"strconv"

	"crypto/sha256"
	"math"
)


// takes in the number of connections for a genome and returns a randomly created genome
func createGenome(numConnections int, seed int64) []string {
	genome := []string{}
	rand.Seed(seed)

	for i := 0; i < numConnections; i++ {
		randomGene := rand.Uint32()
		genome = append(genome, fmt.Sprintf("%X", randomGene))
	}

	// fmt.Println(genome)
	// fmt.Println()

	return genome
}

func generateRandomGene() string {
    gene := rand.Uint32() // Random 32-bit unsigned integer
    return fmt.Sprintf("%X", gene)
}

// takes in a genome, creates a copy, and mutates it
func mutateGenome(genome []string) []string {
	mutationRate := 0.0005 // mutationRate / 15000 for prob

	expectedSwaps := mutationRate * float64(len(genome))
	insertDeleteProbability := float64(expectedSwaps) / float64(len(genome))
	// fmt.Println(insertDeleteProbability, expectedSwaps)

	var mutatedGenome []string
	for _, gene := range genome {
		if rand.Float64() < insertDeleteProbability { // delete gene chance
            continue // skip adding this gene to mutatedGenome
        }


		gBinCopy, _ := strconv.ParseInt(gene, 16, 64)
		gCopy := fmt.Sprintf("%032b", gBinCopy)

		newG := []rune(gCopy)
		for i, base := range gCopy{
			if rand.Float64() < mutationRate{
				if base == '0' {
					newG[i] = '1'
				} else { newG[i] = '0' }
			}
		}
		gCopyInt, _ := strconv.ParseUint(string(newG), 2, 32)
		mutatedGenome = append(mutatedGenome, fmt.Sprintf("%X", gCopyInt))
	}

	if rand.Float64() < insertDeleteProbability { // 5% chance to add a new gene
        newGene := generateRandomGene()
        mutatedGenome = append(mutatedGenome, newGene)
    }

	return mutatedGenome
}

func LevenshteinDistance(arr1, arr2 []string) int {
	len1 := len(arr1)
	len2 := len(arr2)

	// Create a 2D slice to store the distances
	dist := make([][]int, len1+1)
	for i := range dist {
		dist[i] = make([]int, len2+1)
	}

	// Initialize the distance matrix
	for i := 0; i <= len1; i++ {
		dist[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		dist[0][j] = j
	}

	// Compute the distance matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if arr1[i-1] != arr2[j-1] {
				cost = 1
			}
			dist[i][j] = min(
				dist[i-1][j]+1,     // Deletion
				dist[i][j-1]+1,     // Insertion
				dist[i-1][j-1]+cost, // Substitution
			)
		}
	}

	return dist[len1][len2]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	return int(math.Min(float64(a), math.Min(float64(b), float64(c))))
}

// Similarity calculates the similarity as a percentage based on Levenshtein distance
func genomeSimilarity(arr1, arr2 []string) float64 {
	distance := LevenshteinDistance(arr1, arr2)
	maxLen := math.Max(float64(len(arr1)), float64(len(arr2)))

	// If both arrays are empty, they are identical
	if maxLen == 0 {
		return 100.0
	}

	// Calculate similarity percentage
	similarity := (1.0 - float64(distance)/maxLen) * 100.0
	return similarity
}

func extractEdgeList(nnet *NeuralNetwork) error {
	file, err := os.Create("neural_output.brain")
	if err != nil {
		fmt.Println("Unable to create output neural network file:", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	neuralCsvData := [][]string{
		{"source", "target", "value", "color1", "color2"},
	}

	for nodeHash, _ := range nnet.NetworkMap {
		for _, outgoingEdge := range nnet.NetworkMap[nodeHash][0] {
			newEntry := []string{}

			inputNode := nnet.NodeMap[outgoingEdge.Source]
			newEntry = append(newEntry, getNeuronString(inputNode))
			inputColor := neuronGetColor(inputNode)

			destNode := nnet.NodeMap[outgoingEdge.Destination]
			newEntry = append(newEntry, getNeuronString(destNode))
			destColor := neuronGetColor(destNode)

			weight := fmt.Sprintf("%f", outgoingEdge.Weight)
			newEntry = append(newEntry, weight)

			newEntry = append(newEntry, inputColor)
			newEntry = append(newEntry, destColor)

			neuralCsvData = append(neuralCsvData, newEntry)

		}
	}

	for _, wire := range neuralCsvData {
		if err := writer.Write(wire); err != nil {
			fmt.Println("Error writing wire!", err)
			return err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil{
		fmt.Println("Error flushing output", err)
		return err
	}
	return nil
}

func signalHandler(rpy_flag bool){
	// Block until a signal is received
	sig := <-sigChan

	if rpy_flag {
		fmt.Println("\nReceived signal:", sig)
		// fmt.Println(evolutionTree[0].Children)

		saveReplayToFile()
		// newRpy, err := loadFromFile(filename)
		// if err != nil {
		// 	fmt.Println("Error loading array:", err)
		// 	return
		// }
		// fmt.Println(len(newRpy), len(rpy))


		// if reflect.DeepEqual(newRpy, rpy) {
		// 	println("Slices are equal")
		// } else {
		// 	println("Slices are not equal")
		// 	fmt.Println(*newRpy[0][10])
		// }
		// Exit the program after cleanup
	}

	os.Exit(0)
}

// rep check
// func checkGenomeRep(child []string, parent []string) {
// 	for i := range child {
// 		if child[i] != parent[i]{
// 			fmt.Println("not equal")
// 			os.Exit(1)
// 		}
// 	}
// }
//
//


func generateGenusName() string{
	genusRoots := []string{"Aero", "Crypto", "Neo", "Proto", "Hydro", "Pyro", "Thermo", "Xylo", "Lacto", "Phyllo", "Macro", "Micro"}
	suffixes := []string{"us", "a", "um", "is", "os"}

	return genusRoots[rand.Intn(len(genusRoots))] + suffixes[rand.Intn(len(suffixes))]
}

func generateSpeciesName() string{
	speciesRoots := []string{"phagus", "tropus", "donta", "cera", "morpha", "soma", "ptera", "genes", "lithos", "derma", "carpa", "rhiza"}
    suffixes := []string{"is", "a", "us", "um"}

	return speciesRoots[rand.Intn(len(speciesRoots))] + suffixes[rand.Intn(len(suffixes))]
}

// Function to convert color.Color to a hex string
func colorToHex(c color.Color) string {
    r, g, b, _ := c.RGBA() // Get RGBA values
    return fmt.Sprintf("#%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

func getIndivColor(genome string) (int64,int64,int64){
	// genomeString := ""
	// for _, gene := range genome {
	// 	genomeString += gene
	// }

	// Calculate the SHA-256 hash of the input string
	// hash := sha256.Sum256([]byte(genomeString))
	hash := sha256.Sum256([]byte(genome))
	hexHash := fmt.Sprintf("%x", hash)

	// Extract components for the RGB color
	r, _ := strconv.ParseInt(hexHash[0:2], 16, 32)
	g, _ := strconv.ParseInt(hexHash[2:4], 16, 32)
	b, _ := strconv.ParseInt(hexHash[4:6], 16, 32)

	// Normalize the components to the [0, 255] range
	r = int64(math.Round(float64(r) * 255.0 / 255.0))
	g = int64(math.Round(float64(g) * 255.0 / 255.0))
	b = int64(math.Round(float64(b) * 255.0 / 255.0))

	return r, g, b
}

// function euclideanDistanceRGB(color1, color2) {
//     const r1 = color1.r;
//     const g1 = color1.g;
//     const b1 = color1.b;

//     const r2 = color2.r;
//     const g2 = color2.g;
//     const b2 = color2.b;

//     const dr = r2 - r1;
//     const dg = g2 - g1;
//     const db = b2 - b1;

//     return Math.sqrt(dr * dr + dg * dg + db * db);
// }
