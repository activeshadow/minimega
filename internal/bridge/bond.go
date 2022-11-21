// Copyright 2016-2022 National Technology & Engineering Solutions of Sandia, LLC (NTESS).
// Under the terms of Contract DE-NA0003525 with NTESS, the U.S. Government retains certain
// rights in this software.

package bridge

import (
	"fmt"

	log "github.com/sandia-minimega/minimega/v2/pkg/minilog"
)

func (b *Bridge) AddBond(name, mode string, interfaces []string, vlan int) error {
	bridgeLock.Lock()
	defer bridgeLock.Unlock()

	log.Info("adding bond %v on bridge %v", name, b.Name)

	// need to remove the taps from ovs first before we can bond
	for _, iface := range interfaces {
		args := []string{"del-port", b.Name, iface}

		if _, err := ovsCmdWrapper(args); err != nil {
			return fmt.Errorf("failed to delete tap %v from ovs with error: %v", iface, err)
		}
	}

	// ovs-vsctl add-bond <bridge name> <bond name> <list of interfaces>
	args := []string{"add-bond", b.Name, name}

	args = append(args, interfaces...)
	args = append(args, "--", "set", "port", name)
	args = append(args, fmt.Sprintf("tag=%v", vlan))

	// https://www.man7.org/linux/man-pages/man5/ovs-vswitchd.conf.db.5.html#Port_TABLE
	switch mode {
	case "active-backup", "balance-slb":
		args = append(args, fmt.Sprintf("bond_mode=%s", mode))
	case "balance-tcp":
		// https://vstellar.com/2019/01/ahv-networking-part-3-change-ovs-bond-mode/
		args = append(args, "lacp=active", "other_config:lacp-fallback-ab=true")
		args = append(args, fmt.Sprintf("bond_mode=%s", mode))
	default:
		return fmt.Errorf("unsupported bond mode provided: %s", mode)
	}

	if _, err := ovsCmdWrapper(args); err != nil {
		return fmt.Errorf("add bond failed: %v", err)
	}

	b.bonds[name] = interfaces
	return nil
}
