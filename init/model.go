package init

import (
	"fmt"
	"time"
)

type mongoSetting struct {
	URI        string        `yaml:"-" json:"uri"`
	Username   string        `yaml:"username" json:"-"`
	PWD        string        `yaml:"pwd" json:"-"`
	Domains    []domain      `yaml:"domains" json:"-"`
	DBName     string        `yaml:"db_name" json:"dbName"`
	ReplicaSet string        `yaml:"replica_set" json:"-"`
	Timeout    time.Duration `yaml:"timeout" json:""`
}

type domain struct {
	Bind string `yaml:"bind" json:"bind"`
	Host string `yaml:"host" json:"host"`
	Port string `yaml:"port" json:"port"`
	TLS  tls    `yaml:"tls" json:"tls"`
}

type tls struct {
	CertFile string /* file *,pem */
	KeyFile  string /* file *.key */
}

// Config object
type Config struct {
	BlockService domain       `yaml:"block_service" json:"blockService"`
	HostService  domain       `yaml:"host_service" json:"hostService"`
	MongoDB      mongoSetting `yaml:"mongodb" json:"mongodb"`
	LogToFile    bool         `yaml:"log_to_file" json:"logToFile"`
	DebugMode    bool         `yaml:"debug_mode" json:"debugMode"`
}

func (ins *domain) ParseURL() string {
	if len(ins.Host) == 0 {
		ins.Host = "127.0.0.1"
	}
	if len(ins.Port) == 0 || ins.Port == "80" || ins.Port == "443" {
		return ins.Host
	}
	return fmt.Sprintf("%s:%s", ins.Host, ins.Port)
}

func (ins *domain) ParseAddr() string {
	if len(ins.Bind) == 0 {
		ins.Bind = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%s", ins.Bind, ins.Port)
}

func (ins *mongoSetting) ParseURI() string {

	ins.URI = "mongodb://"

	if len(ins.Username) > 0 {
		var pwd string
		ins.URI += ins.Username + ":"
		if len(ins.PWD) > 0 {
			ins.URI += ins.PWD + "@"
		} else {
			fmt.Println("\r\nEnter MONGO_PWD")
			fmt.Scanf("%s", &pwd)
			ins.URI += ins.PWD + "@"
		}
	}

	for i, domain := range ins.Domains {
		if len(domain.Port) == 0 {
			domain.Port = "27017"
		}
		if i != 0 {
			ins.URI += ","
		}
		ins.URI += domain.Host + ":" + domain.Port
	}

	ins.URI += "/" + ins.ReplicaSet

	if len(ins.ReplicaSet) > 0 {
		ins.URI += "?replicaSet=" + ins.ReplicaSet
	}

	return ins.URI
}
