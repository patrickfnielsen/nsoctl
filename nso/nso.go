/*
Copyright © 2022 Patrick Falk Nielsen <git@patricknielsen.dk>
*/
package nso

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type NSO struct {
	Server   string
	Username string
	Password string
	Timeout  time.Duration
}

type NSOResponse struct {
	StatusCode int
	Data       string
}

func (nso NSO) getBasicAuthString() string {
	auth := nso.Username + ":" + nso.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (nso NSO) request(method string, url string, body io.Reader) (resp *http.Response, err error) {
	url = fmt.Sprintf("https://%s/%s", nso.Server, url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatalf("HTTP Request: %s", err)
		return nil, err
	}

	req.Header.Add("Authorization", nso.getBasicAuthString())
	req.Header.Add("Accept", "application/yang-data+json")

	var httpClient = &http.Client{
		Timeout: time.Second * nso.Timeout,
	}

	resp, err = httpClient.Do(req)
	if err != nil {
		log.Fatalf("HTTP Client: %s", err)
		return nil, err
	}

	return resp, err
}

func (nso NSO) getBody(resp *http.Response) (body string, err error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s", err)
		return "", err
	}

	return string(bodyBytes), err
}

func (nso NSO) Post(url string, body io.Reader) (response *NSOResponse, err error) {
	resp, err := nso.request("POST", url, body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := nso.getBody(resp)
	return &NSOResponse{resp.StatusCode, data}, err
}

func (nso NSO) Get(url string) (response *NSOResponse, err error) {
	resp, err := nso.request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := nso.getBody(resp)
	return &NSOResponse{resp.StatusCode, data}, err
}

func (nso NSO) Redeploy(service string, id string) (err error) {
	url := fmt.Sprintf("/restconf/data/tailf-ncs:services/%s:%s=%s/re-deploy", service, service, id)
	resp, err := nso.Post(url, nil)
	if err != nil {
		return fmt.Errorf("Failed to re-deploy service, error: %s", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("Failed to re-deploy service, unexpected status code: %d", resp.StatusCode)
	}

	return err
}

func (nso NSO) Undeploy(service string, id string) (err error) {
	url := fmt.Sprintf("/restconf/data/tailf-ncs:services/%s:%s=%s/un-deploy", service, service, id)
	resp, err := nso.Post(url, nil)
	if err != nil {
		return fmt.Errorf("Failed to un-deploy service, error: %s", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("Failed to un-deploy service, unexpected status code: %d", resp.StatusCode)
	}

	return err
}

func (nso NSO) SyncFromDevice(deviceName string) (err error) {
	url := fmt.Sprintf("/restconf/data/tailf-ncs:devices/device=%s/sync-from", deviceName)
	resp, err := nso.request("POST", url, nil)
	defer resp.Body.Close()

	if err != nil {
		return fmt.Errorf("Failed to sync-from device, error: %s", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to sync-from device, unexpected status code: %d", resp.StatusCode)
	}

	return err
}
