//go:build ignore

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: go run seed.go <count>")
	}

	n, _ := strconv.Atoi(os.Args[1])
	store, err := NewStore("history.db")
	if err != nil {
		panic(err)
	}

	store.db.Exec("PRAGMA synchronous = OFF")
	store.db.Exec("PRAGMA journal_mode = MEMORY")

	tx, _ := store.db.Begin()

	batchSize := 1000
	q := strings.Builder{}
	vals := make([]any, 0, batchSize*2)
	baseTime := time.Now().UTC()

	fmt.Printf("generating %d rows...\n", n)

	for i := 0; i < n; i++ {
		ip := strconv.Itoa(rand.IntN(254)+1) + "." +
			strconv.Itoa(rand.IntN(255)) + "." +
			strconv.Itoa(rand.IntN(255)) + "." +
			strconv.Itoa(rand.IntN(255))

		ts := baseTime.Add(time.Duration(-i) * time.Hour)

		if len(vals) == 0 {
			q.WriteString("INSERT INTO ip_history (ip, ts) VALUES ")
		} else {
			q.WriteString(",")
		}
		q.WriteString("(?,?)")
		vals = append(vals, ip, ts)

		if len(vals) >= batchSize*2 {
			if _, err := tx.Exec(q.String(), vals...); err != nil {
				panic(err)
			}
			q.Reset()
			vals = vals[:0]
		}
	}

	if len(vals) > 0 {
		tx.Exec(q.String(), vals...)
	}

	tx.Commit()
	fmt.Println("done")
}
