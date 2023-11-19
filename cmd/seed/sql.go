package main

import (
	"embed"
	_ "embed"
	"log"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

//go:embed schema.sql
var sqlFile embed.FS

func main() {
	db, err := sqlx.Connect("sqlite", "./junjo.db")
	if err != nil {
		log.Fatalln(err)
	}

	sqlBytes, err := sqlFile.ReadFile("schema.sql")
	if err != nil {
		panic(err)
	}

	sqlContent := string(sqlBytes)

	db.MustExec(sqlContent)
}
