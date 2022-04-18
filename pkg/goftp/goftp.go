package goftp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/secsy/goftp"
	"github.com/whytehack/goftp/pkg/logs"
	"github.com/whytehack/goftp/pkg/panik"
)

type SSFTP struct {
	Cli *goftp.Client
}

func (s *SSFTP) Close() {
	err := s.Cli.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *SSFTP) Files(path string, isForwardSlash bool) ([]FileInfo, error) {
	objects, err := s.Cli.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	files := make([]FileInfo, 0)

	var slash string
	if isForwardSlash {
		slash = "/"
	} else {
		slash = "\\"
	}

	if path == "/" {
		slash = ""
	}

	for _, f := range objects {
		if f.IsDir() {
			f, err := s.Files(fmt.Sprintf("%s%s%s", path, slash, f.Name()), isForwardSlash)
			if err != nil {
				logs.ERROR.Printf("%v \n", err)
			}

			files = append(files, f...)

		} else {
			files = append(files, FileInfo{
				Name: fmt.Sprintf("%s%s%s", path, slash, f.Name()),
				Date: f.ModTime(), //f.ModTime().Format("02-01-2006"),
				Size: f.Size(),
			})
		}
	}

	return files, nil
}

func New(user, password, host string) (*SSFTP, error) {
	config := goftp.Config{
		User:               user,
		Password:           password,
		ConnectionsPerHost: 10,
		Timeout:            10 * time.Second,
	}

	client, err := goftp.DialConfig(config, host)
	if err != nil {
		logs.ERROR.Printf("Can not dial ftp: %v \n", err)
		return nil, err
	}

	return &SSFTP{Cli: client}, nil
}

type FileInfo struct {
	Name string
	Size int64
	Date time.Time
}

func (s *SSFTP) DownloadWorker(i int, remoteFileChan <-chan *DownloadFileInfo, downloadedFileChan chan<- *DownloadFileInfo, wait *sync.WaitGroup) {
	defer wait.Done()
	defer panik.Catch()

	for dfi := range remoteFileChan {
		err := s.download(i, dfi.Source, dfi.Destination)
		if err != nil {
			logs.ERROR.Printf("%v \n", err)
		} else {
			downloadedFileChan <- &DownloadFileInfo{
				Source:      dfi.Destination,
				Destination: filepath.Dir(dfi.Destination),
			}
		}
	}
}

type DownloadFileInfo struct {
	Source      string
	Destination string
}

func (s *SSFTP) download(i int, source, destination string) error {
	f, err := s.Cli.Stat(source)
	if err != nil {
		logs.ERROR.Printf("Failed to get file stat for %s: %v \n", source, err)
		return err
	}

	dir := filepath.Dir(destination)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0700)
	}

	createdFile, err := os.Create(destination)
	if err != nil {
		logs.ERROR.Printf("Failed to create file for %s: %v \n", destination, err)
		return err
	}
	defer createdFile.Close()

	log.Println(source, "is downloading from remote server... worker:", i)

	err = s.Cli.Retrieve(source, createdFile)
	if err != nil {
		logs.ERROR.Printf("Failed to retrieve file for %s: %v \n", source, err)
		return err
	}

	floc, err := os.Stat(destination)
	if err != nil {
		logs.ERROR.Printf("Failed to get file stat for %s: %v \n", destination, err)
		return err
	}

	if floc.Size() != f.Size() {
		logs.ERROR.Printf("%s could not be downloaded properly. \n", source)
		return err
	}

	log.Println(source, "has been downloaded from remote server... worker:", i)
	return nil
}
