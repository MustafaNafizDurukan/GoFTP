package file

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func getFileSize(FileName string) int64 {

	info, err := os.Stat(FileName)
	if err != nil {
		log.Fatal(err)
	}
	fileSize := info.Size()

	return fileSize
}

func GetLocalFiles(path string) map[string]int64 {
	fileNameList := make(map[string]int64)

	filepath.WalkDir(path, func(FilePath string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileName := filepath.Base(FilePath)

		fileNameList[fileName] = getFileSize(FilePath)

		return nil
	})
	return fileNameList
}
