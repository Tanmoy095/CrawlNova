package crawler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func StartWorker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() //tell main the worker is done when exists
	client := &http.Client{}
	for url := range jobs {

		//create a context that auto cancel after 2sec....

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		//create a	 GET req with using that context

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
