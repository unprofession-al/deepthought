package main

import "encoding/json"

type Configuration struct {
	Port                    string             `json:"port"`
	DbHost                  string             `json:"db_host"`
	DbPort                  string             `json:"db_port"`
	DbName                  string             `json:"db_name"`
	NodevarsProvidersString string             `json:"-"`
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
