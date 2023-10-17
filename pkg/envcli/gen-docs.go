// Copyright 2023 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package envcli

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func genDocsCmd() *cobra.Command {
	genDocsCmd := &cobra.Command{
		Use:   "gen-docs <path>",
		Short: "[Internal] Generate documentation for the CLI",
		Long: "[Internal] Generates the documentation for the CLI's Cobra commands. " +
			"Docs are placed in the directory specified by <path>.",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return errors.WithStack(err)
			}
			docsPath := filepath.Join(wd, args[0] /* relative path */)

			// We clear out the existing directory so that the doc-pages for
			// commands that have been deleted in the CLI will also be removed
			// after we re-generate the docs below
			if err := clearDir(docsPath); err != nil {
				return err
			}

			rootCmd := cmd
			for rootCmd.HasParent() {
				rootCmd = rootCmd.Parent()
			}

			// Removes the line in the generated docs of the form:
			// ###### Auto generated by spf13/cobra on 18-Jul-2022
			rootCmd.DisableAutoGenTag = true

			return errors.WithStack(doc.GenMarkdownTree(rootCmd, docsPath))
		},
	}

	return genDocsCmd
}

func clearDir(dir string) error {
	// if the dir doesn't exist, use default filemode 0755 to create it
	// if the dir exists, use its own filemode to re-create it
	var mode os.FileMode
	f, err := os.Stat(dir)
	if err == nil {
		mode = f.Mode()
	} else if errors.Is(err, fs.ErrNotExist) {
		mode = 0o755
	} else {
		return errors.WithStack(err)
	}

	if err := os.RemoveAll(dir); err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(os.MkdirAll(dir, mode))
}
