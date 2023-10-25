package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Design 1 (Sync): Design a request dispatcher for a web-server that accepts and processes incoming web requests concurrently and responds synchronously.
// func main() {

// 	r := http.NewServeMux()
// 	var wg sync.WaitGroup

// 	r.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
// 		for i := 1; i <= 5; i++ {
// 			wg.Add(1)

// 			go func(id int) {
// 				defer wg.Done()
// 				fmt.Printf("Worker %d starting\n", id)

// 				time.Sleep(time.Second * 10)
// 				fmt.Printf("Worker %d done\n", id)
// 			}(i)
// 		}
// 	})

// 	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
// 		wg.Wait()
// 		fmt.Fprintln(w, "All requests processed")
// 	})

// 	http.Handle("/", r)
// 	http.ListenAndServe(":8080", nil)
// }

// Design 1 (Async): Modify design to support asynchronous response
func main() {
	r := http.NewServeMux()
	var mu sync.Mutex
	var processing bool

	r.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Fprintln(w, "Processing...")

		if processing {
			fmt.Fprintln(w, "Processing is already in progress")
			return
		}
		
		processing = true

		go func() {
			mu.Lock()
			// Simulate processing
			time.Sleep(5 * time.Second)

			processing = false
			mu.Unlock()

			fmt.Fprintln(w, "Processing complete")
		}()
	})

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if processing {
			fmt.Fprintln(w, "Processing is in progress")
		} else {
			mu.Lock()
			fmt.Fprintln(w, "No processing is currently happening")
			defer mu.Unlock()
		}
	})

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
