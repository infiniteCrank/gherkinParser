package main

import (
	"fmt"
	"reflect"
	"testing"
)

// Test to verify the generation of a feature file from structured data
// including appending common steps to an existing background
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
	fmt.Println("*************** TEST ONE Generated file ********************")
	fmt.Println(generatedFile)
	generateFeature := parseFeatureFile(generatedFile)
	fmt.Println("*************** TEST ONE generated struct ********************")
	fmt.Println(fmt.Printf("%+v\n", generateFeature))

	control := parseFeatureFile(expected)

	if !reflect.DeepEqual(generateFeature, control) {
		t.Errorf("Expected %v, got %v", control, generateFeature)
	} else {
		fmt.Println("*************** TEST ONE PASS! ********************")
	}

}

// Test to verify the generation when no background exists
func TestFeatureFileGenerationNoBackground(t *testing.T) {
	feature := Feature{
		Name: "User Login",
		Scenarios: []Scenario{
			{
				Name:  "Login attempts",
				Steps: []string{"the user enters credentials", "the user clicks login", "the user sees the dashboard"},
				Tags:  []string{},
			},
		},
	}

	expected := `Feature: User Login

Background:
  Given the user enters credentials
  Given the user clicks login
  Given the user sees the dashboard

Scenario: Login attempts
`

	generatedFile := generateFeatureFile(feature)
	fmt.Println("*************** TEST TWO Generated file ********************")
	fmt.Println(generatedFile)
	generateFeature := parseFeatureFile(generatedFile)
	fmt.Println("*************** TEST TWO generated struct ********************")
	fmt.Println(fmt.Printf("%+v\n", generateFeature))
	control := parseFeatureFile(expected)
	if !reflect.DeepEqual(generateFeature, control) {
		t.Errorf("Expected %v, got %v", control, generateFeature)
	} else {
		fmt.Println("*************** TEST TWO PASS! ********************")
	}
}
