package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/tools/hrms/model"
	"github.com/tools/hrms/util"
)

func Test_InsertEmpGenerel(t *testing.T) {
	t.Skip()
	conJSON := `{
		"dbhost":"localhost",
		"port":5432,
		"dnname":"erp",
		"uid":"hrmapp",
		"password":"cnp4test",
		"timeout": 10

	}`
	dbUtil, err := util.NewPGSqlDBUtil([]byte(conJSON), true)
	if err != nil {
		t.FailNow()

	}

	empID := fmt.Sprintf("AUTO-%d", time.Now().Unix())
	var empGenInfo model.PreBaplie
	empGenInfo.Empcode = empID
	empName := fmt.Sprintf("Name %s", empID)
	lname := fmt.Sprintf("LName %s", empID)
	empGenInfo.Empname = &empName
	empGenInfo.Lastname = &lname
	sql, params := empGenInfo.GetInsertStatement()
	//var id int
	err = dbUtil.InsertOne(sql, params, nil)
	if err != nil {
		t.FailNow()
	}
	//t.Logf("Generated ID %v\n", id)

	rslts, err := dbUtil.Query("select * from hrm.employee_general_info")
	if err != nil {
		t.FailNow()
	}
	defer rslts.Close()

	for rslts.Next() {

		if rowValues, err := rslts.Values(); err == nil {
			for colNumber, value := range rowValues {
				t.Logf("Col number %d=(%v)%v\n", colNumber, reflect.TypeOf(value), value)
			}
		}
	}

	dbUtil.Shutdown()

}
func Test_CheckJSON(t *testing.T) {
	conJSON := `{
		"dbhost":"localhost",
		"port":5432,
		"dnname":"erp",
		"uid":"hrmapp",
		"password":"cnp4test",
		"timeout": 10

	}`
	dbUtil, err := util.NewPGSqlDBUtil([]byte(conJSON), true)
	if err != nil {
		t.FailNow()

	}

	sql := `INSERT INTO hrm.mytable(
		category, empid, extra)
		VALUES ($1, $2, $3)
		`
	params := make([]interface{}, 0)
	params = append(params, fmt.Sprintf("AUTO-%d", time.Now().Unix()))
	params = append(params, "0001")
	var jsnObj interface{}
	jsnObj = []map[string]interface{}{map[string]interface{}{
		"x": "Y",
		"z": fmt.Sprintf("%d", time.Now().Unix()),
	}}
	params = append(params, jsnObj)

	err = dbUtil.InsertOne(sql, params, nil)
	if err != nil {
		t.FailNow()
	}
	//t.Logf("Generated ID %v\n", id)

	rslts, err := dbUtil.Query("select * from hrm.mytable")
	if err != nil {
		t.FailNow()
	}
	defer rslts.Close()

	for rslts.Next() {

		if rowValues, err := rslts.Values(); err == nil {
			for colNumber, value := range rowValues {
				t.Logf("Col number %d=(%v)%v\n", colNumber, reflect.TypeOf(value), value)
			}
		}
	}

	dbUtil.Shutdown()

}
