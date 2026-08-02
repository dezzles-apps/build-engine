Feature: Repositories API

  Scenario: Get repositories by organisation
    Given url 'http://localhost:8087'
    And path 'api/v1/repositories/dezzles-apps'
    When method get
    Then status 200
    And match response[0].organisation == 'dezzles-apps'
    And match response[0].repository == 'test'
    And match response[1].organisation == 'dezzles-apps'
    And match response[1].repository == 'no-channel'

  Scenario: Get repository with channel
    Given url 'http://localhost:8087'
    And path 'api/v1/repositories/dezzles-apps/test'
    When method get
    Then status 200
    And match response.organisation == 'dezzles-apps'
    And match response.repository == 'test'
    And match response.channel == 'dezzles-apps'

  Scenario: Get repository without channel
    Given url 'http://localhost:8087'
    And path 'api/v1/repositories/dezzles-apps/no-channel'
    When method get
    Then status 200
    And match response.organisation == 'dezzles-apps'
    And match response.repository == 'no-channel'
    And match response.channel == 'dezzles-apps'

  