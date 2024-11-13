package main

import (
	"fmt"
    "encoding/json"
    "encoding/gob"
	"bytes"
	"image/color"
	"github.com/klauspost/compress/zstd"
	"os"
	"sync"
	// "reflect"
	"log"
)

type EventType int
const (
	Kill EventType = iota // actual murder
	Death // just a random death
	Birth // a birth event
)


var EVENT_ID = 0

// stores an event for replay logging
type Event struct{
	Id int 				// event id
	Type EventType 		// the type of event
	Source int  // who is the event referring to
	Target int  // the target of event (for killing) null if no target
}

type Replay struct{
	ReplayGrid [][][]int
	Organisms map[int]*Mite
	Events map[int][]Event
}

var replayMu sync.Mutex

type BrainNodeJson struct {
	Id string `json:"id"`
	Group int `json:"group"`
}

type BrainEdgeJson struct{
	Source string `json:"source"`
	Target string `json:"target"`
	Value float64 `json:"value"`
}


func init() {
	gob.Register(&Replay{})
    gob.Register(&Mite{})
	gob.Register(&Event{})
	gob.Register(&NeuralNetwork{})
	gob.Register(color.RGBA{})
	gob.Register(map[int]*Mite{})
}

func addEvent(eventType EventType, source *Mite, target *Mite) Event {

	sourceId := -1
	targetId := -1
	if source != nil { sourceId = source.Id }
	if target != nil { targetId = target.Id }

	event := Event{
		Id: EVENT_ID,
		Type: eventType,
		Source: sourceId,
		Target: targetId,
	}


	EVENT_ID++


	replayMu.Lock()
	_, ok := replay.Events[CURR_STEP]
	if !ok {
		replay.Events[CURR_STEP] = []Event{}
	}
	replay.Events[CURR_STEP] = append(replay.Events[CURR_STEP], event)
	replayMu.Unlock()

	return event
}

var replay *Replay

func initReplay(){
	replay = &Replay{
		Organisms: make(map[int]*Mite),
		Events: make(map[int][]Event),
	}
}


func addMiteToReplay(mite *Mite){
	_, exists := replay.Organisms[mite.Id]
	if !exists{
		replay.Organisms[mite.Id] = mite
	}
}

func saveGrid(){
	newGrid := make([][]int,ROWS)

	for i, _ := range gridOccupy {
		newGrid[i] = make([]int, COLS)
		for j, mite := range gridOccupy[i]{
			if mite != nil{
				addMiteToReplay(mite)
				newGrid[i][j] = mite.Id
			} else{
				newGrid[i][j] = -1
			}
		}
	}

	replay.ReplayGrid = append(replay.ReplayGrid, newGrid)
}



// func loadReplayFromFile(filename string) (*Replay,error) {
// 	// Open the file for reading
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil,err
// 	}
// 	defer file.Close()


// 	zstdReader, err := zstd.NewReader(file)
//     if err != nil {
//         return nil,err
//     }
//     defer zstdReader.Close()

// 	var data Replay

// 	// var data Replay
//     // Decode data with gob after decompression
//     decoder := gob.NewDecoder(zstdReader)
// 	if err := decoder.Decode(&data); err != nil {
// 		fmt.Println(err);
// 		return nil,err
// 	}

// 	fmt.Println("loaded ", len(replay.ReplayGrid), " steps")

// 	return &data,nil
// }

func loadReplayFromFile(filename string) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()


	zstdReader, err := zstd.NewReader(file)
    if err != nil {
        return err
    }
    defer zstdReader.Close()

	// var data Replay

	// var data Replay
    // Decode data with gob after decompression
    decoder := gob.NewDecoder(zstdReader)
	if err := decoder.Decode(&replay); err != nil {
		fmt.Println(err);
		return err
	}

	fmt.Println("loaded ", len(replay.ReplayGrid), " steps")

	return nil
}

func DeepCopy(src, dst interface{}) error {
    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(src); err != nil {
        return err
    }
    return gob.NewDecoder(&buf).Decode(dst)
}

var births = 0
var deaths = 0

func saveReplayToFile() error{
	gridLockMu.Lock()
	// Perform any cleanup or shutdown tasks here
	fmt.Println("saving replay...")
	filename := "organisms.rpy"
	// if err := saveToFile(filename, replay); err != nil {
	// 	return err
	// }

	chartFile, err := os.Create("charts.json")
	if err != nil {
		fmt.Println("Error saving:", err)
		return err
	}
	defer chartFile.Close()


    data := map[string]interface{}{
		"pop_data": genData,
		"death_data": killingData,
		"tree_data": *treeData,
	} // Example of dynamic data
	chartEncoder := json.NewEncoder(chartFile)
	err = chartEncoder.Encode(data)
	if err != nil {
        log.Fatalf("Failed to encode to JSON: %v", err)
    }

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error saving:", err)
		return err
	}
	defer file.Close()

    writer, err := zstd.NewWriter(file, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
	if err != nil {
		fmt.Println("error", err)
        return err
    }

	// Create a new encoder and encode the data
	encoder := gob.NewEncoder(writer)
	if err := encoder.Encode(replay); err != nil {
		fmt.Println("Error saving:", err)
		return err
	}
	writer.Close()

	// var newReplay Replay
	// DeepCopy(&replay, &newReplay)

	// loadReplayFromFile("organisms.rpy")


	// newReplay, err := loadReplayFromFile("organisms.rpy")
    // if err != nil {
    //     log.Fatalf("Failed to load gob file: %v", err)
    // }


	// if reflect.DeepEqual(replay, newReplay) {
    //     fmt.Println("The gob file was saved and loaded correctly.")
	// 	// fmt.Println(reflect.DeepEqual(replay.ReplayGrid[2], newReplay.ReplayGrid[2]))
	// 	pop := 0
	// 	newBirths := 0
	// 	newDeaths := 0

	// 	for i, _ := range replay.Events{
	// 		for _, event := range replay.Events[i]{
	// 			if event.Type == Birth{
	// 				newBirths++
	// 				pop++
	// 			} else if event.Type == Death{
	// 				pop--
	// 				newDeaths++
	// 			}
	// 		}
	// 	}

	// 	fmt.Println(len(organisms), pop)
	// 	fmt.Println("births compared", births, newBirths-1000)
	// 	fmt.Println("deaths compared", deaths, newDeaths)

    // } else {
    //     fmt.Println("Data mismatch: The gob file was not saved correctly.")

    //     // fmt.Printf("Loaded: %+v\n", newReplay)
    // }




	gridLockMu.Unlock()
	os.Exit(0)
	return nil

}


// outputs a brain to a file
func outputBrain(nnet *NeuralNetwork) error{
	file, err := os.Create("brain.json")
	if err != nil {
		fmt.Println("Error saving:", err)
		return err
	}
	defer file.Close()

	nodes := []BrainNodeJson{}
	links := []BrainEdgeJson{}

	// loop over neurons
	for id := range nnet.NodeMap {
		nodes = append(nodes, BrainNodeJson{
			Id: getNeuronString(nnet.NodeMap[id]),
			Group: int(nnet.NodeMap[id].NeuronType),
 		})
	}

	for node, _ := range nnet.NetworkMap {
		for _, link := range nnet.NetworkMap[node][0]{
			input := getNeuronString(nnet.NodeMap[link.Source])
			target := getNeuronString(nnet.NodeMap[link.Destination])
			links = append(links, BrainEdgeJson{
				Source: input,
				Target: target,
				Value:	link.Weight,
			})
		}
	}


	data := map[string]interface{}{
		"nodes": nodes,
		"links": links,
	}


	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
        log.Fatalf("Failed to encode to JSON: %v", err)
    }

	return nil
}

// gets the brains of last iteration
func getBrains(){
	// get brain
	last := replay.ReplayGrid[len(replay.ReplayGrid)-1]

	for i, _ := range last{
		for _, c := range last[i]{
			if c != -1 {
				// output brain
				mite := replay.Organisms[c]
				outputBrain(mite.Nnet)
				return
			}
		}
	}
}
