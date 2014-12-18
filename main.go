package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sontags/env"
)

var config = &Configuration{}
var data = &Collections{}

// http://docs.ansible.com/developing_inventory.html

func init() {
	env.Var(&config.Port, "PORT", "8089", "Port to bind to")
	env.Var(&config.DbName, "DB_NAME", "deepthought", "name of the database")
	env.Var(&config.DbPort, "DB_PORT", "27017", "database port")
	env.Var(&config.DbHost, "DB_HOST", "localhost", "database host")
	env.Var(&config.NodevarsProvidersString, "NODEVARS_PROVIDERS", "", "comma-separated list of nodevars providers")
	env.Var(&config.LdapConnString, "LDAP_CONN", "", "comma-separated list of LDAP Connections")
}

func main() {
	env.Parse("DEEP", false)
	err := config.ParseProviders()
	if err != nil {
		panic(err)
	}

	err = config.ParseLdapConn()
	if err != nil {
		panic(err)
	}

	data = initDatabase()

	gin.SetMode(gin.ReleaseMode)

	g := gin.New()

	g.Use(LogJSON())
	g.Use(SetCORS())
	g.Use(gin.Recovery())

	g.OPTIONS("*path", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := g.Group("/api", BasicAuthLDAP())
	{
		api.GET("/nodes", ListNodes)
		api.POST("/nodes/:node", AddNode)
		api.GET("/nodes/:node", GetNode)
		api.DELETE("/nodes/:node", DeleteNode)
		api.GET("/nodes/:node/vars", GetMergedNodevars)
		api.PUT("/nodes/:node/vars/:var", UpdateNodevars)
		api.GET("/nodes/:node/vars/:var", GetNodevars)
		api.POST("/nodes/:node/providers/:provider", TriggerProvider)

		api.GET("/roles", ListRoles)
		api.POST("/roles/:role", AddRole)
		api.GET("/roles/:role", GetRole)
		api.DELETE("/roles/:role", DeleteRole)
		api.GET("/roles/:role/vars", GetMergedRolevars)
		api.PUT("/roles/:role/vars/:var", UpdateRolevars)
		api.GET("/roles/:role/vars/:var", GetRolevars)

		api.POST("/nodes/:node/roles/:role", LinkNodeWithRole)
		api.POST("/roles/:role/nodes/:node", LinkNodeWithRole)
		api.DELETE("/nodes/:node/roles/:role", UnlinkNodeWithRole)
		api.DELETE("/roles/:role/nodes/:node", UnlinkNodeWithRole)

		api.GET("/_config", GetConfig)
	}

	g.GET("/inventory", GetInventory)

	g.Run(":" + config.Port)

}
