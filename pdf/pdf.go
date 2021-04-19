package pdf

import (
	"github.com/audrenbdb/deiz/intl"
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

type Pdf struct {
	intl *intl.Parser

	fontFamily string
	fontFile   string
	fontDir    string

	blueTheme rgb
}

type ServiceDeps struct {
	Intl *intl.Parser

	FontFamily string
	FontFile   string
	FontDir    string
}

func NewService(deps ServiceDeps) *Pdf {
	blueTheme := rgb{
		red:   0,
		green: 0,
		blue:  70,
	}
	return &Pdf{
		intl:       deps.Intl,
		fontFamily: deps.FontFamily,
		fontFile:   deps.FontFile,
		fontDir:    deps.FontDir,

		blueTheme: blueTheme,
	}
}

func initPDF(doc *gofpdf.Fpdf, headerFunc func(), footerFunc func()) {
	doc.SetHeaderFunc(headerFunc)
	doc.SetFooterFunc(footerFunc)
	doc.AliasNbPages("{nb}")
	doc.AddPage()
}

func (pdf *Pdf) createPDF(o orientation, u unitScale, s pageSize) *gofpdf.Fpdf {
	doc := gofpdf.New(string(o), string(u), string(s), pdf.fontDir)
	doc.SetMargins(20, 80, 20)
	doc.AddUTF8Font(pdf.fontFamily, "", pdf.fontFile)
	doc.SetFont(pdf.fontFamily, "", 12)
	return doc
}
