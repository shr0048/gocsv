package main

import (
	"fmt"
	"github.com/shr0048/gocsv"
)


func main() {
	myCsv := gocsv.CSV{}
	err := myCsv.LoadCSV("./meta.csv", "\t", 2)
	if err != nil {
		fmt.Print("Load error")
	}

	// Print all rows
	row, ok := myCsv.NextRow()
	for ok {
		fmt.Println(row)
		row, ok = myCsv.NextRow()
	}
}

