package main

import (
	"flag"
	"gitlabPatches"
)

func main() {
	perPage := flag.Int("per", 500, "Results to grab per page from GitLab API")
	maxCommits := flag.Int("max", 500, "Maximum number of commits to grab from GitLab API")
	outputDir := flag.String("output", gitlabPatches.GitlabCacheDir, "the directory to save files to. "+gitlabPatches.GitlabCacheDir+" will be made locally if not specified")

	flag.Parse()

	if *outputDir != "" {
		gitlabPatches.GitlabCacheDir = *outputDir
	}

	_ = gitlabPatches.GetGitlabCommits(*perPage, *maxCommits)
}
