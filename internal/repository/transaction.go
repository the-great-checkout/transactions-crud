package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/the-great-checkout/transactions-crud/internal/database"
	"github.com/the-great-checkout/transactions-crud/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db         *gorm.DB
	collection *mongo.Collection
}

func NewTransactionRepository(postgresDB database.Postgres, mongoDB database.Mongo) *TransactionRepository {
	return &TransactionRepository{postgresDB.DB, mongoDB.Collection}
}

func (r *TransactionRepository) Create(transaction *entity.Transaction) error {
	var status entity.Status
	r.db.Where("name = ?", "created").First(&status)

	transaction.Status = status

	if err := r.db.Create(transaction).Error; err != nil {
		return err
	}

	_, err := r.collection.InsertOne(context.TODO(), transaction)
	return err
}

func (r *TransactionRepository) FindByID(id uuid.UUID) (*entity.Transaction, error) {
	var transaction entity.Transaction
	if err := r.db.Preload("Status").First(&transaction, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) FindAll() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	if err := r.db.Preload("Status").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) Update(transaction *entity.Transaction) error {
	existingTransaction := &entity.Transaction{}
	if err := r.db.First(existingTransaction, "id = ?", transaction.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("transaction not found")
		}
		return err
	}

	var status entity.Status
	r.db.Where("name = ?", transaction.Status.Name).First(&status)

	existingTransaction.Status = status
	existingTransaction.Value = transaction.Value

	if err := r.db.Save(existingTransaction).Error; err != nil {
		return err
	}

	filter := bson.M{"_id": existingTransaction}

	update := bson.M{"$set": existingTransaction}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *TransactionRepository) Delete(id uuid.UUID) (*entity.Transaction, error) {
	var status entity.Status
	r.db.Where("name = ?", "deleted").First(&status)

	// Manually update the is_deleted and deleted_at fields
	var transaction entity.Transaction
	result := r.db.Model(&transaction).Where("id = ?", id).Updates(map[string]any{
		"is_deleted": true,
		"deleted_at": time.Now(), // Set current time for deleted_at
		"status":     status,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedTransaction entity.Transaction
	err := r.db.Unscoped().Where("id = ?", id).First(&updatedTransaction).Error
	if err != nil {
		return nil, err
	}

	return &updatedTransaction, nil
}
