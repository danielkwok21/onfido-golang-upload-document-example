# How to run

1. Inject auth token into env var
```bash
# first inject authorization token
export ONFIDO_AUTH_TOKEN=xxx
```

2. Update fields in code according to use case
```golang
// Refer to https://documentation.onfido.com/#upload-document-request-body
formData := map[string]string{
    "type":         "national_identity_card",
    "applicant_id": "xxxx",
}
```

3. Run program
```bash
# run program
go run main.go

# expected output
Status Code: 201
Response: xxxx
```

# Hard lessons learnt
To attach a "file" field to formData, additional handling is required. This learning is captured in the `prepareFileField()` function.
```golang
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
```

If this wasn't handled this way, Onfido endpoint would respond with this.
```bash
Status Code: 422
Response: {"error":{"message":"There was a validation error on this request","type":"validation_error","fields":{"file":["the content_type of the file has been spoofed"]}}}
```

All credits go to https://github.com/uw-labs/go-onfido