package controller

import (
	"io"
	"net/http"
	"log"
)

func Compile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}
	response, err := http.Post("https://golang.org/compile", r.Header.Get("Content-type"), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error sending code to compile:", err)
		return
	}
	defer response.Body.Close()
	if _, err := io.Copy(w, response.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error sending code to compile:", err)
		return
	}
}
