package sysctl

import (
	"github.com/prometheus/procfs"
	"strconv"
)

type SysctlBackend struct {
}

func NewSysctlBackend() *SysctlBackend {
	return &SysctlBackend{}
}

type SysctlConfig struct {
	IPV4 SysctlNetConfig
}

type SysctlNetConfig struct {
	Forward bool
}

func (b *SysctlBackend) Fetch() (SysctlConfig, error) {
	config := SysctlConfig{}
	fs, _ := procfs.NewFS("/proc")
	data, err := fs.SysctlStrings("net.ipv4.ip_forward")
	if err != nil {
		return config, err
	}

	fw, err := strconv.ParseBool(data[0])
	if err != nil {
		return config, err
	}

	config.IPV4.Forward = fw

	return config, nil
}
