package netbox

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/wishlist"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/jon4hz/wishbox/config"
	"github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/dcim"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/netbox-community/go-netbox/netbox/client/virtualization"
)

const sshServiceName = "ssh"

type Client struct {
	c   *client.NetBoxAPI
	cfg *config.Netbox
}

// GetInventory generates endpoints based on the netbox inventory
func GetInventory(cfg *config.Netbox) ([]*wishlist.Endpoint, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.IgnoreTLS}},
	}
	transport := httptransport.NewWithClient(cfg.Host, client.DefaultBasePath, []string{"https"}, httpClient)
	transport.DefaultAuthentication = httptransport.APIKeyAuth("Authorization", "header", "Token "+cfg.Token)

	c := Client{
		c:   client.New(transport, nil),
		cfg: cfg,
	}

	vDevices, err := c.getVirtualDevices()
	if err != nil {
		return nil, err
	}

	pDevices, err := c.getDevices()
	if err != nil {
		return nil, err
	}
	vDevices = append(vDevices, pDevices...)
	return vDevices, nil
}

func (c Client) getVirtualDevices() ([]*wishlist.Endpoint, error) {
	req := virtualization.NewVirtualizationVirtualMachinesListParams()
	req.Role = &c.cfg.FilterRole
	if c.cfg.OnlyActive {
		req.Status = &statusActive
	}

	res, err := c.c.Virtualization.VirtualizationVirtualMachinesList(req, nil)
	if err != nil {
		return nil, err
	}

	endpoints := make([]*wishlist.Endpoint, 0, *res.GetPayload().Count)
	for _, v := range res.GetPayload().Results {

		port, err := c.getSSHService(strconv.Itoa(int(v.ID)), sshServiceName, true)
		if err != nil {
			return nil, err
		}
		if port == 0 {
			port = 22
		}

		ip := strings.Split(*v.PrimaryIP.Address, "/")[0]
		endpoints = append(endpoints, &wishlist.Endpoint{
			Name:         *v.Name,
			Address:      fmt.Sprintf("%s:%d", ip, port),
			ForwardAgent: c.cfg.ForwardAgent,
			User:         c.cfg.User,
		})
	}

	return endpoints, nil
}

var statusActive = "active"

func (c Client) getDevices() ([]*wishlist.Endpoint, error) {
	req := dcim.NewDcimDevicesListParams()
	req.Role = &c.cfg.FilterRole
	if c.cfg.OnlyActive {
		req.Status = &statusActive
	}

	res, err := c.c.Dcim.DcimDevicesList(req, nil)
	if err != nil {
		return nil, err
	}

	endpoints := make([]*wishlist.Endpoint, 0, *res.GetPayload().Count)
	for _, v := range res.GetPayload().Results {

		port, err := c.getSSHService(strconv.Itoa(int(v.ID)), sshServiceName, true)
		if err != nil {
			return nil, err
		}
		if port == 0 {
			port = 22
		}

		ip := strings.Split(*v.PrimaryIP.Address, "/")[0]
		endpoints = append(endpoints, &wishlist.Endpoint{
			Name:         *v.Name,
			Address:      fmt.Sprintf("%s:%d", ip, port),
			ForwardAgent: c.cfg.ForwardAgent,
			User:         c.cfg.User,
		})
	}

	return endpoints, nil
}

func (c Client) getSSHService(deviceID string, name string, isVirtual bool) (int64, error) {
	req := ipam.NewIpamServicesListParams()

	if isVirtual {
		req.VirtualMachineID = &deviceID
	} else {
		req.DeviceID = &deviceID
	}

	req.Name = &name
	res, err := c.c.Ipam.IpamServicesList(req, nil)
	if err != nil {
		return 0, err
	}
	for _, v := range res.GetPayload().Results {
		return v.Ports[0], nil
	}
	return 0, nil
}
