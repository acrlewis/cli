package serviceauthtoken

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
)

type ListServiceAuthTokens struct {
	ui            terminal.UI
	config        core_config.Reader
	authTokenRepo api.ServiceAuthTokenRepository
}

func init() {
	command_registry.Register(&ListServiceAuthTokens{})
}

func (cmd *ListServiceAuthTokens) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "service-auth-tokens",
		Description: i18n.T("List service auth tokens"),
		Usage:       i18n.T("CF_NAME service-auth-tokens"),
	}
}

func (cmd *ListServiceAuthTokens) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. No argument required\n\n") + command_registry.Commands.CommandUsage("service-auth-tokens"))
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *ListServiceAuthTokens) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.authTokenRepo = deps.RepoLocator.GetServiceAuthTokenRepository()
	return cmd
}

func (cmd *ListServiceAuthTokens) Execute(c flags.FlagContext) {
	cmd.ui.Say(i18n.T("Getting service auth tokens as {{.CurrentUser}}...",
		map[string]interface{}{
			"CurrentUser": terminal.EntityNameColor(cmd.config.Username()),
		}))
	authTokens, apiErr := cmd.authTokenRepo.FindAll()
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}
	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{i18n.T("label"), i18n.T("provider")})

	for _, authToken := range authTokens {
		table.Add(authToken.Label, authToken.Provider)
	}

	table.Print()
}
