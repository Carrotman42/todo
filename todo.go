package main

import (
	"fmt"
	"net/http"
	"io"
)

func write(w io.Writer, s string) {
	io.WriteString(w, s)
}

func main() {
	fmt.Println("This is your to-do list")
	
	http.HandleFunc("/", MainPage)
	http.ListenAndServe(":16005", nil)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	write(w, "Hello!")
}
