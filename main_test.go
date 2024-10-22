package main

import (
	"fmt"
	"reflect"
	"testing"
)

// Test to verify the generation of a feature file from structured data.
func TestFeatureFileGeneration(t *testing.T) {
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
	//generate a feature file from content
	generatedFile := generateFeatureFile(feature)
	// Parse the input content.
	generatedFeature := parseFeatureFile(generatedFile)
	//generate control
	generatedControl := parseFeatureFile(expected)
	//compar to original feature struct
	if !reflect.DeepEqual(feature, generatedFeature) {
		fmt.Println(feature)
		fmt.Println("*********************")
		fmt.Println(generatedFeature)
		t.Errorf("Expected original feature content to equal itself after being parsed to a file and back but it did not")
	}
	//positive control
	if !reflect.DeepEqual(feature, generatedControl) {
		t.Errorf("positive control failed")
	}
	//negative control
	if !reflect.DeepEqual(generatedFeature, generatedControl) {
		t.Errorf("negative control failed")
	}
}

// Unit test to parse and then regenerate a feature file.
func TestParseAndGenerateFeatureFile(t *testing.T) {
	inputContent := `Feature: User Login

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

	// Parse the input content.
	feature := parseFeatureFile(inputContent)

	// Regenerate the feature file.
	generated := generateFeatureFile(feature)

	// parse after generated
	parseGenerated := parseFeatureFile(generated)

	if !reflect.DeepEqual(feature, parseGenerated) {
		t.Errorf("Expected structs to be equal, but they are not")
	}
}
