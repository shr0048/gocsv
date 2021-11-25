//MIT License
//Copyright (c) 2020 HeraldSim

package gocsv

// CSV data converter as Go object type to support FP.

import (
	"bufio"
	"os"
	"strings"
)

const InvalidIntVal = -1

// Cell represents a single cell in the csv file.
type Cell map[string]string

// Record is a row of csv file.
type Record []Cell

// CSV is converted data of input csv file
type CSV struct {
	// index for iteration of Records
	index     int
	HeaderNum int
	RowNum    int
	Header    []string
	Records   []Record
}

// RecordMapper custom mapper function type.
type RecordMapper func(value Record) Record

// RecordFilter custom filter function type.
type RecordFilter func(value Record) bool

// RecordReducer custom reducer function type.
type RecordReducer func(a, b interface{}) interface{}

// ParseCSV csv parser function.
type ParseCSV func(string, string) []string

// SetParser return specific parser according to the csv file separator.
func SetParser(separator string) ParseCSV {
	if separator == "," {
		return func(r, s string) []string {
			res := ParseLine(strings.Split(r, s), s)
			return res
		}
	} else {
		return func(r, s string) []string {
			res := strings.Split(r, s)
			return res
		}
	}
}

// LoadCSV load csv file and convert each row and cell.
func (csv *CSV) LoadCSV(filePath string, separator string, startHeader int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	rowidx := 0

	myParser := SetParser(separator)

	for scanner.Scan() {
		if rowidx < startHeader {
			rowidx++
			continue
		} else if rowidx == startHeader {
			csv.Header = myParser(scanner.Text(), separator)
			csv.HeaderNum = len(csv.Header)
			rowidx++
			continue

		} else {
			tempRow := myParser(scanner.Text(), separator)
			var tempRecord Record
			for i := 0; i < csv.HeaderNum; i++ {
				cell := make(map[string]string)
				cell[csv.Header[i]] = tempRow[i]

				tempRecord = append(tempRecord, cell)
			}
			csv.Records = append(csv.Records, tempRecord)
			rowidx++
		}
	}
	csv.RowNum = rowidx - (startHeader + 1)
	csv.index = InvalidIntVal

	return nil
}

// NextRow iterate method of CSV
func (csv *CSV) NextRow() (record Record, ok bool) {
	csv.index++
	if csv.index >= csv.RowNum {
		return Record{}, false
	}
	return csv.Records[csv.index], true
}

// Map 'CSV's 'map' function of FP
func (csv *CSV) Map(mapper RecordMapper) *CSV {
	newRecords := make([]Record, 0, csv.RowNum)
	for _, r := range csv.Records {
		newRecords = append(newRecords, mapper(r))
	}

	return &CSV{
		index:     InvalidIntVal,
		HeaderNum: csv.HeaderNum,
		RowNum:    csv.RowNum,
		Header:    csv.Header,
		Records:   newRecords}
}

// Filter 'CSV's 'filter' function of FP
func (csv *CSV) Filter(filter RecordFilter) *CSV {
	newRecords := make([]Record, 0, csv.HeaderNum)
	rowNum := 0
	for _, v := range csv.Records {
		if filter(v) {
			newRecords = append(newRecords, v)
			rowNum++
		}
	}
	return &CSV{
		index:     InvalidIntVal,
		HeaderNum: csv.HeaderNum,
		RowNum:    rowNum,
		Header:    csv.Header,
		Records:   newRecords}
}

// Reduce 'CSV's 'Reduce' function of FP
func (csv *CSV) Reduce(identity interface{}, reducer RecordReducer) interface{} {
	res := identity
	for _, record := range csv.Records {
		res = reducer(res, record)
	}

	return res
}
