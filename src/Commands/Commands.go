package commands

import (
	"fmt"
	"log"
	"os"

	migration "docker.go/src/Migrations"
	seeder "docker.go/src/Migrations/Seeder"
	"docker.go/src/functions"
	"github.com/jmoiron/sqlx"
)

// Commands é utilizado para utilizar os comandos no terminal
func Commands(db *sqlx.DB) {
	fmt.Println("Commandos")
	fmt.Println("migration faz a migração das tabelas")
	fmt.Println("seed faz o cadastro do adm no banco de dados")
	fmt.Println(" ")
	// Pega os dados do terminal
	arg := os.Args[1:]
	// se tiver a palavra migration ele faz a migration
	if functions.Contains(arg, "migration") {
		MakeMigration(db)
	}

	// se tiver a palavra migration ele faz a migration
	if functions.Contains(arg, "--seed") {
		MakeSeed(db)
	}
}

func MakeMigration(db *sqlx.DB) {
	fmt.Println("Fazendo a migration")
	schema := migration.Schema()
	migrationExecutation := db.MustExec(schema)
	if migrationExecutation != nil {
		fmt.Println("Erro migration")

		log.Fatal(migrationExecutation)
	}
}

func MakeSeed(db *sqlx.DB) {
	fmt.Println("Fazendo a seeed")
	seeder.Seed(db)
	// migrationExecutation := db.MustExec(schema)
	// if migrationExecutation != nil {
	// 	fmt.Println("Erro migration")

	// 	log.Fatal(migrationExecutation)
	// }
}
