package test

import (
	"reflect"
	"testing"
	"time"

	"github.com/tools/hrms/util"
)

func Test_SQLMap(t *testing.T) {
	var emp TestEntity
	emp.ID = "10001"
	emp.Age = 67
	emp.Salary = 10000.00
	emp.DOB = time.Now()
	name := "Sudip"
	emp.Name = &name
	workExp := 10
	emp.WorkExp = &workExp
	ls := float64(5000.00)
	emp.LastSalary = &ls
	lwd := time.Now().Add(100 * time.Hour)
	emp.LWD = &lwd
	sqlFieldMap := util.NewSQLFieldMap("hrm.employee_general_info", emp)
	sql, paramValues := sqlFieldMap.GenerateInsertScript(emp)
	t.Log(sql)
	t.Log(paramValues)
	for index, value := range paramValues {
		t.Logf("Postion %d %v", index+1, reflect.TypeOf(value))
	}
}

type TestEntity struct {
	ID         string     `json:"xid"`
	Name       *string    `json:"name"`
	Age        int        `json:"age"`
	WorkExp    *int       `json:"workexp"`
	Salary     float64    `json:"sal"`
	LastSalary *float64   `json:"lastsal"`
	DOB        time.Time  `json:"dob"`
	LWD        *time.Time `json:"lwd"`
}
