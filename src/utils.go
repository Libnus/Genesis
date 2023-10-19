package main

import (
	"fmt"
	"math/rand"
	"encoding/csv"
	"os"
	"strconv"
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

// takes in a genome, creates a copy, and mutates it
func mutateGenome(genome []string) []string {
	mutationRate := 15 // mutationRate / 15000 for prob

	var mutatedGenome []string
	for _, gene := range genome {
		gBinCopy, _ := strconv.ParseInt(gene, 16, 64)
		gCopy := fmt.Sprintf("%032b", gBinCopy)

		newG := []rune(gCopy)
		for i, base := range gCopy{
			if rand.Intn(15001) == mutationRate{
				fmt.Println("mutation occurred!")
				if base == '0' {
					newG[i] = '1'
				} else { newG[i] = '0' }
			}
		}
		gCopyInt, _ := strconv.ParseUint(string(newG), 2, 32)
		mutatedGenome = append(mutatedGenome, fmt.Sprintf("%X", gCopyInt))
	}

	return mutatedGenome
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

	for nodeHash, _ := range nnet.networkMap {
		for _, outgoingEdge := range nnet.networkMap[nodeHash][0] {
			newEntry := []string{}

			inputNode := nnet.nodeMap[outgoingEdge.source]
			newEntry = append(newEntry, getNeuronString(inputNode))
			inputColor := neuronGetColor(inputNode)

			destNode := nnet.nodeMap[outgoingEdge.destination]
			newEntry = append(newEntry, getNeuronString(destNode))
			destColor := neuronGetColor(destNode)

			weight := fmt.Sprintf("%f", outgoingEdge.weight)
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

// rep check
// func checkGenomeRep(child []string, parent []string) {
// 	for i := range child {
// 		if child[i] != parent[i]{
// 			fmt.Println("not equal")
// 			os.Exit(1)
// 		}
// 	}
// }
