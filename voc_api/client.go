package voc_api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type VocApi struct {
	baseUrl string
	token string
}

func New(baseUrl string, token string) *VocApi {
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &VocApi{
		baseUrl: baseUrl,
		token: token,
	}
}

func (api *VocApi) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", api.token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed: %s: %s", resp.Status, string(data))
	}
	return resp, nil
}

func parseBody(resp *http.Response, responseObject any) error {
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return err
	}

	return nil
}

func (api *VocApi) GetTalk(conference string, event string) (*Talk, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/events/%s", api.baseUrl, conference, event), nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &Talk{}
	err = parseBody(resp, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *VocApi) GetConference(conference string) (*Conference, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", api.baseUrl, conference), nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &Conference{}
	err = parseBody(resp, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *VocApi) UploadFile(
	conference string,
	event string,
	fileName string,
	fileMimeType string,
	fileContent []byte,
	meta map[string]any,
) error {
	// Prepare multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// File part
	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	_, err = filePart.Write(fileContent)
	if err != nil {
		return err
	}

	// Meta part
	metaJson, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	metaPart, err := writer.CreateFormField("meta")
	if err != nil {
		return err
	}
	_, err = metaPart.Write(metaJson)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	// Prepare request
	url := fmt.Sprintf("%s/%s/events/%s/file", api.baseUrl, conference, event)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := api.doRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (api *VocApi) UploadVtt(
	conference string,
	event string,
	vttContent []byte,
	language string,
) error {
	meta := map[string]any{
		"recording": map[string]any{
			"language":  language,
			"mime_type": "text/vtt",
		},
	}
	return api.UploadFile(
		event,
		conference,
		"dummy.vtt",
		"text/vtt",
		vttContent,
		meta,
	)
}
