package main

import (
	"fmt"
	"log"
	"os"

	commands "docker.go/src/Commands"
	connection "docker.go/src/Connections"
	routes "docker.go/src/Routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// carregar env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	db := connection.CreateConnection()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fez conexão")

	// Comandos no terminal
	commands.Commands(db)

	db.Close()

	fmt.Println("Successfully connected!")

	router := gin.Default()
	// ligar as rotas para o gin
	routes.UsersRoutes(router)
	routes.AuthRoutes(router)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	router.Run(":" + os.Getenv("PORT"))
	// go func() {
	//endless.ListenAndServe(":4242", router)
	// }()

}
