package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {

	defer wg.Done() //tell main the worker is done when exists
	client := &http.Client{}
	for url := range jobs {

		//create a context that auto cancel after 2sec....

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		//create a   GET req with using that context

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			results <- fmt.Sprintf("Worker %d failed to create request : %v", id, err)
			cancel() //must call to avoid context memory leaks...
			continue //skip the job and try next

		}

		//measure request duration..........

		// 3. Measure request duration
		start := time.Now()
		resp, err := client.Do(req) //  Executes request with timeout built-in
		if err != nil {
			results <- fmt.Sprintf("[Worker %d] Domain: %s - Timeout/Error: %v", id, url, err)
			cancel() //must call to avoid context memory leaks...
			continue
		}
		resp.Body.Close() //  Always close response body

		// 4. Send result to results channel
		duration := time.Since(start)
		results <- fmt.Sprintf("[Worker %d] Domain: %s - Time: %v", id, url, duration)
	}

}

func main() {
	// 📝 List of domains (jobs to perform)
	domains := []string{
		"https://google.com",
		"https://amazon.com",
		"https://facebook.com",
		"https://invalid.domain", // ❌ Will cause error (bad domain)
	}

	jobs := make(chan string)    // channel for sending jobs(domains)
	results := make(chan string) //channel for receiving results from worksers.......

	var wg sync.WaitGroup

	//sends all domains into the jobs channel

	go func() {
		for _, domain := range domains {
			jobs <- domain
		}

		close(jobs) //tell worker no more job after this ...............................

	}()

	//start workers............

	workcount := 5 //we will use 5 workers

	for w := 0; w < workcount; w++ {
		wg.Add(1)

		go worker(w, jobs, results, &wg) //start worker go routine.........

	}
	//waits for all workers to finish ,then closes results channel
	go func() {
		wg.Wait()
		close(results)

	}()
	for res := range results {
		fmt.Println(res)
	}
}
