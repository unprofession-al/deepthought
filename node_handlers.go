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
			} else if resp.StatusCode != http.StatusOK {
				s, _ := getBodyAsString(resp.Body)
				z := strings.SplitN(s, "\"", 3)
				renderHttp(resp.StatusCode, z[1], c)
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
			renderHttp(http.StatusOK, node, c)
			return
		}
	}

	renderHttp(http.StatusNotFound, "provider not found", c)
}

func DeleteNode(c *gin.Context) {
	nn := c.Params.ByName("node")

	node := &Node{}
	err := data.Nodes.Find(bson.M{"name": nn}).One(node)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	err = data.Nodes.RemoveId(node.Id)
	if err != nil {
		renderHttp(http.StatusInternalServerError, err.Error(), c)
		return
	}

	renderHttp(http.StatusOK, "", c)
}
