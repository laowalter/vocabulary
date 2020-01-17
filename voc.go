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

var dbFullPath string                                                //sqlite database
var fib = [14]int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233} //review statergy
var voc string                                                       //vocabulary.txt file

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
		fmt.Printf("Vocabulary file %s was not created by now, user %s or %s to create it.\n", Cyan(voc), Red("translate"), Red("newtrans"))
		os.Exit(1)
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

	var lines []string
	for n := 0; n < len(txtlines); n++ { //Purge multiple continued ----split----
		if n+1 < len(txtlines) {
			if txtlines[n] == SPLITLINE && txtlines[n+1] == SPLITLINE {
				continue
			} else {
				lines = append(lines, txtlines[n])
			}
		}
	}

	for n := 0; n < len(lines); n++ { //Parse
		if lines[n] == SPLITLINE {
			if tag {
				words = append(words, word)
				word.trans = ""
				tag = false
			}
			n++
			word.name = lines[n]
			tag = true
			continue
		}
		word.trans += lines[n]
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

func storeDB(db *sql.DB, words []Word) {
	/* store words list to word.db */
	for _, word := range words {
		if checkRecord(word.name) {
			fmt.Printf("Word: %s already exist\n", word.name)
			resetWord(db, word.name)
			fmt.Printf("Word: %s is %s as today's new word.\n", Cyan(word.name), Red("Reset"))
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
	}
	removeVoc(voc)
	return
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

func totalWordsDB(db *sql.DB) int {
	/* total words from table word of current date */
	rows, err := db.Query("select * from words")
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
	return len(wordRecords)
}

func updateNextReviewDate(db *sql.DB, rec WordTableRow) {
	rec.reviewStatus += 1 //enter next reivew series
	nextDay := time.Now().AddDate(0, 0, fib[rec.reviewStatus])
	fmt.Printf("Will review on %s\n", nextDay.Format("2006-01-02"))
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

func resetWord(db *sql.DB, word string) {
	reviewStatus := 0                       //reset to the first review status
	date := time.Now().Format("2006-01-02") // reset the nextreviewdate to today
	stmt, err := db.Prepare(`UPDATE words SET nextreviewdate = ?, reviewstatus = ?  WHERE word = ?`)
	if err != nil {
		fmt.Println("Update Prepare Error")
		panic(err)
	}
	_, err = stmt.Exec(date, reviewStatus, word)
	if err != nil {
		fmt.Printf("Can not update %s's nextreview and reviewstatus.\n", word)
		return
	}
	fmt.Printf("Word %s was reset as a NEW %s.\n", Red(word), Cyan("word"))
}

func modifyWordRecord(db *sql.DB, rec WordTableRow) {
	/* Only change a word in DB, for example change pl gaffes to gaffe */
	stmt, err := db.Prepare(`UPDATE words SET word = ? WHERE word = ?`)
	if err != nil {
		fmt.Printf("Can not prepare modify the word %s\n", Cyan(rec.word))
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please input the new word for %s: ", rec.word)
	newWord, _ := reader.ReadString('\n')
	newWord = strings.TrimSuffix(newWord, "\n")
	newWord = strings.Trim(newWord, " ")

	if newWord == "" {
		fmt.Println("You did not input anything.")
		return
	}

	_, err = stmt.Exec(newWord, rec.word)
	if err != nil {
		fmt.Printf("Can not modify the word %s \n", Cyan(rec.word))
		panic(err)
	}

	fmt.Printf("The word %s replaced by %s \n", Cyan(rec.word), Red(newWord))
	return
}

func deleteRecord(db *sql.DB, rec WordTableRow) {
	stmt, err := db.Prepare(`DELETE FROM words WHERE word == ?`)
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(rec.word)
	if err != nil {
		fmt.Printf("Can not delete record: %s", rec.word)
		panic(err)
	}
	fmt.Printf("Word: %s removed from remember database\n", Cyan(rec.word))
	return
}

func review(db *sql.DB, wordList []WordTableRow) {
	wordsLength := len(wordList)

	if wordsLength == 0 {
		return
	}

	index := 0
	for {

		fmt.Printf("\n(%d/%d): %s | %s\n\n", Cyan(index+1), Cyan(wordsLength), Red(wordList[index].word), Red(strings.ToUpper(wordList[index].word)))
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		if char == '\x00' {
			fmt.Printf("%s\n-----------------\n\n", Green(wordList[index].trans))
			for {
				char, _, err = keyboard.GetSingleKey()
				if err != nil {
					panic(err)
				}
				if char == 'p' { // pass after trans displayed
					updateNextReviewDate(db, wordList[index]) //change nextReviewDate
					if index >= wordsLength-1 {
						return
					} else {
						index += 1
						break
					}
				} else if char == 'd' { //delete the word from db
					deleteRecord(db, wordList[index])
					if index >= wordsLength-1 {
						return
					} else {
						index += 1
					}
				} else if char == 'r' { //reset the old word to new word
					resetWord(db, wordList[index].word)
					if index >= wordsLength-1 {
						return
					} else {
						index += 1
					}

				} else if char == 'm' { //modify word's spell.
					modifyWordRecord(db, wordList[index])
					if index >= wordsLength-1 {
						return
					} else {
						index += 1
					}
				} else if char == '\x00' {
					if index >= wordsLength-1 {
						return
					} else {
						index += 1
						break
					}
				} else if char == 'q' { // exit at any time
					return
				}
				break
			}

		} else if char == 'p' { // pass before trans displayed
			updateNextReviewDate(db, wordList[index]) //change nextReviewDate
			if index >= wordsLength-1 {
				return
			} else {
				index += 1
			}

		} else if char == 'm' { // modified the current word after a the word displayed
			modifyWordRecord(db, wordList[index])
			if index >= wordsLength-1 {
				return
			} else {
				index += 1
			}
		} else if char == 'd' { // delete the current word (after the word displayed)
			deleteRecord(db, wordList[index])
			if index >= wordsLength-1 {
				return
			} else {
				index += 1
			}

		} else if char == 'r' { //reset the old word to new word
			resetWord(db, wordList[index].word)
			if index >= wordsLength-1 {
				return
			} else {
				index += 1
			}
		} else if char == 'q' { // exit
			return
		}
	}
}

func main() {
	storePtr := flag.Bool("store", false, "Store new words to Database")
	listPtr := flag.Bool("list", false, "List words in ~/.words/vocabulary.txt ")
	initPtr := flag.Bool("init", false, "Init Local database in ~/.word/words.db")
	totalPtr := flag.Bool("total", false, "Return the total number of words in  database in ~/.word/words.db")
	flag.Parse()

	homeFullPath := os.Getenv("HOME") + "/.word"
	dbFullPath = homeFullPath + "/words.db"
	voc = homeFullPath + "/vocabulary.txt"

	if *initPtr { // Fist time use, build a new words.db in ~/.word
		err := os.MkdirAll(homeFullPath, os.ModePerm)
		if err != nil {
			fmt.Printf("Can not create directory: %s ", homeFullPath)
		}
		db := initDB(dbFullPath)
		createTable(db)
		creatVoc(voc)

	} else if *storePtr { // store all the vocabulary from voc.txt to database
		words := readVoc(voc)
		db := openDB(dbFullPath)
		storeDB(db, words)

	} else if *listPtr {
		words := readVoc(voc)
		for index, word := range words {
			fmt.Printf("Index: %2d, Word: %s\n", index, word.name)
		}
	} else if *totalPtr {
		db := openDB(dbFullPath)
		totalNumber := totalWordsDB(db)
		fmt.Printf("Total numbers of words in DB is  %d\n", Red(totalNumber))
	}
}
