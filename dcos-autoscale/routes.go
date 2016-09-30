package main

import (
	"github.com/dcos/dcos-autoscale/common"
)

var routes = map[string]map[string]common.Handler{

	"POST": {
        "/autoscale/api/v1/auto_query_list":                 auto_query_list,
        "/autoscale/api/v1/auto_query_details":              auto_query_details,
		"/autoscale/api/v1/elastic_expansion_configuration": elastic_expansion_configuration,
	},
	"PATCH": {
		"/autoscale/api/v1/updating_elastic_expansion":      updating_elastic_expansion,
		"/autoscale/api/v1/suspend":                         suspend,
		"/autoscale/api/v1/turn_on":                         turn_on,
	},
}
