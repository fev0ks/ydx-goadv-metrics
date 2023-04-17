package middlewares

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
)

type SubNetChecker struct {
	trustedSubnet string
}

func NewSubNetChecker(trustedSubnet string) *SubNetChecker {
	return &SubNetChecker{trustedSubnet: trustedSubnet}
}

func (sc *SubNetChecker) CheckRealIPRequestHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sc.trustedSubnet == "" {
			next.ServeHTTP(w, r)
			return
		}

		if clientIPVal := sc.getClientIp(r); clientIPVal != "" {
			ok, err := sc.isTrustedIp(clientIPVal)
			if err != nil {
				log.Printf("failed to check clientIP: %v", err)
			} else if !ok {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (sc *SubNetChecker) getClientIp(r *http.Request) (clientIPVal string) {
	if ipHeader := r.Header.Get("X-Real-Ip"); ipHeader != "" {
		clientIPVal = ipHeader
	} else {
		ips := r.Header.Get("X-Forwarded-For")
		ipStrs := strings.Split(ips, ",")
		clientIPVal = ipStrs[0]
	}
	return
}

func (sc *SubNetChecker) isTrustedIp(ipVal string) (bool, error) {
	_, ipNet, err := net.ParseCIDR(ipVal)
	if err != nil {
		return false, err
	}
	trustedSubnet := net.ParseIP(sc.trustedSubnet)
	if trustedSubnet == nil {
		return false, errors.New("failed to parse trustedSubnet")
	}
	return ipNet.Contains(trustedSubnet), nil
}
