/*
Package usecase references all usecases to be implemented
*/
package usecase

import (
	"context"
	"github.com/audrenbdb/deiz"
	"time"
)

type (
	BillingUsecases struct {
		InvoiceCreater       InvoiceCreater
		InvoiceCanceler      InvoiceCanceler
		InvoiceMailer        InvoiceMailer
		InvoicesGetter       InvoicesGetter
		StripeSessionCreater StripeSessionCreater
		UnpaidBookingsGetter UnpaidBookingsGetter
	}
)

type (
	InvoiceCreater interface {
		CreateInvoice(ctx context.Context, invoice *deiz.BookingInvoice, sendToPatient bool) error
	}
	InvoiceCanceler interface {
		CancelInvoice(ctx context.Context, invoice *deiz.BookingInvoice) error
	}
	InvoiceMailer interface {
		MailInvoice(invoice *deiz.BookingInvoice, recipient string) error
		MailInvoicesSummary(ctx context.Context, start, end time.Time, recipient string, clinicianID int) error
	}
	InvoicesGetter interface {
		GetPeriodInvoices(ctx context.Context, start, end time.Time, clinicianID int) ([]deiz.BookingInvoice, error)
	}
	StripeSessionCreater interface {
		CreateStripePaymentSession(ctx context.Context, amount int64, clinicianID int) (string, error)
	}
	UnpaidBookingsGetter interface {
		GetUnpaidBookings(ctx context.Context, clinicianID int) ([]deiz.Booking, error)
	}
)
