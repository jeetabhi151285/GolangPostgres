package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var _typeCache map[string]*SQLFieldMap
var _sqlHelperLog *logrus.Logger

var _strPtrType = reflect.PtrTo(reflect.TypeOf(string("")))
var _intPtrType = reflect.PtrTo(reflect.TypeOf(int(0)))
var _floatPtrType = reflect.PtrTo(reflect.TypeOf(float64(0.0)))
var _dtPtrType = reflect.PtrTo(reflect.TypeOf(time.Now()))
var _intType = reflect.TypeOf(int(0))
var _floatType = reflect.TypeOf(float64(0))
var _strType = reflect.TypeOf(string(""))
var _dtType = reflect.TypeOf(time.Now())

//Datatype constants
const (
	INT uint32 = iota
	STR
	TIME
	FlOAT
	INTERFACE
)

func init() {
	_typeCache = make(map[string]*SQLFieldMap)
	_sqlHelperLog = logrus.New()
}

//SQLFieldMap Sql field mapping
type SQLFieldMap struct {
	TypeName  string
	TableName string
	FieldMap  map[string]uint32 //sql field name to data type

}

//NewSQLFieldMap creates a new SQLFieldMap . Returns from cache if already existing
func NewSQLFieldMap(tableName string, obj interface{}) *SQLFieldMap {
	typ := reflect.TypeOf(obj)
	_sqlHelperLog.Infof("Probing type %s", typ.Name())
	if sqlMap, isExisting := _typeCache[typ.Name()]; isExisting {
		return sqlMap
	}
	fmap := new(SQLFieldMap)
	fmap.TypeName = typ.Name()
	fmap.TableName = tableName
	fieldTypeMap := make(map[string]uint32)
	for index := 0; index < typ.NumField(); index++ {
		fieldDetails := typ.Field(index)
		jsonTagValue := fieldDetails.Tag.Get("json")
		if jsonTagValue != "" {
			_sqlHelperLog.Infof("Processing %s", jsonTagValue)
			tagParts := strings.Split(jsonTagValue, ",")
			switch fieldDetails.Type {
			case _strType:
				fallthrough
			case _strPtrType:
				fieldTypeMap[tagParts[0]] = STR
			case _intType:
				fallthrough
			case _intPtrType:
				fieldTypeMap[tagParts[0]] = INT
			case _floatType:
				fallthrough
			case _floatPtrType:
				fieldTypeMap[tagParts[0]] = FlOAT
			case _dtPtrType:
				fallthrough
			case _dtType:
				fieldTypeMap[tagParts[0]] = TIME
			default:
				_sqlHelperLog.Infof("Could not probe %v", fieldDetails.Type)
				fieldTypeMap[tagParts[0]] = INTERFACE

			}
			//TODO: Fix this
		}
	}
	fmap.FieldMap = fieldTypeMap
	_typeCache[typ.Name()] = fmap
	return fmap
}

//GenerateInsertScript table insert script
func (sfm *SQLFieldMap) GenerateInsertScript(obj interface{}) (string, []interface{}) {
	genericObj := make(map[string]interface{})
	jb, _ := json.Marshal(obj)
	json.Unmarshal(jb, &genericObj)
	var fieldNameBuffer bytes.Buffer
	var valueBuffer bytes.Buffer
	paramValues := make([]interface{}, 0)
	pos := 1
	fieldNameBuffer.WriteString(" ( ")
	valueBuffer.WriteString(" ( ")
	for key, value := range genericObj {
		if key != "_id" {
			if dataType, isExisting := sfm.FieldMap[key]; isExisting {
				fieldNameBuffer.WriteString(key)
				fieldNameBuffer.WriteString(",")
				valueBuffer.WriteString(fmt.Sprintf("$%d", pos))
				valueBuffer.WriteString(",")
				pos++
				_dbUtillog.Infof("Data type %s->%v", key, reflect.TypeOf(value))
				paramValues = append(paramValues, sfm.typeCast(value, dataType))
			}
		}
	}
	fieldNameBuffer.WriteString(") ")
	valueBuffer.WriteString(") ")
	fNames := fieldNameBuffer.String()
	posParams := valueBuffer.String()
	sql := fmt.Sprintf("insert into %s%s values%s", sfm.TableName, strings.Replace(fNames, ",)", ")", -1), strings.Replace(posParams, ",)", ")", -1))
	_sqlHelperLog.Infof("SQL Generated %s", sql)
	return sql, paramValues
}

//GenerateUpdateScript generates an update script
func (sfm *SQLFieldMap) GenerateUpdateScript(obj interface{}, filter map[string]interface{}) (string, []interface{}) {
	genericObj := make(map[string]interface{})
	jb, _ := json.Marshal(obj)
	json.Unmarshal(jb, &genericObj)
	var setFieldsBuf bytes.Buffer
	pos := 1
	setFieldsBuf.WriteString(" set ")
	paramValues := make([]interface{}, 0)
	isFirst := true
	for key, value := range genericObj {
		if key != "_id" {
			if dataType, isExisting := sfm.FieldMap[key]; isExisting {
				if _, isUPDField := filter[key]; !isUPDField {
					if isFirst {
						isFirst = false
						setFieldsBuf.WriteString(key)
					} else {
						setFieldsBuf.WriteString(", ")
						setFieldsBuf.WriteString(key)
					}
					setFieldsBuf.WriteString("= ")
					setFieldsBuf.WriteString(fmt.Sprintf("$%d ", pos))
					paramValues = append(paramValues, sfm.typeCast(value, dataType))
					pos++

				}
			}
		}
	}

	fNames := setFieldsBuf.String()
	var whereClause bytes.Buffer
	whereClause.WriteString(" where ")
	isFirst = true
	for key, value := range filter {
		if isFirst {
			isFirst = false
			whereClause.WriteString(fmt.Sprintf("%s=$%d", key, pos))
		} else {
			whereClause.WriteString(fmt.Sprintf("and %s=$%d", key, pos))
		}

		paramValues = append(paramValues, value)
		pos++
	}
	sql := fmt.Sprintf("update  %s %s %s", sfm.TableName, fNames, whereClause.String())
	_sqlHelperLog.Infof("Update SQL Generated %s", sql)

	return sql, paramValues
}

func (sfm *SQLFieldMap) typeCast(objValue interface{}, targetType uint32) interface{} {

	if targetType == INT {
		//Convert to an int
		//We need a convertion to integer as json.Unmarshal to an generic object
		//always put an int filed to float64
		intVal, _ := strconv.ParseInt(fmt.Sprintf("%v", objValue), 10, 64)
		return intVal
	}
	return objValue
}
