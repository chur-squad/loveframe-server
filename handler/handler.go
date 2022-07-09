package handler

import (
	photos "github.com/chur-squad/loveframe-server/photos"
)

type (
	// Handler is handler struct for using in Echo.
	Handler struct {
		Cfg			*Config
		Photo		photos.Manager		
	}
)

// NewHandler is to create a handler object.
func NewHandler(opts ...Option) (*Handler, error) {
	h := &Handler{}
	//handelr is struct 
	
	return h, nil
}
