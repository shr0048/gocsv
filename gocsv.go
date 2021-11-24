package gocsv

import (
	"bufio"
	"os"
	"strings"
)

const INVALID_INT_VAL = -1

// Cell
type Cell map[string]string

// Record
type Record []Cell

type RecordMapper func(value Record) Record
type RecordFilter func(value Record) bool
type RecordReducer func(a, b interface{}) interface{}
type ParseCSV func(string, string) []string

// CSV
type CSV struct {
	index     int
	HeaderNum int
	RowNum    int
	Header    []string
	Records   []Record
}

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

// LoadCSV
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
			rowidx ++
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
	csv.index = INVALID_INT_VAL

	return nil
}

func (csv *CSV) NextRow() (record Record, ok bool) {
	csv.index++
	if csv.index >= csv.RowNum {
		return Record{}, false
	}
	return csv.Records[csv.index], true
}

func (csv *CSV) Map(mapper RecordMapper) *CSV {
	newRecords := make([]Record, 0, csv.RowNum)
	for _, r := range csv.Records {
		newRecords = append(newRecords, mapper(r))
	}

	return &CSV{
		index:     INVALID_INT_VAL,
		HeaderNum: csv.HeaderNum,
		RowNum:    csv.RowNum,
		Header:    csv.Header,
		Records:   newRecords}
}

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
		index:     INVALID_INT_VAL,
		HeaderNum: csv.HeaderNum,
		RowNum:    rowNum,
		Header:    csv.Header,
		Records:   newRecords}
}

func (csv *CSV) Reduce(identity interface{}, reducer RecordReducer) interface{} {
	res := identity
	for _, record := range csv.Records {
		res = reducer(res, record)
	}

	return res
}
