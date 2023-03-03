package pdfgenerator

import (
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type PdfGenerator struct {
	pdf pdf.Maroto
}

func NewPdfGenerator() PdfGenerator {
	instance := pdf.NewMaroto(consts.Portrait, consts.A4)

	instance.AddUTF8Font("CustomArial", consts.Normal, "pkg/pdfGenerator/fonts/arial-unicode-ms.ttf")
	instance.AddUTF8Font("CustomArial", consts.Italic, "pkg/pdfGenerator/fonts/arial-unicode-ms.ttf")
	instance.AddUTF8Font("CustomArial", consts.Bold, "pkg/pdfGenerator/fonts/arial-unicode-ms.ttf")
	instance.AddUTF8Font("CustomArial", consts.BoldItalic, "pkg/pdfGenerator/fonts/arial-unicode-ms.ttf")
	instance.SetDefaultFontFamily("CustomArial")

	return PdfGenerator{pdf: instance}
}

func GenerateCollectionPdf(contents [][]string, tableName, pathToFile string) ([]byte, error) {
	instance := NewPdfGenerator()

	instance.pdf.SetPageMargins(20, 20, 20)

	instance.generateHeader()

	instance.generateWordsList(contents, tableName)

	res, err := instance.pdf.Output()
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func (p *PdfGenerator) generateHeader() {
	p.pdf.RegisterHeader(func() {
		p.pdf.Row(10, func() {
			p.pdf.Col(12, func() {
				p.pdf.Text("Vocabulary")
			})
		})
	})
}

func (p *PdfGenerator) generateWordsList(contents [][]string, tableName string) {
	tableHeadings := []string{"Word", "Translation"}
	lightPurpleColor := color.Color{
		Red:   210,
		Green: 200,
		Blue:  230,
	}

	p.pdf.SetBackgroundColor(color.Color{
		Red:   3,
		Green: 166,
		Blue:  166,
	})

	p.pdf.Row(10, func() {
		p.pdf.Col(12, func() {
			p.pdf.Text(tableName, props.Text{
				Top:    2,
				Size:   15,
				Color:  color.NewWhite(),
				Family: consts.Courier,
				Style:  consts.Bold,
				Align:  consts.Center,
			})
		})
	})

	p.pdf.SetBackgroundColor(color.NewWhite())

	p.pdf.TableList(tableHeadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      12,
			GridSizes: []uint{6, 6},
		},
		ContentProp: props.TableListContent{
			Size:      12,
			GridSizes: []uint{6, 6},
		},
		Align:                consts.Left,
		AlternatedBackground: &lightPurpleColor,
		HeaderContentSpace:   2,
		Line:                 false,
	})
}
