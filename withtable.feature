Feature: User Login

Background:
  Given the user has opened the login page

Scenario: Successful login with valid credentials
  Given I ate lunch
  Given the user has entered a valid username
  Given the user has entered a valid password
  When the user clicks the login button
  Then the user should be redirected to the dashboard
  And a welcome message should be displayed

Scenario: Unsuccessful login with invalid credentials
  Given I ate lunch
  Given the user has entered an invalid username
  Given the user has entered an invalid password
  When the user clicks the login button
  Then the user should see an error message

Scenario: Unsuccessful login with invalid credentials
  Given I ate lunch
  Given the user has entered an invalid username
  Given the user has entered an valid password
  When the user clicks the login button
  Then the user should see an error message

Scenario: Unsuccessful login with invalid credentials
  Given I ate lunch
  Given the user has entered an valid username
  Given the user has entered an invalid password
  When the user clicks the login button
  Then the user should see an error message

Scenario Outline: Login attempts with different usernames and passwords
  Given I ate lunch
  Given the user has entered <username>
  And the user has entered <password>
  When the user clicks the login button
  Then the user should see <result>

  Examples:
    | username      | password     | result              |
    | valid_user    | valid_pass   | redirected to dashboard |
    | invalid_user  | invalid_pass | error message       |