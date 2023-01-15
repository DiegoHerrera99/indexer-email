package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func main() {

	start := time.Now()
	defer timeTrack(start, "genPathChunks")

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

	chunks := chunkSlice(fileList, nChunks)
	fmt.Printf("NO. OF CHUNKS: %v    ----------\nNO. OF FILES: %v ----------\n", len(chunks), len(fileList))

}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s", name, elapsed)
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
