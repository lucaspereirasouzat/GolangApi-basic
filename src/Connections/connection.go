package connection

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	redis "github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientRedis *redis.Client
var clientMongo *mongo.Client

type Log struct {
	Time    string
	Latency string
	Status  int
	Url     string
	User    interface{}
	Method  string
}

//RedisConnection faz a coneção com o redis
func RedisConnection() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})
	fmt.Println("redis connection")
	clientRedis = client
	return client
}

// SetItemRedis Salva no Redis
func SetItemRedis(key string, value interface{}) error {
	err := clientRedis.Set(key, value, 0).Err()
	fmt.Println("error", err)
	return err
}

// GetItemRedis pega o item no redis
func GetItemRedis(key string) (string, error) {
	val, err := clientRedis.Get(key).Result()
	return val, err
}

func MongoConnection() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		fmt.Println("cliente mongo", err)
	}
	clientMongo = client
}

func InsertMongoDB(db string, cole string, value Log) *mongo.InsertOneResult {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err := clientMongo.Connect(ctx)
	if err != nil {
		fmt.Println("error mongo", err)
	}
	// result, err := json.Marshal(value)
	// resultString := string(result)

	// //filter := bson.D{{"hello", "world"}}
	// fmt.Println(value)
	// fmt.Println(result)
	// //fmt.Println(filter)
	// fmt.Println("result String", resultString)
	collection := clientMongo.Database(db).Collection(cole)

	res, err := collection.InsertOne(context.Background(), bson.M{
		"time":    value.Time,
		"latency": value.Latency,
		"status":  value.Status,
		"url":     value.Url,
		"user":    value.User,
		"method":  value.Method})
	if err != nil {
		fmt.Println(res, err)
	}
	return res
}
func ShowMongoDB() {

}

func SelectMongoDB(db string, cole string, search string, skip int64, limit int64) (inface []bson.M, total int64, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = clientMongo.Connect(ctx)

	collection := clientMongo.Database(db).Collection(cole)

	options := options.Find()

	// Sort by `_id` field descending
	options.SetSort(bson.D{{"_id", -1}})

	// Limit by 10 documents only
	options.SetLimit(limit)

	options.SetSkip(skip)

	//	filter := bson.M{"url": bson.M{"$search": search}}

	cur, err := collection.Find(ctx, bson.M{}, options) //.Sort("_id").Limit()
	if err != nil {
		fmt.Println("error mongo", err)
	}

	count, err := clientMongo.Database(db).Collection(cole).CountDocuments(context.Background(), bson.D{})

	defer cur.Close(ctx)
	//var inface []bson.M
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		inface = append(inface, result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return inface, count, nil
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

// ConnectionResult String result connection
var ConnectionResult string = ""

// CreateConnection  Cria connexão
func CreateConnection() *sqlx.DB {
	// faz a junção de string para conectar ao banco
	db, err := sqlx.Open("postgres", ConnectionResult)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

//ConnectionDB Cria a string de connexão
func ConnectionDB() {
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

	ConnectionResult = fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), portN, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	//	return ConnectionResult
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
