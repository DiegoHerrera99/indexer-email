package helpers

import (
	"email-indexer/globals"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Email struct {
	Date    string `json:"date"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Cc      string `json:"cc"`
	Body    string `json:"body"`
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("---------- TIME OF EXEC(%v): %s \n", name, elapsed)
}

func ChunkSlice(slice []string, nChunks int) [][]string {
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

func UploadChunk(chunk []string, i int, wg *sync.WaitGroup) {

	tempFile := globals.TEMPDIR + "/bulk" + strconv.FormatInt(int64(i+1), 10) + ".json"

	f, _ := os.OpenFile(tempFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	for idx, path := range chunk {
		email, err := GenEmail(path)

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
	BulkFile(tempFile)

	wg.Done()
}

func GenEmail(path string) (Email, error) {

	content, _ := os.ReadFile(path) //Ruta email

	r := strings.NewReader(string(content))
	m, err := mail.ReadMessage(r)
	if err != nil {
		formattedDate := time.Now().Format(time.RFC3339)
		msg := fmt.Sprintf("PATH: %v, ERROR: %v", path, err)
		email := Email{
			Date: formattedDate,
		}

		return email, errors.New(msg)
	}

	header := m.Header
	body, _ := io.ReadAll(m.Body)

	date := header.Get("Date")
	parseDate, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700 (MST)", date)
	if err != nil {
		fmt.Println(err)
	}

	formattedDate := parseDate.Format(time.RFC3339)

	email := Email{
		Date:    formattedDate,
		From:    header.Get("From"),
		To:      header.Get("To"),
		Subject: header.Get("Subject"),
		Cc:      header.Get("Cc"),
		Body:    string(body),
	}

	return email, nil
}

func BulkFile(path string) {
	bulkFile, _ := os.Open(path)
	defer bulkFile.Close()
	defer os.Remove(path)

	req, err := http.NewRequest("POST", globals.ZINC_ENDPOINT, bulkFile)
	if err != nil {
		fmt.Println(err)
	}

	req.SetBasicAuth(globals.ZINC_USER, globals.ZINC_PWD)
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("UPLOADING CHUNK: %v \n", string(body))
}

func CreateIndex(url string, method string, payload string) error {
	payloadReader := strings.NewReader(payload)

	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(globals.ZINC_USER, globals.ZINC_PWD)
	//req.Header.Add("Authorization", "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("CREATING INDEX: %v \n", string(body))

	return nil
}
