package main

import (
	"strings"
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

	generated := generateFeatureFile(feature)
	if strings.TrimSpace(generated) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, generated)
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

	// Ensure that the generated content matches the original input.
	if strings.TrimSpace(generated) != strings.TrimSpace(inputContent) {
		t.Errorf("Expected generated file to match the original input:\nExpected:\n%s\nGot:\n%s", inputContent, generated)
	}
}
