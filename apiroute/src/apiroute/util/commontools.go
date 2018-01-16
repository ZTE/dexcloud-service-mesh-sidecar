package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigFilePath string = "conf"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetGoPath() []string {
	goPath := os.Getenv("GOPATH")
	fmt.Println(goPath)
	if strings.Contains(goPath, ";") { //windows
		return strings.Split(goPath, ";")
	} else if strings.Contains(goPath, ":") { //linux
		return strings.Split(goPath, ":")
	} else { //only one
		path := make([]string, 1, 1)
		path[0] = goPath
		return path
	}
}
func GetCfgFilePath() string {
	var err error
	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", AppPath)
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", workPath)
	appConfigPath := filepath.Join(workPath, ConfigFilePath)
	if !FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, ConfigFilePath)
		if !FileExists(appConfigPath) {
			goPath := GetGoPath()
			for _, val := range goPath {
				appConfigPath = filepath.Join(val, "src", "apiroute", ConfigFilePath)
				fmt.Println(appConfigPath)
				if FileExists(appConfigPath) {
					return appConfigPath
				}
			}
			appConfigPath = "/"
		}
	}

	return appConfigPath
}

func ReadJsonfile(datafile string) []byte {
	f, err := os.Open(datafile)
	if err != nil {
		return nil
	}
	b, err1 := ioutil.ReadAll(f)
	if err1 != nil {
		return nil
	}
	return b

}

func HTTPGet(base, query string) (b []byte, err error) {
	url := base + query
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	buf, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func HTTPGetWithIndex(base, query, index string) (b []byte, tag string, err error) {
	url := base + query
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Data-Tag", index)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, "", err
	}

	if res.StatusCode != 200 {
		return nil, "", fmt.Errorf(res.Status)
	}

	buf, err := ioutil.ReadAll(res.Body)
	dataTag := res.Header.Get("Data-Tag")
	res.Body.Close()

	if err != nil {
		return nil, "", err
	}

	return buf, dataTag, nil
}

func HTTPPost(base, query string, body []byte) error {
	contentType := "application/json"
	url := base + query
	res, err := http.Post(url, contentType, bytes.NewReader(body))

	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return fmt.Errorf(res.Status)
	}

	res.Body.Close()
	return nil
}

func HTTPDelete(base, query string) error {
	url := base + query
	req, _ := http.NewRequest("DELETE", url, nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 204 {
		return fmt.Errorf(res.Status)
	}

	res.Body.Close()
	return nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Printf("Dial errors:%v in GetOutboundIP", err)
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
