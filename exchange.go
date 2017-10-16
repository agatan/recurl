package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SessionID int

type Exchange struct {
	SessionID *SessionID
	Request   *Request
	Response  *Response
}

type Request struct {
	Method string
	URL    *url.URL
	Proto  string // "HTTP/1.0"
	Header http.Header
	Body   []byte
	Host   string
}

// NewRequest returns Request object from http.Request.
func NewRequest(r *http.Request) (*Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	u := *r.URL
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	return &Request{
		Method: r.Method,
		URL:    &u,
		Proto:  r.Proto,
		Header: r.Header,
		Body:   body,
		Host:   r.Host,
	}, nil
}

type Response struct {
	StatusCode int    // e.g. "200 OK"
	Proto      string // e.g. "HTTP/1.0"
	Header     http.Header
	Body       []byte
}

// NewResponse returns Response object from http.Response.
func NewResponse(r *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	return &Response{
		StatusCode: r.StatusCode,
		Proto:      r.Proto,
		Header:     r.Header,
		Body:       body,
	}, nil
}
