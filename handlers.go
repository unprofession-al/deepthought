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
	err := parseBody(c, &vars)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
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
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}
	if count > 0 {
		renderHttp(http.StatusInternalServerError, "node already exists", c)
		return
	}

	for _, provider := range config.NodevarsProviders {
		url := strings.Replace(provider.Url, "{{nodename}}", node.Name, -1)
		resp, err := http.Get(url)
		if err != nil {
			renderHttp(http.StatusInternalServerError, err.Error(), c)
			return
		}

		var pvars map[string]interface{}
		err = getJSONBodyAsStruct(resp.Body, &pvars)
		if err != nil {
			renderHttp(http.StatusInternalServerError, err.Error(), c)
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
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, node, c)
}

func ListNodes(c *gin.Context) {
	nodes := &[]Node{}

	err := data.Nodes.Find(bson.M{}).All(nodes)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}
	renderHttp(http.StatusOK, nodes, c)
}

func UpdateNodevars(c *gin.Context) {
	nn := c.Params.ByName("node")
	vn := c.Params.ByName("var")

	var vars map[string]interface{}
	err := parseBody(c, &vars)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
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
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	node.Vars.AddOrReplace(v)

	err = data.Nodes.UpdateId(node.Id, node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, node, c)
}

func GetNode(c *gin.Context) {
	nn := c.Params.ByName("node")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nn}).One(node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, node, c)
}

func GetMergedNodevars(c *gin.Context) {
	nn := c.Params.ByName("node")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nn}).One(node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	vars := node.Vars.Merge()
	
	renderHttp(http.StatusOK, vars, c)
}

func GetNodevars(c *gin.Context) {
	nn := c.Params.ByName("node")
	vn := c.Params.ByName("var")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nn}).One(node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	for _, vars := range node.Vars {
		if vars.Source == vn {
			renderHttp(http.StatusOK, vars.Vars, c)
			return
		}
	}
	renderHttp(http.StatusNotFound, "vars not found", c)
}

func TriggerProvider(c *gin.Context) {
	p := c.Params.ByName("provider")
	n := c.Params.ByName("node")

	node := &Node{}

	err := data.Nodes.Find(bson.M{"name": n}).One(node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	for _, provider := range config.NodevarsProviders {
		if provider.Name == p {
			url := strings.Replace(provider.Url, "{{nodename}}", node.Name, -1)
			resp, err := http.Get(url)
			if err != nil {
				renderHttp(http.StatusInternalServerError, err.Error(), c)
				return
			}

			var pvars map[string]interface{}
			err = getJSONBodyAsStruct(resp.Body, &pvars)
			if err != nil {
				renderHttp(http.StatusInternalServerError, err.Error(), c)
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
				renderHttp(http.StatusInternalServerError, err.Error(), c)
				return
			}

			c.JSON(http.StatusOK, node)
			renderHttp(http.StatusOK, node, c)
			return
		}
	}

	renderHttp(http.StatusNotFound, "provider not found", c)
}

func AddRole(c *gin.Context) {
	name := c.Params.ByName("role")

	var vars map[string]interface{}
	err := parseBody(c, &vars)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
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
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}
	if count > 0 {
		renderHttp(http.StatusInternalServerError, "role already exists", c)
		return
	}

	err = data.Roles.Insert(role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, role)
}

func ListRoles(c *gin.Context) {
	roles := &[]Role{}

	err := data.Roles.Find(bson.M{}).All(roles)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}
	renderHttp(http.StatusOK, roles, c)
}

func UpdateRolevars(c *gin.Context) {
	rn := c.Params.ByName("role")
	vn := c.Params.ByName("var")

	var vars map[string]interface{}
	err := parseBody(c, &vars)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
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
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	role.Vars.AddOrReplace(v)

	err = data.Roles.UpdateId(role.Id, role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, role, c)
}

func GetRolevars(c *gin.Context) {
	rn := c.Params.ByName("role")
	vn := c.Params.ByName("var")

	role := &Role{}
	err := data.Roles.Find(bson.M{"name": rn}).One(role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	for _, vars := range role.Vars {
		if vars.Source == vn {
			renderHttp(http.StatusOK, vars.Vars, c)
			return
		}
	}
	renderHttp(http.StatusNotFound, "vars not found", c)
}

func GetRole(c *gin.Context) {
	rn := c.Params.ByName("role")

	role := &Role{}
	err := data.Roles.Find(bson.M{"name": rn}).One(role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, role, c)
}

func LinkNodeWithRole(c *gin.Context) {
	roleName := c.Params.ByName("role")
	nodeName := c.Params.ByName("node")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nodeName}).One(node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	role := &Role{}
	err = data.Roles.Find(bson.M{"name": roleName}).One(role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
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
			renderHttp(http.StatusInternalServerError, err.Error(), c)
			return
		}
	}

	renderHttp(http.StatusOK, node, c)
}

func GetInventory(c *gin.Context) {
	i, err := NewAnsibleInventory()
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}
	renderHttp(http.StatusOK, i, c)
}

func GetConfig(c *gin.Context) {
	renderHttp(http.StatusOK, config, c)
}
