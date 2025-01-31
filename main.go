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

type Step struct {
	Text   string // The actual step text
	Prefix string // The keyword prefix (Given, When, Then)
}

type Scenario struct {
	Name     string
	Steps    []Step    // List of steps with keywords.
	Tags     []string  // List of tags.
	Examples []Example // List of example tables.
}
type Example struct {
	Title string
	Rows  []Row
}

type Row struct {
	Cells []string
}

type ScenarioOutline struct {
	Name     string
	Steps    []Step    // List of steps with prefixes
	Tags     []string  // List of tags.
	Examples []Example // Example tables
}

type Feature struct {
	Name             string
	Scenarios        []Scenario
	ScenarioOutlines []ScenarioOutline
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

			//check to see if scenario outline
			if child.Scenario.Keyword == "Scenario Outline" {
				outline := child.Scenario
				var newOutline ScenarioOutline
				newOutline.Name = outline.Name

				// Collect Tags if any
				for _, tag := range outline.Tags {
					newOutline.Tags = append(newOutline.Tags, tag.Name)
				}

				// Collect Steps from the Scenario Outline
				for _, step := range outline.Steps {
					newOutline.Steps = append(newOutline.Steps, Step{Text: step.Text, Prefix: step.Keyword})
				}

				// Collect Example Tables from the Scenario Outline
				for _, example := range outline.Examples {
					var newExample Example
					newExample.Title = example.Description // Set the description as the title

					// Add header
					if example.TableHeader != nil {
						headerCells := make([]string, len(example.TableHeader.Cells))
						for i, cell := range example.TableHeader.Cells {
							headerCells[i] = cell.Value
						}
						newExample.Rows = append(newExample.Rows, Row{Cells: headerCells})
					}

					// Add rows in the TableBody
					for _, row := range example.TableBody {
						rowCells := make([]string, len(row.Cells))
						for i, cell := range row.Cells {
							rowCells[i] = cell.Value
						}
						newExample.Rows = append(newExample.Rows, Row{Cells: rowCells})
					}

					newOutline.Examples = append(newOutline.Examples, newExample)
				}

				// Add the outline to feature's ScenarioOutlines
				feature.ScenarioOutlines = append(feature.ScenarioOutlines, newOutline)

			} else {

				var newScenario Scenario
				newScenario.Name = scenario.Name

				// Collect Tags if any
				for _, tag := range scenario.Tags {
					newScenario.Tags = append(newScenario.Tags, tag.Name)
				}

				// Collect Steps with their appropriate prefixes
				for _, step := range scenario.Steps {
					newScenario.Steps = append(newScenario.Steps, Step{Text: step.Text, Prefix: step.Keyword})
				}

				// Add to feature's scenarios
				feature.Scenarios = append(feature.Scenarios, newScenario)
			}

		}
	}

	return feature
}

// Helper function to identify common steps across scenarios
func findCommonSteps(scenarios []Scenario, scenarioOutlines []ScenarioOutline) ([]string, []Scenario) {
	stepCount := make(map[string]int)

	// Count occurrences of each step in scenarios
	for _, scenario := range scenarios {
		for _, step := range scenario.Steps {
			if step.Prefix != "When " { // Don't include steps that have when prefix in the background
				stepCount[step.Text]++ // Count based on the step text
			}
		}
	}

	// Count occurrences of each step in scenario Outlines
	for _, scenarioOutline := range scenarioOutlines {
		for _, step := range scenarioOutline.Steps {
			if step.Prefix != "When " { // Don't include steps that have when prefix in the background
				stepCount[step.Text]++ // Count based on the step text
			}
		}
	}

	// Identify steps that are common to all scenarios
	var commonSteps []string
	for step, count := range stepCount {
		if count == len(scenarios) || count == len(scenarios)+len(scenarioOutlines) {
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
			if !contains(commonSteps, step.Text) {
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

	// Identify and append common background steps
	commonSteps, newScenarios := findCommonSteps(feature.Scenarios, feature.ScenarioOutlines)

	// Append common steps to the existing background
	if len(feature.Background) > 0 {
		for _, step := range commonSteps {
			if !contains(feature.Background, step) && !strings.HasPrefix(step, "When") {
				// Append new unique steps to the existing background, exclude "When" steps
				feature.Background = append(feature.Background, step)
			}
		}
	}

	// If no existing background, create one with common steps
	if len(feature.Background) > 0 {
		builder.WriteString("Background:\n")
		for _, step := range feature.Background {
			builder.WriteString("  Given " + step + "\n")
		}
		builder.WriteString("\n")
	}

	// Include the existing scenario outlines if they exist
	for _, outline := range feature.ScenarioOutlines {
		builder.WriteString("Scenario Outline: " + outline.Name + "\n")

		// Include Steps from Scenario Outline, excluding common steps
		for _, step := range outline.Steps {
			if !contains(commonSteps, step.Text) {
				builder.WriteString("  " + step.Prefix + " " + step.Text + "\n")
			}
		}

		// Include Example Table
		builder.WriteString("\n  Examples:\n")
		for _, example := range outline.Examples {
			for _, row := range example.Rows {
				builder.WriteString("    | ")
				builder.WriteString(strings.Join(row.Cells, " | "))
				builder.WriteString(" |\n")
			}
		}

		builder.WriteString("\n")
	}

	// Include new scenarios
	for _, scenario := range newScenarios {
		builder.WriteString("Scenario: " + scenario.Name + "\n")

		// Include Tags if present
		for _, tag := range scenario.Tags {
			builder.WriteString("  @" + tag + "\n")
		}

		// Include Steps with their prefixes, excluding common steps
		for _, step := range scenario.Steps {
			if !contains(commonSteps, step.Text) {
				builder.WriteString("  " + step.Prefix + " " + step.Text + "\n")
			}
		}
		builder.WriteString("\n") // Separate scenarios with a newline
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
