package main

import (
	"github.com/charliekim2/songsleuths/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	conn, err := db.Connect()
	if err != nil {
		panic(err)
	}

	conn.AutoMigrate(
		&db.Game{},
		&db.Player{},
		&db.Tierlist{},
		&db.Tier{},
		&db.Submission{},
		&db.Song{},
		&db.Ranking{},
	)
}
