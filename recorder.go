package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type Recorder struct {
	Exchanges     []*Exchange
	Sessions      map[string]SessionID
	lastSessionID SessionID
	target        *url.URL
	mu            sync.Mutex
}

func NewRecorder(target *url.URL) *Recorder {
	return &Recorder{
		Sessions: make(map[string]SessionID),
		target:   target,
	}
}

func (rec *Recorder) AddExchange(req *Request, resp *Response) error {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	ex := &Exchange{Request: req, Response: resp}
	if cookie := req.Header.Get("Cookie"); cookie != "" {
		// cookie exists.
		sess, ok := rec.Sessions[cookie]
		if !ok {
			return errors.Errorf("unknown cookie: %v", cookie)
		}
		ex.SessionID = &sess
	}
	if cookie := resp.Header.Get("Set-Cookie"); cookie != "" {
		sess := rec.registerSession(cookie)
		ex.SessionID = &sess
	}
	rec.Exchanges = append(rec.Exchanges, ex)
	return nil
}

func (rec *Recorder) registerSession(cookie string) SessionID {
	rec.lastSessionID++
	rec.Sessions[cookie] = rec.lastSessionID
	return rec.lastSessionID
}

func (rec *Recorder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := rec.target
	reverseProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.URL.Path = u.Path + r.URL.Path
			r.Host = u.Host
		},
	}
	var response *Response
	reverseProxy.ModifyResponse = func(resp *http.Response) error {
		var err error
		response, err = NewResponse(resp)
		return err
	}
	request, err := NewRequest(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	reverseProxy.ServeHTTP(w, r)
	if err := rec.AddExchange(request, response); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
