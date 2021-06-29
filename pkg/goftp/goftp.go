package goftp

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/jlaffaye/ftp"
	. "github.com/whytehack/goftp/pkg/constants"
)

type SSFTP struct {
	client *ftp.ServerConn
}

func (s *SSFTP) Close() {
	err := s.client.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func New(user, password, host string) (*SSFTP, error) {
	host = fmt.Sprintf("%s:21", host)
	c, err := ftp.Dial(host, ftp.DialWithTimeout(0))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection succeed")

	err = c.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Authetication succeed")

	binsftp := &SSFTP{
		client: c,
	}
	return binsftp, nil
}

func (s *SSFTP) GetRemoteFileList(source string) map[string]int64 {
	fileNames := make(map[string]int64)

	w := s.client.Walk("/")
	for w.Next() {

		if err := w.Err(); err != nil {
			fmt.Println(err.Error())
			continue
		}

		fi := w.Stat()
		if fi.Type == ftp.EntryTypeFolder {
			continue // Skip dirx
		}

		if w.Path() != "" {
			fmt.Println(w.Path(), fmt.Sprint(fi.Size))
			fileNames[w.Path()] = int64(fi.Size)
		}

	}

	return fileNames
}

func (self *SSFTP) Copy(source, destination string, wg *sync.WaitGroup) {
	defer wg.Done()

	srcFile, err := self.client.Retr(source)
	if err != nil {
		log.Fatal(ERROR + source + " could not be read ")
	}
	defer srcFile.Close()

	buf, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Fatal("ASd")
	}
	println(string(buf))

	fmt.Println(srcFile)

}
