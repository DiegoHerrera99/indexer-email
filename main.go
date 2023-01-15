package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

/*type Email struct {
	Date    string `json:"date"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Cc      string `json:"cc"`
	Body    string `json:"body"`
}

type Bulk struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}*/

const path string = "/Users/diegoherrera/Downloads/enron_mail_20110402/"
const nChunks int = 10

// const index string = "enron"

var wg = sync.WaitGroup{}

func main() {

	start := time.Now()
	defer timeTrack(start, "gen&proccessChunks GO")

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
	chunks := chunkSlice(fileList, nChunks)
	fmt.Printf("---------- NO. OF CHUNKS: %v\n---------- NO. OF FILES:  %v\n", len(chunks), len(fileList))

	//TODO: Procesar con Go Rutina!!!!
	wg.Add(nChunks)
	for idx, chunk := range chunks {
		go processChunk(chunk, idx)
	}

	wg.Wait()

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("---------- TIME OF EXEC(%v): %s \n", name, elapsed)
}

func chunkSlice(slice []string, nChunks int) [][]string {
	var chunks [][]string

	chunkSize := len(slice) / (nChunks)
	remainder := len(slice) % nChunks

	end := 0
	for i := 0; i < nChunks; i++ {
		if i == (nChunks-1) && remainder != 0 {
			chunkSize = (len(slice) / (nChunks)) + remainder
		}

		end += chunkSize
		start := end - chunkSize

		chunks = append(chunks, slice[start:end])
	}

	return chunks
}

/*func genEmail(path string) Email {

	content, _ := os.ReadFile(path) //Ruta email

	r := strings.NewReader(string(content))
	m, err := mail.ReadMessage(r)
	if err != nil {
		fmt.Printf("PATH: %v, ERROR: %v \n", path, err)

		var email Email

		return email
	}

	header := m.Header
	body, _ := io.ReadAll(m.Body)

	email := Email{
		Date:    header.Get("Date"),
		From:    header.Get("From"),
		To:      header.Get("To"),
		Subject: header.Get("Subject"),
		Cc:      header.Get("Cc"),
		Body:    string(body),
	}

	return email
}*/

func processChunk(chunk []string, i int) {
	time.Sleep(5 * time.Second)
	f, _ := os.OpenFile("temp"+strconv.FormatInt(int64(i+1), 10)+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Nombre archivo ndjson
	defer f.Close()

	for _, path := range chunk {
		f.WriteString(path + "\n")
	}

	wg.Done()
}
