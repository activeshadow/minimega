package app

import (
	"fmt"
	"os"
	"path/filepath"

	"phenix/internal/common"
	"phenix/tmpl"
	"phenix/types"
	ifaces "phenix/types/interfaces"
)

type Startup struct {
	options Options
}

func (this *Startup) Init(opts ...Option) error {
	this.options = NewOptions(opts...)

	return nil
}

func (Startup) Name() string {
	return "startup"
}

func (this *Startup) Configure(exp *types.Experiment) error {
	return nil
}

func (this Startup) PreStart(exp *types.Experiment) error {
	var (
		startupDir = exp.Spec.BaseDir() + "/startup"
		imageDir   = common.PhenixBase + "/images/"
	)

	if err := os.MkdirAll(startupDir, 0755); err != nil {
		return fmt.Errorf("creating experiment startup directory path: %w", err)
	}

	for _, node := range exp.Spec.Topology().Nodes() {
		// Check if user provided an absolute path to image. If not, prepend path
		// with default image path.
		imagePath := node.Hardware().Drives()[0].Image()

		if !filepath.IsAbs(imagePath) {
			imagePath = imageDir + imagePath
		}

		// check if the disk image is present, if not set do not boot to true
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			node.General().SetDoNotBoot(true)
		}

		// if type is router, skip it and continue
		if node.Type() == "Router" {
			continue
		}

		switch node.Hardware().OSType() {
		case "linux":
			timeZone := "Etc/UTC"

			if this.options.UseC2 {
				filePath := common.PhenixBase + "/images/" + exp.Spec.ExperimentName() + "/" + node.General().Hostname()

				os.MkdirAll(filePath, 0755)

				if err := tmpl.CreateFileFromTemplate("linux_interfaces.tmpl", node, filePath+"/interfaces"); err != nil {
					return fmt.Errorf("generating linux interfaces config: %w", err)
				}

				data := map[string]interface{}{
					"namespace": exp.Spec.ExperimentName(),
					"hostname":  node.General().Hostname(),
					"timezone":  timeZone,
					"os":        node.Hardware().OSType(),
				}

				if err := tmpl.CreateFileFromTemplate("linux-init.tmpl", data, filePath+"/init.sh"); err != nil {
					return fmt.Errorf("generating Linux c2 init script: %w", err)
				}

				command := fmt.Sprintf("send %s/%s/*", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)

				command = fmt.Sprintf("exec bash /tmp/miniccc/files/%s/%s/init.sh", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)
			} else {
				var (
					hostnameFile = startupDir + "/" + node.General().Hostname() + "-hostname.sh"
					timezoneFile = startupDir + "/" + node.General().Hostname() + "-timezone.sh"
					ifaceFile    = startupDir + "/" + node.General().Hostname() + "-interfaces"
				)

				node.AddInject(
					hostnameFile,
					"/etc/phenix/startup/1_hostname-start.sh",
					"0755", "",
				)

				node.AddInject(
					timezoneFile,
					"/etc/phenix/startup/2_timezone-start.sh",
					"0755", "",
				)

				node.AddInject(
					ifaceFile,
					"/etc/network/interfaces",
					"", "",
				)

				if err := tmpl.CreateFileFromTemplate("linux_hostname.tmpl", node.General().Hostname(), hostnameFile); err != nil {
					return fmt.Errorf("generating linux hostname config: %w", err)
				}

				if err := tmpl.CreateFileFromTemplate("linux_timezone.tmpl", timeZone, timezoneFile); err != nil {
					return fmt.Errorf("generating linux timezone config: %w", err)
				}

				if err := tmpl.CreateFileFromTemplate("linux_interfaces.tmpl", node, ifaceFile); err != nil {
					return fmt.Errorf("generating linux interfaces config: %w", err)
				}
			}
		case "rhel", "centos":
			timeZone := "Etc/UTC"

			if this.options.UseC2 {
				filePath := common.PhenixBase + "/images/" + exp.Spec.ExperimentName() + "/" + node.General().Hostname()

				os.MkdirAll(filePath, 0755)

				ifaces := make([]string, len(node.Network().Interfaces()))

				for idx := range node.Network().Interfaces() {
					iface := fmt.Sprintf("ifcfg-eth%d", idx)
					dst := filePath + "/" + iface

					ifaces[idx] = iface

					if err := tmpl.CreateFileFromTemplate("linux_interfaces.tmpl", node, dst); err != nil {
						return fmt.Errorf("generating linux interfaces config: %w", err)
					}
				}

				data := map[string]interface{}{
					"namespace": exp.Spec.ExperimentName(),
					"hostname":  node.General().Hostname(),
					"timezone":  timeZone,
					"os":        node.Hardware().OSType(),
					"ifaces":    ifaces,
				}

				if err := tmpl.CreateFileFromTemplate("linux-init.tmpl", data, filePath+"/init.sh"); err != nil {
					return fmt.Errorf("generating Linux c2 init script: %w", err)
				}

				command := fmt.Sprintf("send %s/%s/*", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)

				command = fmt.Sprintf("exec bash /tmp/miniccc/files/%s/%s/init.sh", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)
			} else {
				var (
					hostnameFile = startupDir + "/" + node.General().Hostname() + "-hostname.sh"
					timezoneFile = startupDir + "/" + node.General().Hostname() + "-timezone.sh"
				)

				node.AddInject(
					hostnameFile,
					"/etc/phenix/startup/1_hostname-start.sh",
					"0755", "",
				)

				node.AddInject(
					timezoneFile,
					"/etc/phenix/startup/2_timezone-start.sh",
					"0755", "",
				)

				if err := tmpl.CreateFileFromTemplate("linux_hostname.tmpl", node.General().Hostname(), hostnameFile); err != nil {
					return fmt.Errorf("generating linux hostname config: %w", err)
				}

				if err := tmpl.CreateFileFromTemplate("linux_timezone.tmpl", timeZone, timezoneFile); err != nil {
					return fmt.Errorf("generating linux timezone config: %w", err)
				}

				for idx := range node.Network().Interfaces() {
					ifaceFile := fmt.Sprintf("%s/interfaces-%s-eth%d", startupDir, node.General().Hostname(), idx)

					node.AddInject(
						ifaceFile,
						fmt.Sprintf("/etc/sysconfig/network-scripts/ifcfg-eth%d", idx),
						"", "",
					)

					if err := tmpl.CreateFileFromTemplate("linux_interfaces.tmpl", node, ifaceFile); err != nil {
						return fmt.Errorf("generating linux interfaces config: %w", err)
					}
				}
			}
		case "windows":
			// Temporary struct to send to the Windows Startup template.
			data := struct {
				Node     ifaces.NodeSpec
				Metadata map[string]interface{}
			}{
				Node:     node,
				Metadata: make(map[string]interface{}),
			}

			// Check to see if a scenario exists for this experiment and if it
			// contains a "startup" app. If so, see if this node has a metadata entry
			// in the scenario app configuration.
			for _, app := range exp.Apps() {
				if app.Name() == "startup" {
					for _, host := range app.Hosts() {
						if host.Hostname() == node.General().Hostname() {
							data.Metadata = host.Metadata()
						}
					}
				}
			}

			if this.options.UseC2 {
				filePath := common.PhenixBase + "/images/" + exp.Spec.ExperimentName() + "/" + node.General().Hostname()

				os.MkdirAll(filePath, 0755)

				if err := tmpl.CreateFileFromTemplate("windows_startup.tmpl", data, filePath+"/init.ps1"); err != nil {
					return fmt.Errorf("generating windows startup config: %w", err)
				}

				command := fmt.Sprintf("send %s/%s/*", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)

				command = fmt.Sprintf("exec powershell.exe -file /tmp/miniccc/files/%s/%s/init.ps1", exp.Spec.ExperimentName(), node.General().Hostname())
				node.AddCommand(command)
			} else {
				var (
					startupFile = startupDir + "/" + node.General().Hostname() + "-startup.ps1"
					schedFile   = startupDir + "/startup-scheduler.cmd"
				)

				node.AddInject(
					startupFile,
					"startup.ps1",
					"0755", "",
				)

				node.AddInject(
					schedFile,
					"ProgramData/Microsoft/Windows/Start Menu/Programs/StartUp/startup_scheduler.cmd",
					"0755", "",
				)

				if err := tmpl.CreateFileFromTemplate("windows_startup.tmpl", data, startupFile); err != nil {
					return fmt.Errorf("generating windows startup config: %w", err)
				}

				if err := tmpl.RestoreAsset(startupDir, "startup-scheduler.cmd"); err != nil {
					return fmt.Errorf("restoring windows startup scheduler: %w", err)
				}
			}
		}
	}

	return nil
}

func (Startup) PostStart(exp *types.Experiment) error {
	return nil
}

func (Startup) Cleanup(exp *types.Experiment) error {
	return nil
}
