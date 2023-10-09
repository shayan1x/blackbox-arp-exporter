package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mdlayher/arp"
)

const IP_REGEX = "^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$"

func GenerateMetrics(Duration float64, Success int) string {
	return fmt.Sprintf(`
# HELP arp_duration_seconds Returns how long the arp request took to complete in seconds
# TYPE arp_duration_seconds gauge
arp_duration_seconds %f
# HELP arp_success Displays whether or not the arp request was a success
# TYPE arp_success gauge
arp_success %d
`, Duration, Success)
}

func main() {
	Address := flag.String("listen-address", ":9230", "The Address used to listen the exporter")
	InterfaceName := flag.String("interface", "eth0", "The interface name to send packets from")
	flag.Parse()

	fmt.Printf("[Info] Listening on %s", *Address)

	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		Query := r.URL.Query()

		if !Query.Has("target") {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Target parameter is missing")
			return
		}

		if !Query.Has("module") {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Module parameter is missing (currently the only supported module is arp)")
			return
		}

		Timeout := 1000

		if Query.Has("timeout") {
			Result, err := strconv.Atoi(Query.Get("timeout"))
			if Result >= 1000 && err == nil {
				Timeout = Result
			}
		}

		Target := Query.Get("target")

		IsIP, _ := regexp.Match(IP_REGEX, []byte(Target))

		if !IsIP {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Target parameter should be an valid IPv4 address like (192.168.10.100)")
			return
		}

		if strings.ToLower(Query.Get("module")) != "arp" {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Invalid module parameter (Set this value to arp)")
			return
		}

		inf, err := net.InterfaceByName(*InterfaceName)

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "[Error] %s", err.Error())
			return
		}

		client, err := arp.Dial(inf)

		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "[Error] %s", err.Error())
			return
		}

		client.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(Timeout)))

		start := time.Now()
		result, err := client.Resolve(netip.MustParseAddr(Target))
		duration := time.Since(start).Seconds()
		Success := 0

		if len(result) > 0 {
			Success = 1
		} else if err != nil {
			Success = 0
		}

		w.WriteHeader(200)
		fmt.Fprint(w, GenerateMetrics(duration, Success))
	})

	err := http.ListenAndServe(*Address, nil)

	if err != nil {
		fmt.Printf("[Error] Could not listen on %s", *Address)
		os.Exit(1)
	}
}
