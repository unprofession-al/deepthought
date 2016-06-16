package main

import "encoding/json"

type Configuration struct {
	Port                    string             `json:"port"`
	DbHosts                 string             `json:"db_hosts"`
	DbUser                  string             `json:"ub_user" yaml:"db_user"`
	DbPass                  string             `json:"-" yaml:"-"`
	DbName                  string             `json:"db_name"`
	NodevarsProvidersString string             `json:"-" yaml:"-"`
	NodevarsProviders       []NodevarsProvider `json:"nodevars_providers"`
}

func (c *Configuration) ParseProviders() error {
	return json.Unmarshal([]byte(c.NodevarsProvidersString), &c.NodevarsProviders)
}

type NodevarsProvider struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Prio int    `json:"prio"`
}
