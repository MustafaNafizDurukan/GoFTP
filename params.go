package main

type paramsT struct {
	RemotePath     string   `short:"x" long:"RemotePath" required:"false" description:"(Optional) Remote directory to get files"`
	LocalPath      string   `short:"f" long:"LocalPath" required:"false" description:"(Optional) Local directory to save zip files and extracted from zip files"`
	Host           string   `short:"h" long:"host" required:"false" description:"(Optional) Host to be used when connecting"`
	UserName       string   `short:"u" long:"userName" required:"false" description:"(Optional) Username to be used when connecting"`
	Password       string   `short:"p" long:"password" required:"false" description:"(Optional) Password to be used when connecting "`
	Config         string   `short:"c" long:"config" required:"false" description:"(Optional) Config file name" default:"config.yml"`
	PathsNotLookAt string   `short:"n" long:"pathsnotlookat" required:"false" description:"(Optional) Comma separated list of path names that not to look at (all,example,Users)"`
	ZipPassword    string   `short:"z" long:"zippassword" required:"false" description:"(Optional) Downloaded zip password"`
	Args           []string // Positional arguments except application path (1st param) will be here
}

var params paramsT
