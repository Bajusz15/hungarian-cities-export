package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Data struct {
	Features Feature `json:"features"`
}

type Feature []struct {
	Type       string `json:"type"`
	Properties struct {
		ID                   string `json:"@id"`
		KshRef               string `json:"ksh_ref"`
		Name                 string `json:"name"`
		NameRu               string `json:"name:ru"`
		Place                string `json:"place"`
		Population           string `json:"population"`
		PostalCode           string `json:"postal_code"`
		SourcePopulation     string `json:"source:population"`
		SourcePopulationDate string `json:"source:population:date"`
		SourcePostalCode     string `json:"source:postal_code"`
		SourcePostalCodeDate string `json:"source:postal_code:date"`
	} `json:"properties"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	ID string `json:"id"`
}

type City struct {
	Name      string  `json:"name"`
	Zip       int     `json:"zip"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func main() {
	file, err := os.Open("source.json")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	var d Data
	byteValue, _ := ioutil.ReadAll(file)
	_ = json.Unmarshal(byteValue, &d)

	var cities []City

	for _, f := range d.Features {
		zip, _ := strconv.Atoi(f.Properties.PostalCode)
		c := City{
			Name:      f.Properties.Name,
			Zip:       zip,
			Longitude: f.Geometry.Coordinates[0],
			Latitude:  f.Geometry.Coordinates[1],
		}
		cities = append(cities, c)
	}
	//log.Println(cities)
	exportToJSON(cities)
	exportToPostgres(cities)
}

func exportToJSON(cities []City) {
	export := struct {
		Cities []City `json:"cities"`
	}{Cities: cities}
	JSONExport, _ := json.Marshal(export)
	_ = ioutil.WriteFile("export.json", JSONExport, os.ModePerm)
}

func exportToPostgres(cities []City) {
	query := ""
	for _, val := range cities {
		query += fmt.Sprintf("INSERT INTO cities (name,zip,longitude,latitude) ('%s', %d, %v, %v);", val.Name, val.Zip, val.Longitude, val.Latitude)
	}
	data := []byte(query)
	_ = ioutil.WriteFile("export.sql", data, os.ModePerm)
}
