package main

import (
	"github.com/zhughes3/website/db"
	"github.com/zhughes3/website/user"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	router := NewRouter()

	db.DB = db.SetupDB()
	db.DB.AutoMigrate(&user.User{})
	defer db.DB.Close()

	server := &http.Server{
		Handler: router,
		Addr: "localhost:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}