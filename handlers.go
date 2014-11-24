package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

func AddNode(c *gin.Context) {
	name := c.Params.ByName("node")

	var vars map[string]interface{}
	err := getJSONBodyAsStruct(c.Request.Body, &vars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	v := Vars{
		Prio:   0,
		Vars:   vars,
		Source: "native",
	}

	node := &Node{Name: name}
	node.Vars.AddOrReplace(v)

	count, err := data.Nodes.Find(bson.M{"name": name}).Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if count > 0 {
		c.JSON(http.StatusInternalServerError, "node already exists")
		return
	}

	for _, provider := range config.NodevarsProviders {
		url := strings.Replace(provider.Url, "{{nodename}}", node.Name, -1)
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var pvars map[string]interface{}
		err = getJSONBodyAsStruct(resp.Body, &pvars)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		providerVars := Vars{
			Prio:   provider.Prio,
			Source: provider.Name,
			Vars:   pvars,
		}

		node.Vars.AddOrReplace(providerVars)
	}

	err = data.Nodes.Insert(node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, node)
}

func ListNodes(c *gin.Context) {
	nodes := &[]Node{}

	err := data.Nodes.Find(bson.M{}).All(nodes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func UpdateNodevars(c *gin.Context) {
	nn := c.Params.ByName("node")
	vn := c.Params.ByName("var")

	var vars map[string]interface{}
	err := getJSONBodyAsStruct(c.Request.Body, &vars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	v := Vars{
		Prio:   0,
		Vars:   vars,
		Source: vn,
	}

	node := &Node{}
	err = data.Nodes.Find(bson.M{"name": nn}).One(node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	node.Vars.AddOrReplace(v)

	err = data.Nodes.UpdateId(node.Id, node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, node)
}

func GetNodevars(c *gin.Context) {
	nn := c.Params.ByName("node")
	vn := c.Params.ByName("var")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nn}).One(node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, vars := range node.Vars {
		if vars.Source == vn {
			c.JSON(http.StatusOK, vars.Vars)
			return
		}
	}
	c.JSON(http.StatusNotFound, "vars not found")
}

func TriggerProvider(c *gin.Context) {
	p := c.Params.ByName("provider")
	n := c.Params.ByName("node")

	node := &Node{}

	err := data.Nodes.Find(bson.M{"name": n}).One(node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, provider := range config.NodevarsProviders {
		if provider.Name == p {
			url := strings.Replace(provider.Url, "{{nodename}}", node.Name, -1)
			resp, err := http.Get(url)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			var pvars map[string]interface{}
			err = getJSONBodyAsStruct(resp.Body, &pvars)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			providerVars := Vars{
				Prio:   provider.Prio,
				Source: provider.Name,
				Vars:   pvars,
			}
			node.Vars.AddOrReplace(providerVars)

			err = data.Nodes.UpdateId(node.Id, node)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			c.JSON(http.StatusOK, node)
			return
		}
	}

	c.JSON(http.StatusNotFound, "provider not found")
}

func AddRole(c *gin.Context) {
	name := c.Params.ByName("role")

	var vars map[string]interface{}
	err := getJSONBodyAsStruct(c.Request.Body, &vars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	v := Vars{
		Prio:   0,
		Vars:   vars,
		Source: "native",
	}

	role := &Role{Name: name}
	role.Vars.AddOrReplace(v)

	count, err := data.Roles.Find(bson.M{"name": name}).Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if count > 0 {
		c.JSON(http.StatusInternalServerError, "role already exists")
		return
	}

	err = data.Roles.Insert(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, role)
}

func ListRoles(c *gin.Context) {
	roles := &[]Role{}

	err := data.Roles.Find(bson.M{}).All(roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, roles)
}

func UpdateRolevars(c *gin.Context) {
	rn := c.Params.ByName("role")
	vn := c.Params.ByName("var")

	var vars map[string]interface{}
	err := getJSONBodyAsStruct(c.Request.Body, &vars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	v := Vars{
		Prio:   0,
		Vars:   vars,
		Source: vn,
	}

	role := &Role{}
	err = data.Roles.Find(bson.M{"name": rn}).One(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	role.Vars.AddOrReplace(v)

	err = data.Roles.UpdateId(role.Id, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, role)
}

func GetRolevars(c *gin.Context) {
	rn := c.Params.ByName("role")
	vn := c.Params.ByName("var")

	role := &Role{}
	err := data.Roles.Find(bson.M{"name": rn}).One(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for _, vars := range role.Vars {
		if vars.Source == vn {
			c.JSON(http.StatusOK, vars.Vars)
			return
		}
	}
	c.JSON(http.StatusNotFound, "vars not found")
}

func LinkNodeWithRole(c *gin.Context) {
	roleName := c.Params.ByName("role")
	nodeName := c.Params.ByName("node")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nodeName}).One(node)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	role := &Role{}
	err = data.Roles.Find(bson.M{"name": roleName}).One(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	hasRole := false
	for _, r := range node.Roles {
		if r == role.Id {
			hasRole = true
		}
	}

	if !hasRole {
		node.Roles = append(node.Roles, role.Id)
		err = data.Nodes.UpdateId(node.Id, node)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, node)
}

func GetInventory(c *gin.Context) {
	i, err := NewAnsibleInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, i)
}

func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, config)
}
