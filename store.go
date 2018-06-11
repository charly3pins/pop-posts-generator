package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
)

func db() *sql.DB {
	db, err := sql.Open("mysql", "root:root@/op-figures")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	return db
}

func store() {
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

	db := db()
	stmt, err := db.Prepare("INSERT figures SET `number`=?,`name`=?,`character`=?,`category`=?,`subcategory`=?,`sculptor`=?,`officialprice`=?,`preorderdate`=?,`releasedate`=?,`reedition1`=?,`reedition2`=?,`height`=?,`weight`=?,`boxsize`=?,`observations`=?")
	if err != nil {
		panic(err)
	}

	for _, f := range figures {
		res, err := stmt.Exec(f.Number, f.Name, f.Character, f.Category, f.Subcategory, f.Sculptor, f.OfficialPrice, f.PreorderDate, f.ReleaseDate, f.Reedition1, f.Reedition2, f.Height, f.Weight, f.BoxSize, f.Observations)
		if err != nil {
			panic(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		log.Printf("Inserted %s with id %d\n", f.Number, id)
	}

	log.Println("all figures inserted")
}

func read() {
	db := db()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM figures")
	if err != nil {
		panic(err)
	}
	var figures []Figure
	for rows.Next() {
		var f Figure
		err = rows.Scan(&f.ID, &f.Number, &f.Name, &f.Character, &f.Category, &f.Subcategory, &f.Sculptor, &f.OfficialPrice, &f.PreorderDate, &f.ReleaseDate, &f.Reedition1, &f.Reedition2, &f.Height, &f.Weight, &f.BoxSize, &f.Observations)
		if err != nil {
			panic(err)
		}
		figures = append(figures, f)
	}

	// Define a template.
	poster := `[[items]]
number = "{{.Number}}"
name = "{{.Name}}"
character = "{{.Character}}"
category = "{{.Category}}"
subcategory = "{{.Subcategory}}"
sculptor = "{{.Sculptor}}"
officialprice = "{{.OfficialPrice}} ¥"
preorderdate = "{{.PreorderDate}}"
releasedate = "{{.ReleaseDate}}"
reedition1 = "{{.Reedition1}}"
reedition2 = "{{.Reedition2}}"
height = "{{.Height}} (cm)"
weight = "{{.Weight}} (g)"
boxsize = "{{.BoxSize}} (cm)"
observations = "{{.Observations}}"
image = "/onepiecefigures/images/poster/{{.Number}}.jpg"
thumb = "/onepiecefigures/images/poster/{{.Number}}.jpg"
alt = "{{.Name}}"
class = "{{if eq .Subcategory "Limited Lawson"}}limited-lawson{{else if eq .Category "Original Series"}}original-series{{else if eq .Category "Neo"}}neo{{else if eq .Category "Neo EX"}}neo-ex{{else if eq .Category "Neo DX"}}neo-dx{{else if eq .Category "CB"}}cb{{else if eq .Category "Mugiwara Theater"}}mugiwara-theater{{else if eq .Category "Strong Edition"}}strong-edition{{else if eq .Category "Stuffed Collection"}}stuffed-collection{{else if eq .Category "Strong Edition Limited Lawson"}}strong-edition-limited-lawson{{else if eq .Category "Limited Edition"}}limited-edition{{else if eq .Category "Sailing Again"}}sailing-again{{else if eq .Category "Maximum"}}maximum{{else if eq .Category "Edition Z"}}edition-z{{else if eq .Category "M.A.S"}}mas{{else if eq .Category "Kabuki Edition"}}kabuki-edition{{else if eq .Category "I.R.O"}}iro{{else if eq .Category "S.O.C"}}soc{{end}}"
`
	t := template.Must(template.New("poster").Parse(poster))
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
