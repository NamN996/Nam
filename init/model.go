package init

import "time"

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
	KeyFole  string /* file *.key */
}

type Config struct {
}
