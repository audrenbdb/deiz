package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (r *repo) GetPaymentMethods(ctx context.Context) ([]deiz.PaymentMethod, error) {
	const query = `SELECT id, name FROM payment_method`
	rows, err := r.conn.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var methods []deiz.PaymentMethod
	for rows.Next() {
		var m deiz.PaymentMethod
		err := rows.Scan(&m.ID, &m.Name)
		if err != nil {
			return nil, err
		}
		methods = append(methods, m)
	}
	return methods, nil
}
