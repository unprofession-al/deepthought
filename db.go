package main

import "gopkg.in/mgo.v2"

type Collections struct {
	Roles mgo.Collection
	Nodes mgo.Collection
}

func initDatabase() *Collections {
	session, err := mgo.Dial(config.DbHost + ":" + config.DbPort)
	if err != nil {
		panic(err)
	}

	nodes := session.DB(config.DbName).C("nodes")
	roles := session.DB(config.DbName).C("roles")

	return &Collections{
		Nodes: *nodes,
		Roles: *roles,
	}
}
