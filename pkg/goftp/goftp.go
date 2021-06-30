package goftp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

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
		log.Fatal(FAIL + "Failed to dial: " + err.Error())
	}

	err = c.Login(user, password)
	if err != nil {
		log.Fatal(FAIL + "Failed to login to client: " + err.Error())
	}

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

func (self *SSFTP) Copy(source, destination string) {

	read, err := self.client.Retr(source)
	if err != nil {
		log.Fatal(err)
	}
	defer read.Close()

	buf, err := ioutil.ReadAll(read)
	if err != nil {
		log.Printf(ERROR + "Failed to read file: " + err.Error())
	}

	var localFileName = path.Base(source)
	dstFilePath := path.Join(destination, localFileName)

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		log.Printf(FAIL + "Failed to create file: " + err.Error())
	}
	defer dstFile.Close()
	/*
		_, err = dstFile.Write([]byte(buf))
		if err != nil {
			log.Printf(ERROR + "Failed to write file: " + err.Error())
		} //Burada tekrar kaç byte indirildiğine bakıp daha sonra da bu byte sayısını kontrol edebilirim

		//Bir de bunu dene
	*/
	err = os.WriteFile(dstFilePath, []byte(buf), 0666)
	if err != nil {
		log.Printf(FAIL + "Failed to write file: " + err.Error())
	}
	log.Printf(SUCCESS+"%s file has been downloaded ", localFileName)

}
