package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//port := ":3333"
	//Routes
	router.HandleFunc("/v1/users", getUsers).Methods("get")
	router.HandleFunc("/v1/courses", getCourses).Methods("get")

	srv := &http.Server{
		//Handler:      http.TimeoutHandler(router, time.Second*3, "¡¡¡ Timeout !!!"),
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	time.Sleep(7 * time.Second)
	fmt.Println("Sigue")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"ok": "course"})
}
