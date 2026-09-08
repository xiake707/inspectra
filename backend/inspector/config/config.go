package config

import (
	mysql "inspector/my_mysql"
)

// HostConfig 单台主机配置
type HostConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Password string `yaml:"password,omitempty"`
	KeyFile  string `yaml:"key_file,omitempty"`
}
var Servers []HostConfig
// GlobalConfig 全局默认配置
/*type GlobalConfig struct {
	User     string `yaml:"user"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password,omitempty"`
	KeyFile  string `yaml:"key_file,omitempty"`
}*/

// Config 整体配置
/*type Config struct {
	Global  GlobalConfig `yaml:"global"`
	Servers []HostConfig `yaml:"servers"`
}*/

// LoadConfig 从 YAML 文件加载配置
func LoadConfig() (*[]HostConfig, error) {

	hosts,err := mysql.GetHost()
	if err != nil {
		return nil,err
	}
	for _,host := range hosts {
		Servers = append(Servers,HostConfig{
			Host: host.Host,
			User: host.User,
			Port: host.Port,
			Password: host.Password,
			KeyFile: host.KeyFile,
		})
	}
	return &Servers, nil
}
