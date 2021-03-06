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
	"fmt"
	"io"
	"os"

	"github.com/autopp/spexec/internal/errors"
	"github.com/autopp/spexec/internal/matcher/status"
	"github.com/autopp/spexec/internal/matcher/stream"
	"github.com/autopp/spexec/internal/reporter"
	"github.com/autopp/spexec/internal/runner"
	"github.com/autopp/spexec/internal/spec/parser"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type options struct {
	filename string
	output   string
	color    string
	format   string
}

// Main is the entrypoint of command line
func Main(version string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	opts := &options{}

	const versionFlag = "version"
	const outputFlag = "output"
	const colorFlag = "color"
	const formatFlag = "format"

	cmd := &cobra.Command{
		Use:           "spexec file",
		SilenceUsage:  true,
		SilenceErrors: false,
		Args:          cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if v, err := cmd.Flags().GetBool(versionFlag); err != nil {
				return err
			} else if v {
				cmd.Println(version)
				return nil
			}

			if err := opts.complete(cmd, args); err != nil {
				return err
			}

			return opts.run()
		},
	}

	cmd.Flags().Bool(versionFlag, false, "print version")
	cmd.Flags().StringVarP(&opts.output, outputFlag, "o", "", "output to file")
	cmd.Flags().StringVar(&opts.color, colorFlag, "auto", "color output")
	cmd.RegisterFlagCompletionFunc(colorFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	})
	cmd.Flags().StringVar(&opts.format, formatFlag, "simple", "format")
	cmd.RegisterFlagCompletionFunc(formatFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"simple", "documentation"}, cobra.ShellCompDirectiveDefault
	})

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	return cmd.Execute()
}

func (o *options) complete(cmd *cobra.Command, args []string) error {
	o.filename = args[0]

	if err := validateEnumFlag(o.color, "always", "never", "auto"); err != nil {
		return err
	}

	if err := validateEnumFlag(o.format, "simple", "documentation"); err != nil {
		return err
	}

	return nil
}

func validateEnumFlag(value string, validValues ...string) error {
	for _, v := range validValues {
		if value == v {
			return nil
		}
	}

	return fmt.Errorf("invalid --color flag: %s", value)
}

func (o *options) run() error {
	statusMR := status.NewStatusMatcherRegistryWithBuiltins()
	streamMR := stream.NewStreamMatcherRegistryWithBuiltins()
	tests, err := parser.New(statusMR, streamMR).ParseFile(o.filename)
	if err != nil {
		return err
	}

	runner := runner.NewRunner()
	reporterOpts := make([]reporter.Option, 0)
	out := os.Stdout
	if len(o.output) != 0 {
		out, err = os.Create(o.output)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	reporterOpts = append(reporterOpts, reporter.WithWriter(out))

	if err != nil {
		return err
	}
	var colorMode bool
	switch o.color {
	case "always":
		colorMode = true
	case "never":
		colorMode = false
	case "auto":
		colorMode = isatty.IsTerminal(out.Fd())
	}
	reporterOpts = append(reporterOpts, reporter.WithColor(colorMode))

	var formatter reporter.ReportFormatter
	switch o.format {
	case "simple":
		formatter = &reporter.SimpleFormatter{}
	case "documentation":
		formatter = &reporter.DocumentationFormatter{}
	}
	reporterOpts = append(reporterOpts, reporter.WithFormatter(formatter))

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
		return errors.New(errors.ErrTestFailed, "test failed")
	}
	return nil
}
