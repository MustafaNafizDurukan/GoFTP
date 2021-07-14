package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/whytehack/goftp/pkg/config"
	"github.com/whytehack/goftp/pkg/constants"
	"github.com/whytehack/goftp/pkg/file"
	"github.com/whytehack/goftp/pkg/flags"
	"github.com/whytehack/goftp/pkg/goftp"
	"github.com/whytehack/goftp/pkg/zipman"
)

func (p *paramsT) equalizeToConfig() {
	p.Host = config.Config.Host
	p.LocalPath = config.Config.LocalPath
	p.Password = config.Config.Password
	p.RemotePath = config.Config.RemotePath
	p.UserName = config.Config.UserName

	if p.Host == "" {
		log.Println(constants.FAIL + "Both parameters and config file are empty. Exiting the program...")
		os.Exit(1)
	}
	log.Println(constants.STATUS + "Host has been set to  " + p.Host)
	log.Println(constants.STATUS + "Local path has been set to  " + p.LocalPath)
	log.Println(constants.STATUS + "Password has been set to  " + p.Password)
	log.Println(constants.STATUS + "Remote path has been set to  " + p.RemotePath)
	log.Println(constants.STATUS + "Username has been set to  " + p.UserName)
}

func (p *paramsT) areFilesCorrect(remoteFileList map[string]int64) {
	localFileList := file.GetLocalFiles(p.LocalPath)
	for element := range remoteFileList {
		if localFileList[path.Base(element)] != remoteFileList[element] {
			log.Printf(constants.FAIL+"Something wrong with %s. Maybe overwrited!", element)
		}
	}

	if len(localFileList) == 0 {
		log.Printf(constants.FAIL + "Local path is empty! Check your local directory.")
	}
}

func Process() {
	var wg sync.WaitGroup

	err := config.Init(params.Config)
	if err != nil {
		log.Println(nil, "[!] Can not load config file: %s, %v", params.Config, err)
		return
	}

	if params.Host == "" {
		params.equalizeToConfig()
	}

	client, err := goftp.New(
		params.UserName,
		params.Password,
		params.Host,
	)
	if err != nil {
		log.Println(constants.ERROR + "Client could not be initialized")
		return
	}
	defer client.Close()

	remoteFileList := client.GetRemoteFileList(params.RemotePath)
	if len(remoteFileList) == 0 {
		log.Println(constants.FAIL + "Remote path is empty!")
		return
	}

	localFileList := file.GetLocalFiles(params.LocalPath)

	for element := range remoteFileList {
		if localFileList[path.Base(element)] != remoteFileList[element] {

			dstFilePath, err := client.Copy(element, params.LocalPath)
			if err != nil {
				log.Println(constants.ERROR + "Something went wrong when copy file from remote server!")
				continue
			}

			wg.Add(1)

			go zipman.Unzip(dstFilePath, params.LocalPath, &wg)
		}
	}
	wg.Wait()
	params.areFilesCorrect(remoteFileList)
}

func main() {
	var err error

	params.Args, err = flags.Parse(&params, os.Args)
	if err != nil {
		msg := fmt.Sprintf("Error parsing command line arguments: %v", err)
		log.Println(msg)
		return
	}

	logFile, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf(constants.ERROR+"error opening file: %v \n", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	ticker := time.NewTicker(time.Second * 2)
	for ; true; <-ticker.C {
		Process()
	}
}
