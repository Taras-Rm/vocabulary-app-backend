package pdfgenerator

import (
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
)

type PdfGenerator struct {
	pdf        *gofpdf.Fpdf
	tr         func(string) string
	outputPath string
}

func NewPdf() *PdfGenerator {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pdf := gofpdf.New("P", "mm", "A4", pwd+"/pkg/pdfGenerator/font")

	pdf.AddFont("Helvetica-1251", "", "helvetica_1251.json")
	pdf.SetFont("Helvetica-1251", "", 14)

	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	outputPah := pwd + "/pkg/pdfGenerator/test.pdf"

	return &PdfGenerator{pdf: pdf, tr: tr, outputPath: outputPah}
}

func GenerateCollectionPdf(contents [][]string, tableName string) ([]byte, error) {
	headings := []string{"Word", "Translation"}

	g := NewPdf()

	g.pdf.AddPage()

	widthOfCell := 100.0

	leftMargin := (210.0 - 2*widthOfCell) / 2
	g.pdf.SetX(leftMargin)

	g.generateWordsListTable(contents, headings, widthOfCell, leftMargin)

	err := g.pdf.OutputFileAndClose(g.outputPath)
	if err != nil {
		return nil, err
	}

	bytes, err := g.prepareFile()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (g *PdfGenerator) generateWordsListTable(contents [][]string, headings []string, cellWidth, leftMargin float64) {
	g.pdf.SetFontSize(16)
	for _, str := range headings {
		g.pdf.CellFormat(cellWidth, 7, g.tr(str), "1", 0, "C", false, 0, "")
	}

	g.pdf.Ln(-1)

	g.pdf.SetFontSize(14)
	for _, c := range contents {
		g.pdf.SetX(leftMargin)
		g.pdf.CellFormat(cellWidth, 6, g.tr(c[0]), "1", 0, "", false, 0, "")
		g.pdf.CellFormat(cellWidth, 6, g.tr(c[1]), "1", 0, "", false, 0, "")
		g.pdf.Ln(-1)
	}
}

func (g *PdfGenerator) prepareFile() ([]byte, error) {
	bytes, err := os.ReadFile(g.outputPath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(g.outputPath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
