package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func Check(Host, Username, Password string, Port int) error {
	dataSourceName := fmt.Sprintf("mongodb://%v:%v@%v:%v/?authMechanism=SCRAM-SHA-1", Username, Password, Host, Port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 建立mongodb连接
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dataSourceName))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}
