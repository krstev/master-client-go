package service

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Info",
		"GET",
		"/info",
		Info,
	},
	Route{
		"Health",
		"POST",
		"/health",
		Health,
	}, Route{
		"PrimeNumbers",
		"POST",
		"/primeNumbers",
		v1PrimeNumbers,
	},
	Route{
		"CountVowels",
		"POST",
		"/countVowels",
		v1CountVowels,
	},
	Route{
		"GoogleSearch",
		"POST",
		"/googleQuery",
		GoogleSearch,
	},
	Route{
		"googleQueryTimeout",
		"POST",
		"/googleQueryTimeout",
		googleQueryTimeout,
	},
	Route{
		"googleQueryTimeoutReplica",
		"POST",
		"/googleQueryTimeoutReplica",
		googleQueryTimeoutReplica,
	},
	Route{
		"v2PrimeNumbers",
		"POST",
		"/v2/primeNumbers",
		v2PrimeNumbers,
	},
	Route{
		"v2CountVowels",
		"POST",
		"/v2/countVowels",
		v2CountVowels,
	},
}
