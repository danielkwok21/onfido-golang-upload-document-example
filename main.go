package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func main() {
	const (
		apiEndpoint = "https://api.eu.onfido.com/v3.4/documents"
		filePath    = "./sample_driving_licence.png"
	)

	authToken := os.Getenv("ONFIDO_AUTH_TOKEN")
	if authToken == "" {
		fmt.Println("API key not set.")
	} else {
		fmt.Println("API key:", authToken)
	}

	// Refer to https://documentation.onfido.com/#upload-document-request-body
	formData := map[string]string{
		"type":         "national_identity_card",
		"applicant_id": "xxxx",
	}

	// Create a new multipart buffer to store the form data.
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the form data to the multipart buffer.
	for key, value := range formData {
		multipartWriter.WriteField(key, value)
	}

	prepareFileField(multipartWriter, filePath)

	// Close the multipart writer
	_ = multipartWriter.Close()

	// Create a new POST request to the API endpoint.
	request, err := http.NewRequest("POST", apiEndpoint, &requestBody)
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return
	}
	// Set the Authorization header with the token.
	request.Header.Set("Authorization", "Token token="+authToken)

	// set the Content-Type header to "multipart/form-data; boundary=xxxxx".
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// make the request using the default HTTP client.
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return
	}
	defer response.Body.Close()

	// preview response
	fmt.Println("Status Code:", response.StatusCode)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response:", string(body))
}

// prepareFileField helps us attach a file to the formData
func prepareFileField(multipartWriter *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// inject mime headers
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filepath.Base(filePath)))
	h.Set("Content-Type", "image/png")
	part, _ := multipartWriter.CreatePart(h)

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file data to multipart writer:", err)
		return err
	}

	return nil
}
