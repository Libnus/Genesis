package main

import (
	"fmt"
	"strconv"
	"sync"
	"os"
	"math"
	"math/rand"
)


func convertTwos(binStr string) int {
	// Determine the sign of the binary number
	isNegative := binStr[0] == '1'

	// Invert the bits and add 1 if it's a negative number
	if isNegative {
		inverted := ""
		for _, bit := range binStr {
			if bit == '0' {
				inverted += "1"
			} else {
				inverted += "0"
			}
		}

		// Add 1 to the inverted value
		carry := 1
		result := ""
		for i := len(inverted) - 1; i >= 0; i-- {
			bit := int(inverted[i]-'0') + carry
			result = strconv.Itoa(bit%2) + result
			carry = bit / 2
		}

		binStr = result
	}

	// Convert the binary number to decimal
	decimalValue, _ := strconv.ParseInt(binStr, 2, 64)

	// Return the result, negating if it was originally negative
	if isNegative {
		return -int(decimalValue)
	}
	return int(decimalValue)
}

// base id for all input and output
type InputNeuron int
const (
	// input
	Px InputNeuron = iota // position x
	Py // position y
	BD // closest border distance
	BDx // distance to closest border in the x-direction
	BDy // distance to closest border in the y-direction
	Rnd // some random number
	Pop // population density in the immediate area
	LMx // last movement in the x-direction
	LMy // last movement in the y-direction
	N   // nutrition level
	// TODO oscillator timed with seasons? Or could this be developed naturally
	// TODO temperature input
	// TODO hunger level / nutritional levels (required vs current level)
	// TODO emit pheromone (non-damaging, how much is determined by genes)
	// TODO emit poison gas (damaging lol)
)

func printInputEnumId(id int) string {
	switch id{
		case int(Px):
			return "Px"
		case int(Py):
			return "Py"
		case int(BD):
			return "BD"
		case int(BDx):
			return "BDx"
		case int(BDy):
			return "BDy"
		case int(Rnd):
			return "Rnd"
		case int(Pop):
			return "Pop"
		case int(LMx):
			return "LMx"
		case int(LMy):
			return "LMy"
	}
	fmt.Println("error couldn't find", id)
	return ""
}

type OutputNeuron int
const (
	// output
	Mr OutputNeuron = iota // move randomly
	Mx // move +/- x
	My // move +/- y
	Mrv // move in the opposite direction currently going type  interface {
    Mfd // move fwd
	Res // responsiveness of the creature (higher output means weights are heavier (which weights are heavier is determined by the genome) )
	Eat // eat creature directly in front of creatures chompers (fwd direction of creature)
	Ko // kill but dont eat
	Dr // drink lean
)

func printOutputEnumId(id int) string {
	switch id{
		case int(Mr):
			return "Mr"
		case int(Mx):
			return "Mx"
		case int(My):
			return "My"
		case int(Mrv):
			return "Mrv"
		case int(Mfd):
			return "Mfd"
		case int(Res):
			return "Res"
		case int(Eat):
			return "Eat"
		case int(Dr):
			return "Dr"
	}
	return ""
}

type NeuronType int
const (
	InputType NeuronType = iota
	OutputType
	InternalType
)

func getNeuronString(node *Neuron) string{
	if node.NeuronType == InputType{
		return printInputEnumId(node.Id)
	} else if node.NeuronType == InternalType{
		return fmt.Sprintf("N%d", node.Id)
	} else{
		return printOutputEnumId(node.Id)
	}
}

func neuronGetColor(node *Neuron) string{
	if node.NeuronType == InputType{
		return "white"
	} else if node.NeuronType == InternalType{
		return "gray"
	} else{
		return "red"
	}
}


// any value greater than 3 represents an internal neuron

//ioSize := 4 // number of input sensory neurons and output neurons
var NUM_INPUTS = 9
var NUM_OUTPUTS = 8
var MAX_INTERNAL = 127 // the maximum number of internal neurons allowed
var MAX_NEURONS = NUM_INPUTS + NUM_OUTPUTS + MAX_INTERNAL

/*
	//name: name of the neuron from the int map. Example: position x, position y, internal neuron, move east etc.
	id: id of the nueron used for storage (can be mapped with the enum above to a specific neuron)
	neuronType: 0 for input, 1 for internal, 2 for output
	currInputSum: float which is the current sum of all inputs
	currOutput: float whcih is the last calculated output of this neuron (may not be up to date and that is okay)
	inputs: number of incoming connections it has
	outputs: number of outgoing connections it has
*/

// a neuron to represent the neuron/node in our network graph
type Neuron struct {
	//name int
	Id int
	NeuronType NeuronType
	CurrInputSum float64
	CurrOutput float64
	Inputs int
	Outputs int
}	

// func createInputNode(name int) *InputSensor{
// 	newInput = InputSensor{name: name}
// 	new

// 	return &InputSensor{}
// }

// a wire struct which contains the informaiton for an edge in our network graph
type Wire struct {
	Weight float64
	Source int // hash of the source node
	Destination int // hash of the destination node
}

// the graph structure itself
// we use id as key to avoid references and to make processing the adding/deleting of genomes as data is primitively stored as an int id
type NeuralNetwork struct {
	//numInternal int // number of internal neurons
	NetworkMap map[int][2][]*Wire // the network where the key is a node id and the key is a list of connections
								 // index 0 represents outgoing connections and index 1 represents incoming connections
	NodeMap map[int]*Neuron // map of all nodes where the key is an id and the value is the the neuron itself for quick access
	InputNeurons []int // nodes sorted topologically (least number of inputs first / input sensory neurons first)
	OutputNeurons []int // all output nodes
}


func hashNeuron(id int, neuronType int) int{
	return neuronType*100 + id
}

// adds a neuron to the neural network
// return the id of the neuron just added
func addNeuron(nnet *NeuralNetwork, id int, neuronType int) *Neuron {
	// if this is an input lets keep track of it for faster processing when doing DFS
	
	neuronHashVal := hashNeuron(id, neuronType)
	_, ok := nnet.NodeMap[neuronHashVal]
	if ok{ // the neuron already exists in the network so return it
		return nnet.NodeMap[neuronHashVal]
	}


	// otherwise add the neuron to the network
	nnet.NodeMap[neuronHashVal] = &Neuron{
		Id: id,
		NeuronType: NeuronType(neuronType),
		CurrInputSum: 0.0,
		CurrOutput: 0.0,
		Inputs: 0,
		Outputs: 0,
	}
	// id int
	// neuronType int
	// currInputSum float64
	// currOutput float64
	// inputs int
	// outputs int  

	if neuronType == 0 { // if input neuron
		nnet.InputNeurons = append(nnet.InputNeurons, neuronHashVal)
	} else if neuronType == 1 { // if output neuron
		nnet.OutputNeurons = append(nnet.OutputNeurons, neuronHashVal)
	}
	return nnet.NodeMap[neuronHashVal]
}	

// strong precondition: NEVER ATTEMPT TO ADD A DUPLICATE WIRE
func addWire(nnet *NeuralNetwork, sourceType int, sourceId int, destType int, destId int, weight float64) *Wire{
	// check to make sure the source and dest neuron exists in the network
	addNeuron(nnet, sourceId, sourceType)
	addNeuron(nnet, destId, destType)

	// never check for duplicate wires because we hate error checking of any kind
	// this is biology for god's sake!

	// TOneverDO hash value is calculcated already in the addNeuron function so we are calling this function uneccessarily here
	sourceHashVal := hashNeuron(sourceId, sourceType)
	destHashVal := hashNeuron(destId, destType)

	nnet.NodeMap[sourceHashVal].Outputs++
	nnet.NodeMap[destHashVal].Inputs++

	// weight float64
	// source hash
	// destination int // hash of the destination node
	newWire := &Wire {
		Weight: weight,
		Source: sourceHashVal,
		Destination: destHashVal,
	}

	tempSlice := nnet.NetworkMap[sourceHashVal]
	tempSlice[0] = append(nnet.NetworkMap[sourceHashVal][0], newWire)
	nnet.NetworkMap[sourceHashVal] = tempSlice

	tempSlice = nnet.NetworkMap[destHashVal]
	tempSlice[1] = append(nnet.NetworkMap[destHashVal][1], newWire)
	nnet.NetworkMap[destHashVal] = tempSlice

	return newWire
}

/*
   takes in a genome and creates a neural network from it
   a genome looks something like this:
   		
   		1 52562f78 3c396612 4989b501 039c5fbd

   the first hex represents the number of internal neurons in the network (to easily map any random id number to a given input sensor or output)
   the rest of the genome represents a 32-bit binary number which represents a brain wiring/connection in our neural network

   will return nil if no error in reading genome otherwise returns the error
*/
func processGenome(genome []string) (*NeuralNetwork, error) {
	// process the number of internal neurons
	// numInternal, err := strconv.ParseInt(genome[0], 16, 64)
	// if err != nil {
	// 	fmt.Printf("[ ERR ] couldn't parse number of internal neurons! %s\n", err)
	// 	return err
	// }

	// create a new network
	newNetwork := new(NeuralNetwork)
	// newNetwork.numInternal = numInternal
	newNetwork.NetworkMap = make(map[int][2][]*Wire)
	newNetwork.NodeMap = make(map[int]*Neuron)
	newNetwork.InputNeurons = []int{}


	for i := 0; i < len(genome); i++ {
		// fmt.Printf("\n\n[ GENEOTRON ] Processing gene %d\n", i)

		geneInt, err := strconv.ParseUint(genome[i], 16, 64)
		if err != nil {
			fmt.Printf("[ ERR ] couldn't parse genome! %s\n", err)
			fmt.Println()
			fmt.Println(genome[i])
			os.Exit(1)
		}
		gene := fmt.Sprintf("%032b", geneInt)

		// aquire the sourceId by modding the value parsed by the number of genes
		// we mod to allow any bit in the gene to be modified through mutation
		sourceType, _ := strconv.Atoi(string(gene[0]))

		sourceIdInt, err := strconv.ParseUint(gene[1:8], 2, 7)

		if err != nil{
			fmt.Printf("[ ERR ] couldn't parse source ID of gene %d; %s\n", i, err)
		} 

		var sourceId int

		// if from an input sensory neuron
		if sourceType == 0 {
			sourceId = int(sourceIdInt) % NUM_INPUTS
			// fmt.Printf("[ G-TRON ] Wire source is of type input with ID %d\n", sourceId)
		} else{
			sourceId = int(sourceIdInt) % MAX_INTERNAL
			// fmt.Printf("[ G-TRON ] Wire source is of type internal neuron with ID %d\n", sourceId)
			sourceType = int(InternalType)
		}

		// add the source neuron
		addNeuron(newNetwork, sourceId, sourceType)

		destType, _ := strconv.Atoi(string(gene[8]))
		// fmt.Println(string(gene))
		// destType++ // destination neurons can only be another internal neuron or an output neuron so we offset this binary int by 1
		// 				  // 0 represents an input neuron, 1 and internal neuron, 2 an output neuron

		destIdInt, err := strconv.ParseUint(gene[9:16], 2, 7)
		if err != nil{
			fmt.Printf("[ ERR ] couldn't parse dest ID of gene %d; %s\n", i, err)
		} 
		//f1351fe3

		// DESTINATION NEURON
		var destId int
		if destType == 0 {
			destId = int(destIdInt) % NUM_OUTPUTS
			// fmt.Printf("[ G-TRON ] Wire destination is of type output with ID %d destType %d\n", destId, destType)
			destType = int(OutputType)
		} else{
			destId = int(destIdInt) % MAX_INTERNAL
			// fmt.Printf("[ G-TRON ] Wire destination is of type internal neuron with ID %d destType %d\n", destId, destType)
			destType = int(InternalType)
		}

		// add the dest neuron
		addNeuron(newNetwork, destId, destType)

		// WEIGHT OF WIRE
		weightInt := convertTwos(gene[16:])

		weight := float64(weightInt) / float64(8192) // 8192 (creates the range -4.0 to 4.0)

		// fmt.Println("[ G-TRON ] Wire weight is", weight)

		addWire( newNetwork, sourceType, sourceId, destType, destId, weight )
	}	

	return newNetwork, nil
}


func printNode(node *Neuron) {
	switch node.NeuronType{
	case InputType:
		fmt.Printf("   Input neuron: %s\n", printInputEnumId(node.Id))
	case OutputType:
		fmt.Printf("   Output neuron: %s\n", printOutputEnumId(node.Id))
	case InternalType:
		fmt.Println("   Internal neuron:", node.Id)
	}
}

func printWire(nnet *NeuralNetwork, wire *Wire) {
	fmt.Println(" Starting neuron:")
	sourceNode := nnet.NodeMap[wire.Source]
	printNode(sourceNode)

	fmt.Println(" Ending neuron:", )
	endNode := nnet.NodeMap[wire.Destination]
	printNode(endNode)

	fmt.Println(" Wire Weight:", wire.Weight)
	fmt.Println()
}

func printNeural(nnet *NeuralNetwork) {
	fmt.Println("[ PRINT ] printing neural network nodes...")

	for _, node := range nnet.NodeMap {
		printNode(node)
	}

	fmt.Println()

	fmt.Println("[ PRINT ] printing wires...")
	for nodeHash, _ := range nnet.NetworkMap {
		for _, wire := range nnet.NetworkMap[nodeHash][0]{
			printWire(nnet, wire)
		}
	}
}

func calculateInput(indiv *Mite, id int) float64{
	switch InputNeuron(id) {
	case Px:
		return float64(indiv.X) / float64(ROWS)
	case Py:
		return float64(indiv.Y) / float64(COLS)
	case Rnd:
		return rand.Float64()
	case N:
		return indiv.Nutrition
	}

	return 0.0
}



// ============== BRAIN TRAVERSAL

func traverseNode(neuralMu map[int]*sync.Mutex, visitedMu map[int]*sync.Mutex, visitedNeurons map[int]*bool, indiv *Mite, nnet *NeuralNetwork, nodeHash int) {
	//visitedMu[nodeHash].Lock()
	neuralMu[nodeHash].Lock()
	if *visitedNeurons[nodeHash] {
		//fmt.Println(" [ DFS ] NODE ALREADY VISITED", nodeHash)
		neuralMu[nodeHash].Unlock()
		return
	} // we have already visited this node


	*visitedNeurons[nodeHash] = true
	neuralMu[nodeHash].Unlock()


	currNeuron := nnet.NodeMap[nodeHash] // current neuron we are traversing over


	/* CUDA
	 * 1.) pre-traverse no computation just calculate paths
	 * 2.) computation []
	 *
	 */

	neuronOutput := 0.0
	// if on an input node calculate sensory inputs
	if currNeuron.NeuronType == InputType {
		neuronOutput = calculateInput(indiv, currNeuron.Id)
	} else{
		// sum all incoming edges
		for _, incomingEdge := range nnet.NetworkMap[nodeHash][1] {
			// incomingEdge is of type wire
			//fmt.Println("[ DFS ]", nodeHash, "calculating input from", incomingEdge.source, "::", nnet.nodeMap[incomingEdge.source].currOutput)
			neuralMu[incomingEdge.Source].Lock()
			neuronOutput += nnet.NodeMap[incomingEdge.Source].CurrOutput * incomingEdge.Weight
			neuralMu[incomingEdge.Source].Unlock()
		}
		neuronOutput = math.Tanh(neuronOutput)
	}

	neuralMu[nodeHash].Lock()
	currNeuron.CurrOutput = neuronOutput
	neuralMu[nodeHash].Unlock()

	//fmt.Println("   calculated output as:", neuronOutput)

	//visitedMu[nodeHash].Lock()
	//visitedMu[nodeHash].Unlock()


	// for each neighbor create a new routine
	var wgTraverse sync.WaitGroup
	for _, outgoingEdge := range nnet.NetworkMap[nodeHash][0] {
		wgTraverse.Add(1)

		go func(outgoingEdge *Wire){
			defer wgTraverse.Done()
			traverseNode(neuralMu, visitedMu, visitedNeurons, indiv, nnet, outgoingEdge.Destination)
		}(outgoingEdge)
	}

	wgTraverse.Wait()
}

// traverse over neural network using concurrent DFS and return a map where the
// key is a output neuron id and the value is a float64 represneting the output
//
// TODO input indiv not just the network
func calcNeuralPotential(indiv *Mite, nnet *NeuralNetwork) map[OutputNeuron]float64{
	var wgNeural sync.WaitGroup // wait group for brain traversal
	neuralMu := make(map[int]*sync.Mutex)

	// fmt.Println("[ NNET ] Calculating neural output potential...\n		Starting DFS...")

	for nodeHash, _ := range nnet.NodeMap{
		neuralMu[nodeHash] = &sync.Mutex{}
	}

	// for each input node create a go routine and then for each neighbor start another routine
	for _, nodeHash := range nnet.InputNeurons {
		// start a go routine
		wgNeural.Add(1)

		// since we are basically creating a DFS for every input node lets create a visited nodes set for each search
		visitedNeurons := make(map[int]*bool)
		visitedMu := make(map[int]*sync.Mutex)

		for node, _ := range nnet.NodeMap{ // for every key
			bar := false
			visitedNeurons[node] = &bar
			//visitedMu[node] = &sync.Mutex{}
		}

		go func(hash int) {
			//fmt.Println("   Traversing input node", hash)
			defer wgNeural.Done()
			traverseNode(neuralMu, visitedMu, visitedNeurons, indiv, nnet, hash)
		}(nodeHash)
	}

	wgNeural.Wait()

	// fmt.Println("Done")

	netOutput := make(map[OutputNeuron]float64)

	// wait for DFS to finish then loop over all output nodes and organize results
	for _, nodeHash := range nnet.OutputNeurons {
		outputNeuron := nnet.NodeMap[nodeHash]
		netOutput[OutputNeuron(outputNeuron.Id)] = outputNeuron.CurrOutput
	}

	// fmt.Println(netOutput)
	return netOutput
}



// TODO move to different file as this pretains to game state rather than the actual neural network
// DOES NOT LOCK POSITION
// func checkCollisions(x int, y int) bool {
// 	if x < 0 || x >= 128 { return true }
// 	if y < 0 || y >= 128 { return true }

// 	// check for other mites in the new square
// 	gridMu[x][y].Lock()
// 	return gridOccupy[x][y];
// }

/*
	this function performs an action for an organism in a single step which includes calculating the neural
	output to get the actions performed

	NOTE: multiple actions can be performed by an organism in a single step even if the action directly goes
	against a previous action (i.e. moving right then deciding to move left)
 .
	however, once a neuron has been fired, it cannot be refired
*/
func stepOrganism(indiv *Mite) {
	if CURR_STEP - indiv.Birth >= 200{ // too old
		indiv.Dead = true
	}

	if indiv.Dead{
		return
	}

	netOutput := calcNeuralPotential(indiv, indiv.Nnet)

	// didMove := false
	indiv.didMove = false

	// perform the actions for each output
	// for each output neuron determine the neuron type and perform action
	for neuronId, output := range netOutput {
		if rand.Float64() > math.Abs(output) { continue }

		switch neuronId {
		case Mr:
				// pick a random direction and move
			xOffset, yOffset := rand.Intn(3) - 1, rand.Intn(3) - 1
			newX := indiv.X + xOffset
			newY := indiv.Y + yOffset

			// x move
			moveMite( indiv, newX, indiv.Y )

			// y move
			moveMite( indiv, indiv.X, newY )
			indiv.didMove = true
		case Mx:
			newX := indiv.X +  int(output / math.Abs(output))

			moveMite( indiv, newX, indiv.Y )
			indiv.didMove = true
		case My:
			newY := indiv.Y +  int(output / math.Abs(output))

			moveMite( indiv, indiv.X, newY )
			indiv.didMove = true
		case Ko:
			//TODO in forward direction rather than just random
			xOffset, yOffset := rand.Intn(3) - 1, rand.Intn(3) - 1
			killX, killY := indiv.X + xOffset, indiv.Y + yOffset

			// dont kill myself ?
			if killX == indiv.X && killY == indiv.Y { continue }

			if (killX >= 0 && killX < 128) && (killY >= 0 && killY < 128) {
				gridMu[ killX ][ killY ].Lock()
				// its possible that after locking the grid position is now empty
				// the kill is ignored or considered a "failed" attempt
				if gridOccupy[ killX ][ killY ] != nil {
					// TODO kill organism there
					gridOccupy[killX][killY].Dead = true
					indiv.Nutrition += 2*gridOccupy[killX][killY].Nutrition
					// fmt.Println("KILLINGNGNGNGNGNNGNG")

					// if *rpy_flag{
					// 	addEvent(Kill, indiv, gridOccupy[killX][killY])
					// }
					addKill()
				}
				gridMu[ killX ][ killY ].Unlock()
			}
		}
	}

	if indiv.Sick > 0{
		indiv.Nutrition -= 0.5

		if indiv.Nutrition > 0.0 {
			indiv.Sick--
		}

		neighbors := [][2]int{
			{-1, -1}, // Top-left
			{-1, 0},  // Top
			{-1, 1},  // Top-right
			{0, -1},  // Left
			{0, 1},   // Right
			{1, -1},  // Bottom-left
			{1, 0},   // Bottom
			{1, 1},   // Bottom-right
		}

		for _, neighbor := range neighbors {
			x := indiv.X + neighbor[0]
			y := indiv.Y + neighbor[1]

			if (x < 0 || x >= 128) || (y < 0 || y >= 128) { continue }


			// TODO fresh sick mite spread easier
			if rand.Float64() < 0.25{ // it spreads 0.1
				gridMu[ x ][ y ].Lock()
				if gridOccupy[x][y] != nil && gridOccupy[x][y].Sick == 0 && !gridOccupy[x][y].Dead  {
					gridOccupy[ x ][ y ].Sick = rand.Intn(50) + 1
				}
				gridMu[ x ][ y ].Unlock()
			}
		}
	}

	if indiv.Nutrition <= 0.0 {
		// if indiv.Sick > 0{
		// 	removeSick()
		// 	// CURR_SICK--
		// }
		indiv.Dead = true
		return
	}

	// if didMove {
	// 	didMove = false
	// }
	if indiv.didMove {
		indiv.Nutrition -= 0.05
	}else { indiv.Nutrition += 0.1 }

}
