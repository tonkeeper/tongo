package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/startfellows/tongo/connect"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
)

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{Jar: jar}
	hash, err := getHash(client)
	if err != nil {
		panic(err)
	}
	req, tempSession, err := getConnectRequest(client, hash)
	if err != nil {
		panic(err)
	}
	fmt.Println(req, tempSession)
}
func getConnectRequest(client *http.Client, hash string) (connect.Request, string, error) {
	resp, err := client.Post("https://fragment.com/api?hash="+hash, "application/x-www-form-urlencoded", strings.NewReader("method=getTonAuthLink"))
	if err != nil {
		return connect.Request{}, "", err
	}

	defer resp.Body.Close()
	var respBody struct {
		Ok          bool   `json:"ok"`
		Link        string `json:"link"`
		CheckMethod string `json:"check_method"`
		CheckParams struct {
			TempSession string `json:"temp_session"`
		} `json:"check_params"`
	}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return connect.Request{}, "", err
	}
	if !respBody.Ok || !strings.HasPrefix(respBody.Link, "https://app.tonkeeper.com/ton-login/") {
		return connect.Request{}, "", fmt.Errorf("invalid fragment auth link")
	}
	url := "https://" + strings.TrimPrefix(respBody.Link, "https://app.tonkeeper.com/ton-login/")
	resp, err = http.Get(url)
	if err != nil {
		return connect.Request{}, "", err
	}
	defer resp.Body.Close()
	var req connect.Request
	err = json.NewDecoder(resp.Body).Decode(&req)
	return req, respBody.CheckParams.TempSession, err
}

func getHash(client *http.Client) (string, error) {
	resp, err := client.Get("https://fragment.com/")

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	found := regexp.MustCompile(`hash=([a-fA-F0-9]{15,32})`).FindSubmatch(b)
	if len(found) != 2 {
		return "", errors.New("hash not found")
	}
	return string(found[1]), nil
}
