package servicebroker

import (
	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/simonleung8/flags"
	"github.com/simonleung8/flags/flag"
)

type DeleteServiceBroker struct {
	ui     terminal.UI
	config core_config.Reader
	repo   api.ServiceBrokerRepository
}

func init() {
	command_registry.Register(&DeleteServiceBroker{})
}

func (cmd *DeleteServiceBroker) MetaData() command_registry.CommandMetadata {
	fs := make(map[string]flags.FlagSet)
	fs["f"] = &cliFlags.BoolFlag{Name: "f", Usage: i18n.T("Force deletion without confirmation")}

	return command_registry.CommandMetadata{
		Name:        "delete-service-broker",
		Description: i18n.T("Delete a service broker"),
		Usage:       i18n.T("CF_NAME delete-service-broker SERVICE_BROKER [-f]"),
		Flags:       fs,
	}
}

func (cmd *DeleteServiceBroker) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 1 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. Requires an argument\n\n") + command_registry.Commands.CommandUsage("delete-service-broker"))
	}

	reqs = append(reqs, requirementsFactory.NewLoginRequirement())
	return
}

func (cmd *DeleteServiceBroker) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.config = deps.Config
	cmd.repo = deps.RepoLocator.GetServiceBrokerRepository()
	return cmd
}

func (cmd *DeleteServiceBroker) Execute(c flags.FlagContext) {
	brokerName := c.Args()[0]
	if !c.Bool("f") && !cmd.ui.ConfirmDelete(i18n.T("service-broker"), brokerName) {
		return
	}

	cmd.ui.Say(i18n.T("Deleting service broker {{.Name}} as {{.Username}}...",
		map[string]interface{}{
			"Name":     terminal.EntityNameColor(brokerName),
			"Username": terminal.EntityNameColor(cmd.config.Username()),
		}))

	broker, apiErr := cmd.repo.FindByName(brokerName)

	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn(i18n.T("Service Broker {{.Name}} does not exist.", map[string]interface{}{"Name": brokerName}))
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.repo.Delete(broker.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	return
}
