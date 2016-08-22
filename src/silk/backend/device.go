package backend

import (
	"fmt"
	"net"
	"os"
	"syscall"

	"github.com/vishvananda/netlink"
)

type vxlanDeviceAttrs struct {
	vni       int
	name      string
	vtepIndex int
	vtepAddr  net.IP
	vtepPort  int
}

type vxlanDevice struct {
	link *netlink.Vxlan
}

func sysctlSet(path, value string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(value))
	return err
}

func newVXLANDevice(devAttrs *vxlanDeviceAttrs) (*vxlanDevice, error) {
	link := &netlink.Vxlan{
		LinkAttrs: netlink.LinkAttrs{
			Name: devAttrs.name,
		},
		VxlanId:      devAttrs.vni,
		VtepDevIndex: devAttrs.vtepIndex,
		SrcAddr:      devAttrs.vtepAddr,
		Port:         devAttrs.vtepPort,
		Learning:     false,
		GBP:          true,
	}

	link, err := ensureLink(link)
	if err != nil {
		return nil, err
	}
	// this enables ARP requests being sent to userspace via netlink
	sysctlPath := fmt.Sprintf("/proc/sys/net/ipv4/neigh/%s/app_solicit", devAttrs.name)
	if err := sysctlSet(sysctlPath, "3"); err != nil {
		return nil, err
	}

	return &vxlanDevice{
		link: link,
	}, nil
}

func ensureLink(vxlan *netlink.Vxlan) (*netlink.Vxlan, error) {
	err := netlink.LinkAdd(vxlan)
	if err == syscall.EEXIST {
		// it's ok if the device already exists as long as config is similar
		existing, err := netlink.LinkByName(vxlan.Name)
		if err != nil {
			return nil, err
		}

		incompat := vxlanLinksIncompat(vxlan, existing)
		if incompat == "" {
			return existing.(*netlink.Vxlan), nil
		}

		// delete existing
		fmt.Printf("%q already exists with incompatable configuration: %v; recreating device", vxlan.Name, incompat)
		if err = netlink.LinkDel(existing); err != nil {
			return nil, fmt.Errorf("failed to delete interface: %v", err)
		}

		// create new
		if err = netlink.LinkAdd(vxlan); err != nil {
			return nil, fmt.Errorf("failed to create vxlan interface: %v", err)
		}
	} else if err != nil {
		return nil, err
	}

	ifindex := vxlan.Index
	link, err := netlink.LinkByIndex(vxlan.Index)
	if err != nil {
		return nil, fmt.Errorf("can't locate created vxlan device with index %v", ifindex)
	}

	var ok bool
	if vxlan, ok = link.(*netlink.Vxlan); !ok {
		return nil, fmt.Errorf("created vxlan device with index %v is not vxlan", ifindex)
	}

	return vxlan, nil
}

func vxlanLinksIncompat(l1, l2 netlink.Link) string {
	if l1.Type() != l2.Type() {
		return fmt.Sprintf("link type: %v vs %v", l1.Type(), l2.Type())
	}

	v1 := l1.(*netlink.Vxlan)
	v2 := l2.(*netlink.Vxlan)

	if v1.VxlanId != v2.VxlanId {
		return fmt.Sprintf("vni: %v vs %v", v1.VxlanId, v2.VxlanId)
	}

	if v1.VtepDevIndex > 0 && v2.VtepDevIndex > 0 && v1.VtepDevIndex != v2.VtepDevIndex {
		return fmt.Sprintf("vtep (external) interface: %v vs %v", v1.VtepDevIndex, v2.VtepDevIndex)
	}

	if len(v1.SrcAddr) > 0 && len(v2.SrcAddr) > 0 && !v1.SrcAddr.Equal(v2.SrcAddr) {
		return fmt.Sprintf("vtep (external) IP: %v vs %v", v1.SrcAddr, v2.SrcAddr)
	}

	if len(v1.Group) > 0 && len(v2.Group) > 0 && !v1.Group.Equal(v2.Group) {
		return fmt.Sprintf("group address: %v vs %v", v1.Group, v2.Group)
	}

	if v1.L2miss != v2.L2miss {
		return fmt.Sprintf("l2miss: %v vs %v", v1.L2miss, v2.L2miss)
	}

	if v1.Port > 0 && v2.Port > 0 && v1.Port != v2.Port {
		return fmt.Sprintf("port: %v vs %v", v1.Port, v2.Port)
	}

	if v1.GBP != v2.GBP {
		return fmt.Sprintf("gbp: %v vs %v", v1.GBP, v2.GBP)
	}

	return ""
}

func (dev *vxlanDevice) Configure(vtepOverlayIP net.IP, fullOverlayMask net.IPMask, hwAddr net.HardwareAddr) error {
	// vxlan's subnet is that of the whole overlay network (e.g. /16)
	// and not that of the individual host (e.g. /24)
	addr := &net.IPNet{
		IP:   vtepOverlayIP,
		Mask: fullOverlayMask,
	}
	dev.setAddr4(addr)
	netlink.LinkSetHardwareAddr(dev.link, hwAddr)

	if err := netlink.LinkSetUp(dev.link); err != nil {
		return fmt.Errorf("failed to set interface %s to UP state: %s", dev.link.Attrs().Name, err)
	}

	if err := dev.PurgeAllRoutes(); err != nil {
		return err
	}

	// fully mask the address before adding the wide route
	routeSubnet := &net.IPNet{
		IP:   addr.IP.Mask(addr.Mask),
		Mask: addr.Mask,
	}
	return dev.AddRoute(routeSubnet, nil, netlink.SCOPE_LINK, nil)
}

// sets IP4 addr on link removing any existing ones first
func (dev *vxlanDevice) setAddr4(ipn *net.IPNet) error {
	err := dev.PurgeAllRoutes()
	if err != nil {
		return err
	}

	addr := netlink.Addr{IPNet: ipn, Label: ""}
	if err = netlink.AddrAdd(dev.link, &addr); err != nil {
		return fmt.Errorf("failed to add IP address %s to %s: %s", ipn.String(), dev.link.Attrs().Name, err)
	}

	return nil
}
