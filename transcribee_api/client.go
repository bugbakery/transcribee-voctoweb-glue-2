package transcribee_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type TranscribeeApi struct {
	baseURL string
	token   string
}

func New(baseURL string, token string) *TranscribeeApi {
	return &TranscribeeApi{
		baseURL: baseURL,
		token:   token,
	}
}

func (api *TranscribeeApi) doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", api.token))
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

func parseBody(resp *http.Response, out any) error {
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w; raw: %s", err, string(data))
	}
	return nil
}

type DocumentBodyWithFile struct {
	Language         string
	Model            string
	Name             string
	NumberOfSpeakers *int
	FileName         string
	File             io.Reader
}

func (api *TranscribeeApi) CreateDocument(body *DocumentBodyWithFile) (*Document, error) {
	if body == nil {
		return nil, fmt.Errorf("body is nil")
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	_ = writer.WriteField("language", body.Language)
	_ = writer.WriteField("model", body.Model)
	_ = writer.WriteField("name", body.Name)
	if body.NumberOfSpeakers != nil {
		_ = writer.WriteField("number_of_speakers", fmt.Sprintf("%d", *body.NumberOfSpeakers))
	}

	if body.File != nil {
		fw, err := writer.CreateFormFile("file", body.FileName)
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		if _, err := io.Copy(fw, body.File); err != nil {
			_ = writer.Close()
			return nil, err
		}
	}

	_ = writer.Close()

	url := fmt.Sprintf("%s/api/v1/documents/", api.baseURL)
	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}

	var doc Document
	if err := parseBody(resp, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (api *TranscribeeApi) GetTranscribeeDocuments() ([]Document, error) {
	url := fmt.Sprintf("%s/api/v1/documents/", api.baseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}
	var docs []Document
	if err := parseBody(resp, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (api *TranscribeeApi) GetTasksForDocument(docID string) ([]TaskResponse, error) {
	url := fmt.Sprintf("%s/api/v1/documents/%s/tasks/", api.baseURL, docID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}
	var tasks []TaskResponse
	if err := parseBody(resp, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (api *TranscribeeApi) CreateShareToken(docID string, data *CreateShareToken) (*DocumentShareTokenBase, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/api/v1/documents/%s/share_tokens/", api.baseURL, docID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := api.doRequest(req)
	if err != nil {
		return nil, err
	}
	var out DocumentShareTokenBase
	if err := parseBody(resp, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (api *TranscribeeApi) Export(docID string, format string, includeSpeakerNames bool, includeWordTiming bool, maxLineLength int) (string, error) {
	params := url.Values{}
	if format != "" {
		params.Set("format", format)
	}
	if includeSpeakerNames {
		params.Set("include_speaker_names", "true")
	} else {
		params.Set("include_speaker_names", "false")
	}
	if includeWordTiming {
		params.Set("include_word_timing", "true")
	} else {
		params.Set("include_word_timing", "false")
	}
	if maxLineLength != 0 {
		params.Set("max_line_length", fmt.Sprintf("%d", maxLineLength))
	}

	url := fmt.Sprintf("%s/api/v1/documents/%s/export/?%s", api.baseURL, docID, params.Encode())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := api.doRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (api *TranscribeeApi) CreateShareUrl(docID string) (string, error) {
	shareToken, err := api.CreateShareToken(docID, &CreateShareToken{
		CanWrite: true,
		Name:     "voctoglue",
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/document/%s?share_token=%s", api.baseURL, docID, url.QueryEscape(shareToken.Token)), nil
}
