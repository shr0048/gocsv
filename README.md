[![Go Report Card](https://goreportcard.com/badge/github.com/shr0048/gocsv)](https://goreportcard.com/report/github.com/shr0048/gocsv)

# GoCsv

GoCsv is a library written in pure Go to use csv data more comfortable

### Supported Go version
- golang >= 1.13

### Installation

```bash
go get github.com/shr0048/gocsv
```

### Feature

- Freely set the separator (Ex: ',' '\t' ';' etc) 
- Support Collection model
- Include Map, Filter, Reduce
- All cell data is converting Go-Map


### Example
```go
package main

import (
	"fmt"
	"github.com/shr0048/gocsv"
	"strconv"
	"strings"
)

func convertCap(record gocsv.Record) gocsv.Record {
	newRecod := gocsv.Record{}
	for _, cell := range record {
		if val, ok := cell[`"Name"`]; ok {
			capVal := strings.ToUpper(val)
			newCell := gocsv.Cell{"Name": capVal}

			newRecod = append(newRecod, newCell)
		} else {
			newRecod = append(newRecod, cell)
		}
	}
	return newRecod
}

func filteringAge(record gocsv.Record) bool {
	for _, cell := range record {
		if val, ok := cell[`"Age"`]; ok {
			age, _ := strconv.ParseFloat(strings.TrimSpace(val), 64)
			if age >= 21 {
				return true
			} else {
				return false
			}
		}
	}
	// No 'Age' Column
	return false
}

func totalAge(a, b interface{}) interface{} {
	var age float64
	for _, cell := range b.(gocsv.Record) {
		if val, ok := cell[`"Age"`]; ok {
			age, _ = strconv.ParseFloat(strings.TrimSpace(val), 64)
		}
	}
	return a.(float64) + age
}

func main() {
	fmt.Println("##### Example of Using Collection functions #####")

	myCsv := gocsv.CSV{}
	err := myCsv.LoadCSV("./mlb_players.csv", ", ", 0)
	if err != nil {
		fmt.Print("Load error")
	}

	fmt.Println("Convert all Name as uppercase")
	// Apply Map
	mapCsv := myCsv.Map(convertCap)
	// Print all rows
	row, ok := mapCsv.NextRow()
	for ok {
		fmt.Println(row)
		row, ok = mapCsv.NextRow()
	}

	// Apply Filter and Reduce
	filtered := mapCsv.Filter(filteringAge)
	count := filtered.RowNum
	avg := filtered.Reduce(float64(0), totalAge)

	fmt.Printf("Get average age over 21: %f \n", avg.(float64)/float64(count))
}
```

### License
MIT License

Copyright (c) 2020 HeraldSim

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

### Contact
shr0048@protonmail.com