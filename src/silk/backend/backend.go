package backend

import (
	"fmt"
	"net"
	"silk/models"
)

type Controller interface {
	OverlayMTU() int
	ConfigureDevice(fullOverlay net.IPNet, local *models.NetHost) error
	InstallRoutes(remoteHosts []*models.NetHost) error
}

type controller struct {
	ExternalInterface net.Interface
	ExternalIP        net.IP
	VNI               int
	Port              int
	Device            *vxlanDevice
}

func (c *controller) OverlayMTU() int {
	return c.ExternalInterface.MTU - 48
}

func New(vni int, port int, externalIP net.IP, externalInterface net.Interface) (Controller, error) {
	cont := &controller{
		ExternalInterface: externalInterface,
		ExternalIP:        externalIP,
		VNI:               vni,
		Port:              port,
	}

	devAttrs := vxlanDeviceAttrs{
		vni:       vni,
		name:      fmt.Sprintf("silk.%v", vni),
		vtepIndex: externalInterface.Index,
		vtepAddr:  externalIP,
		vtepPort:  port,
		gbp:       true,
	}

	dev, err := newVXLANDevice(&devAttrs)
	if err != nil {
		return nil, err
	}

	cont.Device = dev
	return cont, nil
}

func (c *controller) ConfigureDevice(fullOverlay net.IPNet, local *models.NetHost) error {
	if err := c.Device.Configure(
		local.VtepOverlayIP,
		fullOverlay.Mask,
		local.VtepOverlayMAC); err != nil {
		return fmt.Errorf("configuring vxlan device: %s", err)
	}

	return nil
}

func (c *controller) InstallRoutes(remoteHosts []*models.NetHost) error {
	for _, host := range remoteHosts {
		l3neigh := L3Neigh{
			OverlayIP:  host.VtepOverlayIP,
			OverlayMAC: host.VtepOverlayMAC,
		}
		l2neigh := L2Neigh{
			OverlayMAC: host.VtepOverlayMAC,
			UnderlayIP: host.PublicIP,
		}
		err := c.Device.SetL3(l3neigh)
		if err != nil {
			return fmt.Errorf("set l3: %s", err)
		}
		err = c.Device.SetL2(l2neigh)
		if err != nil {
			return fmt.Errorf("set l2: %s", err)
		}
		err = c.Device.AddRoute(host.OverlaySubnet, host.OverlaySubnet.IP)
		if err != nil {
			return fmt.Errorf("add route %s: %s", host.OverlaySubnet.String(), err)
		}
	}
	return nil
}
