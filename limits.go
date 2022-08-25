package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Init database
func init_limit_db() *sql.DB {
	db, err := sql.Open("sqlite3", "limits.db")
	if err != nil {
		log.Fatal(err)
	}
	cmd := `CREATE TABLE IF NOT EXISTS Limits (id INTEGER PRIMARY KEY AUTOINCREMENT, acct TEXT, ticket INTEGER, time TEXT)`
	stat, err := db.Prepare(cmd)
	if err != nil {
		log.Fatal(err)
	}
	stat.Exec()

	return db
}

// Add account to database
func add_to_db(acct string, limit uint16) {
	db := init_limit_db()
	cmd := `INSERT INTO Limits (acct, ticket) VALUES (?, ?)`
	stat, err := db.Prepare(cmd)
	if err != nil {
		log.Fatal(err)
	}
	stat.Exec(acct, limit)
}

// Take ticket for tooting
func take_ticket(acct string) {
	db := init_limit_db()
	cmd1 := `SELECT ticket FROM Limits WHERE acct = ?`
	cmd2 := `UPDATE Limits SET ticket = ?, time = ? WHERE acct = ?`

	var ticket uint16
	db.QueryRow(cmd1, acct).Scan(&ticket)
	if ticket > 0 {
		ticket = ticket - 1
	}

	stat, err := db.Prepare(cmd2)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	last_toot_at := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local).Format("2006/01/02 15:04:05 MST")

	stat.Exec(ticket, last_toot_at, acct)
}

// Check followed once
func followed(acct string) bool {
	db := init_limit_db()
	cmd := `SELECT acct FROM Limits WHERE acct = ?`
	err := db.QueryRow(cmd, acct).Scan(&acct)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}

		return false
	}

	return true
}

// Check ticket availability
func check_ticket(acct string, ticket uint16, toots_interval uint16) uint16 {
	db := init_limit_db()
	cmd1 := `SELECT ticket FROM Limits WHERE acct = ?`
	cmd2 := `SELECT time FROM Limits WHERE acct = ?`

	var tickets uint16
	var lastS string

	db.QueryRow(cmd1, acct).Scan(&tickets)
	db.QueryRow(cmd2, acct).Scan(&lastS)

	lastT, _ := time.Parse("2006/01/02 15:04:05 MST", lastS)

	since := time.Since(lastT)
	limit := fmt.Sprintf("%dh", toots_interval)
	interval, _ := time.ParseDuration(limit)

	if since >= interval {
		cmd := `UPDATE Limits SET ticket = ? WHERE acct = ?`
		stat, err := db.Prepare(cmd)
		if err != nil {
			log.Fatal(err)
		}
		stat.Exec(ticket, acct)

		return ticket
	}

	return tickets
}
