package dredd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"strings"
)

type IpLimitHandler struct {
	defaultQuota int
	defaultRate  int
	ipNets       map[*net.IPNet]Limit
	ips          map[string]Limit
}

type Limit struct {
	Quota int
	Rate  int
}

type IpLimitCategory struct {
	Quota int      `yaml:"quota"`
	Rate  int      `yaml:"rate"`
	Ips   []string `yaml:"ips"`
}

func NewIpLimitsHandler(defaultQuota int, defaultRate int) *IpLimitHandler {

	limits := &IpLimitHandler{}
	limits.defaultQuota = defaultQuota
	limits.defaultRate = defaultRate
	limits.ipNets = make(map[*net.IPNet]Limit)
	limits.ips = make(map[string]Limit)

	return limits
}

func NewIpLimitsHandlerFromFile(path string, defaultQuota int, defaultRate int) (*IpLimitHandler, error) {

	limits := NewIpLimitsHandler(defaultQuota, defaultRate)

	b, err := ioutil.ReadFile(path)
	limitsFile := make(map[string]IpLimitCategory)

	if err != nil {
		return nil, fmt.Errorf("failed to open limits file: %w", err)
	}

	err = yaml.Unmarshal(b, limitsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	for _, cat := range limitsFile {
		limit := Limit{Quota: cat.Quota, Rate: cat.Rate}

		for _, ipString := range cat.Ips {
			if strings.Contains(ipString, "/") {
				_, ipNet, err := net.ParseCIDR(ipString)

				if err != nil {
					return nil, fmt.Errorf("failed to parse cidr: %w", err)
				}
				limits.ipNets[ipNet] = limit
			} else {
				ip := net.ParseIP(ipString)

				if ip == nil {
					return nil, fmt.Errorf("failed to parse ip: %w", err)
				}
				limits.ips[ip.String()] = limit
			}
		}
	}

	return limits, nil
}

func (w *IpLimitHandler) GetLimits(ipString string) (Limit, error) {

	ip := net.ParseIP(ipString)
	defaultLimits := Limit{Quota: w.defaultQuota, Rate: w.defaultRate}

	if ip == nil {
		return defaultLimits, fmt.Errorf("failed to parse ip: %s", ipString)
	}

	if _, ok := w.ips[ip.String()]; ok {
		return w.ips[ip.String()], nil
	}

	for k, v := range w.ipNets {
		if k.Contains(ip) {
			return v, nil
		}
	}

	if defaultLimits.Rate < 0 || defaultLimits.Quota < 0 {
		return defaultLimits, fmt.Errorf("ip not whitelisted: %s", ipString)
	}

	return defaultLimits, nil
}
