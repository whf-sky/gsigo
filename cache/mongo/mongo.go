package mongo
//
//import (
//	"context"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	//"go.mongodb.org/mongo-driver/mongo/readpref"
//	"time"
//)
//
//
//func NewMongo(gname string) *Mongo{
//	return (&Mongo{}).Using(gname)
//}
//
//type Mongo struct {
//	group *group
//	clinet *mongo.Client
//}
//
////Using use group
//func (m *Mongo) Using(gname string) *Mongo{
//	m.group = using(gname)
//	m.clinet = m.group.client
//	return m
//}
//
////UseDb use database
//func (m *Mongo) UseDb(dbname string) *MongoDb{
//	return &MongoDb{db:m.group.client.Database(dbname)}
//}
//
////Client mongo client
//func (m *Mongo) Client() *mongo.Client{
//	return m.clinet
//}
//
//
//type MongoDb struct {
//	db *mongo.Database
//	collection map[string]*mongo.Collection
//}
//
////Client mongo client
//func (m *MongoDb) Collection() *mongo.Client{
//	m.collection["xxx"]
//	//err = client.Ping(ctx, readpref.Primary())
//	collection := client.Database("testing").Collection("numbers")
//	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
//	res, _ := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
//	_ = res.InsertedID
//	return m.clinet
//}
//
//////collection
//func Connect(uri string, maxIdle uint64, timeout int)  (*mongo.Client, error) {
//	//new mongo client
//	client, err1 := mongo.NewClient(options.Client().ApplyURI(uri).SetMaxPoolSize(maxIdle))
//	if err1 != nil {
//		return nil, err1
//	}
//	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
//	//connect mongo
//	err2 := client.Connect(ctx)
//
//	if err2 != nil {
//		return nil, err2
//	}
//	return client, nil
//}
