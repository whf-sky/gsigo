package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var groups map[string]*group

//get groups
func using(gname string) *group {
	// panic error not database group
	if _, ok := groups[gname]; !ok {
		groups[gname] = &group{}
		groups[gname].get(gname)
	}
	return groups[gname]
}

//database group
type group struct {
	Config 	GroupYml
	client 	*mongo.Client
}

//get gorm group
func (g *group) get(name string) *group {
	//get databases group
	var ok bool
	g.Config, ok = configs[name]
	if !ok {
		panic("database configs be short of '"+name+"'")
	}

	g.client = g.connect(g.Config)
	groups[name] = g
	return g
}

//Open initialize a new db connection
func (g *group) connect(config GroupYml)  *mongo.Client {
	client, err := Connect(config.Uri, config.MaxIdle, config.Timeout)
	if err != nil {
		panic(err)
	}
	return client
}

func init()  {
	groups = map[string]*group{}
}

