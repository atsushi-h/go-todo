package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	fmt.Println("Server starting on :4000...")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		fmt.Println(err)
	}
}
