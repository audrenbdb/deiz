package pdf

import (
	"github.com/jung-kurt/gofpdf"
)

type orientation string
type pageSize string
type unitScale string

const (
	portrait  orientation = "Portrait"
	landscape orientation = "Landscape"
	A4        pageSize    = "A4"
	mm        unitScale   = "mm"
)

type pdf struct {
	fontFamily string
	fontFile   string
	fontDir    string

	blueTheme rgb
}

func NewService(fontFamily, fontFile, fontDir string) *pdf {
	blueTheme := rgb{
		red:   0,
		green: 0,
		blue:  70,
	}
	return &pdf{
		fontFamily: fontFamily,
		fontFile:   fontFile,
		fontDir:    fontDir,

		blueTheme: blueTheme,
	}
}

func initPDF(doc *gofpdf.Fpdf, headerFunc func(), footerFunc func()) {
	doc.SetHeaderFunc(headerFunc)
	doc.SetFooterFunc(footerFunc)
	doc.AliasNbPages("{nb}")
	doc.AddPage()
}

func (pdf *pdf) createPDF(o orientation, u unitScale, s pageSize) *gofpdf.Fpdf {
	doc := gofpdf.New(string(o), string(u), string(s), pdf.fontDir)
	doc.SetMargins(20, 80, 20)
	doc.AddUTF8Font(pdf.fontFamily, "", pdf.fontFile)
	doc.SetFont(pdf.fontFamily, "", 12)
	return doc
}
