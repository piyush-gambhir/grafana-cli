package client

import (
	"fmt"
	"net/url"
)

type PageParams struct {
	Page    int
	PerPage int
}

func DefaultPageParams() PageParams {
	return PageParams{Page: 1, PerPage: 100}
}

func (p PageParams) Apply(v url.Values) {
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("perpage", fmt.Sprintf("%d", p.PerPage))
	}
}

func (p PageParams) QueryString() string {
	v := url.Values{}
	p.Apply(v)
	qs := v.Encode()
	if qs != "" {
		return "?" + qs
	}
	return ""
}

func (p PageParams) AppendToPath(path string) string {
	if p.Page <= 0 && p.PerPage <= 0 {
		return path
	}
	sep := "?"
	for _, c := range path {
		if c == '?' {
			sep = "&"
			break
		}
	}
	v := url.Values{}
	p.Apply(v)
	return path + sep + v.Encode()
}
