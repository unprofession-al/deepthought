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
	env.Var(&config.DbHosts, "DB_HOSTS", "localhost:27017", "database hosts")
	env.Var(&config.DbUser, "DB_USER", "", "database user")
	env.Var(&config.DbPass, "DB_PASS", "", "database password")
	env.Var(&config.NodevarsProvidersString, "NODEVARS_PROVIDERS", "", "comma-separated list of nodevars providers")
}

func main() {
	env.Parse("DEEP", false)
	err := config.ParseProviders()
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

	api := g.Group("/api")
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
