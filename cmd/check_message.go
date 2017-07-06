package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/petems/githelpy/githelpy"
)

// checkMessageCmd represents the command that check a message
var checkMessageCmd = &cobra.Command{
	Use:   "message [message]",
	Short: "Check message",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFileConfig()

		if err != nil {
			failure(err)

			exitError()
		}

		message, err := extractCheckMessageArgs(args)

		if err != nil {
			failure(err)

			exitError()
		}

		q := githelpy.MessageQuery{
			Message:  message,
			Matchers: viper.GetStringMapString("matchers"),
			Options: map[string]bool{
				"check-summary-length":  viper.GetBool("config.check-summary-length"),
				"exclude-merge-commits": viper.GetBool("config.exclude-merge-commits"),
			},
		}

		matching, err := githelpy.MatchMessageQuery(q)

		matchings := &[]*githelpy.Matching{}

		if !githelpy.IsZeroMatching(matching) {
			*matchings = append(*matchings, matching)
		}

		processMatchResult(matchings, err, viper.GetStringMapString("examples"))
	},
}

func extractCheckMessageArgs(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("One argument required : message")
	}

	return args[0], nil
}

func init() {
	checkCmd.AddCommand(checkMessageCmd)
}
