package config

import (
	"boilerplate/utils"
	"boilerplate/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	_ "time/tzdata" // Important!

	clientV3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Context     context.Context        `json:"context"`
	AppName     string                 `json:"app_name"`
	SettingName string                 `json:"setting_name"`
	RootPath    string                 `json:"root_path"`
	Environment string                 `json:"environment"`
	EtcdHost    string                 `json:"etcd_host"`
	Port        int                    `json:"port"`
	Timezone    string                 `json:"timezone"`
	Location    *time.Location         `json:"location"`
	TokenConfig TokenConfig            `json:"token_config"`
	MySQL       MySQLConfig            `json:"mysql"`
	Redis       map[string]RedisConfig `json:"redis"`
}

type TokenConfig struct {
	SecretKey         string `json:"secret_key"`
	TokenHourLifeSpan int    `json:"token_hour_lifespan"`
}

type MySQLConfig struct {
	ConnMaxIdleTime int                `json:"conn_max_idle_time"`
	ConnMaxLifeTime int                `json:"conn_max_life_time"`
	MaxIdleConns    int                `json:"max_idle_conns"`
	MaxOpenConns    int                `json:"max_open_conns"`
	Default         MySQLConnection    `json:"default"`
	Connections     []MySQLConnections `json:"connections"`
}

type MySQLConnections struct {
	Datas  []string          `json:"datas"`
	Writes []MySQLConnection `json:"writes"`
	Reads  []MySQLConnection `json:"reads"`
}

type MySQLConnection struct {
	Network  string `json:"network"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	Charset  string `json:"charset"`
}

type RedisConfig struct {
	Addrs    []string `json:"addrs"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Route    struct {
		Latency bool `json:"latency"`
		Random  bool `json:"random"`
	} `json:"route"`
	DB              int    `json:"db"`
	Prefix          string `json:"prefix"`
	PoolSize        int    `json:"pool_size"`
	PoolTimeout     int    `json:"pool_timeout"`
	MinIdleConns    int    `json:"min_idle_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time"`
	ConnMaxLifeTime int    `json:"conn_max_life_time"`
}

func NewConfig(ctx context.Context, etcdHost, settingName string) *Config {
	config := &Config{
		Context:     ctx,
		SettingName: settingName,
		EtcdHost:    etcdHost,
	}

	config.getEtcd()

	var dir string
	var location *time.Location
	var err error

	if dir, err = os.Getwd(); err != nil {
		panic(err)
	}

	config.RootPath = dir

	if location, err = time.LoadLocation(config.Timezone); err != nil {
		panic(err)
	}

	config.Location = location

	logger.Init(fmt.Sprintf("%s/log/%s_%s.log", config.RootPath, config.AppName, config.GetCurrentTime(location).Format(time.DateOnly)))

	logger.Sugar.Debugf("config: %s", utils.DataString(config))

	return config
}

func (config *Config) GetCurrentTime(location *time.Location) time.Time {
	if location != nil {
		location = config.Location
	}

	return time.Now().In(location)
}

// Get the config for the database from etcd server
func (config *Config) getEtcd() {
	cli, err := clientV3.New(clientV3.Config{
		Endpoints:   []string{config.EtcdHost},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	defer cli.Close()

	resp, err := cli.Get(config.Context, config.SettingName)
	if err != nil {
		panic(err)
	}

	if len(resp.Kvs) == 0 {
		panic(fmt.Errorf("key %s not found", config.SettingName))
	}

	value := resp.Kvs[0].Value
	if err = json.Unmarshal(value, &config); err != nil {
		panic(err)
	}
}
