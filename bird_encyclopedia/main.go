package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	//we have to import the driver, but don't use it in our code
	// so we use the `_` symbol
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Bird struct {
	Species string
	Description string 
}


func main() {
	// The `sql.Open` function opens a new `*sql.DB` instance. We specify the driver name
	// and the URI for our database. Here, we're using a Postgres URI
	db, err := sql.Open("pgx", "postgresql://localhost:5432/bird_encyclopedia")
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(1 * time.Second)
	db.SetConnMaxLifetime(30 * time.Second)
	// To verify the connection to our database instance, we can call the `Ping`
	// method. If no error is returned, we can assume a successful connection
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}
	fmt.Println("database is reachable")


row := db.QueryRow("SELECT bird, Description FROM birds LIMIT 1")

bird := Bird{}

if err := row.Scan(&bird.Species, &bird.Description); err != nil {
	log.Fatalf("could not scan row: %v", err)
}
fmt.Printf("found bird: %+v\n", bird)

rows, err := db.Query("SELECT bird, Description FROM birds limit 10")
if err != nil {
	log.Fatalf("could not execute query: %v", err)
}

birds := []Bird{}



for rows.Next() {
	bird := Bird{}

	if err := rows.Scan(&bird.Species, &bird.Description); err != nil {
		log.Fatalf("could not scan row: %v", err)
	}

	birds = append(birds, bird)

	
}
fmt.Printf("found %d birds: %+v", len(birds), birds)
_, err = db.Exec("DELETE FROM birds WHERE bird=$1", "rooster")
if err != nil {
	log.Fatalf("could not delete row: %v", err)
}

newBird := Bird{
	Species: "rooster",
	Description: "wakes you up in the morning",
}

result, err := db.Exec("INSER INTO birds (bird, description) VALUES ($1, $2)", newBird.Species, newBird.Description)
if err != nil {
	log.Fatalf("could not insert row: %v", err)
}

rowsAffected, err := result.RowsAffected()
if err != nil {
	log.Fatalf("could not get affected rows: %v", err)
}

fmt.Println("inserted", rowsAffected, "rows")

}