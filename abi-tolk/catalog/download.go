package catalog

import (
	"context"
	"fmt"
	"net/http"
)

const (
	releaseBase = "https://github.com/ton-blockchain/abis/releases"
	assetName   = "abi-catalog.json"
)

// ReleaseURL builds the URL of a release, tagged as "v0.1.1".
//
// Pin one. The names a catalog gives its contracts end up in whatever a caller
// builds on them, so moving to a newer catalog should be a decision rather than
// something that happens on the next restart.
func ReleaseURL(tag string) string {
	return fmt.Sprintf("%s/download/%s/%s", releaseBase, tag, assetName)
}

// LatestReleaseURL resolves to whichever release is newest at the time of the
// request. Prefer ReleaseURL.
const LatestReleaseURL = releaseBase + "/latest/download/" + assetName

func Download(ctx context.Context, url string) (*Catalog, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	c, err := Load(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", url, err)
	}
	return c, nil
}
