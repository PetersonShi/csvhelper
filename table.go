package csvhelper

import "sync"

const (
	TableTypeMap = iota
	TableTypeSlice
)

type DataModel interface {
	CSVDataBind(raw map[string]interface{})
}

type CSVTable interface {
	GetDataItem(interface{}) interface{}
	GetMapData() map[interface{}]interface{}
	GetSliceData() []interface{}
	SetKeyField(string)
	GetKeyField() string
	bind([]map[string]interface{})
}

type tableBase struct {
	//mutex
	dataMu sync.Mutex
	//keyField field
	keyField string
	//new data item maker
	newData func() interface{}
}

type mapTable struct {
	tableBase
	//data container
	dataMap map[interface{}]interface{}
}

func (mt *mapTable) bind(source []map[string]interface{}) {
	mt.dataMu.Lock()
	defer mt.dataMu.Unlock()

	mt.dataMap = make(map[interface{}]interface{})
	for i, _ := range source {
		rowData := mt.newData()
		if _, ok := rowData.(DataModel); ok {
			rowData.(DataModel).CSVDataBind(source[i])
		} else {
			StructBind(rowData, source[i])
		}
		mt.addRow(rowData)
	}
}

func (mt *mapTable) addRow(rowData interface{}) {
	v := valueElem(rowData)

	if _, found := v.Type().FieldByName(mt.keyField); found == false {
		panic("DataTable addRow,keyField not found. keyField=" + mt.keyField)
	}

	valueField := v.FieldByName(mt.keyField)
	mt.dataMap[valueField.Interface()] = rowData
}

func (mt *mapTable) SetKeyField(keyField string) {
	keyField = ReorganizeKeyField(keyField)
	if keyField == "" {
		keyField = "Id"
	}
	if mt.keyField != keyField {
		mt.keyField = keyField
		mt.rebuild()
	}
}

func (mt *mapTable) GetKeyField() string {
	return mt.keyField
}

func (mt *mapTable) rebuild() {
	mt.dataMu.Lock()
	defer mt.dataMu.Unlock()

	tempRows := mt.dataMap
	mt.dataMap = make(map[interface{}]interface{})
	for _, row := range tempRows {
		mt.addRow(row)
	}
}

func (mt *mapTable) GetDataItem(id interface{}) interface{} {
	mt.dataMu.Lock()
	defer mt.dataMu.Unlock()
	return mt.dataMap[id]
}

func (mt *mapTable) GetMapData() map[interface{}]interface{} {
	mt.dataMu.Lock()
	defer mt.dataMu.Unlock()
	return mt.dataMap
}

func (mt *mapTable) GetSliceData() []interface{} { return nil }

type sliceTable struct {
	tableBase
	//data container
	dataSlice []interface{}
}

func (st *sliceTable) bind(source []map[string]interface{}) {
	st.dataMu.Lock()
	defer st.dataMu.Unlock()

	st.dataSlice = make([]interface{}, len(source))
	for i, _ := range source {
		rowData := st.newData()
		if _, ok := rowData.(DataModel); ok {
			rowData.(DataModel).CSVDataBind(source[i])
		} else {
			StructBind(rowData, source[i])
		}
		st.dataSlice[i] = rowData
	}
}

func (st *sliceTable) SetKeyField(keyField string) {
	keyField = ReorganizeKeyField(keyField)
	if keyField == "" {
		keyField = "Id"
	}
	st.keyField = keyField
}

func (st *sliceTable) GetKeyField() string {
	return st.keyField
}

func (st *sliceTable) GetDataItem(id interface{}) interface{} {
	if id == nil {
		return nil
	}

	st.dataMu.Lock()
	defer st.dataMu.Unlock()
	for _, dataRow := range st.dataSlice {
		v := valueElem(dataRow)

		if _, found := v.Type().FieldByName(st.keyField); found == false {
			return nil
		}

		if id == v.FieldByName(st.keyField).Interface() {
			return dataRow
		}
	}
	return nil
}

func (st *sliceTable) GetMapData() map[interface{}]interface{} { return nil }

func (st *sliceTable) GetSliceData() []interface{} {
	st.dataMu.Lock()
	defer st.dataMu.Unlock()
	return st.dataSlice
}

func newMapTable(newData func() interface{}) *mapTable {
	mapTable := &mapTable{tableBase: tableBase{newData: newData}, dataMap: make(map[interface{}]interface{})}
	mapTable.SetKeyField("id")
	return mapTable
}

func newSliceTable(newData func() interface{}) *sliceTable {
	sliceTable := &sliceTable{tableBase: tableBase{newData: newData}, dataSlice: make([]interface{}, 0)}
	sliceTable.SetKeyField("id")
	return sliceTable
}
