package api

import (
	"net/url"
	"time"
)

type EventOpts struct {
	Action     string
	TargetType string
	Before     time.Time
	After      time.Time
	RequestOpts
}

func (opts EventOpts) ToValues() url.Values {
	urlValues := opts.RequestOpts.ToValues()
	if len(opts.Action) > 0 {
		urlValues.Set("action", opts.Action)
	}
	if len(opts.TargetType) > 0 {
		urlValues.Set("targetType", opts.TargetType)
	}
	if !opts.Before.IsZero() {
		urlValues.Set("before", opts.Before.Format("2006-01-02"))
	}
	if !opts.After.IsZero() {
		urlValues.Set("after", opts.After.Format("2006-01-02"))
	}
	return urlValues
}
