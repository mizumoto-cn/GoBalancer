package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// GetIP returns the real client IP of the request.
// If the X-Forwarded-For header is set, it returns the first IP of the header.
// If the X-Real-IP header is set, it returns the value of the header.
// Otherwise, it returns the client IP.
func GetIP(r *http.Request) string {
	// try to get the real IP from the X-Forwarded-For header
	if len(r.Header.Get(XForwardedFor)) > 0 {
		xff := r.Header.Get(XForwardedFor)
		// if the header is set, return the first IP
		ips := strings.Split(xff, ",")
		return ips[0]
	}
	// try to get the real IP from the X-Real-IP header
	if len(r.Header.Get(XRealIP)) > 0 {
		return r.Header.Get(XRealIP)
	}
	// if the header is not set, return the client IP
	clientIP, _, _ /*err*/ := net.SplitHostPort(r.RemoteAddr)
	return clientIP

}

// GetHost returns the host name from the URL, in the form of "host:port".
func GetHost(u *url.URL) string {
	if _, _, err := net.SplitHostPort(u.Host); err == nil {
		return u.Host
	}
	if u.Scheme == "http" {
		return fmt.Sprintf("%s:%d", u.Host, 80)
	}
	if u.Scheme == "https" {
		return fmt.Sprintf("%s:%d", u.Host, 443)
	}
	return u.Host
}
