package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"github.com/eiannone/keyboard"
	. "github.com/logrusorgru/aurora"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"time"
)

const (
	SPLITLINE string = "----split----"
)

var dbFullPath string
var fib = [13]int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}

type Word struct {
	name  string
	trans string
}

type WordTableRow struct {
	word           string
	trans          string
	createDate     string
	nextReviewDate string
	reviewStatus   int
}

func initDB(filePath string) *sql.DB {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		fmt.Printf("Oops! Database %s already exist", filePath)
		return db
	}
	if db == nil {
		panic("db nil")
	}
	fmt.Printf("Database %s created!\n", filePath)
	return db
}

func createTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
	CREATE TABLE words (
    word text NOT NULL,
    translation text NOT NULL,
    createdate text DEFAULT (STRFTIME('%Y-%m-%d', 'NOW')),
	nextreviewdate text DEFAULT (STRFTIME('%Y-%m-%d', 'NOW')),
    reviewstatus INT DEFAULT 0,
    PRIMARY KEY(word, lastreviewdate)
	);`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println("Table already exist.")
		return
	}
	fmt.Println("Table words created!")
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
	err = os.Remove(voc)
	if err != nil {
		fmt.Printf("Opps!!! Cannot delete file: %s\n", voc)
	} else {
		fmt.Printf("File: %s removed.\n", voc)
	}
	return words
}

func checkRecord(word string) bool {
	/* check if a word.name exist in dbFullPath */
	db, err := sql.Open("sqlite3", dbFullPath)
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
	db, err := sql.Open("sqlite3", dbFullPath)
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

func readDB() []WordTableRow {
	/* read words from table word of current date */
	db, err := sql.Open("sqlite3", dbFullPath)
	if err != nil {
		panic(err)
	}

	date := time.Now().Format("2006-01-02")
	rows, err := db.Query("select * from words where nextreviewdate <= ?", date)
	if err != nil {
		panic(err)
	}

	var wordRecords []WordTableRow
	for rows.Next() {
		var record WordTableRow
		err = rows.Scan(&record.word, &record.trans, &record.createDate, &record.nextReviewDate, &record.reviewStatus)
		if err != nil {
			panic(err)
		}
		wordRecords = append(wordRecords, record)
	}
	return wordRecords
}

func review(wordList []WordTableRow) {
	for _, rec := range wordList {
		fmt.Printf("New Word: %s  %s\n\n", Red(rec.word), Red(strings.ToUpper(rec.word)))
		//fmt.Print("Press SPACE key for translation\n")
		for {
			char, _, err := keyboard.GetSingleKey()
			if err != nil {
				panic(err)
			}
			if char == '\x00' {
				fmt.Printf("%s\n-----------------\n\n", Green(rec.trans))
				break
			} else if char == 'p' {
				/* change nextReviewDate */
				fmt.Printf("Will change nextReviewDate in DB\n")

			} else if char == 'q' {
				os.Exit(1)
			}
		}
	}

}

func main() {
	storePtr := flag.Bool("store", false, "Store new words to Database")
	initPtr := flag.Bool("init", false, "Init Local database in ~/.word/words.db")
	flag.Parse()

	homeFullPath := os.Getenv("HOME")
	dbFullPath = homeFullPath + "/.word/words.db"

	if *initPtr {
		/* Fist time use, build a new words.db in ~/.word/ */
		err := os.MkdirAll(homeFullPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Can not create directory: %s ", homeFullPath)
		}

		db := initDB(dbFullPath)
		createTable(db)

	} else if *storePtr {
		/* store all the vocabulary from voc.txt to database */
		var words []Word
		words = readVoc(homeFullPath + "/.word/vocabulary.txt")
		dbStore(words)

	} else {
		/* Review voc */
		wordList := readDB()
		review(wordList)

	}
}
