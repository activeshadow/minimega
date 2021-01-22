package app

import (
	"fmt"
	"os"
	"strings"

	"phenix/internal/common"
	"phenix/tmpl"
	"phenix/types"
)

type NTP struct {
	options Options
}

func (this *NTP) Init(opts ...Option) error {
	this.options = NewOptions(opts...)

	return nil
}

func (NTP) Name() string {
	return "ntp"
}

func (NTP) Configure(exp *types.Experiment) error {
	return nil
}

func (this NTP) PreStart(exp *types.Experiment) error {
	ntpServers := exp.Spec.Topology().FindNodesWithLabels("ntp-server")

	if len(ntpServers) != 0 {
		// Just take first server if more than one are labeled.
		node := ntpServers[0]

		var ntpAddr string

		for _, iface := range node.Network().Interfaces() {
			if strings.EqualFold(iface.VLAN(), "mgmt") {
				ntpAddr = iface.Address()
				break
			}
		}

		ntpDir := exp.Spec.BaseDir() + "/ntp"
		ntpFile := ntpDir + "/" + node.General().Hostname() + "_ntp"

		if node.Type() == "Router" {
			if err := os.MkdirAll(ntpDir, 0755); err != nil {
				return fmt.Errorf("creating experiment ntp directory path: %w", err)
			}

			node.AddInject(ntpFile, "/opt/vyatta/etc/ntp.conf", "", "")

			if err := tmpl.CreateFileFromTemplate("ntp_linux.tmpl", ntpAddr, ntpFile); err != nil {
				return fmt.Errorf("generating ntp script: %w", err)
			}
		} else if node.Hardware().OSType() == "linux" {
			if this.options.UseC2 {
				filePath := common.PhenixBase + "/images/" + exp.Spec.ExperimentName() + "/" + node.General().Hostname()

				os.MkdirAll(filePath, 0755)

				if err := tmpl.CreateFileFromTemplate("ntp_linux.tmpl", ntpAddr, filePath+"/ntp.conf"); err != nil {
					return fmt.Errorf("generating ntp script: %w", err)
				}

				command := fmt.Sprintf("send %s/%s/*", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)

				command = fmt.Sprintf("exec cp /tmp/miniccc/files/%s/%s/ntp.conf /etc/ntp.conf", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)
			} else {
				if err := os.MkdirAll(ntpDir, 0755); err != nil {
					return fmt.Errorf("creating experiment ntp directory path: %w", err)
				}

				node.AddInject(ntpFile, "/etc/ntp.conf", "", "")

				if err := tmpl.CreateFileFromTemplate("ntp_linux.tmpl", ntpAddr, ntpFile); err != nil {
					return fmt.Errorf("generating ntp script: %w", err)
				}
			}
		} else if node.Hardware().OSType() == "windows" {
			if this.options.UseC2 {
				filePath := common.PhenixBase + "/images/" + exp.Spec.ExperimentName() + "/" + node.General().Hostname()

				os.MkdirAll(filePath, 0755)

				if err := tmpl.CreateFileFromTemplate("ntp_windows.tmpl", ntpAddr, filePath+"/ntp.ps1"); err != nil {
					return fmt.Errorf("generating ntp script: %w", err)
				}

				command := fmt.Sprintf("send %s/%s/*", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)

				command = fmt.Sprintf("exec powershell.exe -file /tmp/miniccc/files/%s/%s/ntp.ps1", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)
			} else {
				if err := os.MkdirAll(ntpDir, 0755); err != nil {
					return fmt.Errorf("creating experiment ntp directory path: %w", err)
				}

				node.AddInject(ntpFile, "ntp.ps1", "0755", "")

				if err := tmpl.CreateFileFromTemplate("ntp_windows.tmpl", ntpAddr, ntpFile); err != nil {
					return fmt.Errorf("generating ntp script: %w", err)
				}
			}
		}
	}

	return nil
}

func (NTP) PostStart(exp *types.Experiment) error {
	return nil
}

func (NTP) Cleanup(exp *types.Experiment) error {
	return nil
}
