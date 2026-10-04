package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	DefaultKey  string `yaml:"defaultKey"`
	AuthAddress string `yaml:"authAddress"`
	MainAddress string `yaml:"mainAddress"`
	Addr        string `yaml:"addr"`
	FakeServer  bool   `yaml:"fakeServer"`

	MQTTHost     string `yaml:"-"`
	MQTTPort     string `yaml:"-"`
	MQTTUsername string `yaml:"-"`
	MQTTPassword string `yaml:"-"`
	MQTTSSL      bool   `yaml:"-"`
}

type AddonOptions struct {
	AuthAddress string `json:"auth_address"`
	MainAddress string `json:"main_address"`
}

type supervisorServiceResponse struct {
	Result  string         `json:"result"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

type UserConfig struct {
	Name     string `yaml:"name"`
	Phone    string `yaml:"phone"`
	Code     string `yaml:"code"`
	Sender   string `yaml:"sender"`
	Key      string `yaml:"key"`
	FileName string `yaml:"-"`
}

func loadUserConfig(filename string) (*UserConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg UserConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.FileName = filename
	return &cfg, nil
}

func (c *UserConfig) save() {
	if c.FileName != "" {
		data, err := yaml.Marshal(c)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(c.FileName, data, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func (c *UserConfig) UpdateCode(code string) {
	c.Code = code
	c.save()
}

func (c *UserConfig) UpdateKey(key string) {
	c.Key = key
	c.save()
}

func (c *UserConfig) UpdateUserInfo(name, phone string) {
	c.Name = name
	c.Phone = phone
	c.save()
}

func InitUserConfig(filename string) *UserConfig {
	userConfig, err := loadUserConfig(filename)
	if err != nil {
		userConfig = &UserConfig{
			Sender:   uuid.NewString(),
			FileName: filename,
		}
		userConfig.save()
	}
	return userConfig
}

func LoadAppConfig(filename string) (*AppConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var appConfig AppConfig
	err = yaml.Unmarshal(data, &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

func OverwriteAppConfig(config *AppConfig, filename string) *AppConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		return config
	}
	var appConfig AppConfig
	err = yaml.Unmarshal(data, &appConfig)
	if err != nil {
		return config
	}
	if appConfig.MainAddress != "" {
		log.Printf("local value MainAddress: %s", appConfig.MainAddress)
		config.MainAddress = appConfig.MainAddress
	}
	if appConfig.AuthAddress != "" {
		log.Printf("local value AuthAddress: %s", appConfig.AuthAddress)
		config.AuthAddress = appConfig.AuthAddress
	}
	if appConfig.FakeServer {
		log.Printf("local value MockServer: %t", appConfig.FakeServer)
		config.FakeServer = appConfig.FakeServer
	}
	return config
}

func LoadEnvOverrides(config *AppConfig) *AppConfig {
	config.MQTTHost = os.Getenv("MQTT_HOST")
	config.MQTTPort = os.Getenv("MQTT_PORT")
	config.MQTTUsername = os.Getenv("MQTT_USERNAME")
	config.MQTTPassword = os.Getenv("MQTT_PASSWORD")
	config.MQTTSSL = os.Getenv("MQTT_SSL") == "true"
	return config
}

func LoadAddonOptions(config *AppConfig, filename string) *AppConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		return config
	}
	var options AddonOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return config
	}
	if options.AuthAddress != "" {
		log.Printf("addon option AuthAddress: %s", options.AuthAddress)
		config.AuthAddress = options.AuthAddress
	}
	if options.MainAddress != "" {
		log.Printf("addon option MainAddress: %s", options.MainAddress)
		config.MainAddress = options.MainAddress
	}
	return config
}

func LoadSupervisorMQTT(config *AppConfig, token string) *AppConfig {
	req, err := http.NewRequest(http.MethodGet, "http://supervisor/services/mqtt", nil)
	if err != nil {
		return config
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("supervisor mqtt service request failed: %v", err)
		return config
	}
	defer resp.Body.Close()

	var body supervisorServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Printf("supervisor mqtt service decode failed: %v", err)
		return config
	}
	if body.Result != "ok" || len(body.Data) == 0 {
		log.Printf("supervisor mqtt service unavailable: %s", body.Message)
		return config
	}

	if host, _ := body.Data["host"].(string); host != "" {
		log.Printf("supervisor mqtt service: host=%s", host)
		config.MQTTHost = host
	}
	if port, ok := body.Data["port"]; ok {
		config.MQTTPort = fmt.Sprint(port)
	}
	if ssl, ok := body.Data["ssl"].(bool); ok {
		config.MQTTSSL = ssl
	}
	if username, _ := body.Data["username"].(string); username != "" {
		config.MQTTUsername = username
	}
	if password, _ := body.Data["password"].(string); password != "" {
		config.MQTTPassword = password
	}
	return config
}
