package commands

import (
	"github.com/cloudfoundry/cli/cf/api/stacks"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
)

type ListStacks struct {
	ui         terminal.UI
	config     core_config.Reader
	stacksRepo stacks.StackRepository
}

func init() {
	command_registry.Register(&ListStacks{})
}

func (cmd *ListStacks) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "stacks",
		Description: i18n.T("List all stacks (a stack is a pre-built file system, including an operating system, that can run apps)"),
		Usage:       i18n.T("CF_NAME stacks"),
	}
}

func (cmd *ListStacks) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. Requires app name as argument\n\n") + command_registry.Commands.CommandUsage("stacks"))
	}

	reqs = append(reqs, requirementsFactory.NewLoginRequirement())
	return
}

func (cmd *ListStacks) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.stacksRepo = deps.RepoLocator.GetStackRepository()
	return cmd
}

func (cmd *ListStacks) Execute(c flags.FlagContext) {
	cmd.ui.Say(i18n.T("Getting stacks in org {{.OrganizationName}} / space {{.SpaceName}} as {{.Username}}...",
		map[string]interface{}{"OrganizationName": terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
			"SpaceName": terminal.EntityNameColor(cmd.config.SpaceFields().Name),
			"Username":  terminal.EntityNameColor(cmd.config.Username())}))

	stacks, apiErr := cmd.stacksRepo.FindAll()
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{i18n.T("name"), i18n.T("description")})

	for _, stack := range stacks {
		table.Add(stack.Name, stack.Description)
	}

	table.Print()
}
