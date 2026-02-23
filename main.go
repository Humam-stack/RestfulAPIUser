package main

import (
	"log"
	"os"
	"restfulapi/config"
	"restfulapi/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error load .env file")
	}

	// koneksi
	config.ConnectDB()

	//setuproutes
	router := routes.SetupRoutes()

	//jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8888"
	}

	log.Printf("server berjalan di port https://localhost:%s", port)
	log.Fatal(router.Run(":" + port))

}
