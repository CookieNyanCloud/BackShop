package repository

import "github.com/jmoiron/sqlx"

type OrdersRepo struct {
	db *sqlx.DB
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{db: db}
}

func (r *OrdersRepo) Create(ctx context.Context, order domain.Order) error {
	_, err := r.db.InsertOne(ctx, order)

	return err
}

func (r *OrdersRepo) AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error) {
	var order domain.Order

	res := r.db.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status": transaction.Status,
		},
		"$push": bson.M{
			"transactions": transaction,
		},
	})
	if res.Err() != nil {
		return order, res.Err()
	}

	err := res.Decode(&order)

	return order, err
}