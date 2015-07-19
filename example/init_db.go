// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@/go2sql_example")
	if err != nil {
		panic(err)
	}

	must(db.Exec(`
		DROP TABLE IF EXISTS languages;
	`))
	must(db.Exec(`
		Create TABLE languages (
			id int NOT NULL AUTO_INCREMENT,
			name TEXT,
			words_stat BIGINT,
			field1 varchar(255) not null default 'text',
			field2 varchar(255) not null default 'text',
			field3 varchar(255) not null default 'text',
			field4 varchar(255) not null default 'text',
			field5 varchar(255) not null default 'text',
			field6 varchar(255) not null default 'text',
			field7 varchar(255) not null default 'text',
			PRIMARY KEY (id)
		);
	`))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		must(db.Exec(fmt.Sprintf("INSERT INTO languages (name, words_stat) VALUES ('Mr. Tester', %d);", rand.Int())))
	}
}

func must(r sql.Result, err error) {
	if err != nil {
		panic(err)
	}
}
