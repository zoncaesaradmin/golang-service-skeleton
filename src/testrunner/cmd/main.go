package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"katharos/testrunner/internal/config"
	"katharos/testrunner/internal/orchestrator"
	"katharos/testrunner/internal/testdata"
	"katharos/testrunner/internal/types"
	"katharos/testrunner/internal/validation"
	"gopkg.in/yaml.v2"
)

func main() {
	// Command line flags
	var (
		configFile = flag.String("config", "config.yaml", "Configuration file path")
		scenario   = flag.String("scenario", "", "Specific scenario to run (leave empty for all)")
		output     = flag.String("output", "console", "Output format: console, json, junit")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
		generate   = flag.Bool("generate", false, "Generate sample test data and config")
	)
	flag.Parse()

	// Setup logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Handle data generation
	if *generate {
		if err := generateSampleData(); err != nil {
			log.Fatalf("Failed to generate sample data: %v", err)
		}
		log.Println("Sample data and config generated successfully")
		return
	}

	log.Printf("Starting Katharos Test Runner...")
	log.Printf("Config file: %s", *configFile)
	log.Printf("Output format: %s", *output)

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create orchestrator
	orch, err := orchestrator.NewOrchestrator(cfg)
	if err != nil {
		log.Fatalf("Failed to create orchestrator: %v", err)
	}

	// Load test scenarios
	loader := testdata.NewLoader(cfg.Testdata.ScenariosPath)
	scenarios, err := loader.LoadAllScenarios()
	if err != nil {
		log.Fatalf("Failed to load test scenarios: %v", err)
	}

	// Filter scenarios if specific scenario requested
	if *scenario != "" {
		filteredScenarios := make([]testdata.TestScenario, 0)
		for _, s := range scenarios {
			if s.Name == *scenario {
				filteredScenarios = append(filteredScenarios, s)
				break
			}
		}
		if len(filteredScenarios) == 0 {
			log.Fatalf("Scenario '%s' not found", *scenario)
		}
		scenarios = filteredScenarios
	}

	log.Printf("Loaded %d test scenario(s)", len(scenarios))

	// Execute scenarios
	results := make([]types.TestResult, 0)
	for _, scenario := range scenarios {
		log.Printf("Executing scenario: %s", scenario.Name)
		result, err := orch.ExecuteScenario(scenario)
		if err != nil {
			log.Printf("Scenario '%s' execution failed: %v", scenario.Name, err)
			result = types.TestResult{
				ScenarioName: scenario.Name,
				Success:      false,
				Error:        err.Error(),
				Duration:     0,
			}
		}
		results = append(results, result)
	}

	// Generate report
	reporter := validation.NewReporter(*output)
	report := validation.TestReport{
		Timestamp: time.Now(),
		Results:   results,
	}

	if err := reporter.GenerateReport(report); err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Calculate success rate
	successful := 0
	for _, result := range results {
		if result.Success {
			successful++
		}
	}

	successRate := float64(successful) / float64(len(results)) * 100
	log.Printf("Test execution completed: %d/%d scenarios passed (%.1f%%)", 
		successful, len(results), successRate)

	if successful == len(results) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// generateSampleData creates sample configuration and test data files
func generateSampleData() error {
	// Create directories
	dirs := []string{"testdata/scenarios", "testdata/fixtures"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate sample config
	sampleConfig := config.Config{
		Component: config.ComponentConfig{
			BinaryPath: "../component/bin/component",
			Port:       8080,
			Timeout:    30 * time.Second,
		},
		MessageBus: config.MessageBusConfig{
			Type: "local", // Use local for development
		},
		Testdata: config.TestdataConfig{
			ScenariosPath: "testdata/scenarios",
			FixturesPath:  "testdata/fixtures",
		},
		Validation: config.ValidationConfig{
			Timeout:      60 * time.Second,
			MaxRetries:   3,
			RetryDelay:   1 * time.Second,
		},
	}

	configData, err := yaml.Marshal(sampleConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile("config.yaml", configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Generate sample scenario
	sampleScenario := testdata.TestScenario{
		Name:        "user_workflow",
		Description: "Test complete user creation and retrieval workflow",
		Input: map[string]interface{}{
			"user": map[string]interface{}{
				"id":   "test-user-123",
				"name": "Test User",
				"email": "test@example.com",
			},
		},
		ExpectedOutput: map[string]interface{}{
			"status": "success",
			"user": map[string]interface{}{
				"id":   "test-user-123",
				"name": "Test User",
				"email": "test@example.com",
			},
		},
		Timeout: 30 * time.Second,
	}

	scenarioData, err := yaml.Marshal(sampleScenario)
	if err != nil {
		return fmt.Errorf("failed to marshal scenario: %w", err)
	}

	scenarioPath := filepath.Join("testdata", "scenarios", "user_workflow.yaml")
	if err := os.WriteFile(scenarioPath, scenarioData, 0644); err != nil {
		return fmt.Errorf("failed to write scenario file: %w", err)
	}

	// Generate sample fixture data
	sampleFixture := map[string]interface{}{
		"users": []map[string]interface{}{
			{
				"id":   "user-1",
				"name": "John Doe",
				"email": "john@example.com",
			},
			{
				"id":   "user-2", 
				"name": "Jane Smith",
				"email": "jane@example.com",
			},
		},
	}

	fixtureData, err := json.MarshalIndent(sampleFixture, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal fixture: %w", err)
	}

	fixturePath := filepath.Join("testdata", "fixtures", "sample_users.json")
	if err := os.WriteFile(fixturePath, fixtureData, 0644); err != nil {
		return fmt.Errorf("failed to write fixture file: %w", err)
	}

	return nil
}
