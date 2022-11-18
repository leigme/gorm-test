package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
)

var (
	yamlPath   string
	yamlConfig YamlConfig
	db         *gorm.DB
)

func init() {
	flag.StringVar(&yamlPath, "y", "app.yaml", "yaml file absolute path")
	flag.Parse()
	f, err := os.Open(yamlPath)
	if err != nil {
		log.Printf("yaml file path is error: %s", err)
		return
	}
	data, err := io.ReadAll(f)
	if err != nil {
		log.Printf("yaml read content is error: %s", err)
		return
	}
	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		log.Printf("yaml unmarshal is error: %s", err)
		return
	}
}

func main() {
	var err error
	db, err = gorm.Open(mysql.Open(yamlConfig.MySqlConfig.Dsn()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Printf("db load is error: %s", err)
		return
	}
	var podInfos []PodInfo
	if result := db.Find(&podInfos); result.Error != nil {
		log.Printf("db is error: %s", result.Error)
		return
	}
	if j, err := json.Marshal(podInfos); err == nil {
		log.Printf("content: %s", string(j))
	}
}

type PodInfo struct {
	Id         int       `orm:"id" json:"id,omitempty"`
	Uid        string    `orm:"uid" json:"uid,omitempty"`
	PodName    string    `orm:"pod_name" json:"podName,omitempty"`
	HostIp     string    `orm:"host_ip" json:"hostIp,omitempty"`
	PodIp      string    `orm:"pod_ip" json:"podIp,omitempty"`
	Status     int       `orm:"status" json:"status,omitempty"`
	CreateTime time.Time `orm:"create_time" json:"createTime"`
	UpdateTime time.Time `orm:"update_time" json:"updateTime"`
}

type YamlConfig struct {
	AppConfig   App   `yaml:"app"`
	MySqlConfig MySql `yaml:"mysql"`
}

type App struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

type MySql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
}

// Dsn = user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func (m *MySql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.Name)
}
