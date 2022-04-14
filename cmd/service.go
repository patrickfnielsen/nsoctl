/*
Copyright © 2022 Patrick Falk Nielsen <git@patricknielsen.dk>
*/
package cmd

import (
	"fmt"

	"github.com/patrickfnielsen/nsoctl/nso"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	nsoService string = "opennet"
	redeploy   bool
	undeploy   bool
)

func redeployService(nso *nso.NSO, serviceId string) {
	spinner, _ := pterm.DefaultSpinner.Start("Redeploying service...")
	if err := nso.Redeploy(nsoService, serviceId); err == nil {
		spinner.Success("Service re-deployed succesfully")
	} else {
		spinner.Fail(err)
	}
}

func undeployService(nso *nso.NSO, serviceId string) {
	spinner, _ := pterm.DefaultSpinner.Start("Undeploying service...")
	if err := nso.Undeploy(nsoService, serviceId); err == nil {
		spinner.Success("Service un-deployed succesfully")
	} else {
		spinner.Fail(err)
	}
}

func getServiceData(nso *nso.NSO, serviceId string) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting service...")
	url := fmt.Sprintf("/restconf/data/tailf-ncs:services/%s:%s=%s", nsoService, nsoService, serviceId)
	resp, err := nso.Get(url)
	if err != nil {
		spinner.Fail(err)
		return
	}

	if resp.StatusCode == 404 {
		spinner.Warning("No service found...")
		return
	}

	spinner.Success()
	pterm.Println(resp.Data)
}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "service actions",
	Long: `get a service, and performe actions like re-deploy, un-deploy
nsoctl service <id> [flags]
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceId := args[0]
		instance := nso.NSO{
			Server:   cfg.Nso.ServerFqdn,
			Username: cfg.Nso.Username,
			Password: cfg.Nso.Password,
			Timeout:  120,
		}

		// re-deploy a service in NSO
		if redeploy {
			redeployService(&instance, serviceId)
			return
		}

		// un-deploy a service in NSO
		if undeploy {
			undeployService(&instance, serviceId)
			return
		}

		// no flags default is to get the service
		getServiceData(&instance, serviceId)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.Flags().StringVarP(&nsoService, "service", "s", "", "name of service (required)")
	serviceCmd.Flags().BoolVarP(&redeploy, "redeploy", "r", false, "redeploy service")
	serviceCmd.Flags().BoolVarP(&undeploy, "undeploy", "u", false, "undeploy service")

	rootCmd.MarkFlagRequired("service")
}
