package plugin

import (
	"fmt"

	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/plugin_config"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/utils"
	"github.com/simonleung8/flags"
	"github.com/simonleung8/flags/flag"
)

type Plugins struct {
	ui     terminal.UI
	config plugin_config.PluginConfiguration
}

func init() {
	command_registry.Register(&Plugins{})
}

func (cmd *Plugins) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["checksum"] = &cliFlags.BoolFlag{Name: "checksum", Usage: i18n.T("Compute and show the sha1 value of the plugin binary file")}

	return command_registry.CommandMetadata{
		Name:        "plugins",
		Description: i18n.T("list all available plugin commands"),
		Usage:       i18n.T("CF_NAME plugins"),
		Flags:       fs,
	}
}

func (cmd *Plugins) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. No argument required\n\n") + command_registry.Commands.CommandUsage("plugins"))
	}

	return
}

func (cmd *Plugins) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.PluginConfig
	return cmd
}

func (cmd *Plugins) Execute(c flags.FlagContext) {
	var version string

	cmd.ui.Say(i18n.T("Listing Installed Plugins..."))

	plugins := cmd.config.Plugins()

	var table terminal.Table
	if c.Bool("checksum") {
		cmd.ui.Say(i18n.T("Computing sha1 for installed plugins, this may take a while ..."))
		table = terminal.NewTable(cmd.ui, []string{i18n.T("Plugin Name"), i18n.T("Version"), i18n.T("Command Name"), "sha1", i18n.T("Command Help")})
	} else {
		table = terminal.NewTable(cmd.ui, []string{i18n.T("Plugin Name"), i18n.T("Version"), i18n.T("Command Name"), i18n.T("Command Help")})
	}

	for pluginName, metadata := range plugins {
		if metadata.Version.Major == 0 && metadata.Version.Minor == 0 && metadata.Version.Build == 0 {
			version = "N/A"
		} else {
			version = fmt.Sprintf("%d.%d.%d", metadata.Version.Major, metadata.Version.Minor, metadata.Version.Build)
		}

		for _, command := range metadata.Commands {
			args := []string{pluginName, version}

			if command.Alias != "" {
				args = append(args, command.Name+", "+command.Alias)
			} else {
				args = append(args, command.Name)
			}

			if c.Bool("checksum") {
				checksum := utils.NewSha1Checksum(metadata.Location)
				sha1, err := checksum.ComputeFileSha1()
				if err != nil {
					args = append(args, "n/a")
				} else {
					args = append(args, fmt.Sprintf("%x", sha1))
				}
			}

			args = append(args, command.HelpText)
			table.Add(args...)
		}
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	table.Print()
}
