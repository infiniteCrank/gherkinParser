Feature: User Login

Background:
  Given the user has entered a valid username
  Given the user has entered a valid password

Scenario: Successful login with valid credentials
  When the user clicks the login button
  Then the user should be redirected to the dashboard

Scenario: Unsuccessful login with invalid credentials
  When the user clicks the login button
  Then the user should see an error message
