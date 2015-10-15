package help

import "github.com/cloudfoundry/cli/cf/i18n"

func GetHelpTemplate() string {
	return `{{.Title "` + i18n.T("NAME:") + `"}}
   {{.Name}} - {{.Usage}}

{{.Title "` + i18n.T("USAGE:") + `"}}
   ` + i18n.T("[environment variables]") + ` {{.Name}} ` + i18n.T("[global options] command [arguments...] [command options]") + `

{{.Title "` + i18n.T("VERSION:") + `"}}
   {{.Version}}

{{.Title "` + i18n.T("BUILD TIME:") + `"}}
   {{.Compiled}}
   {{range .Commands}}
{{.SubTitle .Name}}{{range .CommandSubGroups}}
{{range .}}   {{.Name}} {{.Description}}
{{end}}{{end}}{{end}}
{{.Title "` + i18n.T("ENVIRONMENT VARIABLES:") + `"}}
   CF_COLOR=false                     ` + i18n.T("Do not colorize output") + `
   CF_HOME=path/to/dir/               ` + i18n.T("Override path to default config directory") + `
   CF_PLUGIN_HOME=path/to/dir/        ` + i18n.T("Override path to default plugin config directory") + `
   CF_STAGING_TIMEOUT=15              ` + i18n.T("Max wait time for buildpack staging, in minutes") + `
   CF_STARTUP_TIMEOUT=5               ` + i18n.T("Max wait time for app instance startup, in minutes") + `
   CF_TRACE=true                      ` + i18n.T("Print API request diagnostics to stdout") + `
   CF_TRACE=path/to/trace.log         ` + i18n.T("Append API request diagnostics to a log file") + `
   HTTP_PROXY=proxy.example.com:8080  ` + i18n.T("Enable HTTP proxying for API requests") + `

{{.Title "` + i18n.T("GLOBAL OPTIONS:") + `"}}
   --version, -v                      ` + i18n.T("Print the version") + `
   --build, -b                        ` + i18n.T("Print the version of Go the CLI was built against") + `
   --help, -h                         ` + i18n.T("Show help") + `

`
}
