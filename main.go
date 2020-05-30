package main

import (
	connection "docker.go/src/Connections"
	routes "docker.go/src/Routes"
	"fmt"

	"github.com/gin-gonic/gin"

	"log"
	"os"

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
	//schema := migration.Schema()
	// migrationExecutation := db.MustExec(schema)
	// if migrationExecutation != nil {
	// 	fmt.Println("Erro migration")

	// 	log.Fatal(migrationExecutation)
	// }
	defer db.Close()

	fmt.Println("Successfully connected!")

	// fmt.Println(psqlInfo)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	router := gin.Default()

	//	routes.AuthRoutes(router)
	routes.UsersRoutes(router)
	routes.AuthRoutes(router)
	// userRoutes.Routes(router)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	router.Run(":" + os.Getenv("PORT"))
	// go func() {
	//endless.ListenAndServe(":4242", router)
	// }()

}
