package main

import (
	"email-indexer/globals"
	"email-indexer/helpers"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const path string = "/Users/diegoherrera/Downloads/enron_mail_20110402/"
const nChunks int = 10

var wg = sync.WaitGroup{}

func main() {

	start := time.Now()
	defer helpers.TimeTrack(start, "BulkDB")

	//GENERAR SLICE CON TODOS LOS PATHS DE EMAILS EN LA DB
	fileList := []string{}

	error := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() && d.Name() != ".DS_Store" && d.Name() != "DELETIONS.txt" {
			fileList = append(fileList, path)
		}
		return nil
	})

	if error != nil {
		fmt.Println(error)
	}

	//GENERAR SLICE DE CHUNKS PARA PROCESAR CON GORUTINA
	chunks := helpers.ChunkSlice(fileList, nChunks)
	fmt.Printf("---------- NO. OF CHUNKS: %v\n---------- NO. OF FILES:  %v\n", len(chunks), len(fileList))

	//SE GENERA DOCUMENTO JSON VALIDO PARA CARGAR BULK POR CHUNKS USANDO CONCURRENCIA
	wg.Add(nChunks)
	for idx, chunk := range chunks {
		go uploadChunk(chunk, idx)
	}

	wg.Wait()
}

func uploadChunk(chunk []string, i int) {

	tempFile := "./temp/bulk" + strconv.FormatInt(int64(i+1), 10) + ".json"

	f, _ := os.OpenFile(tempFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for idx, path := range chunk {
		email, err := helpers.GenEmail(path)

		if err != nil {
			fmt.Println(err)
		}

		emailJSON, _ := json.Marshal(email)
		if idx == 0 {
			f.WriteString(`{"index": "` + globals.ZINC_INDEX + `", "records": [` + string(emailJSON) + "," + "\n")
		} else if idx == len(chunk)-1 {
			f.WriteString(string(emailJSON) + "]}" + "\n")
		} else {
			f.WriteString(string(emailJSON) + "," + "\n")
		}
	}

	f.Close()

	//Realizar Bulk
	helpers.BulkFile(tempFile)

	wg.Done()
}
