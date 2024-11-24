package mongodb

import (
	"context"
	"log"

	"github.com/sera/back-end/worker/internal/config"
	"github.com/sera/back-end/worker/internal/config/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBInterface interface {
	GetCollection(collectionName string) *mongo.Collection
	GetCollectionByName(name string) *mongo.Collection
	CheckDB() (bool, error)
}

type mongodb_pool struct {
	DB                  *mongo.Client
	DBName              string
	DBDefaultCollection string
}

var mdbpool = &mongodb_pool{}
var ctx = context.TODO()

func New(conf *config.Config) MongoDBInterface {

	if mdbpool != nil && mdbpool.DB != nil && mdbpool.DBName != "" {

		return mdbpool

	} else {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MDB_URI))
		if err != nil {
			log.Fatal("Erro to make Connect DB:", err.Error())
			logger.Error("Erro to make Connect DB:"+err.Error(), err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal("Erro to contact DB:", err.Error())
			logger.Error("Erro to contact DB:"+err.Error(), err)
		}

		mdbpool = &mongodb_pool{
			DB:                  client,
			DBName:              conf.MDB_NAME,
			DBDefaultCollection: conf.MDB_DEFAULT_COLLECTION,
		}

	}
	logger.Info("About to start user application")
	return mdbpool
}

func (d *mongodb_pool) GetCollection(collectionName string) *mongo.Collection {
	return d.DB.Database(d.DBName).Collection(collectionName)
}

func (d *mongodb_pool) GetCollectionByName(name string) (DBCollection *mongo.Collection) {
	return d.DB.Database(d.DBName).Collection(name)
}

func ObjectIDFromHex(hex string) (objectID primitive.ObjectID, err error) {
	objectID, err = primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Println(err.Error())
		return objectID, err
	}
	return objectID, nil
}

func (d *mongodb_pool) CheckDB() (bool, error) {
	// Tenta realizar um ping no banco de dados para verificar a conexão
	err := d.DB.Ping(ctx, nil)
	if err != nil {
		logger.Error("Erro ao verificar a conexão com o banco de dados: ", err)
		return false, err
	}
	logger.Info("Conexão com o banco de dados está ativa")
	return true, nil
}
