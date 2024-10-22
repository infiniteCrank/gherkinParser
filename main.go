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
		// Check if child is a Background
		if background := child.Background; background != nil {
			// Collect Steps from the Background if it exists
			for _, step := range background.Steps {
				feature.Background = append(feature.Background, step.Text)
			}
		}

		// Check if child is a Scenario
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

// Helper function to identify common steps across scenarios
func findCommonSteps(scenarios []Scenario) ([]string, []Scenario) {
	stepCount := make(map[string]int)
	for _, scenario := range scenarios {
		for _, step := range scenario.Steps {
			stepCount[step]++
		}
	}

	// Identify steps that are common to all scenarios
	var commonSteps []string
	for step, count := range stepCount {
		if count == len(scenarios) {
			commonSteps = append(commonSteps, step)
		}
	}

	// Prepare new scenarios by removing common steps
	var newScenarios []Scenario
	for _, scenario := range scenarios {
		var newScenario Scenario
		newScenario.Name = scenario.Name
		newScenario.Tags = scenario.Tags

		for _, step := range scenario.Steps {
			if !contains(commonSteps, step) {
				newScenario.Steps = append(newScenario.Steps, step)
			}
		}
		newScenarios = append(newScenarios, newScenario)
	}

	return commonSteps, newScenarios
}

// Helper function to check if a slice contains a given value
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// Generate a Gherkin feature file from the structured Feature data
func generateFeatureFile(feature Feature) string {
	var builder strings.Builder

	// Write the feature name
	builder.WriteString("Feature: " + feature.Name + "\n\n")

	// Find and append common background steps
	commonSteps, newScenarios := findCommonSteps(feature.Scenarios)

	// If a background exists, append new common steps; otherwise, create a new background
	if len(feature.Background) > 0 {
		for _, step := range commonSteps {
			if !contains(feature.Background, step) {
				feature.Background = append(feature.Background, step)
			}
		}
	}

	if len(feature.Background) > 0 {
		builder.WriteString("Background:\n")
		for _, step := range feature.Background {
			builder.WriteString("  Given " + step + "\n")
		}
		builder.WriteString("\n")
	}

	// Include new Scenarios
	for _, scenario := range newScenarios {
		builder.WriteString("Scenario: " + scenario.Name + "\n")

		// Include Tags if present
		for _, tag := range scenario.Tags {
			builder.WriteString("  @" + tag + "\n")
		}

		var stepRange []string
		stepRange = scenario.Steps
		if len(newScenarios) == 1 {
			stepRange = commonSteps
		}
		// Include Steps with appropriate prefixes based on the context
		for _, step := range stepRange {
			if strings.Contains(step, "enters") || strings.Contains(step, "entered") || strings.Contains(step, "opened") {
				builder.WriteString("  Given " + step + "\n") // Treating entering or opening steps as Given
			} else if strings.Contains(step, "clicks") {
				builder.WriteString("  When " + step + "\n") // Using When for clicking actions
			} else if strings.Contains(step, "should") {
				builder.WriteString("  Then " + step + "\n") // Using Then for assertion steps
			} else {
				builder.WriteString("  When " + step + "\n") // Fallback to When for general steps
			}
		}

		builder.WriteString("\n")
	}

	return builder.String()
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
		parsedFeature := parseFeatureFile(string(fileContent))
		regeneratedFeature := generateFeatureFile(parsedFeature)

		// Output the regenerated feature file
		w.Write([]byte(regeneratedFeature))
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
