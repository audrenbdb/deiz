package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type rgb struct {
	red   int
	green int
	blue  int
}

func headerAsLetterFunc(doc *gofpdf.Fpdf, cityAndDate string, sender []string, recipient []string, colour rgb, headerMargin float64) func() {
	return func() {
		currentFontSize, _ := doc.GetFontSize()
		doc.SetFontSize(12)
		doc.SetTextColor(colour.red, colour.green, colour.blue)
		doc.SetY(10)
		doc.CellFormat(0, 5, cityAndDate, "", 0, "R", false, 0, "")
		doc.Ln(headerMargin)
		for _, l := range sender {
			doc.Cell(0, 0, l)
			doc.Ln(6)
		}
		doc.Ln(4)
		for _, l := range recipient {
			doc.CellFormat(0, 0, l, "", 0, "R", false, 0, "")
			doc.Ln(6)
		}
		doc.SetFontSize(currentFontSize)
	}
}

func headerAsPeriodEarningsSummary(doc *gofpdf.Fpdf, startDateStr, endDateStr string, totalBeforeTax, totalAfterTax int64, colour rgb, headerMargin float64) func() {
	return func() {
		p := message.NewPrinter(language.French)
		currentFontSize, _ := doc.GetFontSize()
		doc.SetFontSize(12)
		doc.SetTextColor(colour.red, colour.green, colour.blue)
		doc.SetY(10)
		doc.Ln(headerMargin)
		doc.Cell(0, 0, "Entre le "+startDateStr)
		doc.Ln(6)
		doc.Cell(0, 0, "Et le "+endDateStr)
		doc.Ln(6)
		doc.Cell(0, 0, "Chiffre d'affaire : "+p.Sprintf("%.2f €", float64(totalBeforeTax)/100))
		doc.Ln(6)
		doc.Cell(0, 0, "Revenus net de T.V.A : "+p.Sprintf("%.2f €", float64(totalAfterTax)/100))
		doc.SetFontSize(currentFontSize)

		doc.Ln(12)
		doc.SetFontSize(10)
		doc.SetFillColor(0, 0, 70)
		doc.SetTextColor(255, 255, 255)
		doc.Cell(10, 0, "")
		doc.CellFormat(33, 6, "Date de facture", "", 0, "C", true, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 6, "Date du rdv", "", 0, "C", true, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 6, "Identifiant", "", 0, "C", true, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(67, 6, "Destinataire", "", 0, "C", true, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 6, "Montant TTC", "", 0, "C", true, 0, "")
		doc.Cell(1, 0, "")
		doc.CellFormat(33, 6, "Modalité", "", 0, "C", true, 0, "")
		doc.Ln(10)
	}
}

func footerFunc(doc *gofpdf.Fpdf, colour rgb) func() {
	return func() {
		currentFontSize, _ := doc.GetFontSize()
		doc.SetFontSize(12)
		doc.SetTextColor(colour.red, colour.green, colour.blue)
		doc.SetY(-15)
		doc.CellFormat(0, 10, fmt.Sprintf("page %d/{nb}", doc.PageNo()), "", 0, "C", false, 0, "")
		doc.SetFontSize(currentFontSize)
	}
}
