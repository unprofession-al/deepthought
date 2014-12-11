package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

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

func UnlinkNodeWithRole(c *gin.Context) {
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

	for i, r := range node.Roles {
		if r == role.Id {
			node.Roles = append(node.Roles[:i], node.Roles[i+1:]...)

			err = data.Nodes.UpdateId(node.Id, node)
			if err != nil {
				renderHttp(http.StatusInternalServerError, err.Error(), c)
				return
			}
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
