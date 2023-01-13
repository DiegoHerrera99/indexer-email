package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
)

type Email struct {
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
}

func main() {

	path := "/Users/diegoherrera/Downloads/enron_mail_20110402/" //DB PATH "./db-mail/" //
	index := "enron"

	//Generar DB con direcciones de emails
	f, _ := os.OpenFile("temp.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Nombre archivo ndjson
	defer f.Close()
	getEmailDir(path, f)

	//Leer DB de emails, mapear emails a JSON, crear documento para BulkV2
	jsonFile, _ := os.OpenFile("temp.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) //Nombre archivo ndjson
	defer jsonFile.Close()

	jsonFile.WriteString(`{"index": "` + index + `" , "records": [ `)

	file, _ := os.Open("temp.txt")
	defer file.Close()
	defer os.Remove("temp.txt")

	reader := bufio.NewReader(file)
	i := true
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			jsonFile.WriteString("]}")
			break
		}

		email := genEmail(string(line))
		emailJSON, _ := json.Marshal(email)

		if i {
			jsonFile.WriteString(string(emailJSON) + "\n")
			i = false
		} else {
			jsonFile.WriteString("," + string(emailJSON) + "\n")
		}
	}

	//Realizar Bulk
	bulkFile, _ := os.Open("temp.json")
	defer f.Close()
	defer os.Remove("temp.json")

	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", bulkFile)
	if err != nil {
		fmt.Println(err)
	}

	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getEmailDir(path string, f *os.File) {

	//ABRIR DIRECTORIO
	openDir, _ := os.Open(path)

	//LEER DIRECTORIO PARA OBTENER TODOS LOS ARCHIVOS HIJOS
	dirFiles, _ := openDir.ReadDir(0)

	//SI ES ARCHIVO LEER Y AGREGAR AL ARCHIVO SI NO USAR RECURSIVIDAD
	for _, file := range dirFiles {
		if file.IsDir() || file.Name() == ".DS_Store" {
			getEmailDir(path+file.Name()+"/", f)
		} else {
			f.WriteString(path + file.Name() + "\n")
		}
	}

}

func genEmail(path string) Email {

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
}
