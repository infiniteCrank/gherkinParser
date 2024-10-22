package main

import (
	"reflect"
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

	// Expected output feature file
	expected := `Feature: User Login

Background:
  Given the user has opened the login page

Scenario: Successful login with valid credentials
  Given the user has entered a valid username
  Given the user has entered a valid password
  When the user clicks the login button
  Then the user should be redirected to the dashboard
  And a welcome message should be displayed

Scenario: Unsuccessful login with invalid credentials
  Given the user has entered an invalid username
  Given the user has entered an invalid password
  When the user clicks the login button
  Then the user should see an error message
`

	// Generate a feature file from the structured data
	generatedFile := generateFeatureFile(feature)
	//generate a struct from the generated file
	generateTestFeature := parseFeatureFile(generatedFile)
	//generate expected feature
	generateControl := parseFeatureFile(expected)

	if generateControl.Background[0] != generateTestFeature.Background[0] {
		t.Errorf("Expected background %v, got %v", generateControl.Background[0], generateTestFeature.Background[0])
	}

}

// Unit test to ensure proper identification of common steps for backgrounds.
func TestFindCommonSteps(t *testing.T) {
	scenarios := []Scenario{
		{
			Name:  "Successful login with valid credentials",
			Steps: []string{"the user has entered a valid username", "the user has entered a valid password", "the user clicks the login button", "the user should be redirected to the dashboard"},
			Tags:  []string{},
		},
		{
			Name:  "Unsuccessful login with invalid credentials",
			Steps: []string{"the user has entered an invalid username", "the user has entered an invalid password", "the user clicks the login button", "the user should see an error message"},
			Tags:  []string{},
		},
	}

	expectedCommonSteps := []string{}
	expectedNewScenarios := []Scenario{
		{
			Name:  "Successful login with valid credentials",
			Steps: []string{"the user has entered a valid username", "the user has entered a valid password", "the user clicks the login button", "the user should be redirected to the dashboard"},
			Tags:  []string{},
		},
		{
			Name:  "Unsuccessful login with invalid credentials",
			Steps: []string{"the user has entered an invalid username", "the user has entered an invalid password", "the user clicks the login button", "the user should see an error message"},
			Tags:  []string{},
		},
	}

	commonSteps, newScenarios := findCommonSteps(scenarios)

	if !reflect.DeepEqual(commonSteps, expectedCommonSteps) {
		t.Errorf("Expected common steps %v, got %v", expectedCommonSteps, commonSteps)
	}

	if !reflect.DeepEqual(newScenarios, expectedNewScenarios) {
		t.Errorf("New scenarios do not match expected. Got: %+v, Want: %+v", newScenarios, expectedNewScenarios)
	}
}
