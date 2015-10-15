package buildpack

import (
	"strconv"

	"github.com/simonleung8/flags"

	"github.com/cloudfoundry/cli/cf/api"
	"github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
)

type ListBuildpacks struct {
	ui            terminal.UI
	buildpackRepo api.BuildpackRepository
}

func init() {
	command_registry.Register(&ListBuildpacks{})
}

func (cmd *ListBuildpacks) MetaData() command_registry.CommandMetadata {
	return command_registry.CommandMetadata{
		Name:        "buildpacks",
		Description: i18n.T("List all buildpacks"),
		Usage:       i18n.T("CF_NAME buildpacks"),
	}
}

func (cmd *ListBuildpacks) Requirements(requirementsFactory requirements.Factory, fc flags.FlagContext) (reqs []requirements.Requirement, err error) {
	if len(fc.Args()) != 0 {
		cmd.ui.Failed(i18n.T("Incorrect Usage. No argument required\n\n") + command_registry.Commands.CommandUsage("buildpacks"))
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *ListBuildpacks) SetDependency(deps command_registry.Dependency, pluginCall bool) command_registry.Command {
	cmd.ui = deps.Ui
	cmd.buildpackRepo = deps.RepoLocator.GetBuildpackRepository()
	return cmd
}

func (cmd *ListBuildpacks) Execute(c flags.FlagContext) {
	cmd.ui.Say(i18n.T("Getting buildpacks...\n"))

	table := cmd.ui.Table([]string{"buildpack", i18n.T("position"), i18n.T("enabled"), i18n.T("locked"), i18n.T("filename")})
	noBuildpacks := true

	apiErr := cmd.buildpackRepo.ListBuildpacks(func(buildpack models.Buildpack) bool {
		position := ""
		if buildpack.Position != nil {
			position = strconv.Itoa(*buildpack.Position)
		}
		enabled := ""
		if buildpack.Enabled != nil {
			enabled = strconv.FormatBool(*buildpack.Enabled)
		}
		locked := ""
		if buildpack.Locked != nil {
			locked = strconv.FormatBool(*buildpack.Locked)
		}
		table.Add(
			buildpack.Name,
			position,
			enabled,
			locked,
			buildpack.Filename,
		)
		noBuildpacks = false
		return true
	})
	table.Print()

	if apiErr != nil {
		cmd.ui.Failed(i18n.T("Failed fetching buildpacks.\n{{.Error}}", map[string]interface{}{"Error": apiErr.Error()}))
	}

	if noBuildpacks {
		cmd.ui.Say(i18n.T("No buildpacks found"))
	}
}
