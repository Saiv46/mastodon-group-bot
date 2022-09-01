package main

import (
	"database/sql"
	"fmt"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Init database
func init_limit_db() *sql.DB {
	db, err := sql.Open("sqlite3", *DBPath)
	if err != nil {
		ErrorLogger.Println("Open database")
	}

	cmd1 := `CREATE TABLE IF NOT EXISTS Limits (id INTEGER PRIMARY KEY AUTOINCREMENT, acct TEXT, ticket INTEGER, time TEXT)`
	cmd2 := `CREATE TABLE IF NOT EXISTS MsgHashs (message_hash TEXT)`

	stat1, err := db.Prepare(cmd1)
	if err != nil {
		ErrorLogger.Println("Create database")
	}
	stat1.Exec()

	stat2, err := db.Prepare(cmd2)
	if err != nil {
		ErrorLogger.Println("Create database")
	}
	stat2.Exec()

	return db
}

// Add account to database
func add_to_db(acct string) {
	db := init_limit_db()
	cmd := `INSERT INTO Limits (acct, ticket) VALUES (?, ?)`
	stat, err := db.Prepare(cmd)
	if err != nil {
		ErrorLogger.Println("Add account to databse")
	}
	stat.Exec(acct, Conf.Max_toots)
}

// Save message hash
func save_msg_hash(hash string) {
	db := init_limit_db()

	cmd1 := `SELECT COUNT(*) FROM MsgHashs`
	cmd2 := `DELETE FROM MsgHashs WHERE ROWID IN (SELECT ROWID FROM MsgHashs LIMIT 1)`
	cmd3 := `INSERT INTO MsgHashs (message_hash) VALUES (?)`

	var rows int

	db.QueryRow(cmd1).Scan(&rows)

	if rows >= Conf.Duplicate_buf {
		superfluous := rows - Conf.Duplicate_buf

		for i := 0; i <= superfluous; i++ {
			stat2, err := db.Prepare(cmd2)
			if err != nil {
				ErrorLogger.Println("Delete message hash from database")
			}
			stat2.Exec()
		}
	}

	stat1, err := db.Prepare(cmd3)
	if err != nil {
		ErrorLogger.Println("Add message hash to database")
	}
	stat1.Exec(hash)
}

// Check followed once
func followed(acct string) bool {
	db := init_limit_db()
	cmd := `SELECT acct FROM Limits WHERE acct = ?`
	err := db.QueryRow(cmd, acct).Scan(&acct)
	if err != nil {
		if err != sql.ErrNoRows {
			InfoLogger.Println("Check followed")
		}

		return false
	}

	return true
}

// Check message hash
func check_msg_hash(hash string) bool {
	db := init_limit_db()
	cmd := `SELECT message_hash FROM MsgHashs WHERE message_hash = ?`
	err := db.QueryRow(cmd, hash).Scan(&hash)
	if err != nil {
		if err != sql.ErrNoRows {
			InfoLogger.Println("Check message hash in database")
		}

		return false
	}

	return true
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
		ErrorLogger.Println("Take ticket")
	}

	now := time.Now()
	last_toot_at := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local).Format("2006/01/02 15:04:05 MST")

	stat.Exec(ticket, last_toot_at, acct)
}

// Check ticket availability
func check_ticket(acct string) uint16 {
	db := init_limit_db()
	cmd1 := `SELECT ticket FROM Limits WHERE acct = ?`
	cmd2 := `SELECT time FROM Limits WHERE acct = ?`

	var tickets uint16
	var lastS string

	db.QueryRow(cmd1, acct).Scan(&tickets)
	db.QueryRow(cmd2, acct).Scan(&lastS)

	lastT, _ := time.Parse("2006/01/02 15:04:05 MST", lastS)

	since := time.Since(lastT)
	limit := fmt.Sprintf("%dh", Conf.Toots_interval)
	interval, _ := time.ParseDuration(limit)

	if since >= interval {
		cmd := `UPDATE Limits SET ticket = ? WHERE acct = ?`
		stat, err := db.Prepare(cmd)
		if err != nil {
			ErrorLogger.Println("Check ticket availability")
		}
		stat.Exec(Conf.Max_toots, acct)

		return Conf.Max_toots
	}

	return tickets
}
