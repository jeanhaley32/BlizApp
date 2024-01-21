// Function definition for the response writer
package main

import (
	"fmt"
	"net/http"
)

// Defined Function for http.HandleFunc
func blizh(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test")
}
