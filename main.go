package main

import (
	"flag"

	gitlabPatches "github.com/0x4f53/gitlab-patches"
)

func main() {
	perPage := flag.Int("per", 1000, "Results to grab per page from GitLab API (default: 1000)")
	maxCommits := flag.Int("max", 1000, "Maximum number of commits to grab from GitLab API (default: 1000)")
	outputDir := flag.String("output", gitlabPatches.GitlabCacheDir, "the directory to save files to. "+gitlabPatches.GitlabCacheDir+" will be made locally if not specified")

	flag.Parse()

	if *outputDir != "" {
		gitlabPatches.GitlabCacheDir = *outputDir
	}

	_ = gitlabPatches.GetGitlabCommits(*perPage, *maxCommits)
}
