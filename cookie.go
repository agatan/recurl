package main

import "strings"

type Cookie struct {
	ID     SessionID
	Values map[string]string
}

func NewCookie(cookie string) *Cookie {
	pairs := strings.Split(cookie, ";")
	values := make(map[string]string)
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		sp := strings.IndexRune(pair, '=')
		values[pair[:sp]] = pair[sp+1:]
	}
	return &Cookie{Values: values}
}

type Cookies []*Cookie

func (cs *Cookies) Append(c *Cookie) {
	*cs = append(*cs, c)
}

func (cs Cookies) FindMatch(m *Cookie) (*Cookie, bool) {
	for _, c := range cs {
		match := true
		for k, v := range m.Values {
			if c.Values[k] != v {
				match = false
				break
			}
		}
		if match {
			return c, true
		}
	}
	return nil, false
}
