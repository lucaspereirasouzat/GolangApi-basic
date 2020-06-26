package main

import (
	"fmt"
	"log"
	"os"
	"time"

	commands "docker.go/src/Commands"
	connection "docker.go/src/Connections"
	middleware "docker.go/src/Middleware"
	routes "docker.go/src/Routes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	cors "github.com/itsjamie/gin-cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// use a single instance of Validate, it caches struct info
var Validate *validator.Validate

func main() {

	// carregar env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	//mail.ExampleSendMail("lucaspereiradesouzat@gmail.com", "TEste de email do golang")
	//mail.ExampleSendMail("lucaspereiradesouzat@gmail.com", "TEste de email do golang")
	connection.ConnectionDB()
	connection.RedisConnection()
	connection.MongoConnection()

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
	router.Use(middleware.Logger())
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	// ligar as rotas para o gin
	routes.UsersRoutes(router)
	routes.AuthRoutes(router)
	routes.FileRoutes(router)
	routes.NotificationRoutes(router)

	Validate = validator.New()

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	router.Run(":" + os.Getenv("PORT"))
	// go func() {
	//endless.ListenAndServe(":4242", router)
	// }()

}
