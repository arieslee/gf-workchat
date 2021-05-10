package main

import (
	"gf_workchat/config"
	"gf_workchat/core/token"
	"gf_workchat/user"
)

type WorkChat struct {
}

// NewWorkChat init
func NewWorkChat() *WorkChat {
	return &WorkChat{}
}

func (w *WorkChat) GetUser(cfg *config.Config) *user.User {
	return user.New(cfg)
}

func (w *WorkChat) GetToken(cfg *config.Config) *token.Token {
	return token.New(cfg)
}
