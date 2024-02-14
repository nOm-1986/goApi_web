package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nOm-1986/goApi_web/internal/user"
)

func main() {
	router := mux.NewRouter()

	userSrv := user.NewService()
	userEnd := user.MakeEndpoints(userSrv)

	//Routes
	router.HandleFunc("/v1/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/v1/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/v1/users", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/v1/users", userEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		//Handler:      http.TimeoutHandler(router, time.Second*3, "¡¡¡ Timeout !!!"),
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	fmt.Println("!!!! Starting server in port 127.0.0.1:8000 !!!!")
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
