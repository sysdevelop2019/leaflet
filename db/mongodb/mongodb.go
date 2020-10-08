package mongodb

import (
	"context"
	"fmt"
	"gitee.com/aarlin/leaflet/log"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"time"
)

const (
	DatabaseName        = "gameserver"
	DBTimeOutDuration   = time.Second * 30
)

type DBIndex struct {
	Key  map[string]int
	NS   string
	Name string
}

type MongoOptions struct {
	DBUser string
	DBPwd  string
	DBName string
	DBAddr string
	DBPort string
}

type MongoRepo struct {
	MongoOptions
	Client *mongo.Client

}

func NewMongoRepo(options MongoOptions) (*MongoRepo, error) {
	repo := &MongoRepo{
		MongoOptions: options,
	}
	//mongodb://[username:password@]host1[:port1][,...hostN[:portN]]][/[database][?options]]
	dsn := fmt.Sprintf("mongodb://%v:%v@%v:%v", repo.DBUser, repo.DBPwd, repo.DBAddr, repo.DBPort)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	client, err := mongo.Connect(ctx, dsn)
	repo.Client = client

	if len(repo.DBName) < 0 {
		repo.DBName = DatabaseName
	}

	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *MongoRepo) GetClient()(client *mongo.Client)  {
	return r.Client
}


func (r *MongoRepo) CreateDefaultIndex(tableName string,  keys []string) (err error) {
	log.DebugF("MongoRepo CreateIndex: %v keys%v", tableName,keys)
	indexView := r.Client.Database(r.DBName).Collection(tableName).Indexes()
	if cursor, err := indexView.List(context.Background()); err == nil {
		found := false
		idxName := ""
		var indexKeys bsonx.Doc
		for i := 0; i < len(keys);i++ {
			if i == 0 {
				idxName = keys[i] + "_1"
			}else{
				idxName = idxName + "_" + keys[i] + "_1"
			}
			elem := bsonx.Elem{
				Key:keys[i],
				Value:bsonx.Int32(1),
			}
			indexKeys = append(indexKeys,elem)
		}
		//log.Debug("indexName", idxName)
		for !found && cursor.Next(context.Background()) {
			var idx DBIndex
			if err := cursor.Decode(&idx); err == nil {
				//log.Debug("indexName", idx.Name)
				if idx.Name == idxName {
					found = true
				}
			}
		}
		//log.Debug("found", found)
		if !found {
			//log.Debug("CreateDbIndex")
			indexName, err := indexView.CreateOne(
				context.Background(),
				mongo.IndexModel{
					Keys:    indexKeys,
					Options: options.Index().SetBackground(true),
				},
			)
			if err != nil {
				log.Error("CreateDbIndex", "", "indexName", indexName, "err", err)
			}
			//
		}
	}
	return
}


func (r *MongoRepo) InsertOne(tableName string,  document interface{},opts ...*options.InsertOneOptions) (err error) {
	log.DebugF("MongoRepo InsertOne: %v", document)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)

	if _, err = collection.InsertOne(ctx, document,opts ...); err != nil {
		log.Error("InsertOne error %v", err.Error())
		return
	}
	return
}
func (r *MongoRepo) InsertMany(tableName string,  document []interface{},opts ...*options.InsertManyOptions) (err error) {
	log.DebugF("MongoRepo InsertMany: %v", document)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)

	if _, err = collection.InsertMany(ctx, document,opts ...); err != nil {
		log.Error("InsertOne error ", err.Error())
		return
	}
	return
}

func (r *MongoRepo) FindOne(tableName string,  filter interface{},opts ...*options.FindOneOptions) *mongo.SingleResult {
	log.DebugF("MongoRepo FindOne: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return collection.FindOne(ctx, filter, opts...)
}

func (r *MongoRepo) Find(tableName string,  filter interface{},opts ...*options.FindOptions) (context.Context,mongo.Cursor, error) {
	log.DebugF("MongoRepo Find: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	cur,err := collection.Find(ctx, filter, opts...)
	return ctx,cur,err
}

func (r *MongoRepo) Count(tableName string,  filter interface{}, opts ...*options.CountOptions) (int64, error) {
	log.DebugF("MongoRepo Count: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return collection.Count(ctx, filter, opts...)
}

func (r *MongoRepo) UpdateOne(tableName string,  filter interface{},update interface{},opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	log.DebugF("MongoRepo UpdateOne: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	//if len(opts) == 0 {
	//	isUp := true
	//	opts = append(opts, &options.UpdateOptions{Upsert:&isUp})
	//}
	return collection.UpdateOne(ctx, filter, update, opts ...)
}

func (r *MongoRepo) UpdateOneOrInsert(tableName string,  filter interface{},update interface{},opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	log.DebugF("MongoRepo UpdateOne: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	if len(opts) == 0 {
		isUp := true
		opts = append(opts, &options.UpdateOptions{Upsert:&isUp})
	}
	return collection.UpdateOne(ctx, filter, update, opts ...)
}

func (r *MongoRepo) ReplaceOne(tableName string,  filter interface{},replacement interface{},opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	log.DebugF("MongoRepo ReplaceOne: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return collection.ReplaceOne(ctx, filter, replacement, opts ...)
}

func (r *MongoRepo) FindOneAndReplace(tableName string,  filter interface{},replacement interface{},opts ...*options.FindOneAndReplaceOptions)  *mongo.SingleResult {
	log.DebugF("MongoRepo FindOneAndReplace: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return collection.FindOneAndReplace(ctx, filter, replacement, opts ...)
}

func (r *MongoRepo) DeleteOne(tableName string,  filter interface{}, opts ...*options.DeleteOptions) ( *mongo.DeleteResult, error) {
	log.DebugF("MongoRepo DeleteOne: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return collection.DeleteOne(ctx, filter, opts...)
}

func (r *MongoRepo) DeleteMany(tableName string,  filter interface{}, opts ...*options.DeleteOptions) ( *mongo.DeleteResult, error) {
	log.DebugF("MongoRepo DeleteMany: %v", filter)
	collection := r.Client.Database(r.DBName).Collection(tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return collection.DeleteMany(ctx, filter, opts...)
}

func (r *MongoRepo) Drop(tableName string) error {
	log.DebugF("MongoRepo Drop: %v", tableName)
	ctx, _ := context.WithTimeout(context.Background(), DBTimeOutDuration)
	return r.Client.Database(r.DBName).Collection(tableName).Drop(ctx)
}

