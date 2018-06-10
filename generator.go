package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/tealeg/xlsx"
)

func generate() {
	fileLog, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file")
	}
	defer fileLog.Close()

	log.SetOutput(fileLog)

	excelFileName := "POP-Guide-DB.xlsx"
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
					figure.Number = cell.String()
				case 1:
					figure.Name = cell.String()
				case 2:
					figure.Character = cell.String()
				case 3:
					figure.Category = cell.String()
				case 4:
					figure.Subcategory = cell.String()
				case 5:
					figure.Sculptor = cell.String()
				case 6:
					figure.OfficialPrice = cell.String()
				case 7:
					figure.PreorderDate = cell.String()
				case 8:
					figure.ReleaseDate = cell.String()
				case 9:
					figure.Reedition1 = cell.String()
				case 10:
					figure.Reedition2 = cell.String()
				case 11:
					figure.Height = cell.String()
				case 12:
					figure.Weight = cell.String()
				case 13:
					figure.BoxSize = cell.String()
				case 14:
					figure.Observations = cell.String()
				default:
					fmt.Printf("Column [%d] not parsed. Value of column [%s]\n", i, cell.String())
				}
			}
			figures = append(figures, figure)
		}
	}

	// Define a template.
	post := `[[items]]
title = "{{.Number}} - {{.Name}}"
image = "/onepiecefigures/images/poster/{{.Number}}.jpg"
thumb = "/onepiecefigures/images/poster/{{.Number}}.jpg"
alt = "{{.Name}}"
description = "<b>Number:</b> {{.Number}}<br><b>Name:</b> {{.Name}}<br><b>Character:</b> {{.Character}}<br><b>Category:</b> {{.Category}} {{if .Subcategory}}{{.Subcategory}}{{end}}<br><b>Sculptor:</b> {{.Sculptor}}<br><b>Official price:</b> {{if .OfficialPrice}}{{.OfficialPrice}} ¥{{end}}<br><b>Preorder date:</b> {{.PreorderDate}}<br><b>Release date:</b> {{.ReleaseDate}}{{if .Reedition1}}<br><b>Reeditions:</b> {{.Reedition1}}{{if .Reedition2}}, {{.Reedition2}}{{end}}<br><b>Height:</b> {{if .Height}}{{.Height}} (cm){{end}}<br><b>Weight:</b> {{if .Weight}}{{.Weight}} (g){{end}}<br><b>Box size:</b> {{if .BoxSize}}{{.BoxSize}} (cm){{end}}{{else}}<br><b>Height:</b> {{if .Height}}{{.Height}} (cm){{end}}<br><b>Weight:</b> {{if .Weight}}{{.Weight}} (g){{end}}<br><b>Box size:</b> {{if .BoxSize}}{{.BoxSize}} (cm){{end}}{{end}}{{if .Observations}}<br><b>Bonus:</b> {{.Observations}}{{end}}"
`
	t := template.Must(template.New("post").Parse(post))
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range figures {
		if f.Name == "" {
			continue
		}
		file, err := os.Create(fmt.Sprintf("%s/toml/%s.toml", dir, f.Number))
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

	files, err := ioutil.ReadDir(fmt.Sprintf("%s/toml/", dir))
	if err != nil {
		log.Println("Error reading dir: ", err)
		return
	}

	var buffer bytes.Buffer
	for _, file := range files {
		log.Println(file.Name())
		fileContent, err := ioutil.ReadFile(fmt.Sprintf("%s/toml/%s", dir, file.Name()))
		if err != nil {
			log.Println("Error reading file: ", err)
			continue
		}

		buffer.WriteString(string(fileContent))
	}

	f, err := os.Create(fmt.Sprintf("%s/items.toml", dir))
	if err != nil {
		log.Println("Error writting final file: ", err)
		return
	}
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(buffer.String())
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()

	log.Println(buffer.String())
}
