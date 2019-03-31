package engine

import (
	"easyurl/infra/http/response/render"
	"net"
	"net/http"
	"strings"
)

type Context struct {
	Request *http.Request
	engine  *Engine

	index int8

	handlers HandlersChain

	writermem responseWriter

	Params Params

	Writer ResponseWriter
}

func (c *Context) ClientIP() string {

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func (c *Context) Next() {
	c.index++
	for s := int8(len(c.handlers)); c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
}

func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

func (c *Context) Status(code int) {
	c.writermem.WriteHeader(code)
}
