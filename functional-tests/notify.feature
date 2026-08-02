Feature: Notification tests

Scenario: Sending a notification
  Given url 'http://localhost:3000'
  And path '/api/v1/notify'
  And request { "organisation": "dezzles-apps", "repository": "test", "buildNumber": 23, "message": "Hello there" }
  When method post
  Then status 200
