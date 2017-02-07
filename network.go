package main

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/docker/libcontainer/netlink"
)

// GetDefaultIP returns the machine default IP adress
func GetDefaultIP() (string, error) {
	if len(os.Getenv("DEV_ENV")) > 0 {
		return "127.0.0.1", nil
	}

	defaultIfaceName, err := getDefaultGatewayIfaceName()
	if err != nil {
		// A default route is not required; log it and keep going.
		log.Println(nil, err)
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			// Record IPv4 network settings. Stop at the first IPv4 address
			// found for the interface.
			if err == nil && ip.To4() != nil {
				if defaultIfaceName == iface.Name {
					return ip.String(), nil
				}
				break
			}
		}
	}

	return "", errors.New("not found")
}

func getDefaultGatewayIfaceName() (string, error) {
	routes, err := netlink.NetworkGetRoutes()
	if err != nil {
		return "", err
	}
	for _, route := range routes {
		if route.Default {
			if route.Iface == nil {
				return "", errors.New("found default route but could not determine interface")
			}
			return route.Iface.Name, nil
		}
	}
	return "", errors.New("unable to find default route")
}
