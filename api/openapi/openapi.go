package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create short URL from original URL
	// (POST /)
	CreateShortURL(w http.ResponseWriter, r *http.Request)
	// Get stats about redirects
	// (GET /stats/{short-url})
	GetStats(w http.ResponseWriter, r *http.Request, shortUrl string)
	// Redirect to original URL by short URL
	// (GET /{short-url})
	RedirectURL(w http.ResponseWriter, r *http.Request, shortUrl string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// CreateShortURL operation middleware
func (siw *ServerInterfaceWrapper) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateShortURL(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetStats operation middleware
func (siw *ServerInterfaceWrapper) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "short-url" -------------
	var shortUrl string

	err = runtime.BindStyledParameter("simple", false, "short-url", chi.URLParam(r, "short-url"), &shortUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter short-url: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStats(w, r, shortUrl)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// RedirectURL operation middleware
func (siw *ServerInterfaceWrapper) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "short-url" -------------
	var shortUrl string

	err = runtime.BindStyledParameter("simple", false, "short-url", chi.URLParam(r, "short-url"), &shortUrl)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter short-url: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RedirectURL(w, r, shortUrl)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/", wrapper.CreateShortURL)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/stats/{short-url}", wrapper.GetStats)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/{short-url}", wrapper.RedirectURL)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8RUTW/TQBD9K6OBo5s4beHgExSkqFIPKKEn2sPGniTb2rtmdhyIovx3NOs4H00qFSji",
	"5qx33jy/j6ww91XtHTkJmK0w5HOqTHwc0feGgtyObvRXzb4mFkvxnWc7s86Um5dTz5URzLDhEhOUZU2Y",
	"YRC2bobr9fbETx4oF1wnOKJQexfoJHqYe5YXQScYxEj4YxpjnT4m4JpqRIVlytu3W2Dr5P3lDto6oRlx",
	"5PFizsc89Mi6qdfZ3DsxuegjVcaWcZBq484Xvnz0iw9L4wr62eNGtxYUcra1WO8ww69zG8AGkDnBVO2B",
	"mr2ugKlnGBI9XrGxLkDuGw4Ed3hl8kdyBXymBZW+rsgJ/LAyh6HvwY0ewuAO9RuslEr5dnQDY/1UcsTw",
	"8cs1JrggDi2BQS/tpcrL1+RMbTHDi17aG2CCtZF51LIf5fYhfqKKbpT9dYEZfmIyQuNOyQS5jeCVL5ad",
	"NuTioKnr0uZxtP8QdHmXXX16yzTFDN/0d+Hub5Ld34t1FF53WKYCM+GG4kEbzcj2PB284uZd5uPqQ/di",
	"gqDhEvKoQ6FCXqapwh7enJgCNtLonXen7mg0WSMQiBfEQMyeY/hCU1WGl1u9oV2szk7ZV9B1G1oPxMwC",
	"Zt9w3N3Ce0Xpx+L1V3H4rOFyrRxmdMLWIUnbM00Bm4qEWCFXJwWIPDzDjESsmwFvighhg2H1rsYJE3Sm",
	"om7yrG3boZ/JnjdPe3h/5HX6al63H3zK5SbPKYRpU8JWpdboy2MTndfqNq74K5uHtBEPzMQ3slU07Lsb",
	"+bbOvsTT7u+xzchv2NotVwv+oZcX6cWxWP9B+04nEH9QLJgsd717rmXr7fnT1bG52g7jim1T9vE2Yu7g",
	"1slTjGE39kw0Oog2yffrXwEAAP///VL4ECsIAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
