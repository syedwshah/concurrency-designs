package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Design 1 (Sync):
// func main() {

// 	r := http.NewServeMux()
// 	var wg sync.WaitGroup

// 	r.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()

// 			// Simulate processing
// 			time.Sleep(5 * time.Second)
// 			// fmt.Fprintln(w, "Processing complete")
//             fmt.Printf("Processing complete")
// 		}()
// 	})

// 	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
// 		wg.Wait()
// 		fmt.Fprintln(w, "All requests processed")
// 	})

// 	http.Handle("/", r)
// 	http.ListenAndServe(":8080", nil)
// }

// Design 1 (Async):
func main() {
	r := http.NewServeMux()
	var mu sync.Mutex
	var processing bool

	r.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		if processing {
			mu.Unlock()
			fmt.Fprintln(w, "Processing is already in progress")
			return
		}
		processing = true
		mu.Unlock()

		go func() {
			// Simulate processing
			time.Sleep(5 * time.Second)

			mu.Lock()
			processing = false
			mu.Unlock()

			fmt.Fprintln(w, "Processing complete")
		}()
	})

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		if processing {
			mu.Unlock()
			fmt.Fprintln(w, "Processing is in progress")
		} else {
			mu.Unlock()
			fmt.Fprintln(w, "No processing is currently happening")
		}
	})

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
