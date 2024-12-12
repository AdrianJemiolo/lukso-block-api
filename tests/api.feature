Feature: LUKSO Block API

  Scenario: Check the /health endpoint
    When I send a GET request to "/health"
    Then the response code should be 200
    And the response should contain "status" with value "UP"

  Scenario: Fetch the latest block number
    When I send a GET request to "/block-number"
    Then the response code should be 200
    And the response should contain "blockNumber"

  Scenario: Access the Swagger documentation
    When I send a GET request to "/docs/"
    Then the response code should be 200
    And the response should contain Swagger UI

