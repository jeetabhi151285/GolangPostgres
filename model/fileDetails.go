package model

import "ibm.com/MaerskTar/MaerskTarSdk2/util"

//"github.com/tools/hrms/util"

const _TestingTable = "odsschema.xSampleTestTable"

//PreBaplie represnts the odsschema.prebaplie_entry table entries
type PreBaplie struct {
	Vessel           string `json:"vessel"`
	CallSign         string `json:"callSign"`
	Voyage           string `json:"voyage"`
	Carrier          string `json:"carrier"`
	Operation        string `json:"operation"`
	CreationDate     string `json:"creationDate"`
	ContainerDetails []struct {
		Aaa string `json:"aaa"`
	} `json:"containerDetails"`
}

//GetInsertStatement returns insert sql statement
func (e *PreBaplie) GetInsertStatement() (string, []interface{}) {
	var prebaplieType PreBaplie
	fmap := util.NewSQLFieldMap(_TestingTable, prebaplieType)
	return fmap.GenerateInsertScript(e)
}
