package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

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

func GetMergedRolevars(c *gin.Context) {
	rn := c.Params.ByName("role")

	role := &Role{}
	err := data.Roles.Find(bson.M{"name": rn}).One(role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	vars := role.Vars.Merge()

	renderHttp(http.StatusOK, vars, c)
}

func DeleteRole(c *gin.Context) {
	rn := c.Params.ByName("role")

	role := &Role{}
	err := data.Roles.Find(bson.M{"name": rn}).One(role)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	nodes := &[]Node{}

	err = data.Nodes.Find(bson.M{}).All(nodes)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	for _, node := range *nodes {
		for i, rId := range node.Roles {
			if role.Id == rId {
				node.Roles = append(node.Roles[:i], node.Roles[i+1:]...)
			}
		}

		err = data.Nodes.UpdateId(node.Id, node)
		if err != nil {
			renderHttp(http.StatusInternalServerError, err.Error(), c)
			return
		}
	}

	err = data.Roles.RemoveId(role.Id)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, "", c)
}
