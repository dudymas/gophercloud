package meters

import (
	"github.com/rackspace/gophercloud"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToMeterListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type ListOpts struct {
}

// ToMeterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToMeterListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List makes a request against the API to list meters accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) listResult {
	var res listResult
	url := listURL(client)

	if opts != nil {
		query, err := opts.ToMeterListQuery()
		if err != nil {
			res.Err = err
			return res
		}
		url += query
	}

	_, res.Err = client.Get(url, &res.Body, &gophercloud.RequestOpts{})
	return res
}

// OptsKind describes the mode with which a given set of opts should be tranferred
type OptsKind string

var (
	//BodyContentOpts is a kind of option serialization. The MeterStatisticsOptsBuilder is expected
	//to emit JSON from ToMeterStatisticsQuery()
	BodyContentOpts = OptsKind("Body")
	//QueryOpts is a kind of option serialization. The MeterStatisticsOptsBuilder is expected
	//to emit uri encoded fields from ToMeterStatisticsQuery()
	QueryOpts = OptsKind("Query")
)

// MeterStatisticsOptsBuilder allows extensions to add additional parameters to the
// List request.
type MeterStatisticsOptsBuilder interface {
	Kind() OptsKind
	ToMeterStatisticsQuery() (string, error)
}

// MeterStatisticsOpts allows the filtering and sorting of collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the server attributes you want to see returned.
type MeterStatisticsOpts struct {
	QueryField string `q:"q.field"`
	QueryOp    string `q:"q.op"`
	QueryValue string `q:"q.value"`

	// Optional group by
	GroupBy string `q:"groupby"`

	// Optional number of seconds in a period
	Period int `q:"period"`
}

// Kind returns QueryOpts by default for MeterStatisticsOpts
func (opts MeterStatisticsOpts) Kind() OptsKind {
	return QueryOpts
}

// ToMeterStatisticsQuery formats a StatisticsOpts into a query string.
func (opts MeterStatisticsOpts) ToMeterStatisticsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// MeterStatistics gathers statistics based on filters, groups, and period options
func MeterStatistics(client *gophercloud.ServiceClient, n string, optsBuilder MeterStatisticsOptsBuilder) statisticsResult {
	var (
		res  statisticsResult
		url  = statisticsURL(client, n)
		opts gophercloud.RequestOpts
		err  error
		kind OptsKind
	)

	if optsBuilder != nil {
		kind = optsBuilder.Kind()
	}

	switch kind {
	case QueryOpts:
		query, err := optsBuilder.ToMeterStatisticsQuery()
		url += query
	case BodyContentOpts:
		opts.JSONBody, err = optsBuilder.ToMeterStatisticsQuery()
	}

	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Get(url, &res.Body, &opts)
	return res
}
