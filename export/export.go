package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const SPLITLINE string = "----split----"
const VOCABULARYFILE = "../vocabulary.txt"

type Word struct {
	name    string
	explain string
}

func readVoc(voc string) []Word {
	/*read voc file parse every word and explaination
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
	for n := 0; n < len(txtlines); n++ {
		if txtlines[n] == SPLITLINE {
			if n != 0 {
				words = append(words, word)
			}
			n++
			word.name = txtlines[n]
			continue
		}
		word.explain += txtlines[n]
		word.explain += "\n"
	}
	words = append(words, word)
	return words

}

func main() {
	var words []Word
	words = readVoc(VOCABULARYFILE)
	for _, word := range words {
		fmt.Printf("WORD NAME: %s\nWORD EXPLAIN:\n %s", word.name, word.explain)
	}
}
