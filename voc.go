package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

const (
	SPLITLINE      string = "----split----"
	VOCABULARYFILE string = "../vocabulary.txt"
	DB             string = "../words.db"
)

type Word struct {
	name  string
	trans string
}

func initDB(filePath string) *sql.DB {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func createTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE words (
    word text NOT NULL,
    translation text NOT NULL,
    createdate text DEFAULT (STRFTIME('%Y-%m-%d', 'NOW')),
    lastreviewdate text DEFAULT (STRFTIME('%Y-%m-%d', 'NOW')),
    reviewstatus INT DEFAULT 0,
    PRIMARY KEY(word, lastreviewdate);
	`
	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func readVoc(voc string) []Word {
	/*read voc file parse every word and transation
	then save it to a Word array and return it.*/

	file, err := os.Open(voc)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	file.Close()

	var words []Word
	var word Word
	var tag = false
	for n := 0; n < len(txtlines); n++ {
		if txtlines[n] == SPLITLINE {
			if tag {
				words = append(words, word)
				word.trans = ""
				tag = false
			}
			n++
			word.name = txtlines[n]
			tag = true
			continue
		}
		word.trans += txtlines[n]
		word.trans += "\n"
	}
	words = append(words, word)
	return words

}

func checkRecord(word string) bool {
	/* check if a word.name exist in DB */
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		panic(err)
	}

	stmt := `SELECT word FROM words WHERE word = ?`
	err = db.QueryRow(stmt, word).Scan(&word)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return false
	}
	return true
}

func dbStore(words []Word) {
	/* store words list to word.db */
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		panic(err)
	}

	for _, word := range words {

		if checkRecord(word.name) {
			fmt.Printf("Word: %s already exist\n", word.name)
			continue
		}

		stmt, err := db.Prepare("INSERT INTO words(word, translation, createdate, lastreviewdate, reviewstatus) values(?,?,?,?,?)")
		if err != nil {
			panic(err)
		}

		date := time.Now().Format("2006-01-02")
		res, err := stmt.Exec(word.name, word.trans, date, date, 0)
		if err != nil {
			panic(err)
		}
		rowId, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Inserted Word: %s with RowID: %d\n", word.name, rowId)
	}
}

func readDB() {
	/* read words from table word of current date */
	db, err := sql.Open("sqlite3", DB)
	if err != nil {
		panic(err)
	}

	date := time.Now().Format("2006-01-02")
	rows, err := db.Query("select * from words where lastreviewdate = ?", date)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var word string
		var trans string
		var createdate string
		var lastreviewdate string
		var reviewstatus int
		err = rows.Scan(&word, &trans, &createdate, &lastreviewdate, &reviewstatus)
		if err != nil {
			panic(err)
		}
		fmt.Println(word, trans)
	}
}

func main() {
	storePtr := flag.Bool("store", false, "Store new words to Database")
	initPtr := flag.Bool("init", false, "Init Local database in ~/.word/words.db")
	flag.Parse()
	if *initPtr {
		db := initDB("~/./word/words.db")
		createTable(db)

	} else if *storePtr {
		var words []Word
		words = readVoc(VOCABULARYFILE)
		dbStore(words)
	} else {
		readDB()
	}
}
