package local

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/whytehack/goftp/pkg/logs"
)

func getFileSize(fileName string) (int64, error) {
	info, err := os.Stat(fileName)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	fileSize := info.Size()

	return fileSize, nil
}

func Files(path string, PathsNotLookAtIt []string) map[string]int64 {
	fileNameList := make(map[string]int64)

	filepath.WalkDir(path, func(FilePath string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if len(PathsNotLookAtIt) != 0 {
			for _, value := range PathsNotLookAtIt {
				if strings.Contains(FilePath, value) {
					return filepath.SkipDir
				}
			}
		}

		if f.IsDir() {
			return nil
		}

		fileName := filepath.Base(FilePath)

		fileSize, err := getFileSize(FilePath)
		if err != nil {
			logs.ERROR.Println("Could not get file size", fileName)
			return err
		}

		fileNameList[fileName] = fileSize

		return nil
	})

	return fileNameList
}

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func FileSize(path string) int64 {
	fs, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return fs.Size()
}
