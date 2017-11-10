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
	Sessions      Cookies
	lastSessionID SessionID
	target        *url.URL
	mu            sync.Mutex
}

func NewRecorder(target *url.URL) *Recorder {
	return &Recorder{
		target: target,
	}
}

func (rec *Recorder) AddExchange(req *Request, resp *Response) error {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	ex := &Exchange{Request: req, Response: resp}
	if cookie := req.Header.Get("Cookie"); cookie != "" {
		// cookie exists.
		matched, ok := rec.Sessions.FindMatch(NewCookie(cookie))
		if !ok {
			return errors.Errorf("unknown cookie: %v", cookie)
		}
		ex.SessionID = &matched.ID
	}
	if cookie := resp.Header.Get("Set-Cookie"); cookie != "" {
		c := NewCookie(cookie)
		if ex.SessionID == nil {
			c.ID = rec.newSessionID()
		} else {
			c.ID = *ex.SessionID
		}
		rec.registerCookie(c)
	}
	rec.Exchanges = append(rec.Exchanges, ex)
	return nil
}

func (rec *Recorder) newSessionID() SessionID {
	rec.lastSessionID++
	return rec.lastSessionID
}

func (rec *Recorder) registerCookie(cookie *Cookie) {
	rec.Sessions.Append(cookie)
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
