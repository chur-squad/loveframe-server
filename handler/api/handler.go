package api_handler

import (
	"github.com/chur-squad/loveframe-server/handler"
	"github.com/chur-squad/loveframe-server/internal"
)

type Handler struct {
	parent *handler.Handler
}

func NewHandler(h *handler.Handler) (*Handler, error) {
	if h == nil {
		return nil, internal.ErrInvalidParams
	}

	return &Handler{parent: h}, nil
}
