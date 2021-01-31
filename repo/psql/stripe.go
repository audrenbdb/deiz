package psql

import "context"

type stripeKeys struct {
	id     int
	public string
	secret []byte
}

func getStripeKeysByPersonID(ctx context.Context, db db, personID int) (stripeKeys, error) {
	const query = `SELECT COALESCE(public, ''), secret FROM stripe_keys WHERE person_id = $1`
	row := db.QueryRow(ctx, query, personID)
	var k stripeKeys
	err := row.Scan(&k.public, &k.secret)
	if err != nil {
		return stripeKeys{}, err
	}
	return k, nil
}

func insertStripeKeys(ctx context.Context, db db, keys *stripeKeys, personID int) error {
	const query = `INSERT INTO stripe_keys(person_id, public, secret) VALUES($1, NULLIF($2, ''), $3) RETURNING id`
	row := db.QueryRow(ctx, query, personID, keys.public, keys.secret)
	return row.Scan(&keys.id)
}

func updatePersonStripeKeys(ctx context.Context, db db, k stripeKeys, personID int) error {
	const query = `UPDATE stripe_keys SET public = $1, secret = $2 WHERE person_id = $3`
	cmdTag, err := db.Exec(ctx, query, k.public, k.secret, personID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errNoRowsUpdated
	}
	return nil
}

func (r *repo) GetClinicianStripeSecretKey(ctx context.Context, clinicianID int) ([]byte, error) {
	k, err := getStripeKeysByPersonID(ctx, r.conn, clinicianID)
	if err != nil {
		return nil, err
	}
	return k.secret, nil
}

func (r *repo) EditClinicianStripeKeys(ctx context.Context, pk string, sk []byte, clinicianID int) error {
	k := stripeKeys{
		public: pk,
		secret: sk,
	}
	return updatePersonStripeKeys(ctx, r.conn, k, clinicianID)
}
