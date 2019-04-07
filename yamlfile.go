package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ReadYamlFile(name string, cfg interface{}) error {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}
