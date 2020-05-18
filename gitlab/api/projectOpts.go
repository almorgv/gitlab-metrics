package api

import (
	"net/url"
	"strconv"
)

type ProjectsOpts struct {
	MinAccessLevel uint32
	RequestOpts
}

func (opts ProjectsOpts) ToValues() url.Values {
	urlValues := opts.RequestOpts.ToValues()
	if opts.MinAccessLevel > 0 {
		urlValues.Set("min_access_level", strconv.FormatUint(uint64(opts.MinAccessLevel), 10))
	}
	return urlValues
}
