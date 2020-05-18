package api

import (
	"net/url"
	"strconv"
)

type RequestOpts struct {
	PerPage        uint32
	Page           uint32
}

func (opts RequestOpts) ToValues() url.Values {
	urlValues := url.Values{}

	if opts.PerPage > 0 {
		urlValues.Set("per_page", strconv.FormatUint(uint64(opts.PerPage), 10))
	}

	if opts.Page > 0 {
		urlValues.Set("page", strconv.FormatUint(uint64(opts.Page), 10))
	}

	return urlValues
}