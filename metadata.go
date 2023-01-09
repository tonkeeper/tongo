package tongo

type JettonMetadata struct {
	Uri         string `json:"uri,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
	ImageData   []byte `json:"image_data,omitempty"`
	Symbol      string `json:"symbol,omitempty"`
	Decimals    string `json:"decimals,omitempty"`
}
