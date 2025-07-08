package crawler

import "sync"

// SafeSet is a concurrency-safe set-like structure for tracking visited URLs.
// Internally, it uses a map to store seen URLs and a mutex for protecting access.
type SafeSet struct {
	mu   sync.Mutex      // Mutex to ensure safe concurrent access to the 'seen' map
	seen map[string]bool // Map to record which URLs have been processed
}

// NewSafeSet initializes and returns a pointer to an empty SafeSet instance.
func NewSafeSet() *SafeSet {
	return &SafeSet{seen: make(map[string]bool)}
}

// Add attempts to insert the given URL into the SafeSet.
// It returns true if the URL was not seen before and is now added.
// It returns false if the URL was already present (i.e., duplicate).
func (s *SafeSet) Add(url string) bool {
	s.mu.Lock()         // Lock the mutex before accessing/modifying the map
	defer s.mu.Unlock() // Ensure the mutex is unlocked after this function exits

	if s.seen[url] {
		return false // URL already visited, so skip it
	}
	s.seen[url] = true // Mark the URL as seen
	return true
}
