package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
	"github.com/pardot/oidc/discovery"
)

const wellKnownProtectedResourcePath = "/.well-known/oauth-protected-resource"

type OauthConfig struct {
	resourceUrl string
	issuerUrl   string
	jwksUrl     string
}

func NewOauthConfig(ctx context.Context, issuer, resourceUrl string) (*OauthConfig, error) {
	clnt, err := discovery.NewClient(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("Can not retrieve issuer metadata: %v", err)
	}
	metadata := clnt.Metadata()
	return &OauthConfig{
		resourceUrl: resourceUrl,
		issuerUrl:   issuer,
		jwksUrl:     metadata.JWKSURI,
	}, nil
}

func (oc *OauthConfig) AddRoutes(mux *http.ServeMux) {
	metadata := &oauthex.ProtectedResourceMetadata{
		Resource:              oc.resourceUrl,
		ScopesSupported:       []string{"mcp:tools"},
		AuthorizationServers:  []string{oc.issuerUrl},
		ResourceName:          "Prompts on the Job",
		ResourceDocumentation: "https://www.promptsonthejob.nl/",
	}
	mux.Handle(wellKnownProtectedResourcePath, auth.ProtectedResourceMetadataHandler(metadata))
}
