package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/kshwedha/geodata-go/src/common/config"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func testConnection(db *sql.DB) {
	// Test the connection to the database
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() (Config, error) {
	// Load the configuration from a file or any other source
	// Implement your logic here to load the config
	configData := config.GetConfig()
	config := Config{
		Host:     configData.GetString("db.HOST"),
		Port:     configData.GetInt("db.PORT"),
		User:     configData.GetString("db.USER"),
		Password: configData.GetString("db.PASSWORD"),
		DBName:   configData.GetString("db.DATABASE"),
	}
	return config, nil
}

func InitDB() (*sql.DB, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Create a connection string using the configuration data

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test the connection to the database
	testConnection(db)

	return db, nil
}

func ExecPsqlResult(db *sql.DB, query string) (int64, error) {
	// Execute the query
	results, err := db.Exec(query)
	if err != nil {
		return 0, err
	}
	// Get the number of rows affected
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func ExecPsqlRows(db *sql.DB, query string) *sql.Rows {
	// Create a statement object
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the query
	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Create a results set
	results, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	return results
}

func Sample() {
	db, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT * FROM users"
	results := ExecPsqlRows(db, query)
	defer results.Close()
	err = printData(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Query executed successfully!")
}

func printData(results *sql.Rows) error {
	// Iterate over the results set
	for results.Next() {
		var id int
		var username string
		var mail string
		var password string
		var created_at time.Time
		// Scan the results into variables
		err := results.Scan(&id, &username, &mail, &password, &created_at)
		if err != nil {
			return err
		}
		fmt.Println(id, username)
	}
	return nil
}
