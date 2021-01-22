package app

import (
	"fmt"
	"os"
	"strings"

	"phenix/internal/common"
	"phenix/tmpl"
	"phenix/types"
	ifaces "phenix/types/interfaces"
)

type Serial struct {
	options Options
}

func (this *Serial) Init(opts ...Option) error {
	this.options = NewOptions(opts...)

	return nil
}

func (Serial) Name() string {
	return "serial"
}

func (Serial) Configure(exp *types.Experiment) error {
	return nil
}

func (this Serial) PreStart(exp *types.Experiment) error {
	// loop through nodes
	for _, node := range exp.Spec.Topology().Nodes() {
		// We only care about configuring serial interfaces on Linux VMs.
		// TODO: handle rhel and centos OS types.
		if node.Hardware().OSType() != "linux" {
			continue
		}

		var serial []ifaces.NodeNetworkInterface

		// Loop through interface type to see if any of the interfaces are serial.
		for _, iface := range node.Network().Interfaces() {
			if iface.Type() == "serial" {
				serial = append(serial, iface)
			}
		}

		if serial != nil {
			if this.options.UseC2 {
				filePath := common.PhenixBase + "/images/" + exp.Spec.ExperimentName() + "/" + node.General().Hostname()

				os.MkdirAll(filePath, 0755)

				if err := tmpl.CreateFileFromTemplate("serial_startup.tmpl", serial, filePath+"/serial-startup.sh"); err != nil {
					return fmt.Errorf("generating serial startup script: %w", err)
				}

				command := fmt.Sprintf("send %s/%s/*", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)

				command = fmt.Sprintf("exec bash /tmp/miniccc/files/%s/%s/serial-startup.sh", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)
			} else {
				startupDir := exp.Spec.BaseDir() + "/startup"
				serialFile := startupDir + "/" + node.General().Hostname() + "-serial.bash"

				node.AddInject(serialFile, "/etc/phenix/serial-startup.bash", "0755", "")
				node.AddInject(startupDir+"/serial-startup.service", "/etc/systemd/system/serial-startup.service", "", "")
				node.AddInject(startupDir+"/symlinks/serial-startup.service", "/etc/systemd/system/multi-user.target.wants/serial-startup.service", "", "")

				if err := os.MkdirAll(startupDir, 0755); err != nil {
					return fmt.Errorf("creating experiment startup directory path: %w", err)
				}

				if err := tmpl.CreateFileFromTemplate("serial_startup.tmpl", serial, serialFile); err != nil {
					return fmt.Errorf("generating serial script: %w", err)
				}

				if err := tmpl.RestoreAsset(startupDir, "serial-startup.service"); err != nil {
					return fmt.Errorf("restoring serial-startup.service: %w", err)
				}

				symlinksDir := startupDir + "/symlinks"

				if err := os.MkdirAll(symlinksDir, 0755); err != nil {
					return fmt.Errorf("creating experiment startup symlinks directory path: %w", err)
				}

				if err := os.Symlink("../serial-startup.service", symlinksDir+"/serial-startup.service"); err != nil {
					// Ignore the error if it was for the symlinked file already existing.
					if !strings.Contains(err.Error(), "file exists") {
						return fmt.Errorf("creating symlink for serial-startup.service: %w", err)
					}
				}
			}
		}
	}

	return nil
}

func (Serial) PostStart(exp *types.Experiment) error {
	return nil
}

func (Serial) Cleanup(exp *types.Experiment) error {
	return nil
}
