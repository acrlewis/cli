package featureflag

import (
	"github.com/cloudfoundry/cli/cf/api/feature_flags"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
)

type ListFeatureFlags struct {
	ui       terminal.UI
	config   core_config.ReadWriter
	flagRepo feature_flags.FeatureFlagRepository
}

func init() {
	command_registry.Register(&ListFeatureFlags{})
}

func (cmd *ListFeatureFlags) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "feature-flags",
		Description: i18n.T("Retrieve list of feature flags with status of each flag-able feature"),
		Usage:       i18n.T("CF_NAME feature-flags"),
	}
}

func (cmd *ListFeatureFlags) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. No argument required\n\n") + command_registry.Commands.CommandUsage("feature-flags"))
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return reqs, err
}

func (cmd *ListFeatureFlags) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.flagRepo = deps.RepoLocator.GetFeatureFlagRepository()
	return cmd
}

func (cmd *ListFeatureFlags) Execute(c flags.FlagContext) {
	cmd.ui.Say(i18n.T("Retrieving status of all flagged features as {{.Username}}...", map[string]interface{}{
		"Username": terminal.EntityNameColor(cmd.config.Username())}))

	flags, err := cmd.flagRepo.List()
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{i18n.T("Features"), i18n.T("State")})

	for _, flag := range flags {
		table.Add(
			flag.Name,
			cmd.flagBoolToString(flag.Enabled),
		)
	}

	table.Print()
	return
}

func (cmd ListFeatureFlags) flagBoolToString(enabled bool) string {
	if enabled {
		return "enabled"
	} else {
		return "disabled"
	}
}
