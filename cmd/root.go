package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/adamgoose/ghpatch/lib"
	"github.com/google/go-github/v56/github"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "ghpatch {path} [flags] -- {cmd}",
	Short: "ghpatch is a tool for patching GitHub files",
	Args:  cobra.MinimumNArgs(2),
	RunE: lib.RunE(func(ctx context.Context, args []string, gh *github.Client) error {
		// Get the source file
		r, _, _, err := gh.Repositories.GetContents(
			ctx,
			viper.GetString("github_owner"),
			viper.GetString("github_repo"),
			args[0],
			&github.RepositoryContentGetOptions{
				Ref: viper.GetString("github_branch"),
			},
		)
		if err != nil {
			return err
		}

		src, err := r.GetContent()
		if err != nil {
			return err
		}

		// Apply transformation command
		dst := bytes.NewBuffer(nil)
		cmd := exec.Command(args[1], args[2:]...)
		cmd.Env = os.Environ()
		cmd.Stdin = bytes.NewReader([]byte(src))
		cmd.Stderr = os.Stderr
		cmd.Stdout = dst

		if err := cmd.Run(); err != nil {
			return err
		}

		// Compare src and dst
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(src, dst.String(), false)

		// Check for --dry-run
		if viper.GetBool("dry_run") {
			fmt.Println(dmp.DiffPrettyText(diffs))
			return nil
		}

		// Skip if equal
		if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
			return nil
		}

		// Update the File
		rx, _, err := gh.Repositories.UpdateFile(
			ctx,
			viper.GetString("github_owner"),
			viper.GetString("github_repo"),
			args[0],
			&github.RepositoryContentFileOptions{
				Message: github.String(viper.GetString("commit_message")),
				Content: dst.Bytes(),
				SHA:     r.SHA,
				Branch:  github.String(viper.GetString("github_branch")),
				Committer: &github.CommitAuthor{
					Name:  github.String(viper.GetString("commit_author_name")),
					Email: github.String(viper.GetString("commit_author_email")),
				},
			})
		if err != nil {
			return err
		}

		fmt.Print(rx.GetSHA())
		return nil
	}),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("token", "t", "", "[GHPATCH_GITHUB_TOKEN] GitHub token")
	viper.BindPFlag("github_token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.PersistentFlags().StringP("owner", "o", "", "[GHPATCH_GITHUB_OWNER] GitHub owner")
	viper.BindPFlag("github_owner", rootCmd.PersistentFlags().Lookup("owner"))

	rootCmd.PersistentFlags().StringP("repo", "r", "", "[GHPATCH_GITHUB_REPO] GitHub repo")
	viper.BindPFlag("github_repo", rootCmd.PersistentFlags().Lookup("repo"))

	rootCmd.PersistentFlags().StringP("branch", "b", "", "[GHPATCH_GITHUB_BRANCH] GitHub branch")
	viper.BindPFlag("github_branch", rootCmd.PersistentFlags().Lookup("branch"))

	rootCmd.PersistentFlags().StringP("message", "m", "", "[GHPATCH_COMMIT_MESSAGE] Commit message")
	viper.BindPFlag("commit_message", rootCmd.PersistentFlags().Lookup("commit-message"))

	rootCmd.PersistentFlags().StringP("author-name", "n", "", "[GHPATCH_COMMIT_AUTHOR_NAME] Commit Author name")
	viper.BindPFlag("commit_author_name", rootCmd.PersistentFlags().Lookup("author-name"))

	rootCmd.PersistentFlags().StringP("author-email", "e", "", "[GHPATCH_COMMIT_AUTHOR_EMAIL] Commit Author email")
	viper.BindPFlag("commit_author_email", rootCmd.PersistentFlags().Lookup("author-email"))

	rootCmd.PersistentFlags().Bool("dry-run", false, "[GHPATCH_DRY_RUN] Dry run")
	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))
}
