package main

import (
	"fmt"
	"time"
)

func main() {
	nextDate, err := time.Parse("2006-01-02", "2020-01-07")
	if err != nil {
		panic(err)
	}
	fmt.Println(nextDate)
	fmt.Println(nextDate.AddDate(0, 0, 1))
	fmt.Println(nextDate.Format("2006-01-02"))
}
