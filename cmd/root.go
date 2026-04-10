package cmd

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
	rootCmd     = &cobra.Command{
		Use:   "stress-test",
		Short: "A load testing CLI tool for web services",
		Long:  "stress-test is a command-line tool for performing load tests on web services.",
		RunE:  runLoadTest,
	}
)

type TestResult struct {
	TotalRequests   int64
	SuccessCount    int64
	StatusCodeCount map[int]int64
	TotalTime       time.Duration
	StartTime       time.Time
	EndTime         time.Time
	mu              sync.Mutex
}

type LoadTester struct {
	url         string
	requests    int
	concurrency int
	client      *http.Client
	result      *TestResult
}

func NewLoadTester(url string, requests, concurrency int) *LoadTester {
	return &LoadTester{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		result: &TestResult{
			StatusCodeCount: make(map[int]int64),
		},
	}
}

func (lt *LoadTester) Run() error {
	lt.result.StartTime = time.Now()

	requestChan := make(chan int, lt.concurrency)
	var wg sync.WaitGroup

	for i := 0; i < lt.concurrency; i++ {
		wg.Add(1)
		go lt.worker(requestChan, &wg)
	}

	go func() {
		for i := 0; i < lt.requests; i++ {
			requestChan <- i
		}
		close(requestChan)
	}()

	wg.Wait()

	lt.result.EndTime = time.Now()
	lt.result.TotalTime = lt.result.EndTime.Sub(lt.result.StartTime)
	lt.result.TotalRequests = int64(lt.requests)

	return nil
}

func (lt *LoadTester) worker(requestChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for range requestChan {
		lt.makeRequest()
	}
}

func (lt *LoadTester) makeRequest() {
	resp, err := lt.client.Get(lt.url)
	if err != nil {
		lt.result.mu.Lock()
		lt.result.StatusCodeCount[0]++
		lt.result.mu.Unlock()
		return
	}
	defer resp.Body.Close()

	io.ReadAll(resp.Body)

	statusCode := resp.StatusCode
	if statusCode == 200 {
		atomic.AddInt64(&lt.result.SuccessCount, 1)
	}

	lt.result.mu.Lock()
	lt.result.StatusCodeCount[statusCode]++
	lt.result.mu.Unlock()
}

func (lt *LoadTester) PrintReport() {
	fmt.Println("\n=== Load Test Report ===")
	fmt.Printf("Total Time: %v\n", lt.result.TotalTime)
	fmt.Printf("Total Requests: %d\n", lt.result.TotalRequests)
	fmt.Printf("Successful (HTTP 200): %d\n", lt.result.SuccessCount)
	fmt.Printf("Success Rate: %.2f%%\n", float64(lt.result.SuccessCount)*100/float64(lt.result.TotalRequests))

	if lt.result.TotalTime.Seconds() > 0 {
		rps := float64(lt.result.TotalRequests) / lt.result.TotalTime.Seconds()
		fmt.Printf("Requests per second: %.2f\n", rps)
	}

	fmt.Println("\nStatus Code Distribution:")
	for statusCode, count := range lt.result.StatusCodeCount {
		if statusCode == 0 {
			fmt.Printf("  Errors (Connection Failed): %d\n", count)
		} else {
			fmt.Printf("  HTTP %d: %d\n", statusCode, count)
		}
	}
	fmt.Println("========================")
}

func runLoadTest(cmd *cobra.Command, args []string) error {
	if url == "" {
		return fmt.Errorf("error: --url parameter is required")
	}

	if requests <= 0 {
		return fmt.Errorf("error: --requests must be greater than 0")
	}

	if concurrency <= 0 {
		return fmt.Errorf("error: --concurrency must be greater than 0")
	}

	fmt.Printf("Starting load test...\n")
	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Requests: %d\n", requests)
	fmt.Printf("Concurrency: %d\n\n", concurrency)

	tester := NewLoadTester(url, requests, concurrency)
	if err := tester.Run(); err != nil {
		return fmt.Errorf("error running test: %w", err)
	}

	tester.PrintReport()
	return nil
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the service to test (required)")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Total number of requests to be made (required)")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of simultaneous calls (default: 1)")

	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
}

func Execute() error {
	return rootCmd.Execute()
}
