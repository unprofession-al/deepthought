package main

import (
	"sort"

	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	Id    bson.ObjectId   `bson:"_id,omitempty" json:"id"`
	Name  string          `bson:"name" json:"name"`
	Roles []bson.ObjectId `bson:"roles" json:"roles"`
	Vars  VarsBucket      `bson:"vars,omitempty" json:"vars"`
}

type Role struct {
	Id   bson.ObjectId `bson:"_id,omitempty" json:"id" yaml:"id"`
	Name string        `bson:"name" json:"name" yaml:"name"`
	Vars VarsBucket    `bson:"vars,omitempty"  json:"vars,omitempty" yaml:"vars,omitempty"`
}

type Vars struct {
	Source string                 `bson:"source" json:"source" yaml:"source"`
	Prio   int                    `bson:"prio" json:"prio" yaml:"prio"`
	Vars   map[string]interface{} `bson:"vars,omitempty" json:"vars,omitempty" yaml:"vars,omitempty"`
}

type VarsBucket []Vars

func (v VarsBucket) Len() int           { return len(v) }
func (v VarsBucket) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v VarsBucket) Less(i, j int) bool { return v[i].Prio > v[j].Prio }

func (v *VarsBucket) Merge() map[string]interface{} {
	sort.Sort(v)
	merged := make(map[string]interface{})
	for _, vars := range *v {
		for k, v := range vars.Vars {
			merged[k] = v
		}
	}
	return merged
}

func (v *VarsBucket) AddOrReplace(n Vars) {
	for i, vars := range *v {
		if vars.Source == n.Source {
			(*v)[i] = n
			return
		}
	}
	*v = append(*v, n)
	return
}
