package connect

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
)

type Request struct {
	Protocol string `json:"protocol"`
	V1       struct {
		Session        string `json:"session"`
		SessionPayload string `json:"session_payload"`
		ImageUrl       string `json:"image_url"`
		CallbackUrl    string `json:"callback_url"`
		Items          []struct {
			Type     string `json:"type"`
			Required bool   `json:"required"`
		} `json:"items"`
	} `json:"v1"`
}

type Responce struct {
	Version        string `json:"version"`
	Nonce          string `json:"nonce"`
	Clientid       string `json:"clientid"`
	Authenticator  string `json:"authenticator"`
	SessionPayload string `json:"session_payload"`
}

type AuthPayload struct {
	Items []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"items"`
}

func createTonOwnershipSignature(walletVersion, address, clientId string, secretKey *[64]byte) string {
	message := fmt.Sprintf(`tonlogin/ownership/%v/%v/%v`, walletVersion, address, clientId)

	s := sign.Sign(nil, []byte(message), secretKey)

	return base64.StdEncoding.EncodeToString(s)
}
