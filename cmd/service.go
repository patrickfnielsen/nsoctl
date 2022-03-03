/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
			spinner, _ := pterm.DefaultSpinner.Start("Redeploying service...")
			if err := instance.Redeploy(nsoService, serviceId); err == nil {
				spinner.Success("Service re-deployed succesfully")
			} else {
				spinner.Fail(err)
			}

			return
		}

		// un-deploy a service in NSO
		if undeploy {
			spinner, _ := pterm.DefaultSpinner.Start("Undeploying service...")
			if err := instance.Undeploy(nsoService, serviceId); err == nil {
				spinner.Success("Service un-deployed succesfully")
			} else {
				spinner.Fail(err)
			}

			return
		}

		// no flags default is to get the service
		spinner, _ := pterm.DefaultSpinner.Start("Getting service...")
		url := fmt.Sprintf("/restconf/data/tailf-ncs:services/%s:%s=%s", nsoService, nsoService, serviceId)
		resp, err := instance.Get(url)
		if err != nil {
			spinner.Fail(err)
			return
		}

		data, err := instance.GetBody(resp)
		if err != nil {
			spinner.Fail(err)
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode == 404 {
			spinner.Warning("No service found...")
			return
		}

		spinner.Success()
		pterm.Println(data)
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	serviceCmd.Flags().StringVarP(&nsoService, "service", "s", "", "name of service (required)")
	serviceCmd.Flags().BoolVarP(&redeploy, "redeploy", "r", false, "redeploy service")
	serviceCmd.Flags().BoolVarP(&undeploy, "undeploy", "u", false, "undeploy service")

	rootCmd.MarkFlagRequired("service")
}
