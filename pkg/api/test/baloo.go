package test

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"gopkg.in/h2non/gentleman.v1/context"
	"gopkg.in/h2non/gentleman.v1/plugin"
)

// BalooResponse is a plugin for HTTP response logging
var BalooResponse = plugin.NewResponsePlugin(context.HandlerFunc(
	func(ctx *context.Context, h context.Handler) {
		body, err := ioutil.ReadAll(ctx.Response.Body)
		if err != nil {
			h.Error(ctx, err)
			return
		}

		if err := ctx.Response.Body.Close(); err != nil {
			h.Error(ctx, err)
			return
		}

		ctx.Response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		fmt.Println(string(body))

		h.Next(ctx)
	},
))
