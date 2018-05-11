package main

import (
	"encoding/csv"
	"flag"
	"log"
	
	"github.com/mikedata/dp-drr-result-writer/models"
	"github.com/nu7hatch/gouuid"
	"io"
	"os"
	"encoding/json"
)

var sourceFile = flag.String("filepath", "", "The path to the file being uploaded.")

func main() {

	flag.Parse()

	if *sourceFile == "" {
		log.Fatal("Aborting. No upload file specified.")
	}

	fileIn, err := os.Open(*sourceFile)
	if err != nil {
		log.Fatal("Aborting. Unable to load csv: " + *sourceFile)
	}

	csvReader := csv.NewReader(fileIn) // TODO - cant load a string

	// Scan for header row, this information will need to be sent to the
	// dataset API with the number of observations in a PUT request
	headerRow, err := csvReader.Read()
	if err != nil {
		log.Fatal("Aborting. Encountered error when processing header row")
	}

	// create a UUID for source task
	sourceUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to genrate a UUID to represent source file.")
	}

	// Create the source Msg
	// We will NOT send this until the task queue has been populated.
	sourceMsg := &models.MsgSource{Source: *sourceFile, Uuid: sourceUUID.String()}

	// Populate task messages
	// ----------------------

	// Cache things we've seen before - we only want unique items
	datasetCache := make(map[string][]string)
	var emptyDim []string
	for i := range headerRow {
		datasetCache[headerRow[i]] = emptyDim
	}

	type task struct {

	}

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV file.")
		}

		for i := 0; i < len(headerRow); i++ {

			cellVal := line[i]
			dimName := headerRow[i]
			seenBefore := false

			for x := 0; x < len(datasetCache[dimName]); x++ {
				if line[i] == datasetCache[dimName][x] {
					seenBefore = true
				}
			}

			// if we haven't seen it before. Cache it and add a message to task list
			if !seenBefore {
				datasetCache[dimName] = append(datasetCache[dimName], cellVal)
				task := make(map[string]string)
				task[dimName] = cellVal

				taskJson, err := json.Marshal(task)
				if err != nil {
					log.Fatal("Failed to marshall json response.")
				}
			}
		}

	}

	sourceJson, err := json.Marshal(sourceMsg)
	if err != nil {
		log.Fatal("Failed to marshall json response.")
	}
}
