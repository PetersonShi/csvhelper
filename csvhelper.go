package csvhelper

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func Load(filePath string, tableType int, keyField string, newRow func() interface{}) CSVTable {
	/*
		Create CSVTable
	*/
	var table CSVTable = nil
	if tableType == TableTypeMap {
		table = newMapTable(newRow)
	} else {
		table = newSliceTable(newRow)
	}
	table.SetKeyField(keyField)

	/*
		CSVDataBind csv content
	*/
	file, err := os.Open(filePath)
	if err != nil {
		panic("Open csv file error,path=" + filePath)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		panic("CSV read error" + err.Error())
	}

	//records size check
	if len(records) == 0 || len(records[0]) == 0 {
		return nil
	}

	//records structure
	//records[0]:the csv column name
	//records[1:]:the csv real data

	//1.check keyField exists
	keyFieldFound := false
	for i := 0; i < len(records[0]); i++ {
		field := records[0][i]
		if table.GetKeyField() == ReorganizeKeyField(field) {
			keyFieldFound = true
			break
		}
	}
	//*Not found keyField in columns,Set first column as keyField
	if keyFieldFound == false {
		table.SetKeyField(records[0][0])
	}

	//2.bind data
	csvData := make([]map[string]interface{}, 0)
	for i := 1; i < len(records); i++ {
		csvRow := make(map[string]interface{})
		for j := 0; j < len(records[0]); j++ {
			key := records[0][j]
			csvRow[key] = records[i][j]
		}
		csvData = append(csvData, csvRow)
	}
	table.bind(csvData)
	return table
}

func SaveCSVByTable(filename string, source CSVTable) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer f.Close()

	fileContent := ""
	if _, ok := source.(*mapTable); ok {
		dataMap := source.GetMapData()
		titleAdded := false
		for _, data := range dataMap {
			if titleAdded == false {
				titleAdded = true
				fileContent += makeCSVTitle(data)
			}
			fileContent += makeCSVContent(data)
		}
	} else {
		dataSlice := source.GetSliceData()
		titleAdded := false
		for _, data := range dataSlice {
			if titleAdded == false {
				titleAdded = true
				fileContent += makeCSVTitle(data)
			}
			fileContent += makeCSVContent(data)
		}
	}

	f.WriteString(fileContent)
}

type saveValueType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 | ~string
}

func SaveCSVByData[T saveValueType](filename string, source []map[string]T) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer f.Close()

	fileContent := ""

	csvfields := make([]string, 0)
	for _, data := range source {
		if len(csvfields) == 0 {
			for key, _ := range data {
				csvfields = append(csvfields, key)
				fileContent += strings.ToLower(key) + ","
			}
			fileContent = strings.TrimRight(fileContent, ",")
			fileContent += "\n"
		}

		for i := 0; i < len(data); i++ {
			field := csvfields[i]
			fileContent += ToString(data[field]) + ","
		}
		fileContent = strings.TrimRight(fileContent, ",")
		fileContent += "\n"
	}

	f.WriteString(fileContent)
}
