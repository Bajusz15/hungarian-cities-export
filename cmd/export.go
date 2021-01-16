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

	var locations []City

	for _, f := range d.Features {
		zip, _ := strconv.Atoi(f.Properties.PostalCode)
		c := City{
			Name:      f.Properties.Name,
			Zip:       zip,
			Longitude: f.Geometry.Coordinates[0],
			Latitude:  f.Geometry.Coordinates[1],
		}
		locations = append(locations, c)
	}
	//log.Println(locations)
	exportToJSON(locations)
	exportToPostgres(locations)
}

func exportToJSON(locations []City) {
	export := struct {
		Cities []City `json:"locations"`
	}{Cities: locations}
	JSONExport, _ := json.Marshal(export)
	_ = ioutil.WriteFile("export.json", JSONExport, os.ModePerm)
}

//postgres schema
//CREATE TABLE location (
//id SERIAL PRIMARY KEY ,
//name VARCHAR(25),
//zip INT NOT NULL DEFAULT 0,
//coordinates geometry NOT NULL DEFAULT 'point(0 0)'
//);
func exportToPostgres(locations []City) {
	query := ""
	for _, val := range locations {
		query += fmt.Sprintf("INSERT INTO location (name,zip, coordinates) VALUES ('%s', %d, 'point(%v  %v)' );", val.Name, val.Zip, val.Longitude, val.Latitude)
	}
	data := []byte(query)
	_ = ioutil.WriteFile("export.sql", data, os.ModePerm)
}
