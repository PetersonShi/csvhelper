package csvhelper

import (
	"fmt"
	"testing"
)

type UserAutoBind struct {
	Name    string
	Age     int
	Address string
	Hobby   string
	Gender  string
}

func newUserAutoBind() interface{} {
	return &UserAutoBind{}
}

type UserAutoBindWithTag struct {
	TestName    string `csv-field:"name"`
	TestAge     int    `csv-field:"age"`
	TestAddress string `csv-field:"address"`
	TestHobby   string `csv-field:"hobby"`
	TestGender  string `csv-field:"gender"`
}

func newUserAutoBindWithTag() interface{} {
	return &UserAutoBindWithTag{}
}

type UserManualBind struct {
	name    string
	age     int
	address string
	hobby   string
	gender  string
}

func (s *UserManualBind) CSVDataBind(data map[string]interface{}) {
	s.name = ToString(data["name"])
	s.age = ToInt(data["age"])
	s.address = ToString(data["address"])
	s.hobby = ToString(data["hobby"])
	s.gender = ToString(data["gender"])
}

func newStudentManualBind() interface{} {
	return &UserManualBind{}
}

func TestAll(t *testing.T) {

	/*
		1.加载csv文件
	*/

	//data model auto bind
	fmt.Println("---data model auto bind---")
	table := Load("./data_test.csv", TableTypeSlice, "name", newUserAutoBind)
	for _, v := range table.GetSliceData() {
		fmt.Println(v.(*UserAutoBind))
	}

	fmt.Println("---data model auto bind with tag---")
	table2 := Load("./data_test.csv", TableTypeSlice, "name", newUserAutoBindWithTag)
	for _, v := range table2.GetSliceData() {
		fmt.Println(v.(*UserAutoBindWithTag))
	}

	//data model manual bind
	fmt.Println("---data model manual bind---")
	table = Load("./data_test.csv", TableTypeSlice, "name", newStudentManualBind)
	for _, v := range table.GetSliceData() {
		fmt.Println(v.(*UserManualBind))
	}

	/*
		2.保存csv文件
	*/

	source := make([]map[string]string, 0)
	for i := 0; i < 100; i++ {
		item := make(map[string]string)
		item["book"] = "仙逆" + ToString(i)
		item["author"] = "耳根" + ToString(i)
		item["type"] = "古典仙侠" + ToString(i)
		source = append(source, item)
	}
	SaveCSVByData("./book.csv", source)
}
