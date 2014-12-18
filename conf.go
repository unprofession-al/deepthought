package main

import "encoding/json"

type Configuration struct {
	Port                    string             `json:"port"`
	DbHost                  string             `json:"db_host"`
	DbPort                  string             `json:"db_port"`
	DbName                  string             `json:"db_name"`
	NodevarsProvidersString string             `json:"-"`
	NodevarsProviders       []NodevarsProvider `json:"nodevars_providers"`
	LdapConnString          string             `json:"-"`
	LdapConn                LdapConn           `json:"ldap_connection"`
}

func (c *Configuration) ParseProviders() error {
	return json.Unmarshal([]byte(c.NodevarsProvidersString), &c.NodevarsProviders)
}

type NodevarsProvider struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Prio int    `json:"prio"`
}

func (c *Configuration) ParseLdapConn() error {
	return json.Unmarshal([]byte(c.LdapConnString), &c.LdapConn)
}

type LdapConn struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	BaseDN       string `json:"basedn"`
	BindDNpatter string `json:"binddnpatter"`
}

//"ldap-ha.swisstxt.ch", 389,  "dc=stxt,dc=mpc", "uid={{username}},ou=users,ou=stxt,dc=stxt,dc=mpc"
