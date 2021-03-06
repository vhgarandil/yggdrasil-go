package yggdrasil

import water "github.com/songgao/water"
import "os/exec"
import "strings"
import "fmt"

// This is to catch Windows platforms

func (tun *tunDevice) setup(ifname string, iftapmode bool, addr string, mtu int) error {
	if !iftapmode {
		tun.core.log.Printf("TUN mode is not supported on this platform, defaulting to TAP")
	}
	config := water.Config{DeviceType: water.TAP}
	config.PlatformSpecificParams.ComponentID = "tap0901"
	config.PlatformSpecificParams.Network = "169.254.0.1/32"
	iface, err := water.New(config)
	if err != nil {
		panic(err)
	}
	// Disable/enable the interface to resets its configuration (invalidating iface)
	cmd := exec.Command("netsh", "interface", "set", "interface", iface.Name(), "admin=DISABLED")
	tun.core.log.Printf("netsh command: %v", strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		tun.core.log.Printf("Windows netsh failed: %v.", err)
		tun.core.log.Println(string(output))
		return err
	}
	cmd = exec.Command("netsh", "interface", "set", "interface", iface.Name(), "admin=ENABLED")
	tun.core.log.Printf("netsh command: %v", strings.Join(cmd.Args, " "))
	output, err = cmd.CombinedOutput()
	if err != nil {
		tun.core.log.Printf("Windows netsh failed: %v.", err)
		tun.core.log.Println(string(output))
		return err
	}
	// Get a new iface
	iface, err = water.New(config)
	if err != nil {
		panic(err)
	}
	tun.iface = iface
	tun.mtu = mtu
	err = tun.setupMTU(tun.mtu)
	if err != nil {
		panic(err)
	}
	return tun.setupAddress(addr)
}

func (tun *tunDevice) setupMTU(mtu int) error {
	// Set MTU
	cmd := exec.Command("netsh", "interface", "ipv6", "set", "subinterface",
		fmt.Sprintf("interface=%s", tun.iface.Name()),
		fmt.Sprintf("mtu=%d", mtu),
		"store=active")
	tun.core.log.Printf("netsh command: %v", strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		tun.core.log.Printf("Windows netsh failed: %v.", err)
		tun.core.log.Println(string(output))
		return err
	}
	return nil
}

func (tun *tunDevice) setupAddress(addr string) error {
	// Set address
	cmd := exec.Command("netsh", "interface", "ipv6", "add", "address",
		fmt.Sprintf("interface=%s", tun.iface.Name()),
		fmt.Sprintf("addr=%s", addr),
		"store=active")
	tun.core.log.Printf("netsh command: %v", strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		tun.core.log.Printf("Windows netsh failed: %v.", err)
		tun.core.log.Println(string(output))
		return err
	}
	return nil
}
