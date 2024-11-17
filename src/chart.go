package main

import (
    "encoding/json"
    // "time"
    "fmt"
	"log"
	"runtime"
	"os/exec"
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
	treeNode *TreeNode
}

type RequestData struct {
	Message string `json:"message"`
	Value   string    `json:"value"`
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


func addSpecies(species *Mite, treeNode *TreeNode){
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
			treeNode: treeNode,
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

	if data.num <= 0{
		// fmt.Println("removing ", getName(species))
		// delete (speciesData, getName(species))
		data.num = 0
	}
	speciesData[getName(species)] = data
	speciesMu.Unlock()
}

func appendKillingData(){
	killingData = append(killingData, DataPoint{X: CURR_STEP, Y: float64(CURR_KILL)})
	CURR_KILL = 0
}

func treePathSearch(searchItem string, currPath []string, node *TreeNode) []string{
	if node.Name != "root" {
		currPath = append(currPath, node.Name)
	}

	if node.Name == searchItem{

		// fmt.Println(node.Name, searchItem, node.Name == searchItem, currPath)
		return currPath
	}

	if len(node.Children) == 0{
		return nil
	}

	for _, next := range node.Children{
		newPath := append([]string{}, currPath...)
		resultPath := treePathSearch(searchItem, newPath, next)
		if resultPath != nil{
			return resultPath
		}
	}

	return nil
}

func treeInsert(parent string, child string, color string) *TreeNode{

	parentNode := treeData
	if parent != "root"{
		parentNode = speciesData[parent].treeNode
	}

	newSpecies := &TreeNode {
		Name: child,
		Color: color,
		Children: []*TreeNode{},
	}
	parentNode.Children = append(parentNode.Children, newSpecies)
	return newSpecies
	// if node.Name == parent{

	// 	node.Children = append(node.Children, &newSpecies)
	// 	return true
	// }

	// for _, next := range node.Children{
	// 	if treeInsert(parent, child, color, next){
	// 		return true
	// 	}
	// }

	// return false
}

func openBrowser(url string) {
    var err error
    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    }
    if err != nil {
        log.Printf("Failed to open browser: %v", err)
    }
}

// NOTE do not call when program is running recursively
func openBrainBrowser(x,y int){
	// speciesMu.Lock()
	// mite := speciesData[species].original
	// speciesMu.Unlock()
	// miteBrain := getBrain(mite.Nnet)

	// // get brain evolution data
	// evolutionTree := treePathSearch(species, []string{}, treeData)

	// brainEvolutionData := []map[string]interface{}{}

	// speciesMu.Lock()
	// for _, v := range evolutionTree{
	// 	miteBrain := getBrain(speciesData[v].original.Nnet)
	// 	brainEvolutionData = append(brainEvolutionData, map[string]interface{}{
	// 		"name": getName(speciesData[v].original),
	// 		"brain": miteBrain,
	// 	})
	// }
	// speciesMu.Unlock()

	// // Prepare a response
	// data := map[string]interface{}{
	// 	// Status:  "Success",
	// 	// Details: fmt.Sprintf("Received message: %s with value: %d", requestData.Message, requestData.Value),
	// 	"brain": miteBrain,
	// 	"evolve": brainEvolutionData,
	// }

    // jsonData, err := json.Marshal(data)
    // if err != nil {
    //     log.Fatalf("Failed to marshal data: %v", err)
    // }

    // // Send data via POST
    // resp, err := http.Post("http://localhost:8080/brain", "application/json", bytes.NewBuffer(jsonData))
    // if err != nil {
    //     log.Fatalf("Failed to send data: %v", err)
    // }
    // defer resp.Body.Close()

    // if resp.StatusCode == http.StatusOK {
    //     log.Println("Data sent successfully!")
    // } else {
    //     log.Printf("Failed to send data. Status: %v", resp.Status)
    // }

	// NOTE we do not need to lock because no proccesses should be running when this is called
	gridLockMu.Lock()
	species := gridOccupy[x][y]
	gridLockMu.Unlock()

	if species == nil {
		fmt.Println("No species at position!")
		return
	}

	fmt.Println(getName(species), x, y)
    openBrowser("http://localhost:8000/brain.html?id=" + getName(species))
}

func brainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// if r.Method == http.MethodOptions {
	// 	w.WriteHeader(http.StatusOK) // Return 200 OK for OPTIONS requests
	// 	return
	// }
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 	return
	// }

	id := r.URL.Query().Get("id")
	fmt.Println("species name", id)

	speciesMu.Lock()
	mite := speciesData[id].original
	speciesMu.Unlock()
	miteBrain := getBrain(mite.Nnet)

	// get brain evolution data
	evolutionTree := treePathSearch(id, []string{}, treeData)

	brainEvolutionData := []map[string]interface{}{}

	speciesMu.Lock()
	var prev = speciesData[evolutionTree[0]].original
	fmt.Println("Brain similarity test for", id)
	for i, v := range evolutionTree{
		miteBrain := getBrain(speciesData[v].original.Nnet)
		brainEvolutionData = append(brainEvolutionData, map[string]interface{}{
			"name": getName(speciesData[v].original),
			"brain": miteBrain,
		})

		if i > 0{
			similarity := genomeSimilarity(prev, speciesData[v].original)
			fmt.Println("similarity between", getName(prev), "and", getName(speciesData[v].original), similarity*100)
		}
		// prev = speciesData[v].original

	}
	speciesMu.Unlock()

	// Prepare a response
	data := map[string]interface{}{
		// Status:  "Success",
		// Details: fmt.Sprintf("Received message: %s with value: %d", requestData.Message, requestData.Value),
		"brain": miteBrain,
		"evolve": brainEvolutionData,
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)


	// brain similarity test
	// prev := brainEvolutionDta[]
	// for i, brain := brainEvolutionData{ // for each species
	// 	if i == 0 {continue}
	// }

	// Decode the JSON request body into RequestData struct
	// var requestData RequestData
	// if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
	// 	http.Error(w, "Error decoding JSON", http.StatusBadRequest)
	// 	return
	// }

	// Print received data
	// fmt.Printf("Received data: %+v\n", requestData)
	// openBrain(requestData.Value)

	//get brain
	// speciesMu.Lock()
	// mite := speciesData[requestData.Value].original
	// speciesMu.Unlock()
	// miteBrain := getBrain(mite.Nnet)

	// // get brain evolution data
	// evolutionTree := treePathSearch(getName(mite), []string{}, treeData)

	// brainEvolutionData := []map[string]interface{}{}

	// speciesMu.Lock()
	// for _, v := range evolutionTree{
	// 	miteBrain := getBrain(speciesData[v].original.Nnet)
	// 	brainEvolutionData = append(brainEvolutionData, map[string]interface{}{
	// 		"name": getName(speciesData[v].original),
	// 		"brain": miteBrain,
	// 	})
	// }
	// speciesMu.Unlock()

	// // Prepare a response
	// data := map[string]interface{}{
	// 	// Status:  "Success",
	// 	// Details: fmt.Sprintf("Received message: %s with value: %d", requestData.Message, requestData.Value),
	// 	"brain": miteBrain,
	// 	"evolve": brainEvolutionData,
	// }

	// Encode response as JSON
	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(data); err != nil {
	// 	http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	// }
	//
}


// func openBrain(w http.ResponseWriter, r *http.Request) {
// 	data := struct {
// 		Message string
// 		Value   int
// 	}{
// 		Message: "Hello from Go with localStorage",
// 		Value:   42,
// 	}

// 	tmpl := template.New("js")
// 	tmpl, err := tmpl.Parse(jsTemplate)
// 	if err != nil {
// 		http.Error(w, "Error parsing template", http.StatusInternalServerError)
// 		return
// 	}

// 	// Execute the template with data
// 	w.Header().Set("Content-Type", "text/html")
// 	tmpl.Execute(w, data)
// }

func dataHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
	} // example of dynamic data
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
