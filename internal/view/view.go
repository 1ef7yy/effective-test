package view

import (
	"emobile/pkg/logger"

	"emobile/internal/domain"
)

type View interface {
}

type view struct {
	domain domain.Domain
}

func NewView(log logger.Logger) View {
	return &view{
		domain: domain.NewDomain(log),
	}
}
