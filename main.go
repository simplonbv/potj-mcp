package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
	"github.com/pardot/oidc/discovery"
	"github.com/urfave/cli/v3"
	"promptsonthejob.nl/mcp/internal/auth"
	"promptsonthejob.nl/mcp/internal/mcpserver"
	"promptsonthejob.nl/mcp/internal/model"
)

const wellKnownProtectedResourcePath = "/.well-known/oauth-protected-resource"

func isPostgresUrl(ctx context.Context, c *cli.Command, s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if u.Scheme != "https" {
		return errors.New("URL must use postgresql scheme")
	}
	return nil
}

func isHttpUrl(ctx context.Context, c *cli.Command, s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if !u.IsAbs() {
		return errors.New("URL must be absolute")
	}
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cmd := &cli.Command{
		Name:                   "Prompts on the Job MCP server",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dsn",
				Aliases: []string{"d"},
				Value:   "postgresql:///",
				Sources: cli.EnvVars("DSN"),
				Action:  isPostgresUrl,
			},
			&cli.StringFlag{
				Name:    "issuer",
				Value:   "https://login.simplon.biz",
				Sources: cli.EnvVars("ISSUER"),
				Action:  isHttpUrl,
			},
			&cli.StringFlag{
				Name:    "url",
				Usage:   "Public URL for the server",
				Value:   "https://mcp.promptsonthejob.nl/",
				Sources: cli.EnvVars("URL"),
				Action:  isHttpUrl,
			},
			&cli.StringFlag{
				Name:    "transport",
				Aliases: []string{"t"},
				Value:   "stdio",
				Sources: cli.EnvVars("TRANSPORT"),
				Config: cli.StringConfig{
					TrimSpace: true,
				},
				Action: func(ctx context.Context, c *cli.Command, s string) error {
					if s != "sse" && s != "stdio" {
						return errors.New("transport must be sse or stdio")
					}
					return nil
				},
			},
			&cli.Uint32Flag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8000,
				Sources: cli.EnvVars("PORT"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			pool, err := model.Connect(ctx, cmd.String("dsn"))
			if err != nil {
				return fmt.Errorf("can not connect to database: %w", err)
			}
			promptRepository := model.NewPostgresPromptRepository(pool)
			server := mcpserver.NewMCPServer(logger, promptRepository)

			transport := cmd.String("transport")
			if transport == "stdio" {
				return server.Run(ctx, &mcp.StdioTransport{})

			} else {
				issuer := cmd.String("issuer")
				baseUrl := cmd.String("url")
				bu, _ := strings.CutSuffix(baseUrl, "/")
				mcpUrl := fmt.Sprintf("%s/mcp", bu)
				wellKnownProtectedResourceUrl := bu + wellKnownProtectedResourcePath

				clnt, err := discovery.NewClient(ctx, issuer)
				if err != nil {
					return fmt.Errorf("can not retrieve issuer metadata: %w", err)
				}

				validator, err := auth.NewValidator(ctx, issuer, clnt.Metadata().JWKSURI, baseUrl)
				if err != nil {
					return fmt.Errorf("can not initialize authentication: %w", err)
				}

				metadata := &oauthex.ProtectedResourceMetadata{
					Resource:               mcpUrl,
					ScopesSupported:        []string{"mcp:tools"},
					AuthorizationServers:   []string{issuer},
					BearerMethodsSupported: []string{"header"},
					ResourceName:           "Prompts on the Job",
					ResourceDocumentation:  "https://www.promptsonthejob.nl/",
				}

				mux := http.NewServeMux()
				mux.Handle(wellKnownProtectedResourcePath, sdkauth.ProtectedResourceMetadataHandler(metadata))

				mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
					return server
				}, nil)
				mux.Handle("/mcp", validator.Protect(wellKnownProtectedResourceUrl, []string{"read"}, mcpHandler))

				port := cmd.Uint32("port")
				log.Printf("Listening on port %d", port)
				log.Printf("MCP service available at %s", mcpUrl)
				return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
			}
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
