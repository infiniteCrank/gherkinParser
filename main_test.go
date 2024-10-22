package main

import (
	"strings"
	"testing"
)

// Test to verify the generation of a feature file from structured data
func TestFeatureFileGeneration(t *testing.T) {
	feature := Feature{
		Name:       "Sample Feature",
		Background: []string{"some background step"},
		Scenarios: []Scenario{
			{
				Name:  "Sample Scenario",
				Steps: []string{"step 1", "step 2"},
				Tags:  []string{"tag1", "tag2"},
			},
		},
	}

	expected := `Feature: Sample Feature

Background:
  Given some background step

Scenario: Sample Scenario
  @tag1
  @tag2
  When step 1
  When step 2

`

	generated := generateFeatureFile(feature)
	if strings.TrimSpace(generated) != strings.TrimSpace(expected) {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, generated)
	}
}

// Test to parse and then regenerate a feature file
func TestParseAndGenerateFeatureFile(t *testing.T) {
	inputContent := `
Feature: Sample Feature

Background:
  Given some background step

Scenario: Sample Scenario
  @tag1
  @tag2
  When step 1
  When step 2
`

	// Parse the input content
	feature := parseFeatureFile(inputContent)
	// Regenerate the feature file
	generated := generateFeatureFile(feature)

	// Ensure that the generated content matches the original input
	if strings.TrimSpace(generated) != strings.TrimSpace(inputContent) {
		t.Errorf("Expected generated file to match the original input:\nExpected:\n%s\nGot:\n%s", inputContent, generated)
	}
}
