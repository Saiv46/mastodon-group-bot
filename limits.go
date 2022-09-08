package main

import (
	"database/sql"
	"fmt"

	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db = init_limit_db()
)

// Init database
func init_limit_db() *sql.DB {
	db, err := sql.Open("sqlite3", *DBPath)
	if err != nil {
		ErrorLogger.Println("Open database")
	}

	cmd1 := `CREATE TABLE IF NOT EXISTS Limits (id INTEGER PRIMARY KEY AUTOINCREMENT, acct TEXT, ticket INTEGER, order_msg INTEGER, got_notice INTEGER, posted TEXT)`
	cmd2 := `CREATE TABLE IF NOT EXISTS MsgHashs (message_hash TEXT)`

	stat1, err := db.Prepare(cmd1)
	if err != nil {
		ErrorLogger.Println("Create database and table Limits")
	}
	stat1.Exec()

	stat2, err := db.Prepare(cmd2)
	if err != nil {
		ErrorLogger.Println("Create database and table MsgHashs")
	}
	stat2.Exec()

	return db
}

// Add account to database
func add_to_db(acct string) {
	cmd := `INSERT INTO Limits (acct, ticket, order_msg, got_notice) VALUES (?, ?, ?, ?)`
	stat, err := db.Prepare(cmd)
	if err != nil {
		ErrorLogger.Println("Add account to database")
	}
	stat.Exec(acct, Conf.Max_toots, 0, 0)
}

// Check followed once
func exist_in_database(acct string) bool {
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

// Take ticket for tooting
func take_ticket(acct string) {
	cmd1 := `SELECT ticket FROM Limits WHERE acct = ?`
	cmd2 := `UPDATE Limits SET ticket = ?, posted = ? WHERE acct = ?`

	var ticket uint
	db.QueryRow(cmd1, acct).Scan(&ticket)
	if ticket > 0 {
		ticket = ticket - 1
	}

	now := time.Now()
	last_toot_at := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local).Format("2006/01/02 15:04:05 MST")

	stat, err := db.Prepare(cmd2)
	if err != nil {
		ErrorLogger.Println("Take ticket")
	}

	stat.Exec(ticket, last_toot_at, acct)
}

// Check ticket availability
func check_ticket(acct string) uint {
	cmd1 := `SELECT ticket FROM Limits WHERE acct = ?`
	cmd2 := `SELECT posted FROM Limits WHERE acct = ?`

	var tickets uint
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

// Save message hash
func save_msg_hash(hash string) {
	cmd1 := `SELECT COUNT(*) FROM MsgHashs`
	cmd2 := `DELETE FROM MsgHashs WHERE ROWID IN (SELECT ROWID FROM MsgHashs LIMIT 1)`
	cmd3 := `INSERT INTO MsgHashs (message_hash) VALUES (?)`

	var rows uint

	db.QueryRow(cmd1).Scan(&rows)

	if rows >= Conf.Duplicate_buf {
		superfluous := rows - Conf.Duplicate_buf

		for i := uint(0); i <= superfluous; i++ {
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

// Check message hash
func check_msg_hash(hash string) bool {
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

// Count order
func count_order(acct string) {
	cmd1 := `UPDATE Limits SET order_msg = ? WHERE acct != ?`
	cmd2 := `SELECT order_msg FROM Limits WHERE acct = ?`
	cmd3 := `UPDATE Limits SET order_msg = ? WHERE acct = ?`

	stat1, err := db.Prepare(cmd1)
	if err != nil {
		ErrorLogger.Println("Count order to zero")
	}

	stat1.Exec(0, acct)

	var order uint
	db.QueryRow(cmd2, acct).Scan(&order)
	if order < Conf.Order_limit {
		order = order + 1
	}

	stat2, err := db.Prepare(cmd3)
	if err != nil {
		ErrorLogger.Println("Count order")
	}

	stat2.Exec(order, acct)
}

// Check order
func check_order(acct string) uint {
	cmd := `SELECT order_msg FROM Limits WHERE acct = ?`

	var order uint
	db.QueryRow(cmd, acct).Scan(&order)

	return order
}

// Mark notice
func mark_notice(acct string) {
	cmd1 := `SELECT got_notice FROM Limits WHERE acct = ?`
	cmd2 := `UPDATE Limits SET got_notice = ? WHERE acct = ?`

	var notice uint
	db.QueryRow(cmd1, acct).Scan(&notice)

	if notice == 0 {
		notice = notice + 1
	}

	stat, err := db.Prepare(cmd2)
	if err != nil {
		ErrorLogger.Println("Mark notice")
	}
	stat.Exec(notice, acct)
}

// Reset notice counter
func reset_notice_counter() {
	cmd := `UPDATE Limits SET got_notice = ?`

	stat, err := db.Prepare(cmd)
	if err != nil {
		ErrorLogger.Println("Reset notice counter")
	}
	stat.Exec(0)
}

// Check if got notification
func got_notice(acct string) uint {
	cmd := `SELECT got_notice FROM Limits WHERE acct = ?`

	var notice uint
	db.QueryRow(cmd, acct).Scan(&notice)

	return notice
}
