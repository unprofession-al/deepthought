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
	g.Use(SetCORS())
	g.Use(LogJSON())
	g.Use(gin.Recovery())

	g.OPTIONS("*path", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	g.GET("/nodes", ListNodes)
	g.POST("/nodes/:node", AddNode)
	g.GET("/nodes/:node", GetNode)
	g.DELETE("/nodes/:node", DeleteNode)
	g.GET("/nodes/:node/vars", GetMergedNodevars)
	g.PUT("/nodes/:node/vars/:var", UpdateNodevars)
	g.GET("/nodes/:node/vars/:var", GetNodevars)
	g.POST("/nodes/:node/providers/:provider", TriggerProvider)

	g.GET("/roles", ListRoles)
	g.POST("/roles/:role", AddRole)
	g.GET("/roles/:role", GetRole)
	g.DELETE("/roles/:role", DeleteRole)
	g.GET("/roles/:role/vars", GetMergedRolevars)
	g.PUT("/roles/:role/vars/:var", UpdateRolevars)
	g.GET("/roles/:role/vars/:var", GetRolevars)

	g.POST("/nodes/:node/roles/:role", LinkNodeWithRole)
	g.POST("/roles/:role/nodes/:node", LinkNodeWithRole)
	g.DELETE("/nodes/:node/roles/:role", UnlinkNodeWithRole)
	g.DELETE("/roles/:role/nodes/:node", UnlinkNodeWithRole)

	g.GET("/inventory", GetInventory)

	g.GET("/_config", GetConfig)

	g.Run(":" + config.Port)

}
