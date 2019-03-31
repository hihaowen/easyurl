package engine

import (
	"net/http"
	"sync"
)

type Engine struct {
	RouterGroup

	pool sync.Pool

	trees methodTrees

	UseRawPath bool

	UnescapePathValues bool

	RedirectFixedPath bool

	RedirectTrailingSlash bool

	allNoRoute HandlersChain
}

var _ IRouter = &Engine{}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		UseRawPath:         false,
		UnescapePathValues: false,
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}

func (engine *Engine) Run(address string) (err error) {
	defer func() { debugPrintError(err) }()

	debugPrint("Listening and serving HTTP on %s\n", address)

	err = http.ListenAndServe(address, engine)
	return
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = r
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	path := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		path = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		handlers, params, tsr := root.getValue(path, c.Params, unescape)
		if handlers != nil {
			c.handlers = handlers
			c.Params = params
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != "CONNECT" && path != "/" {
			if tsr && engine.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
				return
			}
		}
		break
	}

	c.handlers = engine.allNoRoute
	serveError(c, http.StatusNotFound, []byte("404 page not found"))
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	path := req.URL.Path
	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != "GET" {
		code = http.StatusTemporaryRedirect
	}

	req.URL.Path = path + "/"
	if length := len(path); length > 1 && path[length-1] == '/' {
		req.URL.Path = path[:length-1]
	}
	http.Redirect(c.Writer, req, req.URL.String(), code)
	c.writermem.WriteHeaderNow()
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	path := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(path), trailingSlash); ok {
		code := http.StatusMovedPermanently // Permanent redirect, request with GET method
		if req.Method != "GET" {
			code = http.StatusTemporaryRedirect
		}
		req.URL.Path = string(fixedPath)
		http.Redirect(c.Writer, req, req.URL.String(), code)
		c.writermem.WriteHeaderNow()
		return true
	}
	return false
}

var mimePlain = []string{"text/plain"}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}

	if c.writermem.Status() == code {
		debugPrint("context Writer: %+v", c.Writer)

		c.writermem.Header()["Content-Type"] = mimePlain
		c.Writer.Write(defaultMessage)
		return
	}
	c.writermem.WriteHeaderNow()

	return
}
