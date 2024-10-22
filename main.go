package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	gherkin "github.com/cucumber/gherkin/go/v27"
	messages "github.com/cucumber/messages/go/v22"
	"github.com/gorilla/mux"
)

type Scenario struct {
	Name     string
	Steps    []string  // List of steps
	Tags     []string  // List of tags
	Examples []Example // List of example tables
}

type Example struct {
	Title string
	Rows  []Row
}

type Row struct {
	Cells []string
}

type Feature struct {
	Name       string
	Scenarios  []Scenario
	Background []string
}

func parseFeatureFile(fileContent string) Feature {
	uuid := &messages.UUID{}
	reader := strings.NewReader(fileContent)
	gherkinDocument, err := gherkin.ParseGherkinDocument(reader, uuid.NewId)
	if err != nil {
		fmt.Println("Error parsing Gherkin document:", err)
		return Feature{}
	}

	var feature Feature
	feature.Name = gherkinDocument.Feature.Name

	// Iterate over children of the Feature directly
	for _, child := range gherkinDocument.Feature.Children {
		if background := child.Background; background != nil {
			// Collect Steps from the Background if it exists
			for _, step := range background.Steps {
				feature.Background = append(feature.Background, step.Text)
			}
		}

		if scenario := child.Scenario; scenario != nil {
			var newScenario Scenario
			newScenario.Name = scenario.Name

			// Collect Tags if any
			for _, tag := range scenario.Tags {
				newScenario.Tags = append(newScenario.Tags, tag.Name)
			}

			// Collect Steps
			for _, step := range scenario.Steps {
				newScenario.Steps = append(newScenario.Steps, step.Text)
			}

			// Add to the feature's scenarios
			feature.Scenarios = append(feature.Scenarios, newScenario)
		}
	}

	return feature
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		file, _, err := r.FormFile("feature")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read the file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			return
		}

		// Parse the Gherkin feature file
		parseFeatureFile(string(fileContent))

		w.Write([]byte("Feature file parsed successfully!"))
	} else {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/upload", uploadHandler).Methods("POST")

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", router)
}
