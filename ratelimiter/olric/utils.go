package ratelimiter

import (
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

func parseURL(configURL string) (olricPeers []string, userRateLimits map[string]int64, whitelistedIps map[string]bool, err error) {
	urlObject, err := url.Parse(configURL)
	if err != nil {
		return
	}
	values := urlObject.Query()

	olricPeers = strings.Split(urlObject.Host, ",")
	if len(olricPeers) == 1 {
		if olricPeers[0] == "" {
			err = fmt.Errorf("missing olric address")
			return
		}
		if olricPeers[0] != "local" {
			err = fmt.Errorf("need at least 2 hosts to create a cluster, use \"local\" for single instances instead")
			return
		}
	} else {
		for _, peer := range olricPeers {
			if !strings.Contains(peer, ":") {
				err = fmt.Errorf("invalid host [%s], needs to be specified as host:port", peer)
				return
			}
		}
	}

	userRateLimitsString := values.Get("rates")
	if userRateLimitsString == "" {
		// rate limits are optional
		userRateLimits = map[string]int64{}
	} else {
		userRateLimits, err = constructRateLimits(userRateLimitsString)
		if err != nil {
			return
		}
	}

	whitelistedIpsString := values.Get("whitelist")
	if whitelistedIpsString == "" {
		// whitelist is optional
		whitelistedIps = map[string]bool{}
	} else {
		whitelistedIps, err = constructWhitelist(whitelistedIpsString)
		if err != nil {
			return
		}
	}

	return
}

func constructRateLimits(in string) (map[string]int64, error) {
	userRateLimits, err := parseRateLimitsString(in)
	if err != nil {
		return nil, err
	}
	return userRateLimits, nil
}

func constructWhitelist(in string) (map[string]bool, error) {
	whitelistEntries := strings.Split(in, ",")
	out := make(map[string]bool)

	// todo remove
	zlog.Info("parsing whitelist",
		zap.String("in", in),
		zap.String("split", strings.Join(whitelistEntries, ",")),
	)

	for _, entry := range whitelistEntries {
		out[entry] = true
		// todo remove
		zlog.Info("adding whitelist",
			zap.String("entry", entry),
		)
		// todo check if valid ip?
	}
	return out, nil
}

func parseRateLimitsString(in string) (map[string]int64, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make(map[string]int64)
	for _, pair := range strings.Split(in, ",") {
		kv := strings.Split(pair, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid value pair for rate limits: %s", pair)
		}
		asInt64, err := strconv.ParseInt(kv[1], 10, 64)
		if err != nil {
			return nil, err
		}
		out[kv[0]] = asInt64
	}

	return out, nil
}

func validateServices(userRateLimits map[string]int64, serviceNames []string) error {
	for providedName := range userRateLimits {
		if !contains(serviceNames, providedName) {
			return fmt.Errorf("invalid service name: %s", providedName)
		}
	}
	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
