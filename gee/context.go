package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// init info
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string

	//response info
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Writer: w, Req: r, Path: r.URL.Path, Method: r.Method}
}

// === 实现获取请求中的参数 ( GET, POST )
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key) // https://www.google.com?q="why"  --> parse the params
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code) // WriteHeader sends an HTTP response header with the provided status code
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// === 实现发送 string, json, html 功能的封装 ===
func (c *Context) String(code int, formatString string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(formatString, values...)))
}

type H map[string]interface{} // obj in JSON

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
