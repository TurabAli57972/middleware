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

		if !allowedIPsmap[ip] {
			http.Error(w, "Not Allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	allowedIPs := []string{
		"39.32.0.0", "110.36.0.0", "119.152.0.0", "182.176.0.0",
	}
	http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println(w, "You're allowed!")
	}), allowedIPs))

	fmt.Println("Server listening on :8080...")

	http.ListenAndServe(":8080", nil)

}
