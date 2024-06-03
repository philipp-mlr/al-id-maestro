package service

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/philipp-mlr/al-id-maestro/model"
)

func ReadLicensedObjectsCSV(path string) *model.Objects {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	return loadLicensedObjects(data)

}

func loadLicensedObjects(data [][]string) *model.Objects {
	var alObjects model.Objects
	for i, line := range data {
		if i > 0 { // omit header line
			var rec model.Object
			for j, field := range line {
				if j == 0 {
					rec.ObjectType.Name = MapIntToType(field)
				} else if j == 1 {
					id, err := strconv.Atoi(field)
					if err != nil {
						log.Fatalf("Error converting ID to int in CSV file")
					}
					rec.ID = uint(id)
				}
			}
			alObjects.Objects = append(alObjects.Objects, rec)
		}
	}
	return &alObjects
}

func MapIntToType(id string) string {
	switch id {
	case "1":
		return "table"
	case "3":
		return "report"
	case "5":
		return "codeunit"
	case "6":
		return "xmlport"
	case "8":
		return "page"
	case "9":
		return "query"
	case "16":
		return "enum"
	case "14":
		return "pageextension"
	case "15":
		return "tableextension"
	case "22":
		return "reportextension"
	case "17":
		return "enumextension"
	case "20":
		return "permissionset"
	case "21":
		return "permissionsetextension"
	}
	return "unknown"
}
