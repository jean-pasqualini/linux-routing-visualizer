//go:build IPVS_BINARY

package ipvs

import (
	"bytes"
	"errors"
	"os/exec"
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
	Backends  []IPVSServiceBE
}
type IPVSServiceBE struct {
	RealIP string
	Port   int
	Mode   string // NAT / Route / Tunnel
	Weight int
}

var serviceRe = regexp.MustCompile(`^(TCP|UDP)\s+([\d\.]+):(\d+)\s+(\w+)`)
var backendRe = regexp.MustCompile(`^\s+->\s+([\d\.]+):(\d+)\s+(\w+)\s+(\d+)\s+(\d+)\s+(\d+)`)

// bytes, _ := os.ReadFile("/proc/net/ip_vs")
// /proc/net/ip_vs
// ipvsadm ioctl
// Netlink IPVS
// ipvsctl
func (b *IPVSBackend) Fetch() (string, []IPVSService, error) {
	services := []IPVSService{}
	raw, err := b.runProces()
	if err != nil {
		return "", []IPVSService{}, err
	}
	var current *IPVSService
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
			current = &services[len(services)-1]
		}
		if m := backendRe.FindStringSubmatch(line); m != nil {
			p, _ := strconv.Atoi(m[2])
			w, _ := strconv.Atoi(m[4])
			back := IPVSServiceBE{
				RealIP: m[1],
				Port:   p,
				Mode:   m[3],
				Weight: w,
			}
			current.Backends = append(current.Backends, back)
		}
	}

	return raw, services, nil
}

func (b *IPVSBackend) runProces() (string, error) {
	// -c add the counters , "-c"
	cmd := exec.Command("ipvsadm", "-Ln")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	if errRun := cmd.Run(); errRun != nil {
		return "", errors.New(err.String() + errRun.Error())
	}

	return out.String(), nil
}
