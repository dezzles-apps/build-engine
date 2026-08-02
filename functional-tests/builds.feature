Feature: Builds API

  Scenario: Get detailed build information
    Given url 'http://localhost:8087'
    And path 'api/v1/builds/dezzles-apps/test/23'
    When method get
    Then status 200
    And match response.organisation == 'dezzles-apps'
    And match response.repository == 'test'
    And match response.buildNumber == 23
    And match response.startTime == '#string'
    And match response.discordThreadId == 23456
