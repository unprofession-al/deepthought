package main

import (
	"encoding/base64"
	"errors"
	"fmt"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marcsauter/ldap"
)

func BasicAuthLDAP() gin.HandlerFunc {
	AddSource("default", config.LdapConn.Host, config.LdapConn.Port, false, config.LdapConn.BaseDN, config.LdapConn.BindDNpatter)

	return func(c *gin.Context) {
		user, pass, _ := validateBasicAuthHeader(c.Request.Header.Get("Authorization"))

		ok := LoginUser(user, pass)

		if !ok {
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.Fail(401, errors.New("Unauthorized"))
		} else {
			c.Set(gin.AuthUserKey, "user")
		}
	}
}

func validateBasicAuthHeader(h string) (user string, pass string, e error) {
	var credentials string

	parts := strings.SplitN(h, " ", 2)
	if len(parts) == 2 && parts[0] == "Basic" {
		credentials = parts[1]
	} else {
		e = errors.New("Not a valid authorization header.")
		return
	}

	if b, err := base64.StdEncoding.DecodeString(credentials); err == nil {
		parts = strings.Split(string(b), ":")
		if len(parts) == 2 {
			user = parts[0]
			pass = parts[1]
			return
		} else {
			e = errors.New("Credentials are malformed.")
			return
		}
	} else {
		e = err
		return
	}
}

type Ldapsource struct {
	Name          string // canonical name (ie. corporate.ad)
	Host          string // LDAP host
	Port          int    // port number
	UseSSL        bool   // Use SSL
	BaseDN        string // Base DN
	BindDNPattern string
}

var (
	Authensource []Ldapsource
)

func AddSource(name string, host string, port int, usessl bool, basedn string, binddnpattern string) {
	ldaphost := Ldapsource{name, host, port, usessl, basedn, binddnpattern}
	Authensource = append(Authensource, ldaphost)
}

func LoginUser(name, passwd string) bool {
	for _, ls := range Authensource {
		err := ls.VerifyCredentials(name, passwd)
		if err == nil {
			return true
		}
	}
	return false
}

func (ls Ldapsource) VerifyCredentials(name, passwd string) error {
	l, err := ldapDial(ls)
	if err != nil {
		return errors.New(err.String())
	}
	defer l.Close()

	binddn := ls.MakeBindDN(name)

	err = l.Bind(binddn, passwd)
	if err != nil {
		return errors.New(err.String())
	}

	return nil
}

func ldapDial(ls Ldapsource) (*ldap.Conn, *ldap.Error) {
	if ls.UseSSL {
		return ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", ls.Host, ls.Port), nil)
	} else {
		return ldap.Dial("tcp", fmt.Sprintf("%s:%d", ls.Host, ls.Port))
	}
}

func (ls Ldapsource) MakeBindDN(name string) string {
	return strings.Replace(ls.BindDNPattern, "{{username}}", name, -1)

}
