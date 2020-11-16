package dredd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"strings"
)

type IpQuotaHandler struct {
	defaultQuota int
	ipNets       map[*net.IPNet]int
	ips          map[string]int
}

type IpQuotaCategory struct {
	Quota int      `yaml:"quota"`
	Ips   []string `yaml:"ips"`
}

func NewIpQuotaHandler(defaultQuota int) *IpQuotaHandler {

	quotas := &IpQuotaHandler{}
	quotas.defaultQuota = defaultQuota
	quotas.ipNets = make(map[*net.IPNet]int)
	quotas.ips = make(map[string]int)

	return quotas
}

func NewIpQuotaHandlerFromFile(path string, defaultQuota int) (*IpQuotaHandler, error) {

	quotas := NewIpQuotaHandler(defaultQuota)

	b, err := ioutil.ReadFile(path)
	quotaFile := make(map[string]IpQuotaCategory)

	if err != nil {
		return nil, fmt.Errorf("failed to open quota file: %w", err)
	}

	err = yaml.Unmarshal(b, quotaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	for _, cat := range quotaFile {
		quota := cat.Quota

		for _, ipString := range cat.Ips {
			if strings.Contains(ipString, "/") {
				_, ipNet, err := net.ParseCIDR(ipString)

				if err != nil {
					return nil, fmt.Errorf("failed to parse cidr: %w", err)
				}
				quotas.ipNets[ipNet] = quota
			} else {
				ip := net.ParseIP(ipString)

				if ip == nil {
					return nil, fmt.Errorf("failed to parse ip: %w", err)
				}
				quotas.ips[ip.String()] = quota
			}
		}
	}

	return quotas, nil
}

func (w *IpQuotaHandler) GetQuota(ipString string) (int, error) {

	ip := net.ParseIP(ipString)

	if ip == nil {
		return w.defaultQuota, fmt.Errorf("failed to parse ip: %s", ipString)
	}

	if _, ok := w.ips[ip.String()]; ok {
		return w.ips[ip.String()], nil
	}

	for k, v := range w.ipNets {
		if k.Contains(ip) {
			return v, nil
		}
	}

	return w.defaultQuota, nil
}
