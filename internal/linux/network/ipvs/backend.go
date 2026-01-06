package ipvs

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

type IPVSBackend struct{}

func NewIPVSBackend() *IPVSBackend {
	return &IPVSBackend{}
}

type IPVSService struct {
	Protocol  string
	VirtualIP string
	Port      int
	Scheduler string
	Backends  []IPVSBackend
}
type IPVSServiceBE struct {
	RealIP string
	Port   int
	Mode   string // NAT / Route / Tunnel
	Weight int
}

var serviceRe = regexp.MustCompile(`^(TCP|UDP)\s+([\d\.]+):(\d+)\s+(\w+)`)
var backendRe = regexp.MustCompile(`^\s+->\s+([\d\.]+):(\d+)\s+(\w+)\s+(\d+)\s+(\d+)\s+(\d+)`)

func (b *IPVSBackend) Fetch() (string, []IPVSService) {
	services := []IPVSService{}
	bytes, _ := os.ReadFile("/proc/net/ip_vs")
	raw := string(bytes)
	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		if m := serviceRe.FindStringSubmatch(line); m != nil {
			p, _ := strconv.Atoi(m[3])
			svc := IPVSService{
				Protocol:  m[1],
				VirtualIP: m[2],
				Port:      p,
				Scheduler: m[4],
			}
			services = append(services, svc)
		}
	}

	return raw, []IPVSService{}
}
