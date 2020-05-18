package api

import (
	"net/url"
)

type UserOpts struct {
	RequestOpts
}

func (opts UserOpts) ToValues() url.Values {
	return opts.RequestOpts.ToValues()
}
