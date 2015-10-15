package environmentvariablegroup

import (
	"github.com/cloudfoundry/cli/cf/api/environment_variable_groups"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
)

type StagingEnvironmentVariableGroup struct {
	ui                           terminal.UI
	config                       core_config.ReadWriter
	environmentVariableGroupRepo environment_variable_groups.EnvironmentVariableGroupsRepository
}

func init() {
	command_registry.Register(&StagingEnvironmentVariableGroup{})
}

func (cmd *StagingEnvironmentVariableGroup) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "staging-environment-variable-group",
		Description: i18n.T("Retrieve the contents of the staging environment variable group"),
		ShortName:   "sevg",
		Usage:       i18n.T("CF_NAME staging-environment-variable-group"),
	}
}

func (cmd *StagingEnvironmentVariableGroup) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) ([]requirements.Requirement, error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. No argument required\n\n") + command_registry.Commands.CommandUsage("staging-environment-variable-group"))
	}

	reqs := []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return reqs, nil
}

func (cmd *StagingEnvironmentVariableGroup) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.environmentVariableGroupRepo = deps.RepoLocator.GetEnvironmentVariableGroupsRepository()
	return cmd
}

func (cmd *StagingEnvironmentVariableGroup) Execute(c flags.FlagContext) {
	cmd.ui.Say(i18n.T("Retrieving the contents of the staging environment variable group as {{.Username}}...", map[string]interface{}{
		"Username": terminal.EntityNameColor(cmd.config.Username())}))

	stagingEnvVars, err := cmd.environmentVariableGroupRepo.ListStaging()
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()

	table := terminal.NewTable(cmd.ui, []string{i18n.T("Variable Name"), i18n.T("Assigned Value")})
	for _, envVar := range stagingEnvVars {
		table.Add(envVar.Name, envVar.Value)
	}
	table.Print()
}
