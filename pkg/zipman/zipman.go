package zipman

import (
	"archive/zip"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/whytehack/goftp/pkg/goftp"
	"github.com/whytehack/goftp/pkg/logs"
	"github.com/whytehack/goftp/pkg/panik"
)

func extractAndWriteFile(f *zip.File, dest string) error {
	fileInZip, err := f.Open()
	if err != nil {
		logs.ERROR.Printf("%s could not be opened: %v \n", f.Name, err)
		return err
	}
	defer fileInZip.Close()

	path := filepath.Join(dest, f.Name)

	// Check for ZipSlip (Directory traversal)
	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		logs.ERROR.Printf("ZipSlip has been found: %v \n", err)
		return err
	}

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
		return nil
	}

	os.MkdirAll(filepath.Dir(path), f.Mode())
	Newfile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		logs.ERROR.Printf("%s could not be opened: %v \n", path, err)
		return err
	}
	defer Newfile.Close()

	fileSize, err := Newfile.ReadFrom(fileInZip)
	if err != nil {
		return err
	}

	log.Printf(" %s was extracted to %s. Its size is %d \n", filepath.Base(Newfile.Name()), Newfile.Name(), fileSize)
	return nil
}

func unzipWithoutPassword(src, dest string) error {
	z, err := zip.OpenReader(src)
	if err != nil {
		logs.ERROR.Printf("%s could not be opened: %v \n", src, err)
		return err
	}
	defer z.Close()

	for _, f := range z.File {
		err := extractAndWriteFile(f, dest)
		if err != nil {
			logs.ERROR.Printf("%s could not be extracted: %v \n", src, err)
			return err
		}
	}

	return nil
}

func unzipWithPassword(src, dest, zipPassword string) error {
	cmd := exec.Command("C:\\Program Files\\7-Zip\\7z.exe", "e", src, "-o"+dest, "-p"+zipPassword, "-aoa")
	err := cmd.Run()

	if err != nil {
		logs.ERROR.Printf("%s could not be extracted: %v \n", src, err)
		return err
	}

	return nil
}

func Unzip(src, dest, zipPassword string) error {
	log.Printf("Unzipping %s to %s\n", filepath.Base(src), dest)
	if zipPassword == "" {
		err := unzipWithoutPassword(src, dest)
		if err != nil {
			logs.ERROR.Printf("%s could not be extracted: %v \n", src, err)
			return err
		}
		return nil
	}

	err := unzipWithPassword(src, dest, zipPassword)
	if err != nil {
		logs.ERROR.Printf("%s could not be extracted: %v \n", src, err)
		return err
	}

	return nil
}

func UnzipWorker(downloadedFileChan <-chan *goftp.DownloadFileInfo, zipPassword string, wait *sync.WaitGroup) {
	defer panik.Catch()
	defer wait.Done()

	for dfi := range downloadedFileChan {
		err := Unzip(dfi.Source, dfi.Destination, zipPassword) // zipPassword
		if err != nil {
			log.Printf("Something went wrong when %s unextracting \r\n", dfi.Source)
			logs.ERROR.Printf("Something went wrong when %s unextracting: %v \r\n", dfi.Source, err)
			continue
		}
	}
}
