package crawler

import (
	"context"
	"crawl-nova/types"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func StartWorker(id int, jobs chan types.CrawlJob, results chan<- string, wg *sync.WaitGroup, deduper *SafeSet, MaxDepth int, rateLimiter *RateLimiter) {
	defer wg.Done() //tell main the worker is done when exists
	client := &http.Client{}
	for job := range jobs {
		// Step 1: Respect rate limit for this domain before crawling
		if err := rateLimiter.Wait(job.URL); err != nil {
			results <- fmt.Sprintf("[Worker %d] RateLimiter error: %v", id, err)
			continue
		}

		//create a context that auto cancel after 2sec timeout to avoid hanging requests

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		//create a	 GET req with using that context

		req, err := http.NewRequestWithContext(ctx, "GET", job.URL, nil)
		if err != nil {
			results <- fmt.Sprintf("Worker %d failed to create request : %v", id, err)
			cancel() //must call to avoid context memory leaks...
			continue //skip the job and try next

		}

		//measure request duration..........

		start := time.Now()

		resp, err := client.Do(req) //  Executes request with timeout built-in
		if err != nil {
			results <- fmt.Sprintf("[Worker %d] Domain: %s - Timeout/Error: %v", id, job.URL, err)
			cancel() //must call to avoid context memory leaks...
			continue
		}

		//must extract links before closing the resp.body

		links, err := ExtractLinks(resp.Body, job.URL)
		resp.Body.Close() //  Always close response body

		if err != nil {
			results <- fmt.Sprintf("[Worker %d]  %s - Parse Error: %v", id, job.URL, err)

		} else if job.Depth < MaxDepth {
			// Step 6: Process discovered links
			for _, link := range links {
				// Only queue unseen URLs
				if deduper.Add(link) {
					// 👇 NOTE: No rateLimiter here — because this is just enqueuing, not crawling yet
					jobs <- types.CrawlJob{
						URL:   link,
						Depth: job.Depth + 1,
					}
					results <- fmt.Sprintf("[Worker %d] ➕ Discovered: %s", id, link)
				}
			}
		}
		//Log how long the crawl took
		duration := time.Since(start)
		results <- fmt.Sprintf("[Worker %d]  Crawled: %s in %v", id, job.URL, duration)
	}
}
