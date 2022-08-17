/*
Copyright © 2022 chez-shanpu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/chez-shanpu/kubectl-template/pkg/resource"
	"github.com/spf13/cobra"
)

var name string
var opts map[string]string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "kubectl-template",
	Short:        "kubernetes manifest template generator",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("number of args must be one")
		}

		var template resource.Template
		var err error
		switch args[0] {
		case "deploy":
			fallthrough
		case "deployment":
			if template, err = resource.NewDeploymentTemplate(name, opts); err != nil {
				return err
			}
		default:
			if template, err = resource.NewCustomTemplate(args[0], name, opts); err != nil {
				return err
			}
		}
		manifest, err := template.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate a manifest: %w", err)
		}
		fmt.Println(manifest)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()

	flags.StringVar(&name, "name", "", "resource name")
	flags.StringToStringVar(&opts, "opts", nil, "options (e.g. --opts key=val)")

	if err := rootCmd.MarkFlagRequired("name"); err != nil {
		fmt.Printf("failed to init: %v", err)
		os.Exit(1)
	}
}
