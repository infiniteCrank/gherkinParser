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

type ScenarioOutline struct {
	Name     string
	Steps    []string  // List of steps, same as in scenario
	Examples []Example // Example tables
}

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
	Name             string
	Scenarios        []Scenario
	ScenarioOutlines []ScenarioOutline // Newly added field
	Background       []string
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
			for _, step := range background.Steps {
				feature.Background = append(feature.Background, step.Text)
			}
		}

		// Check if child is a Scenario
		if scenario := child.Scenario; scenario != nil {
			var newScenario Scenario
			newScenario.Name = scenario.Name

			for _, tag := range scenario.Tags {
				newScenario.Tags = append(newScenario.Tags, tag.Name)
			}

			for _, step := range scenario.Steps {
				newScenario.Steps = append(newScenario.Steps, step.Text)
			}

			feature.Scenarios = append(feature.Scenarios, newScenario)
		}
	}

	// Generate outlines after collecting scenarios
	feature.ScenarioOutlines = findScenarioOutlines(feature.Scenarios)

	return feature
}

// helper function to find scenario outlines
func findScenarioOutlines(scenarios []Scenario) []ScenarioOutline {
	var outlines []ScenarioOutline
	stepMappings := make(map[string][]Scenario) // Map of normalized step strings to scenarios

	// Group scenarios by steps (ignoring tags and names)
	for _, scenario := range scenarios {
		// Create a key based on the steps, ignoring tags or names
		key := strings.Join(scenario.Steps, "|")
		stepMappings[key] = append(stepMappings[key], scenario)
	}

	// Create Scenario Outlines for groups with more than one scenario
	for _, group := range stepMappings {
		if len(group) > 1 {
			var outline ScenarioOutline
			outline.Name = group[0].Name // Base name; you may want to adjust this logic

			// Build the sample table from grouped scenarios
			var exampleRows []Row
			for _, scenario := range group {
				// Assume the first step contains the test data key-value pairs
				row := make([]string, len(scenario.Steps))
				for i := range scenario.Steps {
					row[i] = "DATA_VALUE" // Placeholder; you would replace with actual values
				}
				exampleRows = append(exampleRows, Row{Cells: row})
			}

			outline.Steps = group[0].Steps
			outline.Examples = []Example{{Title: "Example", Rows: exampleRows}}
			outlines = append(outlines, outline)
		}
	}

	return outlines
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

	// Check if common steps exist and append them to the existing background
	commonSteps, newScenarios := findCommonSteps(feature.Scenarios)

	if len(feature.Background) > 0 {
		for _, step := range commonSteps {
			if !contains(feature.Background, step) {
				feature.Background = append(feature.Background, step) // Append new unique steps to the existing background
			}
		}
	} else if len(commonSteps) > 0 {
		// If no existing background, create one
		builder.WriteString("Background:\n")
		for _, step := range commonSteps {
			builder.WriteString("  Given " + step + "\n") // Treating all common steps as Given
		}
		builder.WriteString("\n")
	}

	// Output the existing background if present
	if len(feature.Background) > 0 {
		builder.WriteString("Background:\n")
		for _, step := range feature.Background {
			builder.WriteString("  Given " + step + "\n")
		}
		builder.WriteString("\n")
	}

	// Include the new Scenarios
	for _, scenario := range newScenarios {
		builder.WriteString("Scenario: " + scenario.Name + "\n")

		// Include Tags if present
		for _, tag := range scenario.Tags {
			builder.WriteString("  @" + tag + "\n")
		}

		// Include Steps while determining their prefixes
		for _, step := range scenario.Steps {
			if strings.Contains(step, "entered") || strings.Contains(step, "opened") {
				builder.WriteString("  Given " + step + "\n")
			} else if strings.Contains(step, "clicks") {
				builder.WriteString("  When " + step + "\n")
			} else if strings.Contains(step, "should") {
				builder.WriteString("  Then " + step + "\n")
			} else {
				builder.WriteString("  When " + step + "\n") // Fallback for other step types
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
