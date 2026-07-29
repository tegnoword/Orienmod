// internal/adapters/output/storage/memory.go
package storage

import (
	"context"
	"os"
	"sync"

	"golang.org/x/oauth2"
)

type MemoryTokenStore struct {
	mu     sync.RWMutex
	tokens map[string]*oauth2.Token
}

func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{
		tokens: make(map[string]*oauth2.Token),
	}
}

func (s *MemoryTokenStore) SaveToken(ctx context.Context, email string, token *oauth2.Token) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[email] = token
	return nil
}

func (s *MemoryTokenStore) GetToken(ctx context.Context, email string) (*oauth2.Token, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	token, exists := s.tokens[email]
	if !exists {
		return nil, os.ErrNotExist
	}
	return token, nil
}

func (s *MemoryTokenStore) DeleteToken(ctx context.Context, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, email)
	return nil
}
