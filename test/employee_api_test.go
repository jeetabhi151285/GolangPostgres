package test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/tools/hrms/util"
)

func TestEmployeeExtnInfoAPI(t *testing.T) {
	apiTester := util.NewAPITester("http://localhost:7070/api", "", "", true)
	//Create a contactInfo
	contactInfo := []map[string]interface{}{
		map[string]interface{}{
			"relation": "Y",
			"name":     fmt.Sprintf("Y-%d", time.Now().Unix()),
		},
	}
	empCode := fmt.Sprintf("%d", time.Now().Unix())
	empextnInfo := map[string]interface{}{
		"empcode":     empCode,
		"contactinfo": contactInfo,
	}
	respCode, data := apiTester.PerformHTTPCall("PUT", "/employee/extninfo/contactinfo", empextnInfo)
	if respCode != 200 {
		t.FailNow()
	}
	x, isFound := data["isSuccess"]
	if !isFound || fmt.Sprintf("%v", x) != "true" {
		t.FailNow()
	}
	//Add one more input
	newContact := map[string]interface{}{
		"relation": "X",
		"name":     fmt.Sprintf("X-%d", time.Now().Unix()),
	}
	contactInfo = append(contactInfo, newContact)
	empextnInfo["contactinfo"] = contactInfo
	respCode, data = apiTester.PerformHTTPCall("PUT", "/employee/extninfo/contactinfo", empextnInfo)
	if respCode != 200 {
		t.FailNow()
	}
	x, isFound = data["isSuccess"]
	if !isFound || fmt.Sprintf("%v", x) != "true" {
		t.FailNow()
	}
	//Now Retrive the records inserted
	empDetails := map[string]interface{}{
		"empcode": empCode,
	}

	respCode, data = apiTester.PerformHTTPCall("POST", "/employee/extninfo/contactinfo", empDetails)
	if respCode != 200 {
		t.FailNow()
	}
	x, isFound = data["isSuccess"]
	if !isFound || fmt.Sprintf("%v", x) != "true" {
		t.FailNow()
	}
	payload, isFound := data["payload"]
	jsb, _ := json.MarshalIndent(payload, "", " ")
	t.Logf("Inserted records \n %s\n", string(jsb))
	//Remove all contact info
	empextnInfo["contactinfo"] = []interface{}{}
	respCode, data = apiTester.PerformHTTPCall("PUT", "/employee/extninfo/contactinfo", empextnInfo)
	if respCode != 200 {
		t.FailNow()
	}
	x, isFound = data["isSuccess"]
	if !isFound || fmt.Sprintf("%v", x) != "true" {
		t.FailNow()
	}
	respCode, data = apiTester.PerformHTTPCall("POST", "/employee/extninfo/contactinfo", empDetails)
	if respCode != 200 {
		t.FailNow()
	}
	x, isFound = data["isSuccess"]
	if !isFound || fmt.Sprintf("%v", x) != "true" {
		t.FailNow()
	}
	payload, isFound = data["payload"]
	jsb, _ = json.MarshalIndent(payload, "", " ")
	t.Logf("updated records \n %s\n", string(jsb))
}
