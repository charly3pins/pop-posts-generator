package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/tealeg/xlsx"
)

type Figure struct {
	ID            int
	Number        string
	Name          string
	Character     string
	Category      string
	Subcategory   string
	Sculptor      string
	OfficialPrice string
	PreorderDate  string
	ReleaseDate   string
	Reedition1    string
	Reedition2    string
	Height        string
	Weight        string
	BoxSize       string
	Observations  string
}

func main() {
	excelFileName := "Libro1.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Println("error opening excel", err)
	}
	figures := []Figure{}
	for _, sheet := range xlFile.Sheets {
		for r, row := range sheet.Rows {
			if r == 0 {
				continue
			}
			figure := Figure{}
			for i, cell := range row.Cells {
				switch i {
				case 0:
					figure.Name = cell.String()
				case 1:
					figure.Character = cell.String()
				case 2:
					figure.Category = cell.String()
				case 3:
					figure.Subcategory = cell.String()
				case 4:
					figure.Sculptor = cell.String()
				case 5:
					figure.OfficialPrice = cell.String()
				case 6:
					figure.PreorderDate = cell.String()
				case 7:
					figure.ReleaseDate = cell.String()
				case 8:
					figure.Reedition1 = cell.String()
				case 9:
					figure.Reedition2 = cell.String()
				case 10:
					figure.Height = cell.String()
				case 11:
					figure.Weight = cell.String()
				case 12:
					figure.BoxSize = cell.String()
				case 13:
					figure.Observations = cell.String()
				default:
					fmt.Printf("Column [%d] not parsed. Value of column [%s]\n", i, cell.String())
				}
			}
			figures = append(figures, figure)
		}
	}
	//fmt.Println("final figures\n", figures)

	// Define a template.
	const post = `+++
banner = ""
categories = ["figure"]
date = "2018-01-05T00:00:00Z"
description = ""
images = []
tags = ["onepiece", "portrait of pirates", "{{.Sculptor}}"]
title = "{{.Name}}"
+++

**Name:** {{.Name}}

**Character:** {{.Character}}

**Category:** {{.Category}} {{if .Subcategory}} {{.Subcategory}} {{end}}

**Sculptor:** {{.Sculptor}}

**Official price:** {{if .OfficialPrice}}{{.OfficialPrice}} ¥{{end}}

**Preorder date:** {{.PreorderDate}}

**Release date:** {{.ReleaseDate}}

{{if .Reedition1}}**Reeditions:** {{.Reedition1}}{{if .Reedition2}}, {{.Reedition2}}{{end}}

**Height:** {{if .Height}}{{.Height}} (cm){{end}}

**Weight:** {{if .Weight}}{{.Weight}} (g){{end}}

**Box size:** {{if .BoxSize}}{{.BoxSize}} (cm){{end}}

{{else}}**Height:** {{if .Height}}{{.Height}} (cm){{end}}

**Weight:** {{if .Weight}}{{.Weight}} (g){{end}}

**Box size:** {{if .BoxSize}}{{.BoxSize}} (cm){{end}}{{end}}
{{if .Observations}}

**Bonus:** {{.Observations}}{{end}}
`
	t := template.Must(template.New("post").Parse(post))
	for _, f := range figures {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		file, err := os.Create(fmt.Sprintf("%s/generated/%s.md", dir, f.Name))
		if err != nil {
			log.Println("Error creating file: ", err)
			return
		}

		err = t.Execute(file, f)
		if err != nil {
			log.Print("Error executing template: ", err)
			return
		}

		file.Close()

	}

}
