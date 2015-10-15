package user

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/api/spaces"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
)

type UnsetSpaceRole struct {
	ui        terminal.UI
	config    core_config.Reader
	spaceRepo spaces.SpaceRepository
	userRepo  api.UserRepository
	userReq   requirements.UserRequirement
	orgReq    requirements.OrganizationRequirement
}

func init() {
	command_registry.Register(&UnsetSpaceRole{})
}

func (cmd *UnsetSpaceRole) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "unset-space-role",
		Description: i18n.T("Remove a space role from a user"),
		Usage: i18n.T("CF_NAME unset-space-role USERNAME ORG SPACE ROLE\n\n") +
			i18n.T("ROLES:\n") +
			i18n.T("   SpaceManager - Invite and manage users, and enable features for a given space\n") +
			i18n.T("   SpaceDeveloper - Create and manage apps and services, and see logs and reports\n") +
			i18n.T("   SpaceAuditor - View logs, reports, and settings on this space\n"),
	}
}

func (cmd *UnsetSpaceRole) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 4 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. Requires USERNAME, ORG, SPACE, ROLE as arguments\n\n") + command_registry.Commands.CommandUsage("unset-space-role"))
	}

	cmd.userReq = requirementsFactory.NewUserRequirement(fc.Args()[0])
	cmd.orgReq = requirementsFactory.NewOrganizationRequirement(fc.Args()[1])

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.userReq,
		cmd.orgReq,
	}

	return
}

func (cmd *UnsetSpaceRole) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.spaceRepo = deps.RepoLocator.GetSpaceRepository()
	cmd.userRepo = deps.RepoLocator.GetUserRepository()
	return cmd
}

func (cmd *UnsetSpaceRole) Execute(c flags.FlagContext) {
	spaceName := c.Args()[2]
	role := models.UserInputToSpaceRole[c.Args()[3]]

	user := cmd.userReq.GetUser()
	org := cmd.orgReq.GetOrganization()
	space, apiErr := cmd.spaceRepo.FindByNameInOrg(spaceName, org.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Say(i18n.T("Removing role {{.Role}} from user {{.TargetUser}} in org {{.TargetOrg}} / space {{.TargetSpace}} as {{.CurrentUser}}...",
		map[string]interface{}{
			"Role":        terminal.EntityNameColor(role),
			"TargetUser":  terminal.EntityNameColor(user.Username),
			"TargetOrg":   terminal.EntityNameColor(org.Name),
			"TargetSpace": terminal.EntityNameColor(space.Name),
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))

	apiErr = cmd.userRepo.UnsetSpaceRole(user.Guid, space.Guid, role)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
