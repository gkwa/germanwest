package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const imgbbURL = "https://api.imgbb.com/1/upload"

type ImageResponse struct {
	Data struct {
		DeleteURL  string `json:"delete_url"`
		DisplayURL string `json:"display_url"`
		Expiration int    `json:"expiration"`
		Height     int    `json:"height"`
		ID         string `json:"id"`
		Image      Image  `json:"image"`
		Medium     Image  `json:"medium"`
		Size       int    `json:"size"`
		Thumb      Image  `json:"thumb"`
		Time       int    `json:"time"`
		Title      string `json:"title"`
		URL        string `json:"url"`
		URLViewer  string `json:"url_viewer"`
		Width      int    `json:"width"`
	} `json:"data"`
	Status  int  `json:"status"`
	Success bool `json:"success"`
}

type Image struct {
	Extension string `json:"extension"`
	Filename  string `json:"filename"`
	Mime      string `json:"mime"`
	Name      string `json:"name"`
	URL       string `json:"url"`
}

func uploadImage(imagePath, apiKey string) (*ImageResponse, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("error copying file: %v", err)
	}

	err = writer.WriteField("key", apiKey)
	if err != nil {
		return nil, fmt.Errorf("error writing field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	req, err := http.NewRequest("POST", imgbbURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var result ImageResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func Run(apiKey string) {
	imagePath := "testdata/image.png"
	response, err := uploadImage(imagePath, apiKey)
	if err != nil {
		fmt.Printf("error uploading image: %v\n", err)
		return
	}

	fmt.Printf("Image uploaded successfully: %s\n", response.Data.URL)
}
