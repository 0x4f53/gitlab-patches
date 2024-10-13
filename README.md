# gitlab-patches

does what is says on the tin

this tool helps you scrape the last x commits pushed to gitlab (and grab their patch files) using the gitlab api. It returns the patch URL, associated domains, repo type and more!

# Usage

```bash
go build
./gitlab-patches
```

# Examples

A list of the last X patches pushed to gitlab in json format

```bash
./gitlab-patches --json
```

This can be combined with a service file or a cron to automate it.