package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"katharos/testrunner/internal/client"
	"katharos/testrunner/internal/tests"
)

func main() {
	// Command line flags
	var (
		serviceURL = flag.String("url", "http://localhost:8080", "Service URL to test")
		mode       = flag.String("mode", "integration", "Test mode: integration, performance, or all")
		timeout    = flag.Duration("timeout", 5*time.Minute, "Test timeout duration")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Setup logging
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	log.Printf("Starting test runner...")
	log.Printf("Service URL: %s", *serviceURL)
	log.Printf("Test Mode: %s", *mode)
	log.Printf("Timeout: %v", *timeout)

	// Create test suite
	testSuite := tests.NewTestSuite(*serviceURL)

	// Set timeout
	done := make(chan bool, 1)
	var success bool

	// Run tests in a goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Test runner panicked: %v", r)
				success = false
			}
			done <- true
		}()

		switch *mode {
		case "integration":
			log.Println("Running integration tests...")
			testSuite.RunAllTests()
			success = true
		case "performance":
			log.Println("Running performance tests...")
			runPerformanceTests(testSuite)
			success = true
		case "all":
			log.Println("Running all tests...")
			testSuite.RunAllTests()
			runPerformanceTests(testSuite)
			success = true
		default:
			log.Printf("Unknown test mode: %s", *mode)
			printUsage()
			success = false
		}
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		if success {
			log.Println("Test runner completed successfully")
			os.Exit(0)
		} else {
			log.Println("Test runner completed with errors")
			os.Exit(1)
		}
	case <-time.After(*timeout):
		log.Printf("Test runner timed out after %v", *timeout)
		os.Exit(1)
	}
}

// runPerformanceTests runs performance-specific tests
func runPerformanceTests(testSuite *tests.TestSuite) {
	log.Println("Starting performance tests...")

	// Run load tests
	runLoadTest(testSuite, 10, "Light Load Test")
	runLoadTest(testSuite, 50, "Medium Load Test")
	runLoadTest(testSuite, 100, "Heavy Load Test")

	log.Println("Performance tests completed")
}

// runLoadTest runs a load test with specified number of concurrent requests
func runLoadTest(_ *tests.TestSuite, concurrency int, name string) {
	log.Printf("Running %s (concurrency: %d)...", name, concurrency)

	start := time.Now()

	// Create channels for coordination
	done := make(chan bool, concurrency)
	errors := make(chan error, concurrency)

	// Launch concurrent goroutines
	for i := 0; i < concurrency; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Perform a simple health check
			client := client.NewClient("http://localhost:8080")
			_, err := client.HealthCheck()
			if err != nil {
				errors <- fmt.Errorf("goroutine %d failed: %w", id, err)
				return
			}
		}(i)
	}

	// Wait for all goroutines to complete
	completed := 0
	errorCount := 0

	for completed < concurrency {
		select {
		case <-done:
			completed++
		case err := <-errors:
			errorCount++
			log.Printf("Load test error: %v", err)
		case <-time.After(30 * time.Second):
			log.Printf("Load test timed out, completed: %d/%d", completed, concurrency)
			return
		}
	}

	duration := time.Since(start)
	successRate := float64(concurrency-errorCount) / float64(concurrency) * 100

	log.Printf("%s completed: %d requests in %v (%.1f%% success rate)",
		name, concurrency, duration, successRate)
}

// printUsage prints usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nTest Modes:\n")
	fmt.Fprintf(os.Stderr, "  integration  - Run integration tests (default)\n")
	fmt.Fprintf(os.Stderr, "  performance  - Run performance tests\n")
	fmt.Fprintf(os.Stderr, "  all          - Run all tests\n")
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s                                    # Run integration tests on localhost:8080\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -url http://api.example.com        # Run tests on remote service\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -mode performance                  # Run performance tests\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -mode all -verbose                 # Run all tests with verbose output\n", os.Args[0])
}
