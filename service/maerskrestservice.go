package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ibm.com/MaerskTar/MaerskTarSdk2/model"
	"ibm.com/MaerskTar/MaerskTarSdk2/util"
	//"github.com/tools/hrms/model"
	//"github.com/tools/hrms/util"
)

var _logger = logrus.New()

//APIResponse returns the service response
type APIResponse struct {
	Message   string      `json:"serviceMessage"`
	Payload   interface{} `json:"payload,omitempty"`
	ServiceTS string      `json:"ts"`
	IsSuccess bool        `json:"isSuccess"`
}

//MAERSKRestService implements the rest API for HRM transactions
type MAERSKRestService struct {
	dbUtil         *util.PGSqlDBUtil
	refDataService *RefDataRESTService
}

//NewMAERSKRestService retuens a new initialized version of the service
func NewMAERSKRestService(config []byte, verbose bool) *MAERSKRestService {
	service := new(MAERSKRestService)
	if err := service.Init(config, verbose); err != nil {
		_logger.Errorf("Unable to intialize service instance %v", err)
		return nil
	}
	return service
}

//Init initializes the service
func (srv *MAERSKRestService) Init(config []byte, verbose bool) error {
	if verbose {
		_logger.SetLevel(logrus.DebugLevel)
	}
	dbUtil, err := util.NewPGSqlDBUtil(config, verbose)
	if err != nil {
		_logger.Errorf("Error in intializing PGSQL util %v", err)
		return err
	}
	srv.dbUtil = dbUtil
	refDS := NewRefDataRestService(srv.dbUtil, verbose)
	if refDS == nil {
		_logger.Errorf("Error in intializing RefDataService ")
		return fmt.Errorf("Error in intializing RefDataService")
	}
	srv.refDataService = refDS
	_logger.Info("Service instance intialized. Waiting to lauch the service")
	return nil
}

//Serve runs the infinite service method.
func (srv *MAERSKRestService) Serve(address string, port int, stopSignal chan bool) {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//TODO: Following to be changed for production

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowMethods("GET", "POST", "DELETE", "PUT", "OPTIONS", "HEAD")
	corsConfig.AddAllowHeaders("Authorization", "WWW-Authenticate", "Content-Type", "Accept", "X-Requested-With")
	corsConfig.AddExposeHeaders("Authorization", "WWW-Authenticate", "Content-Type", "Accept", "X-Requested-With")
	router.Use(cors.New(corsConfig))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, buildResponse(true, "Service Available", nil))
	})
	//prebaplie general insertion
	//Add a new prebaple file
	router.POST("/api/prebaplie/insertPrebaplie", func(c *gin.Context) {
		resp := srv.insertPrebaplie(c)
		c.JSON(http.StatusOK, resp)
	})
	//Reference data
	srv.refDataService.AddRouters(router)
	//Create a server instance and start listening
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: router,
	}
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		// service connections
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_logger.Errorf("Unable to listen: %v\n", err)
			return
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		<-stopSignal
		_logger.Infof("Received stop signal")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_logger.Infof("Sending shutdown to http server")
		httpServer.Shutdown(ctx)
	}()
	time.Sleep(1 * time.Second)
	//Wait indefinitely
	wg.Wait()

}

//Adds a new Prebaplie file.
func (srv *MAERSKRestService) insertPrebaplie(c *gin.Context) APIResponse {
	//TODO: Check with the existing employee code
	var preBap model.PreBaplie
	if !parseInput(c, &preBap) {
		return buildResponse(false, "Unable to parse input json", nil)
	}

	sql, params := preBap.GetInsertStatement()
	err := srv.dbUtil.InsertOne(sql, params, nil)
	if err != nil {
		_logger.Errorf("Error in adding new employee in db %v", err)
		return buildResponse(false, "Employee could not be added ", err.Error())
	}
	return buildResponse(true, "Employee added successfully ", preBap)

}

//All utility methods follows here

func (srv *MAERSKRestService) checkAuth(c *gin.Context) bool {
	//TODO:Need to perform a proper implementation
	return true
}

func parseInput(c *gin.Context, obj interface{}) bool {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_logger.Errorf("Error in reading the request body %v", err)
		return false
	}
	if err = json.Unmarshal(bodyBytes, &obj); err != nil {
		_logger.Errorf("Error in parsing request body to %v %v", reflect.TypeOf(obj), err)
		return false
	}
	return true
}

func buildResponse(isOk bool, msg string, payload interface{}) APIResponse {
	return APIResponse{
		IsSuccess: isOk,
		Message:   msg,
		Payload:   payload,
		ServiceTS: time.Now().Format("2006-01-02-15:04:05.000"),
	}

}
