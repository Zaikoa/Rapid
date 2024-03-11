package database

import (
	"embed"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx"
)

var conn *pgx.Conn
var connMutex sync.Mutex

// Create target connection for the database
func GetConn() (*pgx.Conn, error) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if conn != nil {
		return conn, nil
	}

	connConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "rapidtransfer",
		User:     "swen344",
		Password: "Forzano17**",
	}

	newConn, err := pgx.Connect(connConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect: %v", err)
	}

	conn = newConn
	return conn, nil
}

// Inits all of the tables for the database
func InitializeDatabase() {
	var content embed.FS
	path, _ := content.ReadFile("database.sql")

	// Explicitly initialize the connection
	conn, err := GetConn()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	}

	// Execute the SQL file
	_, err = conn.Exec(string(path))
	if err != nil {
		fmt.Println("Error executing SQL file:", err)
		os.Exit(1)
	}
}
