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

//var fib = [14]string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "144", "233"}
var fib = [14]int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}
var voc string

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
    PRIMARY KEY(word, nextreviewdate)
	);`
	_, err := db.Exec(sql_table)
	if err != nil {
		fmt.Println("Table already exist.")
		return
	} else {
		fmt.Println("Table words created!")
		return
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

func creatVoc(voc string) {
	emptyFile, err := os.Create(voc)
	if err != nil {
		fmt.Printf("Oops, can not creat %s", voc)
		panic(err)
	}
	emptyFile.Close()
	fmt.Printf("%s created!\n", voc)
}

func removeVoc(voc string) {
	err := os.Remove(voc)
	if err != nil {
		fmt.Printf("Opps!!! Cannot delete file: %s\n", voc)
	} else {
		fmt.Printf("File: %s removed.\n", voc)
	}
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

func dbStore(db *sql.DB, words []Word) {
	/* store words list to word.db */
	for _, word := range words {

		if checkRecord(word.name) {
			fmt.Printf("Word: %s already exist\n", word.name)
			continue
		}

		stmt, err := db.Prepare("INSERT INTO words(word, translation, createdate, nextreviewdate, reviewstatus) values(?,?,?,?,?)")
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
		removeVoc(voc)
	}
}

func openDB(dfFullPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbFullPath)
	if err != nil {
		panic(err)
	}
	return db
}

func readDB(db *sql.DB) []WordTableRow {
	/* read words from table word of current date */

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

func updateTable(db *sql.DB, rec WordTableRow) {
	rec.reviewStatus += 1
	nextDay := time.Now().AddDate(0, 0, fib[rec.reviewStatus])
	fmt.Println(nextDay)
	stmt, err := db.Prepare(`UPDATE words SET nextreviewdate = ?, reviewstatus = ?  WHERE word = ?`)
	if err != nil {
		fmt.Println("Update Prepare Error")
		panic(err)
	}
	_, err = stmt.Exec(nextDay.Format("2006-01-02"), rec.reviewStatus, rec.word)
	if err != nil {
		fmt.Printf("Can not update %s's nextreview and reviewstatus", rec.word)
	}
	return
}

func review(db *sql.DB, wordList []WordTableRow) {
	for _, rec := range wordList {
		fmt.Printf("New Word: %s | %s\n\n", Red(rec.word), Red(strings.ToUpper(rec.word)))
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
				fmt.Printf("Pass to the next Round\n")
				updateTable(db, rec)

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

	homeFullPath := os.Getenv("HOME") + "/.word"
	dbFullPath = homeFullPath + "/words.db"
	voc = homeFullPath + "/vocabulary.txt"

	if *initPtr {
		/* Fist time use, build a new words.db in ~/.word/ */
		err := os.MkdirAll(homeFullPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Can not create directory: %s ", homeFullPath)
		}

		db := initDB(dbFullPath)
		createTable(db)
		creatVoc(voc)

	} else if *storePtr {
		/* store all the vocabulary from voc.txt to database */
		var words []Word
		words = readVoc(voc)
		db := openDB(dbFullPath)
		dbStore(db, words)

	} else {
		/* Review voc */
		db := openDB(dbFullPath)
		wordList := readDB(db)
		review(db, wordList)

	}
}
