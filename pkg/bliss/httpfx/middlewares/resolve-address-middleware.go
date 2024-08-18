package middlewares

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/eser/go-service/pkg/bliss/httpfx"
)

const (
	ClientAddr       httpfx.ContextKey = "client-addr"
	ClientAddrIp     httpfx.ContextKey = "client-addr-ip"
	ClientAddrOrigin httpfx.ContextKey = "client-addr-origin"
)

func ResolveAddressMiddleware() httpfx.Handler {
	return func(ctx *httpfx.Context) httpfx.Response {
		addrs := GetClientAddrs(ctx.Request)

		newContext := context.WithValue(
			ctx.Request.Context(),
			ClientAddr,
			addrs,
		)

		isLocal, err := DetectLocalNetwork(addrs[0])
		if err != nil {
			return ctx.Results.Error(http.StatusInternalServerError, err.Error())
		}

		if isLocal {
			newContext = context.WithValue(
				newContext,
				ClientAddrOrigin,
				"local",
			)

			ctx.ResponseWriter.Header().
				Set("X-Request-Origin", "local: "+strings.Join(addrs, ", "))

			ctx.UpdateContext(newContext)

			return ctx.Next()
		}

		// TODO(@eser) add ip allowlist and blocklist implementations

		newContext = context.WithValue(
			newContext,
			ClientAddrOrigin,
			"remote",
		)

		ctx.ResponseWriter.Header().
			Set("X-Request-Origin", strings.Join(addrs, ", "))

		ctx.UpdateContext(newContext)

		return ctx.Next()
	}
}

func DetectLocalNetwork(requestAddr string) (bool, error) {
	requestIp, _, err := net.SplitHostPort(requestAddr)
	if err != nil {
		return false, err
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, err
	}

	requestIpNet := net.ParseIP(requestIp)

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if !ipNet.Contains(requestIpNet) {
			continue
		}

		if requestIpNet.IsLoopback() {
			return true, nil
		}
	}

	return false, nil
}

func GetClientAddrs(req *http.Request) []string {
	// first check the X-Forwarded-For header
	requester := req.Header.Get("True-Client-IP")

	if len(requester) == 0 {
		requester = req.Header.Get("X-Forwarded-For")
	}

	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = req.Header.Get("X-Real-IP")
	}

	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = req.RemoteAddr
	}

	// split comma delimited list into a slice
	// (this happens when proxied via elastic load balancer then again through nginx)
	return strings.Split(requester, ",")
}
