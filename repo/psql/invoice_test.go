package psql

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestInsertInvoice(t *testing.T) {
	invoice := invoice{
		personID:        1,
		identifier:      "000000001",
		sender:          []string{"toto"},
		recipient:       []string{"tata"},
		cityAndDate:     "TEST",
		label:           "TOTO",
		priceAfterTax:   4000,
		deliveryDate:    time.Now(),
		taxFee:          19.6,
		exemption:       "test",
		paymentMethodID: 1,
	}
	ctx := context.Background()
	pool, _ := pgxpool.Connect(ctx, os.Getenv("TEST_DATABASE_URL"))
	err := insertInvoice(ctx, pool, &invoice)
	assert.NoError(t, err)
	assert.Positive(t, invoice.id)
}
