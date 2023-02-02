package main

import (
	"email-indexer/globals"
	"email-indexer/helpers"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

const nChunks int = 10

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Ingrese un path válido!")
	}

	path := os.Args[1]

	start := time.Now()
	defer helpers.TimeTrack(start, "BulkDB")

	//SE ASEGURA QUE EL TEMPDIR SEA ÚNICO
	if _, err := os.Stat(globals.TEMPDIR); err == nil {
		if err := os.RemoveAll(globals.TEMPDIR); err != nil {
			log.Fatal(err)
		}
	}

	err := os.Mkdir(globals.TEMPDIR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

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

	//SE GENERA DOCUMENTO JSON VALIDO Y SE CARGAR BULK POR CHUNKS USANDO CONCURRENCIA
	wg.Add(nChunks)

	for idx, chunk := range chunks {
		go helpers.UploadChunk(chunk, idx, &wg)
	}

	wg.Wait() //Se bloquea la ejecución hasta que se terminen todos los threads

	//SE BORRA EL DIRECTORIO TEMPORAL
	if err := os.RemoveAll(globals.TEMPDIR); err != nil {
		log.Fatal(err)
	}
}
