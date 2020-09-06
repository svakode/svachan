package handler

import (
	"github.com/svakode/svachan/framework"
)

type handler struct {
	ctx *framework.Context
}

// NewHandler is the method for construct a handler
func NewHandler(ctx *framework.Context) Handler {
	return &handler{
		ctx: ctx,
	}
}
