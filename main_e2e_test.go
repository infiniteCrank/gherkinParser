package main

// import (
// 	"bytes"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/gorilla/mux"
// )

// func TestUploadFeature(t *testing.T) {
// 	// Sample feature data that will be uploaded
// 	featureData := `Feature: User Login

// Background:
//   Given the user has opened the login page

// Scenario: Successful login with valid credentials
//   Given the user has entered a valid username
//   Given the user has entered a valid password
//   When the user clicks the login button
//   Then the user should be redirected to the dashboard
//   And a welcome message should be displayed

// Scenario: Unsuccessful login with invalid credentials
//   Given the user has entered an invalid username
//   Given the user has entered an invalid password
//   When the user clicks the login button
//   Then the user should see an error message
// `

// 	// Create a buffer to simulate file input
// 	body := io.NopCloser(bytes.NewBufferString(featureData))

// 	// Create a new HTTP request
// 	req := httptest.NewRequest(http.MethodPost, "/upload", body)
// 	req.Header.Set("Content-Type", "multipart/form-data")

// 	// Create a ResponseRecorder to capture the response
// 	rr := httptest.NewRecorder()

// 	// Create a new router and set the upload handler
// 	router := mux.NewRouter()
// 	router.HandleFunc("/upload", uploadHandler).Methods("POST")

// 	// Serve the HTTP request to the response recorder
// 	router.ServeHTTP(rr, req)

// 	// Check the response code
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
// 	}

// 	// Check the body of the response
// 	expectedResponse := "Feature file parsed successfully!"
// 	if strings.TrimSpace(rr.Body.String()) != expectedResponse {
// 		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
// 	}

// 	// Additionally, you can check if the file was parsed correctly by regenerating the expected output
// 	expectedGenerated := generateFeatureFile(parseFeatureFile(featureData))
// 	if strings.TrimSpace(expectedGenerated) != strings.TrimSpace(featureData) {
// 		t.Errorf("Expected generated feature file to match original input:\nExpected:\n%s\nGot:\n%s", featureData, expectedGenerated)
// 	}
// }
