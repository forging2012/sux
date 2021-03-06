package sux

import (
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
)

/*************************************************************
 * Context
 *************************************************************/

const (
	abortIndex int8 = math.MaxInt8 / 2
)

type IContext interface {
	Req() *http.Request
	Res() http.ResponseWriter
	Init(http.ResponseWriter, *http.Request, HandlersChain)

	Next()
	Reset()
	Params() Params
	SetParams(Params)

	HandlerName() string
}

// Context for http server
type Context struct {
	req *http.Request
	res http.ResponseWriter

	index int8
	// current route params, if route has var params
	params Params
	// context data, you can save some custom data.
	values map[string]interface{}
	// all handlers for current request
	handlers HandlersChain
}

func newContext(res http.ResponseWriter, req *http.Request, handlers HandlersChain) *Context {
	return &Context{
		res: res,
		req: req,

		index:  -1,
		values: make(map[string]interface{}),

		handlers: handlers,
	}
}

func (c *Context) Init(res http.ResponseWriter, req *http.Request, handlers HandlersChain) {
	c.res = res
	c.req = req
	c.values = make(map[string]interface{})
	c.handlers = handlers
}

func (c *Context) HandlerName() string {
	return nameOfFunction(c.handlers.Last())
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

// Values get all values
func (c *Context) Values() map[string]interface{} {
	return c.values
}

// Set a value to context by key
func (c *Context) Set(key string, val interface{}) {
	c.values[key] = val
}

// Get a value from context
func (c *Context) Get(key string) interface{} {
	return c.values[key]
}

// Next call next handler
func (c *Context) Next() {
	c.index++
	s := int8(len(c.handlers))

	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// AppendHandlers to the context
func (c *Context) AppendHandlers(handlers ...HandlerFunc) {
	c.handlers = append(c.handlers, handlers...)
}

// Reset context
func (c *Context) Reset() {
	// c.Writer = &c.writermem
	c.params = nil
	c.handlers = nil
	c.index = -1
	c.values = nil
	// c.Errors = c.Errors[0:0]
	// c.Accepted = nil
}

// Copy a new context
func (c *Context) Copy() *Context {
	var ctx = *c
	ctx.handlers = nil
	ctx.index = abortIndex

	return &ctx
}

/*************************************************************
 * getter/setter methods
 *************************************************************/

// Req get request instance
func (c *Context) Req() *http.Request {
	return c.req
}

// Res get response instance
func (c *Context) Res() http.ResponseWriter {
	return c.res
}

// Params get current route params
func (c *Context) Params() Params {
	return c.params
}

// SetParams to the context
func (c *Context) SetParams(params Params) {
	c.params = params
}

/*************************************************************
 * Context helper methods
 *************************************************************/

// URL get URL instance from request
func (c *Context) URL() *url.URL {
	return c.req.URL
}

// GetRawData return stream data
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.req.Body)
}

// WriteString to response
func (c *Context) WriteString(str string) (n int, err error) {
	return c.res.Write([]byte(str))
}

// Write byte data to response
func (c *Context) Write(bt []byte) (n int, err error) {
	return c.res.Write(bt)
}

// WriteBytes byte data to response
func (c *Context) WriteBytes(bt []byte) (n int, err error) {
	return c.res.Write(bt)
}
