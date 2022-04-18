package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/whytehack/goftp/pkg/config"
	"github.com/whytehack/goftp/pkg/flags"
	"github.com/whytehack/goftp/pkg/goftp"
	"github.com/whytehack/goftp/pkg/local"
	"github.com/whytehack/goftp/pkg/logs"
	"github.com/whytehack/goftp/pkg/panik"
	"github.com/whytehack/goftp/pkg/w32"
	"github.com/whytehack/goftp/pkg/zipman"
)

func (p *paramsT) equalizeToConfig() {
	p.Host = config.Config.Host
	p.LocalPath = config.Config.LocalPath
	p.Password = config.Config.Password
	p.RemotePath = config.Config.RemotePath
	p.UserName = config.Config.UserName
	p.PathsNotLookAt = config.Config.PathsNotLookAt
	p.ZipPassword = config.Config.ZipPassword

	if p.Host == "" {
		log.Println("Both parameters and config file are empty. Exiting the program...")
		os.Exit(1)
	}
	log.Println("Host has been set to  " + p.Host)
	log.Println("Username has been set to  " + p.UserName)
	log.Println("Password has been set to  " + p.Password)
	log.Println("Remote path has been set to  " + p.RemotePath)
	log.Println("Local path has been set to  " + p.LocalPath)
	log.Println("Paths that not will be looked at has been set to  " + p.PathsNotLookAt)
	log.Println("Zip password has been set to  " + p.ZipPassword)
}

func Process() {
	err := config.Init(params.Config)
	if err != nil {
		log.Println(nil, "[!] Can not load config file: %s, %v", params.Config, err)
		return
	}

	if params.Host == "" {
		params.equalizeToConfig()
	}

	server, err := goftp.New(
		params.UserName,
		params.Password,
		params.Host,
	)
	if err != nil {
		logs.ERROR.Println("Client could not be initialized")
		return
	}
	defer server.Close()

	remoteFiles := make(sortedRemoteFileType, 0)
	remoteFiles, err = server.Files(params.RemotePath, true)
	if err != nil {
		logs.ERROR.Println(err)
		return
	}
	if len(remoteFiles) == 0 {
		log.Println("Remote path is empty!")
		return
	}

	sort.Sort(remoteFiles)

	chanSize := 1000
	remoteFileChan := make(chan *goftp.DownloadFileInfo, chanSize)
	downloadedFileChan := make(chan *goftp.DownloadFileInfo, chanSize)

	WORKERNUMBER := 5

	go sendFileListToChan(remoteFiles, remoteFileChan)

	var unzipW sync.WaitGroup
	for i := 0; i < WORKERNUMBER*5; i++ {
		unzipW.Add(1)
		go zipman.UnzipWorker(downloadedFileChan, params.ZipPassword, &unzipW)
	}

	func() {
		var downloadW sync.WaitGroup
		for i := 0; i < WORKERNUMBER; i++ {
			downloadW.Add(1)
			go server.DownloadWorker(i, remoteFileChan, downloadedFileChan, &downloadW)
		}
		downloadW.Wait()

		close(downloadedFileChan)
	}()

	unzipW.Wait()

	log.Println("All files in remote has been downloaded...")
}

func sendFileListToChan(remoteFiles []goftp.FileInfo, remoteFileChan chan<- *goftp.DownloadFileInfo) {
	for _, fi := range remoteFiles { // remoteFiles is FileInfo slice
		filename := filepath.Base(fi.Name)
		filedir := filepath.Dir(fi.Name)

		dstPath := filepath.Join(params.LocalPath, filedir, fi.Date.Format("02-01-2006"))
		// "C:\FileZillaFTPFiles\server\malware\2022-03-03\"

		dstFile := filepath.Join(dstPath, filename)
		// "C:\FileZillaFTPFiles\server\malware\2022-03-03\malware_20220303.txt"

		if !local.IsFileExist(dstFile) {
			remoteFileChan <- &goftp.DownloadFileInfo{
				Source:      fi.Name,
				Destination: dstFile,
			}
		}
	}
	close(remoteFileChan)
}

func main() {
	var err error
	_, err = w32.CreateMutex(w32.MutexName)
	if err != nil {
		fmt.Printf("Error: this process has been created before\n")
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, logs.Set())
	log.SetOutput(multiWriter)

	defer panik.Catch()

	params.Args, err = flags.Parse(&params, os.Args)
	if err != nil {
		log.Printf("Error parsing command line arguments: %v \n", err)
		return
	}

	for {
		Process()
		time.Sleep(30 * time.Minute)
	}
}

type sortedRemoteFileType []goftp.FileInfo

func (p sortedRemoteFileType) Len() int {
	return len(p)
}

func (p sortedRemoteFileType) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p sortedRemoteFileType) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
