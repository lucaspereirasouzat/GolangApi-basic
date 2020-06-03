package main

import (
	"fmt"
	"log"
	"os"

	connection "docker.go/src/Connections"
	migration "docker.go/src/Migrations"
	routes "docker.go/src/Routes"
	functions "docker.go/src/functions"
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

	var db = connection.CreateConnection()

	fmt.Println("Fez conexão")

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Pega os dados do terminal
	arg := os.Args[1:]
	// se tiver a palavra migration ele faz a migration
	if functions.Contains(arg, "migration") {
		fmt.Println("Fazendo a migration")
		schema := migration.Schema()
		migrationExecutation := db.MustExec(schema)
		if migrationExecutation != nil {
			fmt.Println("Erro migration")

			log.Fatal(migrationExecutation)
		}
	}

	defer db.Close()

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
