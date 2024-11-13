package main

import (
    "encoding/json"
    // "time"
    "fmt"
    "sync"
    "net/http"
	// "image/color"
)

type DataPoint struct {
    X int   `json:"X"`  // X axis
    Y float64 `json:"Y"`  // Y axis
}

// converted from popdata point to send to frontend
type BarDataPoint struct {
	X int
	Y string
	Color string
}

// before processing bar
// keep track of population
type PopDataPoint struct{
	num int
	original *Mite
	color string
}


type TreeNode struct {
	Name string	`json:name`
	Color string `json:color`
	Children []*TreeNode `json:"children,omitempty"`
}

var genData []DataPoint
var sickData []DataPoint
var brainData []DataPoint
var killingData []DataPoint
var treeData = &TreeNode{
	Name: "root",
	Children: []*TreeNode{},
}

var speciesData = map[string]PopDataPoint{}

var sickMu sync.Mutex

func appendSickData(){
	sickData = append(sickData, DataPoint{
		X: CURR_STEP,
		Y: float64(CURR_SICK),
	})
}

func appendBrainData(){
	brainData = append(brainData, DataPoint{
		X: CURR_STEP,
		Y: float64(10*AVG_BRAIN_SIZE),
	})
}

func addSick(){
	sickMu.Lock()
	CURR_SICK++
	sickMu.Unlock()
}

func removeSick(){
	sickMu.Lock()
	CURR_SICK--
	fmt.Println("dead", CURR_SICK)
	sickMu.Unlock()
}

func appendGenData(gen int){
	// genData = append(genData, opts.LineData{Value: gen})
	// genDataX = append(genDataX, strconv.Itoa(CURR_STEP))
	genData = append(genData, DataPoint{X: CURR_STEP, Y: float64(gen)})
}


func addSpecies(species *Mite){
	speciesMu.Lock()
	if data, ok := speciesData[getName(species)]; ok {
        // Key exists, increment the num field
        data.num++
        speciesData[getName(species)] = data // Reassign the modified struct back to the map
    } else {
        // Key does not exist, create a new entry
        speciesData[getName(species)] = PopDataPoint{
            num:   1,
			original: species,
            color: colorToHex(species.Color), // Example color
        }
        // fmt.Println("added new species", speciesData[getName(species)])
    }
	speciesMu.Unlock()
}


// assumes the species is already added to speciesData
func removeSpecies(species *Mite){
	speciesMu.Lock()
	data, _ := speciesData[getName(species)]
	data.num--

	if data.num == 0{
		delete (speciesData, getName(species))
	} else{
		speciesData[getName(species)] = data
	}
	speciesMu.Unlock()
}

func appendKillingData(){
	killingData = append(killingData, DataPoint{X: CURR_STEP, Y: float64(CURR_KILL)})
	CURR_KILL = 0
}


func treeInsert(parent string, child string, color string, node *TreeNode) bool{
	if node.Name == parent{
		newSpecies := TreeNode {
			Name: child,
			Color: color,
			Children: []*TreeNode{},
		}

		node.Children = append(node.Children, &newSpecies)
		return true
	}

	for _, next := range node.Children{
		if treeInsert(parent, child, color, next){
			return true
		}
	}

	return false
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == http.MethodOptions {
        return
    }


	parsedSpecies := []BarDataPoint{}

	speciesMu.Lock()
	// fmt.Println(speciesData)
	for key, val := range speciesData{
		parsedSpecies = append(parsedSpecies, BarDataPoint{
			Y: key,
			X: val.num,
			Color: val.color,
		})
	}
	speciesMu.Unlock()

    data := map[string]interface{}{
		"pop_data": genData,
		"death_data": killingData,
		"tree_data": *treeData,
		"species_pop": parsedSpecies,
		"sick_data": sickData,
		"brain_data": brainData,
	} // Example of dynamic data
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
