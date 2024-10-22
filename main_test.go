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

	expected := `Feature: User Login

Background:
  Given the user has opened the login page
  Given the user clicks the login button

Scenario: Successful login with valid credentials
  Given the user has entered a valid username
  Given the user has entered a valid password
  Then the user should be redirected to the dashboard
  Then a welcome message should be displayed

Scenario: Unsuccessful login with invalid credentials
  Given the user has entered an invalid username
  Given the user has entered an invalid password
  Then the user should see an error message`
	// Generate a feature file from the structured data
	generatedFile := generateFeatureFile(feature)
	generateFeature := parseFeatureFile(generatedFile)
	control := parseFeatureFile(expected)

	if !reflect.DeepEqual(generateFeature, control) {
		t.Errorf("Expected %v, got %v", control, generateFeature)
	}

}
