package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/tools/hrms/model"
	"github.com/tools/hrms/util"
)

func TestGenerateModel(t *testing.T) {
	t.SkipNow()
	tableName := "hrm.employee_general_info"
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

	rslts, err := dbUtil.Query(fmt.Sprintf("select * from %s", tableName))
	if err != nil {
		t.FailNow()
	}
	defer rslts.Close()
	attributes := make([]string, 0)
	for rslts.Next() {

		for _, fdetails := range rslts.FieldDescriptions() {
			//fmt.Printf("%d-->[%d]%v\n", fpos, fdetails.DataTypeOID, string(fdetails.Name))
			jsonAttrName := string(fdetails.Name)
			attrName := fmt.Sprintf("%s%s", strings.ToUpper(jsonAttrName[0:1]), jsonAttrName[1:])
			dataType := ""
			switch fdetails.DataTypeOID {
			case 23:
				dataType = "int"
			case 1043:
				dataType = "string"
			case 1082:
				dataType = "time.Time"
			default:
				dataType = "unkown"

			}
			line := fmt.Sprintf("%s *%s `json:\"%s,omitempty\"`", attrName, dataType, jsonAttrName)
			attributes = append(attributes, line)
		}
		break

	}
	for _, attr := range attributes {
		fmt.Println(attr)
	}

	dbUtil.Shutdown()

}

func Test_ModelGenerator(t *testing.T) {
	//Generates an full obj json
	var obj model.EmployeeExtnInfo
	//var intr interface{}
	str := "temp"
	dt := time.Time{}
	typ := reflect.TypeOf(obj)
	//intrType := reflect.TypeOf(intr)
	t.Logf("schemas:")
	t.Logf("  %s:", typ.Name())
	t.Logf("   type: object")
	t.Logf("   properties:")
	for index := 0; index < typ.NumField(); index++ {
		field := typ.Field(index)
		fieldAttrName := field.Tag.Get("json")
		t.Logf(" %s:", strings.Split(fieldAttrName, ",")[0])
		switch field.Type {
		case reflect.TypeOf(str):
			t.Logf("     type: string")
		case reflect.TypeOf(&str):
			t.Logf("     type: string")
		case reflect.TypeOf(10):
			t.Logf("     type: integer")
		case reflect.TypeOf(time.Time{}):
			t.Logf("     type: string")
			t.Logf("     format: \"yyyy-mm-dd\"")
		case reflect.TypeOf(&dt):
			t.Logf("     type: string")
			t.Logf("     format: \"yyyy-mm-dd\"")
		default:
			t.Logf("     type: object")
			t.Logf("     description: \"Object description\"")
		}
	}

}
