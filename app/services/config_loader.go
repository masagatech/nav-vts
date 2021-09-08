package services

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/masagatech/nav-vts/app/models"
	"gopkg.in/yaml.v2"
)

type Confiuration struct {
	conf models.Config
}

func (c *Confiuration) LoadConfig() *models.Config {

	var config map[interface{}]interface{}
	pwd, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(pwd + "/app/config/conf.yml")
	if err != nil {
		log.Panic(err)
	}

	err1 := yaml.Unmarshal([]byte(yamlFile), &config)
	if err1 != nil {
		log.Panic(err1)
	}

	return c.loadEnvConfig(config["env"].(string))
}

func (c *Confiuration) loadEnvConfig(env string) *models.Config {
	pwd, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(pwd + "/app/config/conf-" + env + ".yml")
	if err != nil {
		log.Panic(err)
	}

	err1 := yaml.Unmarshal([]byte(yamlFile), &c.conf)
	if err1 != nil {
		log.Panic(err1)
	}
	return &c.conf
}

func (c *Confiuration) GetConfig() *models.Config {
	return &c.conf
}
