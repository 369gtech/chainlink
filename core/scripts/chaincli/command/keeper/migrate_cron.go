package keeper

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/handler"
)

var inputFile string

var MigrateCronCmd = &cobra.Command{
	Use:   "migrate-cron",
	Short: "Migrate feed util jobs to use Cron upkeep",
	Long:  `This command reads in a list of feed contracts from input file, creates a new cron keeper contract and registers the upkeep`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		hdlr := handler.NewKeeper(cfg)

		if inputFile == "" {
			log.Fatal("Input file should be provided")
		}

		hdlr.MigrateCron(cmd.Context(), inputFile)
	},
}

func init() {
	MigrateCronCmd.Flags().StringVar(&inputFile, "input-file", "", "path to csv file in format: targetAddress,targetFunction,cronSchedule")
}
