package crawler

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

type RateLimiter struct {
	mutex    sync.Mutex
	limiters map[string]*time.Ticker
	delay    time.Duration
}

func NewRateLimiter(delay time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*time.Ticker),
		delay:    delay,
	}
}

func (rl *RateLimiter) Wait(RawUrl string) error {
	domain, err := getDomain(RawUrl) // it will return original domain or hos  of an url....
	if err != nil {
		return fmt.Errorf("invalid url")
	}
	//get or create ticker for the domain....................
	rl.mutex.Lock()
	ticker, exists := rl.limiters[domain] //get the key --> cheack ticker with domain name

	if !exists {
		//create a new ticker for the domain with specific delay..
		ticker = time.NewTicker(rl.delay)
		rl.limiters[domain] = ticker

	}
	rl.mutex.Unlock()

	//wait for the tickers to tick
	<-ticker.C
	return nil

}

func getDomain(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
func (rl *RateLimiter) StopAll() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	for _, ticker := range rl.limiters {
		ticker.Stop()
	}
}
