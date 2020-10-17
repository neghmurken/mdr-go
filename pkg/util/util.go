package util

import (
	"fmt"
	"flag"
	"os"
	"net"
	"errors"
	"strings"
	"unicode"
)

func Quit(message string, err error) {
    fmt.Println(message)
    if nil != err {
        fmt.Println(err)
    }

    flag.PrintDefaults()
    os.Exit(1)
}

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func CleanData(data string) string {
	return strings.TrimFunc(data, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})
}