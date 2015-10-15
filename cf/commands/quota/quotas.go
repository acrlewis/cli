package quota

import (
	"fmt"
	"strconv"

	"github.com/simonleung8/flags"

	"github.com/cloudfoundry/cli/cf/api/quotas"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/formatters"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
)

type ListQuotas struct {
	ui        terminal.UI
	config    core_config.Reader
	quotaRepo quotas.QuotaRepository
}

func init() {
	command_registry.Register(&ListQuotas{})
}

func (cmd *ListQuotas) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "quotas",
		Description: i18n.T("List available usage quotas"),
		Usage:       i18n.T("CF_NAME quotas"),
	}
}

func (cmd *ListQuotas) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. No argument required\n\n") + command_registry.Commands.CommandUsage("quotas"))
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *ListQuotas) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.quotaRepo = deps.RepoLocator.GetQuotaRepository()
	return cmd
}

func (cmd *ListQuotas) Execute(c flags.FlagContext) {
	cmd.ui.Say(i18n.T("Getting quotas as {{.Username}}...", map[string]interface{}{"Username": terminal.EntityNameColor(cmd.config.Username())}))

	quotas, apiErr := cmd.quotaRepo.FindAll()

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}
	cmd.ui.Ok()
	cmd.ui.Say("")

	table := terminal.NewTable(cmd.ui, []string{i18n.T("name"), i18n.T("total memory limit"), i18n.T("instance memory limit"), i18n.T("routes"), i18n.T("service instances"), i18n.T("paid service plans")})

	var megabytes string
	for _, quota := range quotas {
		if quota.InstanceMemoryLimit == -1 {
			megabytes = i18n.T("unlimited")
		} else {
			megabytes = formatters.ByteSize(quota.InstanceMemoryLimit * formatters.MEGABYTE)
		}

		servicesLimit := strconv.Itoa(quota.ServicesLimit)
		if quota.ServicesLimit == -1 {
			servicesLimit = i18n.T("unlimited")
		}

		table.Add(
			quota.Name,
			formatters.ByteSize(quota.MemoryLimit*formatters.MEGABYTE),
			megabytes,
			fmt.Sprintf("%d", quota.RoutesLimit),
			fmt.Sprintf(servicesLimit),
			formatters.Allowed(quota.NonBasicServicesAllowed),
		)
	}

	table.Print()
}
