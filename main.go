package main

import (
	"crawl-nova/crawler"
	"fmt"
	"sync"
)

func main() {
	// Define a list of domains (some are intentionally duplicated for testing deduplication)
	domains := []string{
		"https://google.com",
		"https://amazon.com",
		"https://facebook.com",
		"https://invalid.domain",
		"https://google.com",   // duplicate , just to chech the safeset logic
		"https://facebook.com", // duplicate
		"https://amazon.com",   // duplicate
	}

	// Create channels:
	jobs := make(chan string)    // Channel for sending jobs (URLs to crawl) to workers
	results := make(chan string) // Channel for receiving results from workers

	var wg sync.WaitGroup // WaitGroup to track when all workers are done

	// Initialize a SafeSet for deduplication of URLs
	deduper := crawler.NewSafeSet()

	// Start a goroutine to send unique (non-duplicate) domains to the jobs channel
	go func() {
		for _, domain := range domains {
			if deduper.Add(domain) {
				jobs <- domain // Send domain to job queue if not seen before
			} else {
				fmt.Println("✘ Duplicate skipped:", domain) // Log duplicates that are skipped
			}
		}
		//stop jobs..  logic will implement latter
		//close(jobs) // Important: Close the jobs channel to signal no more incoming tasks
	}()

	// Define how many worker goroutines to run concurrently
	workcount := 5

	// Launch the worker pool
	for w := 0; w < workcount; w++ {
		wg.Add(1)                                              // Increment WaitGroup for each worker
		go crawler.StartWorker(w, jobs, results, &wg, deduper) // Start a worker with an ID
	}

	// Start a goroutine to wait for all workers to finish and then close results channel
	go func() {
		wg.Wait()      // Block until all workers call Done()
		close(results) // Important: Close results channel so main loop below can finish
	}()

	// Continuously read from results channel and print the crawl results
	for res := range results {
		fmt.Println(res)
	}
}
