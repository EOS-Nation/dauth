package redis

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func parseURL(configURL string) (redisAddr string, userRateLimits map[string]int64, err error) {
	urlObject, err := url.Parse(configURL)
	if err != nil {
		return
	}
	values := urlObject.Query()

	redisAddr = urlObject.Host
	if redisAddr == "" {
		err = fmt.Errorf("missing redis address")
		return
	}

	userRateLimitsString := values.Get("rates")
	if userRateLimitsString == "" {
		// rate limits are optional
		userRateLimits = map[string]int64{}
		return
	}

	userRateLimits, err = constructRateLimits(userRateLimitsString)
	if err != nil {
		return
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
