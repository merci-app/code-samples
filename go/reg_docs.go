package main

import (
	"fmt"
	"github.com/merci-app/code-samples/go/authorization"
	"github.com/merci-app/code-samples/go/request"
	"net/http"
)

func main() {
	accessToken := authorization.NewAuthorization("<USERNAME>", "<PASSWORD>")

	var response RegDocsResponse
	req := request.NewRequest(*accessToken)
	resp, body, err := req.Get("https://regdocs.hml.caradhras.io/v1/registration?types=PRIVACY_POLICY&types=TERMS_OF_USE", &response)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		panic(string(body))
	}

	fmt.Println(response.Message)

	var tokens []string
	for _, docs := range response.Result.RegulatoryDocuments {
		tokens = append(tokens, docs.Token)
	}

	agreeRequest := AgreeRegDocsRequest{
		Tokens:      tokens,
		Fingerprint: "fingerprint",
	}
	var agreeResponse AgreeRegDocsResponse
	resp, body, err = req.Post("https://regdocs.hml.caradhras.io/v1/agreement", agreeRequest, &agreeResponse)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		panic(string(body))
	}

	fmt.Println(agreeResponse.Message)
}

type RegDocsResponse struct {
	Message string                `json:"message"`
	Result  RegDocsResultResponse `json:"result"`
}

type RegDocsResultResponse struct {
	RegulatoryDocuments []RegDocsDocumentsResponse `json:"regulatoryDocuments"`
}

type RegDocsDocumentsResponse struct {
	Type      string `json:"type"`
	Token     string `json:"token"`
	RegDocObj string `json:"regDocObj"`
}

type AgreeRegDocsRequest struct {
	Tokens         []string `json:"tokens"`
	IdRegistration string   `json:"idRegistration,omitempty"`
	Fingerprint    string   `json:"fingerprint"`
}

type AgreeRegDocsResponse struct {
	Message string `json:"message"`
}
