package context

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	ctx *gin.Context
}

func NewContext(c *gin.Context) *Context {
	return &Context{ctx: c}
}

func (c *Context) GinContext() *gin.Context {
	return c.ctx
}

// JSON sends a JSON response with the given status code and data.
func (c *Context) JSON(code int, data interface{}) {
	c.ctx.JSON(code, data)
}

// String sends a string response with the given status code.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.ctx.String(code, format, values...)
}

// Status sets the HTTP status code.
func (c *Context) Status(code int) {
	c.ctx.Status(code)
}

// Param returns the value of the URL parameter.
func (c *Context) Param(key string) string {
	return c.ctx.Param(key)
}

// Query returns the value of the query parameter.
func (c *Context) Query(key string) string {
	return c.ctx.Query(key)
}

// DefaultQuery returns the value of the query parameter or a default value if not present.
func (c *Context) DefaultQuery(key, defaultValue string) string {
	return c.ctx.DefaultQuery(key, defaultValue)
}

// GetQuery returns the value of the query parameter and a boolean indicating if it exists.
func (c *Context) GetQuery(key string) (string, bool) {
	return c.ctx.GetQuery(key)
}

// Header returns the value of the request header.
func (c *Context) Header(key string) string {
	return c.ctx.Request.Header.Get(key)
}

// SetHeader sets a response header.
func (c *Context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

// Method returns the HTTP method of the request.
func (c *Context) Method() string {
	return c.ctx.Request.Method
}

// Path returns the URL path of the request.
func (c *Context) Path() string {
	return c.ctx.Request.URL.Path
}

// ClientIP returns the client's IP address.
func (c *Context) ClientIP() string {
	return c.ctx.ClientIP()
}

// Next executes the next handler in the chain.
func (c *Context) Next() {
	c.ctx.Next()
}

// Abort prevents pending handlers from being called.
func (c *Context) Abort() {
	c.ctx.Abort()
}

// AbortWithStatus aborts the request and sets the HTTP status code.
func (c *Context) AbortWithStatus(code int) {
	c.ctx.AbortWithStatus(code)
}

// AbortWithStatusJSON aborts the request and returns a JSON response.
func (c *Context) AbortWithStatusJSON(code int, data interface{}) {
	c.ctx.AbortWithStatusJSON(code, data)
}

// Set stores a key-value pair in the context.
func (c *Context) Set(key string, value interface{}) {
	c.ctx.Set(key, value)
}

// Get retrieves a value from the context by key.
func (c *Context) Get(key string) (interface{}, bool) {
	return c.ctx.Get(key)
}

// MustGet retrieves a value from the context by key or panics if not found.
func (c *Context) MustGet(key string) interface{} {
	return c.ctx.MustGet(key)
}

// StatusCode returns the HTTP status code of the response.
func (c *Context) StatusCode() int {
	return c.ctx.Writer.Status()
}

type Handler func(c *Context)

// Wrap converts a Handler function into a gin.HandlerFunc.
// This is the bridge between our clean API and Gin's middleware system.
func Wrap(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		wrappedCtx := NewContext(c)
		h(wrappedCtx)
	}
}

// Response is a helper struct for standard API responses.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// OK sends a successful JSON response.
func (c *Context) OK(data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: "success",
		Data:   data,
	})
}

// Created sends a 201 Created response.
func (c *Context) Created(data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Status: "success",
		Data:   data,
	})
}

// BadRequest sends a 400 Bad Request response.
func (c *Context) BadRequest(message string) {
	c.JSON(http.StatusBadRequest, Response{
		Status:  "error",
		Message: message,
	})
}

// Unauthorized sends a 401 Unauthorized response.
func (c *Context) Unauthorized(message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Status:  "error",
		Message: message,
	})
}

// Forbidden sends a 403 Forbidden response.
func (c *Context) Forbidden(message string) {
	c.JSON(http.StatusForbidden, Response{
		Status:  "error",
		Message: message,
	})
}

// NotFound sends a 404 Not Found response.
func (c *Context) NotFound(message string) {
	c.JSON(http.StatusNotFound, Response{
		Status:  "error",
		Message: message,
	})
}

// InternalError sends a 500 Internal Server Error response.
func (c *Context) InternalError(message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Status:  "error",
		Message: message,
	})
}
