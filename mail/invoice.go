package mail

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
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
		b.Clinician.Email,
		"Facture de consultation",
		&emailBuffer, plainBody,
		invoicePDF))
}
