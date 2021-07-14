package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var Config *configuration
var path string

type configuration struct {
	RemotePath string `yaml:"remotepath"`
	LocalPath  string `yaml:"localpath"`
	Host       string `yaml:"host"`
	UserName   string `yaml:"username"`
	Password   string `yaml:"password"`
}

func EndsWithAny(haystack string, caseInsensitive bool, needles ...string) bool {
	if caseInsensitive {
		haystack = strings.ToLower(haystack)
	}

	for _, needle := range needles {
		search := needle
		if caseInsensitive {
			search = strings.ToLower(search)
		}

		if strings.HasSuffix(haystack, search) {
			return true
		}
	}
	return false
}

// FileExists reports whether the named file exists as a boolean
func FileExists(filePath string) bool {
	if fi, err := os.Stat(filePath); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// Init loads a config file from the provided path.
// Providing an empty path for filePath parameter will load the default config file
func Init(name string) error {
	var err error
	// Save to global path
	path = name
	if !EndsWithAny(path, true, ".yml") {
		path = path + ".yml"
	}
	// Make sure it exists
	path, err = filepath.Abs(path)
	if err != nil {
		log.Printf("Can not get absolute path for config path: %s, err: %v", name, err)
		return err
	}
	// Create directories
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	// Create default config if not exists
	if !FileExists(path) {
		log.Println("The given file does not exist")
		return err
	}
	return load(path)
}

// Path returns the file path of config file
func Path() string {
	return path
}

// load loads config file from provided path
func load(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Can not read config file: %s, %v", filePath, err)
		return err
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Printf("Can not deserialize yaml file: %s, %v", filePath, err)
		return err
	}
	return nil
}

var lock = &sync.Mutex{}

// save saves the config file file.
// Empty path saves into default config file DRONE.yaml
func (self *configuration) save() error {
	lock.Lock()
	defer lock.Unlock()

	data, err := yaml.Marshal(self)
	if err != nil {
		log.Printf("Can not serialize yaml file: %s, %v", path, err)
		return err
	}
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		log.Printf("Can not write yaml file: %s, %v", path, err)
		return err
	}
	return nil
}
