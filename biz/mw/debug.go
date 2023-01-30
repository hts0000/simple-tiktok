package mw

import (
	"context"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
)

func Debug(c context.Context, ctx *app.RequestContext) {
	log.Println("recv request", string(ctx.Method()), string(ctx.URI().FullURI()))
	ctx.Next(c)
}
