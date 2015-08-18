package main

import "gopkg.in/mgo.v2"

type Collections struct {
	Roles mgo.Collection
	Nodes mgo.Collection
}

func initDatabase() *Collections {
	session, err := mgo.Dial(config.DbHosts)
	if err != nil {
		panic(err)
	}

	sess := session.DB(config.DbName)

	if config.DbPass != "" && config.DbUser != "" {
		err = sess.Login(config.DbUser, config.DbPass)
		if err != nil {
			panic(err)
		}
	}

	nodes := sess.C("nodes")
	roles := sess.C("roles")

	return &Collections{
		Nodes: *nodes,
		Roles: *roles,
	}
}
