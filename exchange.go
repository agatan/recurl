package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"unicode/utf8"
)

type SessionID int

type Exchange struct {
	SessionID *SessionID
	Request   *Request
	Response  *Response
}

type Request struct {
	Method     string
	URL        *url.URL
	Proto      string // "HTTP/1.0"
	Header     http.Header
	Body       []byte `json:"Body,omitempty"`
	BodyString string `json:"BodyString,omitempty"`
	Host       string
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
	req := &Request{
		Method: r.Method,
		URL:    &u,
		Proto:  r.Proto,
		Header: r.Header,
		Host:   r.Host,
	}
	if utf8.Valid(body) {
		req.BodyString = string(body)
	} else {
		req.Body = body
	}
	return req, nil
}

type Response struct {
	StatusCode int    // e.g. "200 OK"
	Proto      string // e.g. "HTTP/1.0"
	Header     http.Header
	Body       []byte `json:"Body,omitempty"`
	BodyString string `json:"BodyString,omitempty"`
}

// NewResponse returns Response object from http.Response.
func NewResponse(r *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	resp := &Response{
		StatusCode: r.StatusCode,
		Proto:      r.Proto,
		Header:     r.Header,
	}
	if utf8.Valid(body) {
		resp.BodyString = string(body)
	} else {
		resp.Body = body
	}
	return resp, nil
}
