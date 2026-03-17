package client

import (
	"fmt"
	"net/url"
)

// PageParams holds pagination parameters.
type PageParams struct {
	Page    int
	PerPage int
}

// DefaultPageParams returns sensible default pagination.
func DefaultPageParams() PageParams {
	return PageParams{Page: 1, PerPage: 100}
}

// Apply adds pagination query parameters to the given URL values.
func (p PageParams) Apply(v url.Values) {
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("perpage", fmt.Sprintf("%d", p.PerPage))
	}
}

// QueryString returns the pagination parameters as a query string fragment.
func (p PageParams) QueryString() string {
	v := url.Values{}
	p.Apply(v)
	qs := v.Encode()
	if qs != "" {
		return "?" + qs
	}
	return ""
}

// AppendToPath appends pagination query params to an existing path that may have query params.
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
