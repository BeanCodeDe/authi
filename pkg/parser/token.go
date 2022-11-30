// Package to verify tokens
package parser

import (
	"github.com/BeanCodeDe/authi/pkg/adapter"
)

type (
	// Interface fo pars and validate tokens
	Parser interface {
		ParseToken(authorizationString string) (*adapter.Claims, error)
	}
)
