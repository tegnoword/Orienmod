package ports

import (
	"context"

	"golang.org/x/oauth2"
)

type TokenRepository interface {
	SaveToken(ctx context.Context, email string, token *oauth2.Token) error
	GetToken(ctx context.Context, email string) (*oauth2.Token, error)
}
