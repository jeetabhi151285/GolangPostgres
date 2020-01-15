package model

import (

	//"github.com/tools/hrms/util"

	"ibm.com/MaerskTar/MaerskTarSdk2/util"
)

const _PreBapliedataInfoTableName = "odsschema.prebaplie"

//PreBaplieDataInfo represents records of odsschema.prebaplie table
type PreBaplieDataInfo struct {
	Vessel           string      `json:"vessel"`
	CallSign         string      `json:"callSign"`
	Voyage           string      `json:"voyage"`
	Carrier          string      `json:"carrier"`
	Operation        string      `json:"operation"`
	CreationDate     string      `json:"creationDate"`
	ContainerDetails interface{} `json:"containerDetails"`
}

//GetInsertStatement returns insert sql statement
func (r *PreBaplieDataInfo) GetInsertStatement() (string, []interface{}) {
	var prebapdataType PreBaplieDataInfo
	fmap := util.NewSQLFieldMap(_PreBapliedataInfoTableName, prebapdataType)
	return fmap.GenerateInsertScript(r)
}

//GetUpdateStatement returns a update sql statement
// func (r *ReferenceDataInfo) GetUpdateStatement() (string, []interface{}) {
// 	var refdataType ReferenceDataInfo
// 	fmap := util.NewSQLFieldMap(_RefdataInfoTableName, refdataType)
// 	filter := map[string]interface{}{
// 		"category": r.Category,
// 		"value":    r.Value,
// 	}
// 	return fmap.GenerateUpdateScript(r, filter)
// }

// //GetTableName retrurns the table name
// func (r *ReferenceDataInfo) GetTableName() string {
// 	return _RefdataInfoTableName
// }

// //BuildRefDataMap converts db results to map of array
// func BuildRefDataMap(objs []interface{}) map[string][]ReferenceDataInfo {
// 	jb, _ := json.Marshal(objs)
// 	records := make([]ReferenceDataInfo, 0)
// 	err := json.Unmarshal(jb, &records)
// 	if err != nil {
// 		return nil
// 	}
// 	//Covert them into a map
// 	refDataMap := make(map[string][]ReferenceDataInfo)
// 	for _, refData := range records {
// 		categoty := refData.Category
// 		if arry, isCreated := refDataMap[categoty]; isCreated {
// 			arry = append(arry, refData)
// 			refDataMap[categoty] = arry
// 		} else {
// 			arry := make([]ReferenceDataInfo, 0)
// 			arry = append(arry, refData)
// 			refDataMap[categoty] = arry
// 		}
// 	}
// 	return refDataMap
// }
