package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func middleware(next http.Handler, allowedIPs []string) http.Handler {
	allowedIPsmap := make(map[string]bool)
	for _, ip := range allowedIPs {
		allowedIPsmap[ip] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		ip = r.Header.Get("X-Real-IP")
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}

		ip = strings.TrimSpace(ip)
		fmt.Printf("Checking IP: %s\n", ip)
		fmt.Println("Allowed IPs:", allowedIPsmap)

		allowed := false
		for _, prefix := range allowedIPs {
			if strings.HasPrefix(ip, prefix) {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "Not Allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	allowedIPs := []string{
		"39.32.", "110.36.", "119.152.", "182.176.",
	}
	http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(w, "You're allowed!")
	}), allowedIPs))

	fmt.Println("Server listening on :8080...")

	http.ListenAndServe(":8080", nil)

}
