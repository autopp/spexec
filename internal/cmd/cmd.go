// Copyright (C) 2021 Akira Tanimura (@autopp)
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/autopp/spexec/internal/config"
	"github.com/autopp/spexec/internal/reporter"
	"github.com/autopp/spexec/internal/runner"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// Main is the entrypoint of command line
func Main(version string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	const versionFlag = "version"
	const outputFlag = "output"
	const colorFlag = "color"

	cmd := &cobra.Command{
		Use:           "spexec file",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if v, err := cmd.Flags().GetBool(versionFlag); err != nil {
				return err
			} else if v {
				cmd.Println(version)
				return nil
			}

			filename := args[0]
			f, err := os.Open(filename)
			if err != nil {
				return err
			}
			format := config.JSONFormat
			ext := filepath.Ext(filename)
			if ext == ".yml" || ext == ".yaml" {
				format = config.YAMLFormat
			}
			tests, err := config.Load(f, format)
			if err != nil {
				return err
			}

			runner := runner.NewRunner()
			reporterOpts := make([]reporter.Option, 0)
			output, err := cmd.Flags().GetString(outputFlag)
			if err != nil {
				return err
			}
			out := os.Stdout
			if len(output) != 0 {
				out, err = os.Create(output)
				if err != nil {
					return err
				}
				defer f.Close()
			}
			reporterOpts = append(reporterOpts, reporter.WithWriter(out))

			color, err := cmd.Flags().GetString(colorFlag)
			if err != nil {
				return err
			}
			var colorMode bool
			switch color {
			case "always":
				colorMode = true
			case "never":
				colorMode = false
			case "auto":
				colorMode = isatty.IsTerminal(out.Fd())
			default:
				return fmt.Errorf("invalid --color flag: %s", color)
			}
			reporterOpts = append(reporterOpts, reporter.WithColor(colorMode))

			reporter, err := reporter.New(reporterOpts...)
			if err != nil {
				return err
			}
			results := runner.RunTests(tests, reporter)

			allGreen := true
			for _, r := range results {
				if !r.IsSuccess {
					allGreen = false
					break
				}
			}

			if !allGreen {
				return errors.New("test failed")
			}
			return nil
		},
	}

	cmd.Flags().Bool(versionFlag, false, "print version")
	cmd.Flags().StringP(outputFlag, "o", "", "output to file")
	cmd.Flags().String(colorFlag, "auto", "color output")
	cmd.RegisterFlagCompletionFunc(colorFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	})

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args[1:])

	return cmd.Execute()
}
