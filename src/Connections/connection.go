package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

//ConnectionDB Cria a string de connexão
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
func QueryTable(table string, page uint64, rowsPerPage uint64, data interface{}) error {
	db := CreateConnection()
	err := db.Select(data, `SELECT * FROM `+table+` LIMIT ($1) OFFSET ($2)`, rowsPerPage, rowsPerPage*page)
	db.Close()
	return err
}

//QueryTotalTable Faz a query em uma tabela para pegar o numero total de itens
func QueryTotalTable(table string) (total uint64, err error) {
	db := CreateConnection()
	err = db.Get(&total, `SELECT COUNT(*) FROM `+table)
	db.Close()
	return total, err
}

//InserIntoTable insere na tabela de acordo com os campos e valores
func InserIntoTable(table string, fields []string, args ...interface{}) (result sql.Result, err error) {
	allfields := strings.Join(fields, ", ")
	var fieldsd []string
	for i := 0; i < len(fields); i++ {
		fieldsd = append(fieldsd, "$"+strconv.FormatInt(int64(i)+1, 16))
	}
	itensvalues := strings.Join(fieldsd, ", ")
	db := CreateConnection()
	result, err = db.Exec("INSERT INTO "+table+"( "+allfields+" ) VALUES ( "+itensvalues+" )", args...)
	db.Close()
	return result, err
}

//ShowRow mostra um item de alguma tabela
func ShowRow(table string, row interface{}, field string, args interface{}) (err error) {
	db := CreateConnection()
	err = db.Get(row, "SELECT * FROM "+table+" WHERE "+field+"=($1)", args)
	fmt.Println(row)
	db.Close()
	return err
}

//InserIntoTable insere na tabela de acordo com os campos e valores
func UpdateRow(table string, fields []string, field string, key interface{}, args ...interface{}) (result sql.Result, err error) {
	allfields := strings.Join(fields, ",")
	db := CreateConnection()
	result, err = db.MustBegin().Exec("INSERT INTO "+table+"( "+allfields+" )", args)
	db.Close()
	return result, err
}

// //Cria Faz a query em uma tabela
// func Cria(db *sqlx.DB, table string, page int, rowsPerPage int, data interface{}) error {
// 	q := fmt.Sprintf(`INSERT INTO "%s" VALUES "%data"`, table, data)
// 	err := db.Select(&data, q)
// 	return err
// }
