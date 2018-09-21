package cache

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
)

type contextKeyType int

var cacheHitContextKey = new(contextKeyType)

type cacheObj struct {
	r       *http.Response
	content []byte
}

func (c *cacheObj) copy() *http.Response {
	r := c.r
	r.Body = ioutil.NopCloser(bytes.NewBuffer(c.content))
	return r
}

// IsCacheHit returns true if the context was used for a cached API request
func IsCacheHit(ctx context.Context) bool {
	return ctx.Value(cacheHitContextKey) != nil
}

// AddCaching adds caching to the Session
func AddCaching(s *session.Session, cacheConfig *Config) {

	// Handle caching
	s.Handlers.Validate.PushFront(func(r *request.Request) {
		// Handle cache flushes on requests that would modify the cache contents
		cacheConfig.flushCaches(r)

		// Get item from cache
		i := cacheConfig.get(r)

		if i != nil && !i.Expired() {
			cacheConfig.incHit(r)

			// Copy the cached response to this requests response
			r.HTTPResponse = i.Value().(*cacheObj).copy()

			// Add cache hit marker to the request context
			r.HTTPRequest = r.HTTPRequest.WithContext(context.WithValue(r.HTTPRequest.Context(), cacheHitContextKey, true))
		}
	})

	// Add an empty handler to the top of Send and short circuit the rest on a cache hit
	s.Handlers.Send.PushFront(func(r *request.Request) {})
	s.Handlers.Send.AfterEachFn = shortCircuitRequestHandler

	// ValidateResponse is the first handler after Send, cache the response if this was not a cached hit
	s.Handlers.ValidateResponse.PushFront(func(r *request.Request) {
		if !IsCacheHit(r.HTTPRequest.Context()) {
			cacheConfig.incMiss(r)

			content, _ := ioutil.ReadAll(r.HTTPResponse.Body)
			r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(content))

			cacheConfig.set(r, &cacheObj{
				r:       r.HTTPResponse,
				content: content,
			})
		}
	})
}

// Returns false when request is a cache hit, used to short circuit request handlers
var shortCircuitRequestHandler = func(item request.HandlerListRunItem) bool {
	return !IsCacheHit(item.Request.HTTPRequest.Context())
}
