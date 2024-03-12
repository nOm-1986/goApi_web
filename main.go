package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nOm-1986/goApi_web/internal/course"
	"github.com/nOm-1986/goApi_web/internal/enrollment"
	"github.com/nOm-1986/goApi_web/internal/user"
	"github.com/nOm-1986/goApi_web/pkg/bootstrap"
)

func main() {
	router := mux.NewRouter()
	_ = godotenv.Load()
	loger := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		loger.Fatal(err)
	}

	userRepo := user.NewRepo(loger, db)
	userSrv := user.NewService(loger, userRepo)
	userEnd := user.MakeEndpoints(userSrv)

	courseRepo := course.NewRepo(loger, db)
	courseSrv := course.NewService(loger, courseRepo)
	courseEnd := course.MakeEndpoints(courseSrv)

	enrollRepo := enrollment.NewRepo(loger, db)
	enrollSrv := enrollment.NewService(loger, enrollRepo)
	enrollEnd := enrollment.MakeEndpoints(enrollSrv)

	//Routes User
	router.HandleFunc("/v1/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/v1/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/v1/users/{id:[0-9a-z/-]+}", userEnd.Get).Methods("GET")
	router.HandleFunc("/v1/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/v1/users/{id:[0-9a-z/-]+}", userEnd.Delete).Methods("DELETE")

	//Routes Courses
	router.HandleFunc("/v1/course", courseEnd.Create).Methods("POST")
	router.HandleFunc("/v1/course", courseEnd.GetAll).Methods("GET")
	router.HandleFunc("/v1/course/{id:[0-9a-z/-]+}", courseEnd.Get).Methods("GET")
	router.HandleFunc("/v1/course/{id}", courseEnd.Update).Methods("PATCH")

	//Routes Enrollment
	router.HandleFunc("/v1/enrollment", enrollEnd.Create).Methods("POST")

	srv := &http.Server{
		//Handler:      http.TimeoutHandler(router, time.Second*3, "¡¡¡ Timeout !!!"),
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	fmt.Println("!!!! Starting server in port 127.0.0.1:8000 !!!!")
	err2 := srv.ListenAndServe()

	if err2 != nil {
		log.Fatal(err2)
	}
}
