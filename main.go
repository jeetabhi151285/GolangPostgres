package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"ibm.com/MaerskTar/MaerskTarSdk2/service"
)

var (
	_version = "default"
)

func main() {
	fmt.Println("Starting Maersk Tar Service ", _version)
	defer fmt.Println("Done....")
	port := flag.Int("p", 7070, "Service listen port")
	bindAddress := flag.String("b", "0.0.0.0", "Bind address")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Parse()
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Config file missing")
		fmt.Println("hrms [flags] <path to config file> ")
		flag.Usage()
		os.Exit(1)
	}
	//Read the config file
	configBytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Println("Unable to read config file ", err)
		os.Exit(1)
	}
	if maerskService := service.NewMAERSKRestService(configBytes, *verbose); maerskService != nil {
		stopSignal := make(chan bool)
		termination := make(chan os.Signal)
		signal.Notify(termination, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-termination
			fmt.Println("SIGTERM/SIGINT received from os")
			stopSignal <- true
		}()
		maerskService.Serve(*bindAddress, *port, stopSignal)
	} else {
		fmt.Println("Unable to start the service ...")
		os.Exit(2)
	}

}
