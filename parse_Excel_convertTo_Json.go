package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"platformOps-EC/models"
	"platformOps-EC/services"
)

func main() {

	var excelFileName, output string

	flag.StringVar(&excelFileName, "i", "", "Input Excel baseline file. If missing, program will exit.")
	flag.StringVar(&output, "o", "manifest.json", "Execution output location.")
	flag.Parse()

	if excelFileName == "" {
		fmt.Println("Missing input Excel file. Program will exit.")
		os.Exit(1)
	}

	if output == "manifest.json" {
		fmt.Println("Default to manifest.json")

	}

	fmt.Println("Loading Excel file ", excelFileName)

	baseline, controls := services.LoadFromExcel(excelFileName)
	var manifest []models.EC_Manifest

	fmt.Println("Converting to Json object")

	for _, c := range controls {

		m := models.EC_Manifest{Title: c.GetCategory(),
			Baseline: baseline.GetName() + " at Req# " + c.GetRow_desc()}
		manifest = append(manifest, m)

	}

	//fmt.Println(models.ToJson(manifest))

	file, err := os.Create(output)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Println("Writing Json Object to file")

	fmt.Fprintf(file, "%v", models.ToJson(manifest))

	fmt.Println("Done writing to output file at ", output)

}