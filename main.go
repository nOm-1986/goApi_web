package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	port := ":3333"
	//Routes
	http.HandleFunc("/v1/users", getUsers)
	http.HandleFunc("/v1/courses", getCourses)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /users")
	io.WriteString(w, "This is my user endpoint!\n")
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /courses")
	io.WriteString(w, "Courses endpoint!\n")
}
