package mail

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

func (m *mailer) MailBookingInvoice(ctx context.Context, invoice *deiz.BookingInvoice, invoicePDF *bytes.Buffer, sendTo string) error {
	b := invoice.Booking
	var emailBuffer bytes.Buffer
	emailData := struct {
		Name   string
		Date   string
		Amount string
	}{
		Name:   b.Clinician.Surname + " " + b.Clinician.Name,
		Date:   invoice.DeliveryDateStr,
		Amount: fmt.Sprintf("%.2f€", float64(invoice.PriceAfterTax)/100),
	}
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "booking-invoice.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Deiz\n
	Nouvelle facture\n
	\n
	De %s\n
	Montant : %s\n
	Faite à %s\n
	\n
	Deiz\n
	\Agenda pour thérapeutes\n
	https://deiz.fr`, emailData.Name, emailData.Amount, emailData.Date)
	return m.sender.Send(ctx, createMail(
		sendTo,
		noReplyAddress,
		"Facture de consultation",
		&emailBuffer, plainBody,
		invoicePDF))
}

func (m *mailer) MailInvoicesSummary(ctx context.Context, summaryPDF *bytes.Buffer, start, end time.Time, tz *time.Location, sendTo string) error {
	var emailBuffer bytes.Buffer
	emailData := struct {
		Start string
		End   string
	}{
		Start: start.In(tz).Format("02/01/2006"),
		End:   end.In(tz).Format("02/01/2006"),
	}
	err := m.tmpl.ExecuteTemplate(&emailBuffer, "invoices-summary.html", emailData)
	if err != nil {
		return err
	}
	plainBody := fmt.Sprintf(`Deiz\n
	Résumé de factures\n
	\n
	Du %s\n
	Au %s\n
	\n
	Deiz\n
	\Agenda pour thérapeutes\n
	https://deiz.fr`, emailData.Start, emailData.End)
	return m.sender.Send(ctx, createMail(
		sendTo,
		noReplyAddress,
		"Résumé de factures",
		&emailBuffer, plainBody,
		summaryPDF))
}
