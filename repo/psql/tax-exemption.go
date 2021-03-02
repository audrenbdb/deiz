package psql

import (
	"context"
	"github.com/audrenbdb/deiz"
)

func (r *repo) GetTaxExemptionCodes(ctx context.Context) ([]deiz.TaxExemption, error) {
	const query = `SELECT id, code FROM tax_exemption`
	rows, err := r.conn.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var exemptions []deiz.TaxExemption
	for rows.Next() {
		var t deiz.TaxExemption
		err := rows.Scan(&t.ID, &t.Code)
		if err != nil {
			return nil, err
		}
		exemptions = append(exemptions, t)
	}
	return exemptions, nil
}
