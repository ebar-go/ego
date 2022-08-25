package component

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	Named
	*mongo.Client
	database string
}

// MongoConfig represents the configuration of a MongoDB connection
type MongoConfig struct {
	Address    string // mongo address, like: mongodb://127.0.0.1:27017
	Username   string // username
	Password   string
	Database   string
	AuthSource string        // default is admin
	Timeout    time.Duration // default is 30 seconds
}

func (c *MongoConfig) complete() {
	if c.AuthSource == "" {
		c.AuthSource = "admin"
	}
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
}

// Connect connects to the server
func (m *Mongo) Connect(config *MongoConfig, timeout time.Duration) error {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(config.Address)
	if len(config.Username) > 0 {
		credential := options.Credential{AuthSource: config.AuthSource, Username: config.Username, Password: config.Password}
		clientOptions = clientOptions.SetAuth(credential)
	}
	clientOptions.ConnectTimeout = &timeout
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	m.Client = client
	m.database = config.Database
	return nil
}

// Close closes the connection to the
func (m *Mongo) Close() error {
	return m.Disconnect(context.Background())
}

// Database use the other database
func (m *Mongo) Database(database string, opts ...*options.DatabaseOptions) *mongo.Database {
	return m.Database(database, opts...)
}

// Collection use current database for the provide collection
func (m *Mongo) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return m.Database(m.database).Collection(name, opts...)
}

func NewMongo() *Mongo {
	return &Mongo{Named: componentMongo}
}
