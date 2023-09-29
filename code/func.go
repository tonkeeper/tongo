package code

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FunCCompiler is not real compiler - it just send code to tonapi.io
// be careful - it behavior can be changed in the future
type FunCCompiler struct {
	url string
}

func NewFunCCompiler() *FunCCompiler {
	return &FunCCompiler{url: "https://code.tonapi.io/func/compile"}
}

// Compile returns FIFT code, BoC with compiled code and error
func (c *FunCCompiler) Compile(files [][2]string) (string, []byte, error) {
	type file struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
	}
	var filess []file
	for _, f := range files {
		filess = append(filess, file{
			Filename: f[0],
			Content:  f[1],
		})
	}
	b, err := json.Marshal(map[string]any{"files": filess})
	fmt.Println(string(b))
	if err != nil {
		return "", nil, err
	}
	resp, err := http.Post(c.url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("invalid compiler status code %v", resp.StatusCode)
	}
	var respBody struct {
		Success          bool
		Hex, Fift, Error string
	}
	b, _ = io.ReadAll(resp.Body)
	fmt.Println(string(b))
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", nil, err
	}
	if !respBody.Success {
		return "", nil, fmt.Errorf(respBody.Error)
	}
	boc, err := hex.DecodeString(respBody.Hex)
	return respBody.Fift, boc, err
}
