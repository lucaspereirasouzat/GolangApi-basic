package connection

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

// CreateConnection  Cria connexão
func CreateConnection() *sqlx.DB {
	// faz a junção de string para conectar ao banco
	db, err := sqlx.Open("postgres", ConnectionDB())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// ConnectionDB Cria a string de connexão
func ConnectionDB() string {
	// carregar env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	// conversão de string para int64
	port := os.Getenv("DB_PORT")

	portN, err := strconv.ParseInt(port, 10, 64)

	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), portN, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
}

//QueryTable Faz a query em uma tabela
func QueryTable(table string, page int64, rowsPerPage int64, data interface{}) error {
	//q := fmt.Sprintf()
	db := CreateConnection()
	err := db.Select(&data, `SELECT * FROM ? LIMIT ? OFFSET ?`, table, strconv.FormatInt(page, 64), strconv.FormatInt(page*rowsPerPage, 64))
	db.Close()
	return err
}

// //Cria Faz a query em uma tabela
// func Cria(db *sqlx.DB, table string, page int, rowsPerPage int, data interface{}) error {
// 	q := fmt.Sprintf(`INSERT INTO "%s" VALUES "%data"`, table, data)
// 	err := db.Select(&data, q)
// 	return err
// }
