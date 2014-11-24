package main

import "gopkg.in/mgo.v2/bson"

type AnsibleInventory map[string]*AnsibleGroup

type AnsibleGroup struct {
	Hosts    []AnsibleHost                     `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Vars     map[string]interface{}            `json:"vars,omitempty" yaml:"vars,omitempty"`
	Hostvars map[string]map[string]interface{} `json:"hostvars,omitempty" yaml:"hostvars,omitempty"`
}

type AnsibleHost string

func NewAnsibleInventory() (*AnsibleInventory, error) {
	i := AnsibleInventory{}

	allGroupName := "all"
	i[allGroupName] = &AnsibleGroup{}
	metaGroupName := "_meta"
	i[metaGroupName] = &AnsibleGroup{}
	i[metaGroupName].Hostvars = make(map[string]map[string]interface{})

	nodes := []Node{}
	err := data.Nodes.Find(bson.M{}).All(&nodes)
	if err != nil {
		return &i, err
	}

	roles := []Role{}
	err = data.Roles.Find(bson.M{}).All(&roles)
	if err != nil {
		return &i, err
	}

	for _, role := range roles {
		if role.Name == allGroupName {
			i[allGroupName].Vars = role.Vars.Merge()
		}
	}

	for _, n := range nodes {
		host := AnsibleHost(n.Name)
		for _, rId := range n.Roles {
			r := Role{}
			for _, role := range roles {
				if role.Id == rId {
					r = role
				}
			}

			if _, ok := i[r.Name]; !ok {
				i[r.Name] = &AnsibleGroup{
					Vars: r.Vars.Merge(),
				}
			}

			if r.Name != allGroupName {
				i[r.Name].Hosts = append(i[r.Name].Hosts, host)
			}
		}
		i[allGroupName].Hosts = append(i[allGroupName].Hosts, host)
		i[metaGroupName].Hostvars[n.Name] = n.Vars.Merge()
	}
	return &i, nil
}
