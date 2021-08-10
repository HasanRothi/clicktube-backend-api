package helpers

import (
	"net"
	"net/url"

	"regexp"

	"github.com/haccer/available"
)

func UrlValidator(link string) bool {
	match, _ := regexp.MatchString("([http | https]+)://([a-z]+)([.]+)([a-z])", link)
	if match == true {
		val, _ := url.Parse(link)
		available := available.Domain(val.Host)
		if !available {
			ips, _ := net.LookupIP(val.Host)
			if len(ips) == 0 {
				return false
			} else {
				return true
			}
		} else {
			return false
		}
	}

	// fmt.Println("IP :", ips)
	// for _, ip := range ips {
	// 	if ipv4 := ip.To4(); ipv4 != nil {
	// 		fmt.Println("IPv4: ", ipv4)
	// 	}
	// }
	// }
	return true
}
