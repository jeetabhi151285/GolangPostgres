package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ibm.com/MaerskTar/MaerskTarSdk2/model"
	"ibm.com/MaerskTar/MaerskTarSdk2/util"
	//"github.com/tools/hrms/model"
	//"github.com/tools/hrms/util"
)

var _rdsLogger = logrus.New()

//RefDataRESTService provides reference data related rest services
type RefDataRESTService struct {
	dbUtil *util.PGSqlDBUtil
}

//NewRefDataRestService retuens a new initialized version of the service
func NewRefDataRestService(dbUtil *util.PGSqlDBUtil, verbose bool) *RefDataRESTService {
	service := new(RefDataRESTService)
	if err := service.Init(dbUtil, verbose); err != nil {
		_rdsLogger.Errorf("Unable to intialize service instance %v", err)
		return nil
	}
	return service
}

//Init intializes the service instance
func (s *RefDataRESTService) Init(dbUtil *util.PGSqlDBUtil, verbose bool) error {
	if verbose {
		_rdsLogger.SetLevel(logrus.DebugLevel)
	}
	if dbUtil == nil {
		return fmt.Errorf("Null DB Util reference passed ")
	}
	s.dbUtil = dbUtil
	return nil
}

//AddRouters add api end points specific to this service
func (s *RefDataRESTService) AddRouters(router *gin.Engine) {
	router.POST("/api/refdata", func(c *gin.Context) {
		resp := s.addNewEntry(c)
		c.JSON(http.StatusOK, resp)
	})
	// router.POST("/api/refdata", func(c *gin.Context) {
	// 	resp := s.updateEntry(c, false)
	// 	c.JSON(http.StatusOK, resp)
	// })
	// router.DELETE("/api/refdata", func(c *gin.Context) {
	// 	resp := s.updateEntry(c, true)
	// 	c.JSON(http.StatusOK, resp)
	// })
	// router.GET("/api/refdata/:searchCriteria", func(c *gin.Context) {
	// 	resp := s.searchEntries(c)
	// 	c.JSON(http.StatusOK, resp)
	// })

}

func (s *RefDataRESTService) addNewEntry(c *gin.Context) APIResponse {
	refDataRecords := make([]model.PreBaplie, 0)
	if !parseInput(c, &refDataRecords) {
		return buildResponse(false, "Unable to parse input json", nil)
	}
	//Validate input
	for _, refData := range refDataRecords {
		if len(strings.TrimSpace(refData.Voyage)) == 0 || len(strings.TrimSpace(refData.Vessel)) == 0 {
			return buildResponse(false, "Voyage/Vessel code is mandatory ", refData)
		}
		//sql, params := s.getInsertStatement(refData)
		sql, params := refData.GetInsertStatement()
		err := s.dbUtil.InsertOne(sql, params, nil)

		if err != nil {
			_logger.Errorf("Error in adding new reference data in db %v", err)
			return buildResponse(false, "Reference data could not be added ", err.Error())
		}
	}
	return buildResponse(true, "All reference data added successfully ", refDataRecords)
}

// func (s *RefDataRESTService) updateEntry(c *gin.Context, isDelete bool) APIResponse {
// 	refDataRecords := make([]model.ReferenceDataInfo, 0)
// 	var sql string
// 	var params []interface{}
// 	if !parseInput(c, &refDataRecords) {
// 		return buildResponse(false, "Unable to parse input json", nil)
// 	}
// 	//Validate input
// 	for _, refData := range refDataRecords {
// 		if len(strings.TrimSpace(refData.Value)) == 0 || len(strings.TrimSpace(refData.Category)) == 0 {
// 			return buildResponse(false, "Value/Category code is mandatory ", refData)
// 		}
// 		if isDelete {
// 			sql, params = s.getDeleteStatement(refData)
// 		} else {
// 			sql, params = refData.GetUpdateStatement()
// 		}
// 		count, err := s.dbUtil.UpdateRecords(sql, params)

// 		if err != nil || count == 0 {
// 			_logger.Errorf("Error in updating reference data in db %v", err)
// 			return buildResponse(false, "Reference data could not be updated ", err.Error())
// 		}
// 	}
// 	return buildResponse(true, "All reference data updated successfully ", refDataRecords)
// }
// func (s *RefDataRESTService) searchEntries(c *gin.Context) APIResponse {
// 	inputCriteria := strings.TrimSpace(c.Param("searchCriteria"))
// 	if len(inputCriteria) == 0 {
// 		return buildResponse(false, "Category not provided", nil)
// 	}
// 	categoryList := strings.Split(inputCriteria, ",")
// 	params := make([]interface{}, 0)
// 	var strBuf bytes.Buffer
// 	for index, category := range categoryList {
// 		if index == 0 {
// 			strBuf.WriteString(fmt.Sprintf(" $%d ", index+1))
// 		} else {
// 			strBuf.WriteString(fmt.Sprintf(" , $%d ", index+1))
// 		}
// 		params = append(params, category)
// 	}
// 	sql := fmt.Sprintf("select category,value,diplaytext,displayorder,extrainfo from hrm.reference_info where category in ( %s ) order by category,displayorder", strBuf.String())
// 	_rdsLogger.Debugf("%s %v", sql, categoryList)
// 	list, err := s.dbUtil.QueryRecords(sql, params...)
// 	if err != nil {
// 		buildResponse(false, "Error in retrieval", err.Error())
// 	}
// 	dataMap := model.BuildRefDataMap(list)
// 	return buildResponse(true, "Sucessfully retrived reference data", dataMap)
// }

/*func (s *RefDataRESTService) getInsertStatement(refData model.ReferenceDataInfo) (string, []interface{}) {
	sql := "insert into hrm.reference_info (category,value,diplaytext,displayorder) values( $1,$2,$3,$4) "
	params := []interface{}{refData.Category, refData.Value, refData.DisplayText, refData.DisplayOrder}
	return sql, params
}
func (s *RefDataRESTService) getUpdateStatement(refData model.ReferenceDataInfo) (string, []interface{}) {
	sql := "update hrm.reference_info set diplaytext = $1 ,displayorder = $2 where category =$3 and value = $4 "
	params := []interface{}{refData.DisplayText, refData.DisplayOrder, refData.Category, refData.Value}
	return sql, params
}*/
// func (s *RefDataRESTService) getDeleteStatement(refData model.ReferenceDataInfo) (string, []interface{}) {
// 	sql := "delete from hrm.reference_info where category =$1 and value = $2 "
// 	params := []interface{}{refData.Category, refData.Value}
// 	return sql, params
// }
