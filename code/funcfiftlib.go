package code

// #cgo darwin LDFLAGS: -L ../lib/darwin/ -Wl,-rpath,../lib/darwin/ -l funcfiftlib
// #cgo linux LDFLAGS: -L ../lib/linux/ -Wl,-rpath,../lib/linux/ -l funcfiftlib
// #include "../lib/funcfiftlib.h"
// #include <stdlib.h>
// #include <stdbool.h>
import "C"
import (
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed fift
var fiftLibs embed.FS

var fileMap map[string]string
var fMutex sync.Mutex

// todo: C.free
//
//export fileReader
func fileReader(_kind *C.char, _data *C.char, content **C.char, errorString **C.char) {
	kind, data := C.GoString(_kind), C.GoString(_data)
	switch kind {
	case "realpath", "source":
		s, prs := fileMap[data]
		if !prs {
			s := C.CString(fmt.Sprintf("file %s not found", data))
			errorString = &s
			return
		}
		*content = C.CString(s)
		return
	default:
		fmt.Println(kind, data)
		return
	}
}

type CompileResult struct {
	Success  bool `json:"success"`
	Boc      []byte
	Fift     string
	Warnings string
}

func Compile(files []string) (CompileResult, error) {
	fMutex.Lock()
	defer fMutex.Unlock()
	fileMap = make(map[string]string)
	var sources struct {
		OptLevel int      `json:"optLevel"`
		Sources  []string `json:"sources"`
	}
	sources.OptLevel = 2
	for i, f := range files {
		filename := fmt.Sprintf("%d.fc", i)
		fileMap[filename] = f
		sources.Sources = append(sources.Sources, filename)
	}
	b, err := json.Marshal(sources)
	if err != nil {
		return CompileResult{}, err
	}
	a := C.compile(C.CString(string(b)))
	var cResult struct {
		Status   string `json:"status"`
		CodeBoc  string `json:"codeBoc"`
		FiftCode string `json:"fiftCode"`
		Warnings string `json:"warnings"`
		Message  string `json:"message"`
	}
	err = json.Unmarshal([]byte(C.GoString(a)), &cResult)
	if err != nil {
		return CompileResult{}, err
	}
	if cResult.Status == "error" {
		return CompileResult{}, fmt.Errorf("compilation error: %v", cResult.Message)
	}
	boc, err := base64.StdEncoding.DecodeString(cResult.CodeBoc)
	if err != nil {
		return CompileResult{}, err
	}
	return CompileResult{
		Success:  cResult.Status == "ok",
		Boc:      boc,
		Fift:     cResult.FiftCode,
		Warnings: cResult.Warnings,
	}, err
}
