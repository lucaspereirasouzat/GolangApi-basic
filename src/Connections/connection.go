package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func MongoConnection() {
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	// err = client.Ping(ctx, readpref.Primary())
	// collection := client.Database("testing").Collection("numbers")
	// fmt.Println("entrou no mongo db")
	// fmt.Println(collection)
	// ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	// res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	// id := res.InsertedID
	//fmt.Println(id)

	// es, err := elasticsearch.NewDefaultClient()
	// if err != nil {
	// 	log.Fatalf("Error creating the client: %s", err)
	// }

	// res, err := es.Info()
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }

	// //defer res.Body.Close()
	// log.Println(res)
}

func Elastic() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	//defer res.Body.Close()
	log.Println(res)
}

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
func QueryTable(db *sqlx.DB, table string, selectFields string, page uint64, rowsPerPage uint64, query string, data interface{}) error {
	err := db.Select(data, `SELECT`+selectFields+` FROM `+table+query+` LIMIT ($1) OFFSET ($2)`, rowsPerPage, rowsPerPage*page)
	return err
}

//QueryTotalTable Faz a query em uma tabela para pegar o numero total de itens
func QueryTotalTable(db *sqlx.DB, table string, query string) (total uint64, err error) {
	err = db.Get(&total, `SELECT COUNT(*) FROM `+table+query)
	return total, err
}

//InserIntoTable insere na tabela de acordo com os campos e valores
func InserIntoTable(db *sqlx.DB, table string, fields []string, args ...interface{}) (result sql.Result, err error) {
	allfields := strings.Join(fields, ", ")
	var fieldsd []string
	for i := 0; i < len(fields); i++ {
		fieldsd = append(fieldsd, "$"+strconv.FormatInt(int64(i)+1, 16))
	}
	itensvalues := strings.Join(fieldsd, ", ")
	result, err = db.Exec("INSERT INTO "+table+"( "+allfields+" ) VALUES ( "+itensvalues+" )", args...)
	return result, err
}

//ShowRow mostra um item de alguma tabela
func ShowRow(db *sqlx.DB, table string, row interface{}, field string, args interface{}) (err error) {
	err = db.Get(row, "SELECT * FROM "+table+" WHERE "+field+"=($1)", args)
	// db.Close()
	return err
}

//UpdateRow insere na tabela de acordo com os campos e valores
func UpdateRow(db *sqlx.DB, table string, fields []string, field string, key interface{}, args ...interface{}) (result sql.Result, err error) {
	var allfields string
	for i, v := range fields {
		allfields += v + "=$" + strconv.Itoa(i+2)
	}
	result, err = db.Exec("UPDATE "+table+" SET "+allfields+" WHERE"+field+"=$1", key, args)
	return result, err
}

//DeleteRow exclui um item da tabela
func DeleteRow(db *sqlx.DB, table string, field string, args interface{}) (err error, result sql.Result) {
	result, err = db.Exec("DELETE FROM "+table+" WHERE "+field+"=($1)", args)
	return err, result
}

// //Cria Faz a query em uma tabela
// func Cria(db *sqlx.DB, table string, page int, rowsPerPage int, data interface{}) error {
// 	q := fmt.Sprintf(`INSERT INTO "%s" VALUES "%data"`, table, data)
// 	err := db.Select(&data, q)
// 	return err
// }
