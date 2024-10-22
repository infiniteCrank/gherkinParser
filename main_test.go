package main

import (
	"fmt"
	"testing"
)

// Test to verify the generation of a feature file from structured data including common steps in background.
func TestFeatureFileGenerationWithBackground(t *testing.T) {
	feature := Feature{
		Name:       "User Login",
		Background: []string{"the user has opened the login page"},
		Scenarios: []Scenario{
			{
				Name:  "Successful login with valid credentials",
				Steps: []string{"the user has entered a valid username", "the user has entered a valid password", "the user clicks the login button", "the user should be redirected to the dashboard", "a welcome message should be displayed"},
				Tags:  []string{},
			},
			{
				Name:  "Unsuccessful login with invalid credentials",
				Steps: []string{"the user has entered an invalid username", "the user has entered an invalid password", "the user clicks the login button", "the user should see an error message"},
				Tags:  []string{},
			},
		},
	}

	// Generate a feature file from the structured data
	generatedFile := generateFeatureFile(feature)
	generateFeature := parseFeatureFile(generatedFile)

	fmt.Println(feature.Scenarios[0].Name)
	fmt.Println(feature.Scenarios[1].Name)
	fmt.Println(feature.Background[0])
	fmt.Println("***********************************")
	fmt.Println(generateFeature.Scenarios[0].Name)
	fmt.Println(generateFeature.Scenarios[1].Name)
	fmt.Println(generateFeature.Background[0])
	if feature.Name != generateFeature.Name {
		fmt.Println("did not pass feature name:" + generateFeature.Name)
	}

	fmt.Println("***********************************")
	fmt.Println(generatedFile)

}

// // Unit test to ensure proper identification of common steps for backgrounds.
// func TestFindCommonSteps(t *testing.T) {
// 	scenarios := []Scenario{
// 		{
// 			Name:  "Successful login with valid credentials",
// 			Steps: []string{"the user has entered a valid username", "the user has entered a valid password", "the user clicks the login button", "the user should be redirected to the dashboard"},
// 			Tags:  []string{},
// 		},
// 		{
// 			Name:  "Unsuccessful login with invalid credentials",
// 			Steps: []string{"the user has entered an invalid username", "the user has entered an invalid password", "the user clicks the login button", "the user should see an error message"},
// 			Tags:  []string{},
// 		},
// 	}

// 	expectedCommonSteps := []string{}
// 	expectedNewScenarios := []Scenario{
// 		{
// 			Name:  "Successful login with valid credentials",
// 			Steps: []string{"the user has entered a valid username", "the user has entered a valid password", "the user clicks the login button", "the user should be redirected to the dashboard"},
// 			Tags:  []string{},
// 		},
// 		{
// 			Name:  "Unsuccessful login with invalid credentials",
// 			Steps: []string{"the user has entered an invalid username", "the user has entered an invalid password", "the user clicks the login button", "the user should see an error message"},
// 			Tags:  []string{},
// 		},
// 	}

// 	commonSteps, newScenarios := findCommonSteps(scenarios)

// 	if !reflect.DeepEqual(commonSteps, expectedCommonSteps) {
// 		t.Errorf("Expected common steps %v, got %v", expectedCommonSteps, commonSteps)
// 	}

// 	if !reflect.DeepEqual(newScenarios, expectedNewScenarios) {
// 		t.Errorf("New scenarios do not match expected. Got: %+v, Want: %+v", newScenarios, expectedNewScenarios)
// 	}
// }
