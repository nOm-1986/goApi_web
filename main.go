package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nOm-1986/goApi_web/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	router := mux.NewRouter()
	_ = godotenv.Load()
	loger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//Modo debug
	db = db.Debug()

	//Crear la tabla
	_ = db.AutoMigrate(&user.User{})

	userRepo := user.NewRepo(loger, db)
	userSrv := user.NewService(loger, userRepo)
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
