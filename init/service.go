package init

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type service struct {
	settingPath string
	fileName    string
	aes256Key   [32]byte
	naclKey     [32]byte
}

func (ins service) GetFilePath() string {
	var (
		filePath string
		arr      []string = strings.Split(ins.settingPath, "/")
	)
	if len(arr) > 0 && len(arr[0]) == 0 {
		filePath = "/"
	}
	if strings.HasSuffix(ins.fileName, ".yaml") == false {
		ins.fileName += ".yaml"
	}
	arr = append(arr, ins.fileName)
	filePath += strings.Join(arr, "/")

	return filePath
}

// Service Repository
type Service interface {
	Init() *Config

	GetAES256Key() [32]byte
	GetNaclKey() [32]byte
}

// NewService function
func NewService() Service {
	return &service{
		settingPath: "config",
		fileName:    "setting.yaml",
		aes256Key:   [32]byte{38, 84, 187, 217, 156, 255, 214, 123, 56, 55, 10, 25, 95, 75, 163, 58, 146, 125, 147, 32, 163, 245, 21, 32, 85, 26, 54, 36, 189, 231, 91, 65},
		naclKey:     [32]byte{22, 25, 98, 123, 65, 14, 189, 148, 26, 59, 48, 163, 235, 247, 15, 42, 82, 93, 54, 168, 215, 193, 167, 14, 167, 65, 148, 128, 39, 248, 149, 64},
	}
}

func (ins *service) GetAES256Key() [32]byte {
	return ins.aes256Key
}

func (ins *service) GetNaclKey() [32]byte {
	return ins.naclKey
}

// Init function
func (ins *service) Init() *Config {

	var (
		filePath = ins.GetFilePath()
		setting  = &Config{
			BlockService: domain{Bind: "0.0.0.0", Host: "127.0.0.1", Port: "8080", TLS: tls{"", ""}},
			HostService:  domain{Bind: "0.0.0.0", Host: "127.0.0.1", Port: "8080", TLS: tls{"", ""}},
			MongoDB: mongoSetting{
				Username: "",
				PWD:      "",
				Domains: []domain{
					{Host: "127.0.0.1", Port: "27017"},
				},
				DBName:     "BlockChain",
				ReplicaSet: "",
				Timeout:    5 * time.Second,
			},
			DebugMode: false,
			LogToFile: false,
		}
	)

	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("\r\nThe file \"setting.yaml\" was not found.\r\nThe system will create this file with default settings.\r\nPlease edit this at the path \"PWD/config/setting.yaml\" and restart service to update the new values!")
		if err := os.MkdirAll(ins.settingPath, 0755); err != nil {
			logrus.Fatal(err)
		}
		body, err := yaml.Marshal(setting)
		if err != nil {
			logrus.Fatal(err)
		}
		err = ioutil.WriteFile(filePath, body, 0664)
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		if err := yaml.Unmarshal(body, setting); err != nil {
			logrus.Fatal(err)
		}
	}

	setting.MongoDB.URI = setting.MongoDB.ParseURI()

	fmt.Printf("\r\nSETTING DETAIL: \r\n%+v\r\n", setting)

	return setting
}
