package api

import (
	"net/url"
	"time"
)

type ProjectMergeRequestOpts struct {
	State         string
	CreatedBefore time.Time
	CreatedAfter  time.Time
	UpdatedBefore time.Time
	UpdatedAfter  time.Time
	RequestOpts
}

func (opts ProjectMergeRequestOpts) ToValues() url.Values {
	urlValues := opts.RequestOpts.ToValues()
	if len(opts.State) > 0 {
		urlValues.Set("state", opts.State)
	}
	if !opts.CreatedBefore.IsZero() {
		urlValues.Set("created_before", opts.CreatedBefore.String())
	}
	if !opts.CreatedAfter.IsZero() {
		urlValues.Set("created_after", opts.CreatedAfter.String())
	}
	if !opts.UpdatedBefore.IsZero() {
		urlValues.Set("updated_before", opts.UpdatedBefore.String())
	}
	if !opts.UpdatedAfter.IsZero() {
		urlValues.Set("updated_after", opts.UpdatedAfter.String())
	}
	return urlValues
}
