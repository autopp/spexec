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
	"io"
	"io/ioutil"

	"github.com/autopp/spexec/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Main is the entrypoint of command line
func Main(version string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	const versionFlag = "version"
	cmd := &cobra.Command{
		Use:          "spexec file",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if v, err := cmd.Flags().GetBool(versionFlag); err != nil {
				return err
			} else if v {
				cmd.Println(version)
				return nil
			}

			cmd.Println(args[0])
			in, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			var c config.Config
			yaml.Unmarshal(in, &c)

			return nil
		},
	}

	cmd.Flags().Bool(versionFlag, false, "print version")

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args[1:])

	return cmd.Execute()
}
