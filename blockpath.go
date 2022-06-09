// Package plugin_blockpath a plugin to block a path.
package plugin_blockpath

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Element for holding response code and regex 
type Element struct {
	Regex string `json:"regex,omitempty"`
	Response int `json:"code,omitempty"`
}

// Config holds the plugin configuration.
type Config struct {
	Elements []Element `json:"elements,omitempty"`
}

// container for compiled targets
type CompiledElement struct {
	Regexp *regexp.Regexp
	Code   int 
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type blockPath struct {
	name    string
	next    http.Handler
	// regexps []*regexp.Regexp
	el      []*CompiledElement
}

// New creates and returns a plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// elms := make([]*CompiledElement, len(config.Elements))

	elms := []*CompiledElement{}

	// for i, _ := range config.Elements {
	for i := 0; i < len(config.Elements); i++ {
		re, err := regexp.Compile(config.Elements[i].Regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", config.Elements[i].Regex, err)
		}

		// set http response code
		code := config.Elements[i].Response
		if code == 0 {
			code = 401
		}

		// regexps[i] = re
		elms = append(elms, &CompiledElement{
			Regexp: re, 
			Code: code,
		})
	}

	return &blockPath{
		name:    name,
		next:    next,
		el:      elms,
		// regexps: regexps,
	}, nil
}

func lookupStatusMessage(code int) (ret string, err error) {
	switch {
	case code == 404:
		ret = "404 page not found"
	default:
		err = fmt.Errorf("Not able to decode status code.")
	}
	return ret, err  
}

func (b *blockPath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.EscapedPath()

	for _, e := range b.el {
		if e.Regexp.MatchString(currentPath) {
			rw.WriteHeader(e.Code)
			sc, err := lookupStatusMessage(e.Code)
			if err != nil {
				rw.Write( []byte(sc) )
			}
			return
		}
	}

	b.next.ServeHTTP(rw, req)
}
