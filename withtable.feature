Feature: User Login

Scenario Outline: Login Scenario Outline
  Given the user has entered <username>
  Given the user has entered <password>
  When the user clicks the login button
  Then the user should see <result>

Examples:
  | username      | password     | result              |
  | valid_user    | valid_pass   | redirected to dashboard |
  | invalid_user  | invalid_pass | error message       |