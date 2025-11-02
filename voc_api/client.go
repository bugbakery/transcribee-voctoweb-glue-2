package voc_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type VocApi struct {
	baseUrl string
}

func New(baseUrl string) *VocApi {
	return &VocApi{
		baseUrl: baseUrl,
	}
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
	resp, err := http.Get(fmt.Sprintf("%s/%s/events/%s", api.baseUrl, conference, event))
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
	resp, err := http.Get(fmt.Sprintf("%s/%s", api.baseUrl, conference))
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
