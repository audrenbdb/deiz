package mail

import (
	"bytes"
	"fmt"
	"github.com/audrenbdb/deiz"
	"time"
)

func (m *Mailer) MailBookingInvoice(invoice *deiz.BookingInvoice, invoicePDF *bytes.Buffer, sendTo string) error {
	details := getInvoiceEmailDetails(invoice)
	template, err := m.htmlTemplate("booking-invoice.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBody()
	return m.client.Send(createMail(mail{
		to:       sendTo,
		from:     noReplyAddress,
		subject:  "Facture de consultation",
		template: template, plainBody: plainBody,
		attachment: invoicePDF}))
}

func (m *Mailer) MailInvoicesSummary(summaryPDF *bytes.Buffer, start, end time.Time, sendTo string) error {
	details := m.getInvoicesEmailDetails(start, end)
	template, err := m.htmlTemplate("invoices-summary.html", details)
	if err != nil {
		return err
	}
	plainBody := details.plainBody()
	return m.client.Send(createMail(mail{
		to:       sendTo,
		from:     noReplyAddress,
		subject:  "Résumé de factures",
		template: template, plainBody: plainBody,
		attachment: summaryPDF,
	}))
}

func (details *invoiceEmailDetails) plainBody() string {
	return fmt.Sprintf(`Deiz\n
	Nouvelle facture\n
	\n
	De %s\n
	Montant : %s\n
	Faite à %s\n
	\n
	Deiz\n
	\Agenda pour thérapeutes\n
	https://deiz.fr`, details.Name, details.Amount, details.Date)
}

func getInvoiceEmailDetails(invoice *deiz.BookingInvoice) invoiceEmailDetails {
	return invoiceEmailDetails{
		Name:   invoice.Booking.Clinician.FullName(),
		Date:   invoice.DeliveryDateStr,
		Amount: fmt.Sprintf("%.2f€", float64(invoice.PriceAfterTax)/100),
	}
}

type invoicesEmailDetail struct {
	Start string
	End   string
}

type invoiceEmailDetails struct {
	Name   string
	Date   string
	Amount string
}

func (m *Mailer) getInvoicesEmailDetails(start, end time.Time) invoicesEmailDetail {
	return invoicesEmailDetail{
		Start: m.intl.Fr.FmtyMd(start),
		End:   m.intl.Fr.FmtyMd(end),
	}
}

func (details *invoicesEmailDetail) plainBody() string {
	return fmt.Sprintf(`Deiz\n
	Résumé de factures\n
	\n
	Du %s\n
	Au %s\n
	\n
	Deiz\n
	\Agenda pour thérapeutes\n
	https://deiz.fr`, details.Start, details.End)
}
