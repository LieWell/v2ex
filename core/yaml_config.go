package core

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var GlobalConfig YamlConfig

type YamlConfig struct {
	Http  YamlHttp  `yaml:"http"`
	Mysql YamlMysql `yaml:"mysql"`
	Zap   YamlZap   `yaml:"zap"`
	V2ex  YamlV2ex  `yaml:"v2ex"`
}

type YamlHttp struct {
	ListenAddr    string `yaml:"listenAddr"`    // http 监听配置
	ListenAddrTLS string `yaml:"listenAddrTLS"` // TLS 配置,留空表示不起用
	CertFile      string `yaml:"certFile"`      // TLS 证书
	KeyFile       string `yaml:"keyFile"`       // TLS 密钥
	Proxy         string `yaml:"proxy"`         // 代理地址
	ExposeAPI     bool   `yaml:"exposeAPI"`     // 是否暴露 API 服务
}

type YamlMysql struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Database           string `yaml:"database"`
	MaxOpenConnections int    `yaml:"maxOpenConnections"`
	MaxIdleConnections int    `yaml:"maxIdleConnections"`
}

type YamlZap struct {
	Level   string `yaml:"level"`   // 日志级别
	File    string `yaml:"file"`    // 日志文件
	MaxSize int    `yaml:"maxSize"` // 日志文件大小,单位 Mi
	MaxAge  int    `yaml:"maxAge"`  // 日志保存时长,单位 天
}

type YamlV2ex struct {
	AvatarDir   string `yaml:"avatarDir"`
	Token       string `yaml:"token"`
	StartSpider bool   `yaml:"startSpider"`
}

func LoadYamlConfig(filepath string) {

	raw, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("read config file error: %v", err)
	}

	var cfg YamlConfig
	err = yaml.Unmarshal(raw, &cfg)
	if err != nil {
		log.Fatalf("unmarshal config file error: %v", err)
	}

	// 记录所有配置的值
	GlobalConfig = cfg
}
