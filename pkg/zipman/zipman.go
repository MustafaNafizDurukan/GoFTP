package zipman

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/whytehack/goftp/pkg/constants"
)

func Unzip(src, dest string, wg *sync.WaitGroup) {
	defer wg.Done()

	dest += "\\all"

	r, err := zip.OpenReader(src)
	if err != nil {
		return
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	log.Printf(constants.RUNNING+"%s is extracting...", src)
	dest += "\\" + strings.TrimRight(filepath.Base(src), ".zip")
	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			log.Printf(constants.ERROR+"illegal file path: %s", path)
			return err
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Printf(constants.ERROR+"%s could not be opened", f.Name())
				return err
			}

			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			fileSize, err := io.Copy(f, rc)
			if err != nil {
				log.Printf(constants.ERROR+"%s could not be copied", f.Name())
				return err
			}

			log.Printf(constants.SUCCESS+" %s was extracted to %s. Its size is %d", filepath.Base(f.Name()), f.Name(), fileSize)
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			log.Println(constants.ERROR + "File could not be extracted")
			return
		}
	}
	log.Printf(constants.SUCCESS+"%s has been completely extracted to %s.", src, dest)
}
