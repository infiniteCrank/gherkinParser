Feature: User Login

  Background:
    Given the user has opened the login page

  Scenario: Successful login with valid credentials
    Given the user has entered a valid username
    Given the user has entered a valid password
    When the user clicks the login button
    Then the user should be redirected to the dashboard

  Scenario: Unsuccessful login with invalid credentials
    Given the user has entered an invalid username
    Given the user has entered an invalid password
    When the user clicks the login button
    Then the user should see an error message




