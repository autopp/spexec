// Copyright (C) 2021-2022	 Akira Tanimura (@autopp)
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
	"github.com/autopp/spexec/internal/model"
	"github.com/autopp/spexec/internal/model/template"
	"github.com/autopp/spexec/internal/reporter"
	"github.com/autopp/spexec/internal/runner"
	"github.com/autopp/spexec/internal/spec"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

type options struct {
	filenames []string
	isStdin   bool
	output    string
	color     string
	format    string
	isStrict  bool
}

const versionFlag = "version"
const outputFlag = "output"
const colorFlag = "color"
const formatFlag = "format"
const strictFlag = "strict"

// Main is the entrypoint of command line
func Main(version string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	opts := &options{}

	cmd := &cobra.Command{
		Use:           "spexec file",
		SilenceUsage:  true,
		SilenceErrors: true,
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
	cmd.Flags().BoolVar(&opts.isStrict, strictFlag, false, "parse spec with strict mode")

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)

	return cmd.Execute()
}

func (o *options) complete(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.Errorf(errors.ErrInvalidSpec, "spec is not given")
	} else if args[0] == "-" {
		o.filenames = []string{"<stdin>"}
		o.isStdin = true
	} else {
		o.filenames = args
		o.isStdin = false
	}

	if err := validateEnumFlag(colorFlag, o.color, "always", "never", "auto"); err != nil {
		return err
	}

	if err := validateEnumFlag(formatFlag, o.format, "simple", "documentation", "json"); err != nil {
		return err
	}

	return nil
}

func validateEnumFlag(flag, value string, validValues ...string) error {
	for _, v := range validValues {
		if value == v {
			return nil
		}
	}

	return fmt.Errorf("invalid --%s flag: %s", flag, value)
}

func (o *options) run() error {
	statusMR := status.NewStatusMatcherRegistryWithBuiltins()
	streamMR := stream.NewStreamMatcherRegistryWithBuiltins()

	p := spec.NewParser(statusMR, streamMR)
	specTemplates := []struct {
		filename      string
		testTemplates []*template.TestTemplate
	}{}
	var err error
	env := model.NewEnv(nil)
	if o.isStdin {
		v, err := model.NewValidator("", o.isStrict)
		if err != nil {
			return err
		}
		testTemplates, err := p.ParseStdin(env, v)
		if err != nil {
			return err
		}
		specTemplates = append(specTemplates, struct {
			filename      string
			testTemplates []*template.TestTemplate
		}{"<stdin>", testTemplates})
	} else {
		for _, filename := range o.filenames {
			v, err := model.NewValidator(filename, o.isStrict)
			if err != nil {
				return err
			}
			testTemplates, err := p.ParseFile(env, v, filename)
			if err != nil {
				return err
			}
			specTemplates = append(specTemplates, struct {
				filename      string
				testTemplates []*template.TestTemplate
			}{filename, testTemplates})
		}
	}

	out := os.Stdout
	if len(o.output) != 0 {
		out, err = os.Create(o.output)
		if err != nil {
			return err
		}
		defer out.Close()
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

	var formatter reporter.ReportFormatter
	switch o.format {
	case "simple":
		formatter = &reporter.SimpleFormatter{}
	case "documentation":
		formatter = &reporter.DocumentationFormatter{}
	case "json":
		formatter = &reporter.JSONFormatter{}
	}

	reporter, err := reporter.New(reporter.WithWriter(out), reporter.WithColor(colorMode), reporter.WithFormatter(formatter))
	if err != nil {
		return err
	}

	specs := []struct {
		filename string
		tests    []*model.Test
	}{}

	for _, st := range specTemplates {
		v, err := model.NewValidator(st.filename, o.isStrict)
		if err != nil {
			return err
		}

		tests := make([]*model.Test, 0)
		for _, tt := range st.testTemplates {
			t, err := tt.Expand(env, v, statusMR, streamMR)
			if err != nil {
				return err
			}
			tests = append(tests, t)
		}
		err = v.Error()
		if err != nil {
			return err
		}
		specs = append(specs, struct {
			filename string
			tests    []*model.Test
		}{filename: st.filename, tests: tests})
	}

	runner := runner.NewRunner()
	var results []*model.TestResult
	for _, spec := range specs {
		rs, err := runner.RunTests(spec.filename, spec.tests, reporter)
		if err != nil {
			return err
		}
		results = append(results, rs...)
	}

	for _, r := range results {
		if !r.IsSuccess {
			return errors.New(errors.ErrTestFailed, "test failed")
		}
	}

	return nil
}
