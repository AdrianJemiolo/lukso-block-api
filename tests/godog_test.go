package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

var response *http.Response
var responseBody map[string]interface{}

func iSendAGETRequestTo(endpoint string) error {
	var err error
	response, err = http.Get("http://localhost:8080" + endpoint)
	if err != nil {
		return fmt.Errorf("failed to send GET request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	defer response.Body.Close()

	// Attempt to parse the response body as JSON
	responseBody = make(map[string]interface{})
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		// If the response is not JSON, ignore the error for non-JSON responses
		responseBody = map[string]interface{}{
			"raw": string(body),
		}
	}

	return nil
}

func theResponseCodeShouldBe(statusCode int) error {
	if response.StatusCode != statusCode {
		return fmt.Errorf("expected status code %d, got %d", statusCode, response.StatusCode)
	}
	return nil
}

func theResponseShouldContain(key string, value string) error {
	actualValue, exists := responseBody[key]
	if !exists {
		return fmt.Errorf("key %s not found in response", key)
	}

	if fmt.Sprintf("%v", actualValue) != value {
		return fmt.Errorf("expected %s to be %s, got %s", key, value, actualValue)
	}

	return nil
}

func theResponseShouldContainSwaggerUI() error {
	// Check if the raw response body contains expected Swagger UI content
	rawBody, ok := responseBody["raw"].(string)
	if !ok {
		return fmt.Errorf("raw response body not found or not a string")
	}

	if !strings.Contains(rawBody, "Swagger UI") {
		return fmt.Errorf("Swagger UI content not found in response")
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send a GET request to "([^"]*)"$`, iSendAGETRequestTo)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response should contain "([^"]*)" with value "([^"]*)"$`, theResponseShouldContain)
	ctx.Step(`^the response should contain Swagger UI$`, theResponseShouldContainSwaggerUI)
}

func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",
	}
	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	if status != 0 {
		m.Run()
	}
}
