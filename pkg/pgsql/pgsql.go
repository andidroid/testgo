package pgsql

import (
	"fmt"

	"database/sql"

	"github.com/andidroid/testgo/pkg/util"
	_ "github.com/lib/pq" // sql behavior modified
)

var connection *sql.DB

func init() {
	var err error
	connection, err = InitDB()
	util.CheckErr(err)
}

func InitDB() (*sql.DB, error) {
	fmt.Println("create database connection")
	var connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// var err error
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		fmt.Println("Error ", err)
	}

	return db, err
}

func main() {

}

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "osm"
	schema   = "public"
)

func GetConnection() (*sql.DB, error) {
	return connection, nil
}

// return db
// age := 21
// rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)

// var userid int
// err := db.QueryRow(`INSERT INTO users(name, favorite_fruit, age)
// 	VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)

// if err != nil {
// 	return nil, err
// }

// //stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS web_url(ID SERIAL PRIMARY KEY, URL TEXT NOT NULL);")

// if err != nil {
// 	return nil, err
// }

// _, err = stmt.Exec()

// if err != nil {
// 	return nil, err
// }

// rows, err := db.Query("SELECT * FROM userinfo")
// checkErr(err)

// for rows.Next() {
// 	var uid int
// 	var username string
// 	var department string
// 	var created time.Time
// 	err = rows.Scan(&uid, &username, &department, &created)
// 	checkErr(err)
// 	fmt.Println("uid | username | department | created ")
// 	fmt.Printf("%3v | %8v | %6v | %6v\n", uid, username, department, created)
// }
