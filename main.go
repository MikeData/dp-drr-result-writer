package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/mikedata/dp-drr-uploader/messager"
	"github.com/mikedata/dp-drr-uploader/models"
	"github.com/mikedata/dp-drr-uploader/s3"
	"github.com/nu7hatch/gouuid"
	"io"
	"os"
)

var (
	sourceFile           = flag.String("filepath", "", "The path to the file being uploaded.")
	sqs_source_queue_url = os.Getenv("SQS_SOURCE_QUEUE_URL")
	sqs_task_queue_url   = os.Getenv("SQS_TASK_QUEUE_URL")
	aws_region           = os.Getenv("AWS_REGION")
	bucket               = os.Getenv("DRR_IMPORT_BUCKET")
)

func main() {

	flag.Parse()

	if *sourceFile == "" {
		log.Fatal("Aborting. No upload file specified.")
	}

	fileIn, err := os.Open(*sourceFile)
	if err != nil {
		log.Fatal("Aborting. Unable to load csv: " + *sourceFile)
	}

	fmt.Println("here")
	// load to s3
	err = s3.UploadSource(fileIn, *sourceFile, aws_region)
	if err != nil {
		log.Fatal("Issue uploading to s3", err)
	}

	fmt.Println("here")

	// read in the header row
	csvReader := csv.NewReader(fileIn)
	headerRow, err := csvReader.Read()
	if err != nil {
		log.Fatal("Aborting. Encountered error when processing header row")
	}

	// create a UUID for source task
	sourceUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to genrate a UUID to represent source file.")
	}

	// create a cache for things we've seen before - we only want unique items
	datasetCache := make(map[string][]string)
	var emptyDim []string
	for i := range headerRow {
		datasetCache[headerRow[i]] = emptyDim
	}

	// marshall all the tasks into a list before sending any
	var taskList [][]byte

	// iterate and create our unique {dimension:item} tasks
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error reading CSV file.", err)
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
					log.Fatal("Error marshalling task message to json.", err)
				}

				taskList = append(taskList, taskJson)
			}
		}
	}

	// Send all the task messages
	for i := 0; i < len(taskList); i++ {
		messager.SendMsg(taskList[i], sqs_task_queue_url)
	}

	// Send the source message
	sourceMsg := &models.MsgSource{Source: *sourceFile, SourceId: sourceUUID.String()}
	sourceJson, err := json.Marshal(sourceMsg)
	if err != nil {
		log.Fatal("Error marshalling source message to json", err)
	}
	messager.SendMsg(sourceJson, sqs_source_queue_url)

}
