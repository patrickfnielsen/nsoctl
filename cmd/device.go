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
	syncFrom bool
)

func syncFromDevice(nso *nso.NSO, deviceName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Starting sync-from on device...")
	if err := nso.SyncFromDevice(deviceName); err == nil {
		spinner.Success(deviceName + " has been synced!")
	} else {
		spinner.Fail(err)
	}
}

func getDeviceConfig(nso *nso.NSO, deviceName string) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting device config...")
	url := fmt.Sprintf("/restconf/data/tailf-ncs:devices/device=%s/config/", deviceName)
	resp, err := nso.Get(url)
	if err != nil {
		spinner.Fail(err)
		return
	}

	spinner.Success()
	pterm.Println(resp.Data)
}

// deviceCmd represents the service command
var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "device actions",
	Long: `get a device, or run a sync-from
nsoctl device <name> [flags]
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deviceName := args[0]
		instance := nso.NSO{
			Server:   cfg.Nso.ServerFqdn,
			Username: cfg.Nso.Username,
			Password: cfg.Nso.Password,
			Timeout:  900,
		}

		if syncFrom {
			syncFromDevice(&instance, deviceName)
			return
		}

		// no flags default is to get the device
		getDeviceConfig(&instance, deviceName)
	},
}

func init() {
	rootCmd.AddCommand(deviceCmd)

	deviceCmd.Flags().BoolVarP(&syncFrom, "sync-from", "s", false, "run sync-from on device")
}
