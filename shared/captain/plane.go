package captain

import (
	"errors"
	"fmt"
	"github.com/ARMmaster17/Captain/shared/ampq"
	"github.com/ARMmaster17/Captain/shared/ipam"
	"github.com/ARMmaster17/Captain/shared/proxmox"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const MaxCpu = 8192
const MaxRam = 3072000
const MaxStorage = 8000

type Plane struct {
	Name	string	`yaml:"name"`
	CPU		int		`yaml:"cpu"`
	RAM		int		`yaml:"ram"`
	Storage	int		`yaml:"storage"`
}

func NewPlane(name string, cpu int, ram int, storage int) (Plane, error) {
	// TODO: Validate name for invalid characters
	if cpu <= 0 || cpu > MaxCpu {
		return Plane{}, errors.New("invalid CPU parameter " + string(cpu))
	}
	if ram <= 0 || ram > MaxRam {
		return Plane{}, errors.New("invalid RAM parameter " + string(ram))
	}
	if storage <= 0 || storage > MaxStorage {
		return Plane{}, errors.New("invalid Storage parameter " + string(storage))
	}
	return Plane{
		Name: name,
		CPU: cpu,
		RAM: ram,
		Storage: storage,
	}, nil
}

func (p *Plane) Create() (string, error) {
	machineConfig, err := buildPlaneConfig(plane)
	if err != nil {
		log.Println(err)
		return "", errors.New("an error occurred while building the plane configuration")
	}
	proxmoxAPI, err := proxmox.NewProxmox()
	lxc, err := proxmoxAPI.LXCCreate(machineConfig)
	return lxc.VMID, nil
}

func (p *Plane) Destroy() error {
	ipamAPI, err := ipam.NewIPAM()
	if err != nil {
		return errors.New("unable to contact IPAM API")
	}
	hostname, err := p.getFQDNHostname()
	if err != nil {
		log.Println(err)
		return errors.New("unable to build FQDN")
	}
	h := ipam.Hostname(hostname)
	err = h.Delete(ipamAPI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to release IP address")
	}
	proxmoxAPI, err := proxmox.NewProxmox()
	if err != nil {
		log.Println(err)
		return errors.New("uanble to contact Proxmox API")
	}
	lxc, err := proxmoxAPI.GetLXCFromHostname(h)
	lxc.Stop(proxmoxAPI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to stop container")
	}
	time.Sleep(30 * time.Second)
	err = lxc.Destroy(proxmoxAPI)
	if err != nil {
		log.Println(err)
		return errors.New("unable to destroy container")
	}
	return nil
}

func (p *Plane) getFQDNHostname() (ipam.Hostname, error) {
	allPlaneConfig, err := getAllPlaneConfig()
	if err != nil {
		log.Println(err)
		return "", errors.New("unable to build cluster-wide plane configuration")
	}
	return ipam.Hostname(fmt.Sprintf("%s.%s", p.Name, allPlaneConfig.Domain)), nil
}
