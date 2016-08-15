package backend

import (
	"fmt"
	"net"
	"syscall"

	"github.com/vishvananda/netlink"
)

type L2Neigh struct {
	OverlayMAC net.HardwareAddr
	UnderlayIP net.IP
}

func (dev *vxlanDevice) SetL2(n L2Neigh) error {
	// bridge fdb replace 0a:00:0a:ff:02:00 dev silk.1 dst 10.244.99.11
	fmt.Printf("calling L2 NeighSet: %s to %s", n.OverlayMAC.String(), n.UnderlayIP.String())
	return netlink.NeighSet(&netlink.Neigh{
		LinkIndex:    dev.link.Index,
		State:        netlink.NUD_PERMANENT,
		Family:       syscall.AF_BRIDGE,
		Flags:        netlink.NTF_SELF,
		IP:           n.UnderlayIP,
		HardwareAddr: n.OverlayMAC,
	})
}

type L3Neigh struct {
	OverlayMAC net.HardwareAddr
	OverlayIP  net.IP
}

func (dev *vxlanDevice) SetL3(n L3Neigh) error {
	// ip neigh replace to 10.255.1.0 dev silk.1 lladdr 0a:00:0a:ff:01:00
	fmt.Printf("calling L3 NeighSet: %s to %s", n.OverlayIP.String(), n.OverlayMAC.String())
	return netlink.NeighSet(&netlink.Neigh{
		LinkIndex:    dev.link.Index,
		State:        netlink.NUD_PERMANENT,
		Type:         syscall.RTN_UNICAST,
		IP:           n.OverlayIP,
		HardwareAddr: n.OverlayMAC,
	})
}

func (dev *vxlanDevice) AddRoute(destNet net.IPNet, gateway net.IP) error {
	// ip route add 10.255.2.0/24 via 10.255.2.0 dev silk.1
	return netlink.RouteAdd(&netlink.Route{
		LinkIndex: dev.link.Index,
		Scope:     netlink.SCOPE_UNIVERSE,
		Dst:       &destNet,
		Gw:        gateway,
	})
}
