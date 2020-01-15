package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/tools/hrms/model"
	"github.com/tools/hrms/util"
)

const connectionJSON = `{
	"dbhost":"localhost",
	"port":5432,
	"dnname":"erp",
	"uid":"hrmapp",
	"password":"cnp4test",
	"timeout": 10

}`

func TestPGInsertAndUpdate(t *testing.T) {
	t.SkipNow()
	dbUtil, err := util.NewPGSqlDBUtil([]byte(connectionJSON), true)
	if err != nil {
		t.FailNow()

	}
	defer dbUtil.Shutdown()
	ts := fmt.Sprintf("%d", time.Now().Unix())
	empJSON := map[string]string{
		"empcode":           ts,
		"empname":           fmt.Sprintf("NAME AUTO %s", ts),
		"dob":               "2001-10-11",
		"originalbirthdate": "2001-06-12",
	}
	jb, _ := json.Marshal(empJSON)
	var employee model.PreBaplie
	json.Unmarshal(jb, &employee)
	t.Log(employee.Dob)
	t.Log(employee.Originalbirthdate)
	sql, params := employee.GetInsertStatement()
	err = dbUtil.InsertOne(sql, params, nil)
	if err != nil {
		t.FailNow()
	}
	//Now update
	newDob := time.Now()
	employee.Dob = &newDob
	sql, params = employee.GetUpdateStatement()
	count, err := dbUtil.UpdateRecords(sql, params)
	if err != nil || count == 0 {
		t.FailNow()
	}
}
func TestPGConnectionBasic(t *testing.T) {
	//t.SkipNow()
	dbUtil, err := util.NewPGSqlDBUtil([]byte(connectionJSON), true)
	if err != nil {
		t.FailNow()

	}
	//empCode := pgtype.Text{}
	//empCode.Set("100001")
	empCode := "100001"
	rslts, err := dbUtil.QueryRecords("select * from hrm.employee_general_info where empcode = $1 ", empCode)
	if err != nil {
		t.FailNow()
	}
	emps := model.BuildEmployeGenInfoArray(rslts)
	jsb, _ := json.MarshalIndent(emps, " ", " ")
	t.Logf("%s", string(jsb))

	dbUtil.Shutdown()

}
