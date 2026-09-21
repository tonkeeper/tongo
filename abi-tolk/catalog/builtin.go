package catalog

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
)

//go:embed builtin/abi-catalog.json.gz
var builtinBundle []byte

//go:embed builtin/version.txt
var builtinTag string

func BuiltinVersion() string { return strings.TrimSpace(builtinTag) }

var builtin = sync.OnceValues(func() (*Catalog, error) {
	r, err := gzip.NewReader(bytes.NewReader(builtinBundle))
	if err != nil {
		return nil, fmt.Errorf("the built-in catalog is not gzip: %w", err)
	}
	c, err := Load(r)
	if err != nil {
		return nil, fmt.Errorf("the built-in catalog does not read: %w", err)
	}
	return c, nil
})

func GetBuiltin() *Catalog {
	c, err := builtin()
	if err != nil {
		panic(err)
	}
	return c
}

func GetLatest(ctx context.Context) (*Catalog, error) {
	tag, err := LatestVersion(ctx)
	if err != nil {
		return nil, err
	}
	if tag == BuiltinVersion() {
		return builtin()
	}
	return Download(ctx, ReleaseURL(tag))
}

func LatestVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, releaseBase+"/latest", nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	location := resp.Header.Get("Location")
	tag := path.Base(location)
	if location == "" || !strings.HasPrefix(tag, "v") {
		return "", fmt.Errorf("%s/latest does not point at a release, it answered %s", releaseBase, resp.Status)
	}
	return tag, nil
}
