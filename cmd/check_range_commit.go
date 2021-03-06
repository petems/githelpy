package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/petems/githelpy/githelpy"
)

// checkRangeCommitCmd represents the check command
var checkRangeCommitCmd = &cobra.Command{
	Use:   "range [commitFrom] [commitTo] [&path]",
	Short: "Check messages in commit range",
	Run: func(cmd *cobra.Command, args []string) {
		err := validateFileConfig()

		if err != nil {
			failure(err)

			exitError()
		}

		from, to, path, err := extractCheckRangeCommitArgs(args)

		if err != nil {
			failure(err)

			exitError()
		}

		q := githelpy.RangeCommitQuery{
			Path:     path,
			From:     from,
			To:       to,
			Matchers: viper.GetStringMapString("matchers"),
			Options: map[string]bool{
				"check-summary-length":  viper.GetBool("config.check-summary-length"),
				"exclude-merge-commits": viper.GetBool("config.exclude-merge-commits"),
			},
		}

		matchings, err := githelpy.MatchRangeCommitQuery(q)

		processMatchResult(matchings, err, viper.GetStringMapString("examples"))
	},
}

func extractCheckRangeCommitArgs(args []string) (string, string, string, error) {
	if len(args) < 2 {
		return "", "", "", fmt.Errorf("Two arguments required : origin commit and end commit")

	}

	if len(args) > 3 {
		return "", "", "", fmt.Errorf("3 arguments must be provided at most")
	}

	var path string
	var err error

	if len(args) == 3 {
		path = args[2]
	}

	path, err = parseDirectory(path)

	if err != nil {
		return "", "", "", err
	}

	return args[0], args[1], path, nil
}

func init() {
	checkCmd.AddCommand(checkRangeCommitCmd)
}
