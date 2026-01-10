package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Record struct {
	ID        int
	Timestamp time.Time
	IP        string
}

type Store struct {
	db *sql.DB
}

func NewStore(path string) (*Store, error) {
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous=NORMAL&_pragma=mmap_size=268435456", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS ip_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ts DATETIME DEFAULT CURRENT_TIMESTAMP,
		ip TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_ts ON ip_history(ts DESC);
	CREATE INDEX IF NOT EXISTS idx_ip ON ip_history(ip COLLATE NOCASE);`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) GetLatest() (string, error) {
	var ip string
	err := s.db.QueryRow("SELECT ip FROM ip_history ORDER BY ts DESC LIMIT 1").Scan(&ip)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return ip, err
}

func (s *Store) Insert(ip string) error {
	_, err := s.db.Exec("INSERT INTO ip_history (ip, ts) VALUES (?, ?)", ip, time.Now().UTC())
	return err
}

func (s *Store) FetchPage(query string, page, pageSize int) ([]Record, bool, error) {
	offset := (page - 1) * pageSize
	limit := pageSize + 1

	q := query + "%"
	if query == "" {
		q = "%"
	}

	rows, err := s.db.Query(
		"SELECT id, ts, ip FROM ip_history WHERE ip LIKE ? ORDER BY ts DESC LIMIT ? OFFSET ?",
		q, limit, offset,
	)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Timestamp, &r.IP); err != nil {
			return nil, false, err
		}
		records = append(records, r)
	}

	hasMore := len(records) > pageSize
	if hasMore {
		records = records[:pageSize]
	}
	return records, hasMore, nil
}
