package pdf

import (
	"bytes"
	"context"
	"fmt"
	"github.com/audrenbdb/deiz"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

func (pdf *pdf) GenerateBookingInvoicePDF(ctx context.Context, i *deiz.BookingInvoice) (*bytes.Buffer, error) {
	doc := pdf.createPDF(portrait, mm, A4)
	p := message.NewPrinter(language.French)
	initPDF(doc,
		headerAsLetterFunc(
			doc,
			i.CityAndDate,
			i.Sender,
			i.Recipient,
			pdf.blueTheme,
			20),
		footerFunc(doc, pdf.blueTheme),
	)

	//title
	doc.Ln(1)
	doc.SetTextColor(pdf.blueTheme.red, pdf.blueTheme.green, pdf.blueTheme.blue)
	doc.Ln(6)
	doc.CellFormat(0, 0, fmt.Sprintf("Identifiant de facture : %s", i.Identifier), "", 0, "R", false, 0, "")
	doc.Ln(20)
	doc.SetFillColor(pdf.blueTheme.red, pdf.blueTheme.green, pdf.blueTheme.blue)
	doc.SetTextColor(255, 255, 255)
	doc.Cell(1, 0, "")
	doc.CellFormat(56, 10, "Date", "", 0, "C", true, 0, "")
	doc.Cell(1, 0, "")
	doc.CellFormat(56, 10, "Prestation", "", 0, "C", true, 0, "")
	doc.Cell(1, 0, "")
	doc.CellFormat(56, 10, "Prix unitaire", "", 0, "C", true, 0, "")
	doc.SetTextColor(pdf.blueTheme.red, pdf.blueTheme.green, pdf.blueTheme.blue)

	//product details
	doc.Ln(12)
	doc.Cell(1, 0, "")
	doc.CellFormat(56, 20, i.DeliveryDateStr, "", 0, "C", false, 0, "")
	doc.Cell(1, 0, "")
	doc.CellFormat(56, 20, i.Label, "", 0, "C", false, 0, "")
	doc.Cell(1, 0, "")
	doc.CellFormat(56, 20, p.Sprintf("%.2f €", float64(i.PriceBeforeTax)/100), "", 1, "C", false, 0, "")

	//total
	if doc.GetY() > 180 {
		doc.AddPage()
		doc.Ln(10)
	}
	doc.SetDrawColor(0, 0, 70)
	doc.Cell(87, 0, "")
	doc.CellFormat(84, 1, "", "B", 0, "L", false, 0, "")
	doc.Ln(10)
	doc.Cell(102, 0, "")
	doc.CellFormat(50, 10, "Total hors taxes :", "", 0, "L", false, 0, "")
	doc.CellFormat(15, 10, p.Sprintf("%.2f", float64(i.PriceBeforeTax)/100), "", 0, "R", false, 0, "")
	doc.CellFormat(5, 10, "€", "", 1, "R", false, 0, "")
	doc.Ln(1)
	doc.Cell(102, 0, "")
	doc.CellFormat(50, 10, "T.V.A :", "", 0, "L", false, 0, "")
	doc.CellFormat(15, 10, p.Sprintf("%.2f", float32(i.PriceBeforeTax/100)*i.TaxFee/100), "", 0, "R", false, 0, "")
	doc.CellFormat(5, 10, "€", "", 1, "R", false, 0, "")
	doc.Ln(1)
	doc.Cell(102, 0, "")
	doc.CellFormat(50, 10, "Total T.T.C :", "", 0, "L", false, 0, "")
	doc.CellFormat(15, 10, p.Sprintf("%.2f", float64(i.PriceAfterTax)/100), "", 0, "R", false, 0, "")
	doc.CellFormat(5, 10, "€", "", 1, "R", false, 0, "")
	doc.Ln(10)
	doc.Cell(87, 0, "")
	doc.CellFormat(84, 1, "", "B", 0, "L", false, 0, "")
	doc.Ln(10)
	doc.CellFormat(171, 10, fmt.Sprintf("Acquitté ce jour via %s", i.PaymentMethod.Name), "", 1, "R", false, 0, "")
	doc.Ln(20)

	//tax exemption
	if i.Exemption != "" {
		doc.CellFormat(171, 10, fmt.Sprintf("TVA non applicable - article %s du CGI", i.Exemption), "", 1, "C", false, 0, "")
	}
	var buffer bytes.Buffer
	err := doc.Output(&buffer)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	return &buffer, nil
}

func (pdf *pdf) GetPeriodBookingInvoicesSummaryPDF(ctx context.Context, invoices []deiz.BookingInvoice, start, end time.Time, totalBeforeTax, totalAfterTax int64, clinicianTz *time.Location) (*bytes.Buffer, error) {
	doc := pdf.createPDF(landscape, mm, A4)
	p := message.NewPrinter(language.French)
	initPDF(doc,
		headerAsPeriodEarningsSummary(doc,
			start.In(clinicianTz).Format("02/01/2006"),
			end.In(clinicianTz).Format("02/01/2006"), totalBeforeTax, totalAfterTax, pdf.blueTheme, 20),
		footerFunc(doc, pdf.blueTheme),
	)

	for _, i := range invoices {
		doc.SetFontSize(8)
		doc.SetTextColor(pdf.blueTheme.red, pdf.blueTheme.green, pdf.blueTheme.blue)
		doc.SetDrawColor(pdf.blueTheme.red, pdf.blueTheme.green, pdf.blueTheme.blue)
		doc.Cell(10, 0, "")
		doc.CellFormat(33, 4, i.CreatedAt.In(clinicianTz).Format("02/01/2006"), "", 0, "C", false, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 4, i.DeliveryDate.In(clinicianTz).Format("02/01/2006"), "", 0, "C", false, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 4, i.Identifier, "", 0, "C", false, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(67, 4, i.Recipient[0], "", 0, "C", false, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 4, p.Sprintf("%.2f €", float64(i.PriceBeforeTax)/100), "", 0, "C", false, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 4, i.PaymentMethod.Name, "", 0, "C", false, 0, "")
		doc.Ln(6)
	}
	var buffer bytes.Buffer
	err := doc.Output(&buffer)
	if err != nil {
		return &bytes.Buffer{}, err
	}
	return &buffer, nil
}
