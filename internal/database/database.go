package database

import (
	"github.com/the-great-checkout/transactions-crud/internal/entity"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Mongo struct {
	Collection *mongo.Collection
}

type Postgres struct {
	DB *gorm.DB
}

func NewMongo(uri, databaseName, collectionName string) Mongo {
	client, _ := mongo.Connect(options.Client().ApplyURI(uri))
	collection := client.Database(databaseName).Collection(collectionName)

	return Mongo{
		collection,
	}
}

func NewPostgres(dsn, schemaName string) Postgres {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   schemaName,
			SingularTable: false,
		}})
	if err != nil {
		panic(err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = db.AutoMigrate(&entity.Status{}, &entity.Transaction{})
	if err != nil {
		panic(err)
	}

	err = db.Exec(`
        INSERT INTO transactions.statuses (id, name) VALUES 
        (uuid_generate_v4(), 'created'), 
        (uuid_generate_v4(), 'pending'), 
        (uuid_generate_v4(), 'completed'), 
        (uuid_generate_v4(), 'deleted')
        ON CONFLICT DO NOTHING;
    `).Error
	if err != nil {
		panic(err)
	}

	return Postgres{
		db,
	}
}
