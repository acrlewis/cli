package securitygroup

import (
	"strings"

	"github.com/cloudfoundry/cli/cf/api/security_groups"
	"github.com/cloudfoundry/cli/cf/api/security_groups/defaults/running"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
)

type unbindFromRunningGroup struct {
	ui                terminal.UI
	configRepo        core_config.Reader
	securityGroupRepo security_groups.SecurityGroupRepo
	runningGroupRepo  running.RunningSecurityGroupsRepo
}

func init() {
	command_registry.Register(&unbindFromRunningGroup{})
}

func (cmd *unbindFromRunningGroup) MetaData() command_registry.CommandMetadata {
	primaryUsage := i18n.T("CF_NAME unbind-running-security-group SECURITY_GROUP")
	tipUsage := i18n.T("TIP: Changes will not apply to existing running applications until they are restarted.")
	return command_registry.CommandMetadata{
		Name:        "unbind-running-security-group",
		Description: i18n.T("Unbind a security group from the set of security groups for running applications"),
		Usage:       strings.Join([]string{primaryUsage, tipUsage}, "\n\n"),
	}
}

func (cmd *unbindFromRunningGroup) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) ([]requirements.Requirement, error) {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("unbind-running-security-group"))
	}

	return []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}, nil
}

func (cmd *unbindFromRunningGroup) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.configRepo = deps.Config
	cmd.securityGroupRepo = deps.RepoLocator.GetSecurityGroupRepository()
	cmd.runningGroupRepo = deps.RepoLocator.GetRunningSecurityGroupsRepository()
	return cmd
}

func (cmd *unbindFromRunningGroup) Execute(context flags.FlagContext) {
	name := context.Args()[0]

	securityGroup, err := cmd.securityGroupRepo.Read(name)
	switch (err).(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(i18n.T("Security group {{.security_group}} {{.error_message}}",
			map[string]interface{}{
				"security_group": terminal.EntityNameColor(name),
				"error_message":  terminal.WarningColor(i18n.T("does not exist.")),
			}))
		return
	default:
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Say(i18n.T("Unbinding security group {{.security_group}} from defaults for running as {{.username}}",
		map[string]interface{}{
			"security_group": terminal.EntityNameColor(securityGroup.Name),
			"username":       terminal.EntityNameColor(cmd.configRepo.Username()),
		}))
	err = cmd.runningGroupRepo.UnbindFromRunningSet(securityGroup.Guid)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}
	cmd.ui.Ok()
	cmd.ui.Say("\n\n")
	cmd.ui.Say(i18n.T("TIP: Changes will not apply to existing running applications until they are restarted."))
}
