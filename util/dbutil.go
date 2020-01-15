package util

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	//"ibm.com/MaerskTar/MaerskTarSdk2/jackc/pgx/v4"
	//"ibm.com/MaerskTar/MaerskTarSdk2/sirupsen/logrus"
	//"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

var _dbUtillog = logrus.New()

//PGSqlDBUtil type implements a postgressql db driver wrapper
type PGSqlDBUtil struct {
	connectionStr string
	dbConnection  *pgx.Conn
	verbose       bool
	mutex         *sync.Mutex
}

type dbConfig struct {
	DBHost   string `json:"dbhost"`
	DBName   string `json:"dnname"`
	UserID   string `json:"uid"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
	Port     int    `json:"port"`
}

//NewPGSqlDBUtil returns a new instance of PGSQLUtil
func NewPGSqlDBUtil(configBytes []byte, verbose bool) (*PGSqlDBUtil, error) {
	_dbUtillog.Info("Creating DBUtil...")
	dbUtil := new(PGSqlDBUtil)
	if err := dbUtil.Init(configBytes); err != nil {
		return nil, err
	}
	if verbose {
		_dbUtillog.SetLevel(logrus.DebugLevel)
	}
	return dbUtil, nil
}

//Init intializes the connection to db
func (dbu *PGSqlDBUtil) Init(configBytes []byte) error {
	var config dbConfig
	err := json.Unmarshal(configBytes, &config)
	if err != nil {
		_dbUtillog.Errorf("Unable to parse the configuration %v", err)
		return err
	}
	dbu.connectionStr = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", config.UserID, config.Password, config.DBHost, config.Port, config.DBName)
	dbu.dbConnection, err = pgx.Connect(context.Background(), dbu.connectionStr)
	if err != nil {
		_dbUtillog.Errorf("Unable to connection to database: %v\n", err)
		return err
	}
	dbu.mutex = new(sync.Mutex)
	go dbu.keepAlive()
	return nil
}

func (dbu *PGSqlDBUtil) keepAlive() {
	for {
		<-time.After(5 * time.Minute)
		isSuccess := dbu.refreshConnection()
		_dbUtillog.Infof("Refresh connection %v", isSuccess)
	}
}

func (dbu *PGSqlDBUtil) refreshConnection() bool {
	dbu.mutex.Lock()
	defer dbu.mutex.Unlock()
	if dbu.dbConnection.Ping(context.Background()) != nil {
		_dbUtillog.Infof("Reconnecting....")

		dbConnection, err := pgx.Connect(context.Background(), dbu.connectionStr)
		if err != nil {
			_dbUtillog.Errorf("Unable to connection to database: %v\n", err)
			return false
		}
		dbu.dbConnection = dbConnection
	}
	return true
}

//Query runs a query in the db and reurns pgx.Rows
func (dbu *PGSqlDBUtil) Query(sql string, params ...interface{}) (pgx.Rows, error) {
	dbu.refreshConnection()
	results, err := dbu.dbConnection.Query(context.Background(), sql, params...)
	if err != nil {
		_dbUtillog.Errorf("Error in executing the query %v", err)
		return nil, err
	}
	_dbUtillog.Infof("Execution result %v", results)
	return results, nil
}

//QueryRecords runs a query in the db and converts them []interface{}
func (dbu *PGSqlDBUtil) QueryRecords(sql string, params ...interface{}) ([]interface{}, error) {
	dbu.refreshConnection()
	rslts, err := dbu.dbConnection.Query(context.Background(), sql, params...)
	if err != nil {
		_dbUtillog.Errorf("Error in executing the query %v", err)
		return nil, err
	}
	defer rslts.Close()
	rows := make([]interface{}, 0)
	var colMap map[int]string
	for rslts.Next() {
		objMap := make(map[string]interface{})

		if colMap == nil {
			colMap = dbu.getColMap(rslts)
		}
		if rowValues, err := rslts.Values(); err == nil {
			for colNumber, value := range rowValues {
				attrName, _ := colMap[colNumber]
				if value != nil {
					objMap[attrName] = value
				}
				if dbu.verbose {
					_dbUtillog.Debugf("Col number %d=(%v)%v\n", colNumber, reflect.TypeOf(value), value)
				}
			}
		}

		rows = append(rows, objMap)
	}

	return rows, err
}

//InsertOne inserts one record in the db
func (dbu *PGSqlDBUtil) InsertOne(sql string, params []interface{}, generatedID interface{}) error {
	dbu.refreshConnection()
	tx, bErr := dbu.dbConnection.Begin(context.Background())
	if bErr != nil {
		_dbUtillog.Errorf("Error in creating trxn %v", bErr)
	}
	//rslt, err := dbu.dbConnection.Exec(context.Background(), sql, params...)
	rslt, err := tx.Exec(context.Background(), sql, params...)
	if err != nil {
		_dbUtillog.Errorf("Error in executing the transaction %v", err)
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil {
			_dbUtillog.Errorf("Unable to roll back %v", rbErr)
		}
		return err
	}
	_dbUtillog.Infof("Execution result %v", rslt)
	if generatedID != nil {
		err = dbu.dbConnection.QueryRow(context.Background(), "SELECT LASTVAL()").Scan(generatedID)
		if err != nil {
			_dbUtillog.Errorf("Error in executing the uids %v", err)
			return err
		}
		_dbUtillog.Infof("Execution generated ids  %v", generatedID)
	}
	if err := tx.Commit(context.Background()); err != nil {
		err = tx.Rollback(context.Background())
		_dbUtillog.Errorf("Unable to roll back %v", err)
	}

	return nil
}

//UpdateRecords inserts one record in the db
func (dbu *PGSqlDBUtil) UpdateRecords(sql string, params []interface{}) (int64, error) {
	dbu.refreshConnection()
	tx, bErr := dbu.dbConnection.Begin(context.Background())
	if bErr != nil {
		_dbUtillog.Errorf("Error in creating trxn %v", bErr)
	}
	rslt, err := tx.Exec(context.Background(), sql, params...)
	if err != nil {
		_dbUtillog.Errorf("Error in executing the transaction %v", err)
		return 0, err
	}
	_dbUtillog.Infof("Execution result %v", rslt)
	if err := tx.Commit(context.Background()); err != nil {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil {
			_dbUtillog.Errorf("Unable to roll back %v", rbErr)
		}
		return 0, err
	}

	return rslt.RowsAffected(), nil
}

//InsertMultiple inserts multple records
func (dbu *PGSqlDBUtil) InsertMultiple(sqls []string, params [][]interface{}) (int64, error) {
	dbu.refreshConnection()
	tx, bErr := dbu.dbConnection.Begin(context.Background())
	if bErr != nil {
		_dbUtillog.Errorf("Error in creating trxn %v", bErr)
	}
	count := int64(0)
	for index, sql := range sqls {
		rslt, err := tx.Exec(context.Background(), sql, params[index]...)
		if err != nil {
			_dbUtillog.Errorf("Error in executing the transaction %v", err)
			rbErr := tx.Rollback(context.Background())
			if rbErr != nil {
				_dbUtillog.Errorf("Unable to roll back %v", rbErr)
			}

			return 0, err
		}
		_dbUtillog.Infof("Execution result %v", rslt)
		count += rslt.RowsAffected()
	}

	if err := tx.Commit(context.Background()); err != nil {
		rbErr := tx.Rollback(context.Background())
		if rbErr != nil {
			_dbUtillog.Errorf("Unable to roll back %v", rbErr)
		}
		return 0, err
	}

	return count, nil

}

//Shutdown close the db connection
func (dbu *PGSqlDBUtil) Shutdown() {
	if dbu.dbConnection != nil {
		dbu.dbConnection.Close(context.Background())
	}
}

//TODO: Do caching here based on the sql
func (dbu *PGSqlDBUtil) getColMap(rslt pgx.Rows) map[int]string {

	colMap := make(map[int]string)
	for pos, fd := range rslt.FieldDescriptions() {
		colMap[pos] = string(fd.Name)
		if string(fd.Name) == "id" {
			colMap[pos] = "_id"
		}
	}
	return colMap
}
