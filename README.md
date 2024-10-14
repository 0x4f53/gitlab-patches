# gitlab-patches

does what is says on the tin

this tool helps you scrape the last x commits pushed to gitlab (and grab their patch files) using the gitlab api. It returns the patch URL, associated domains, repo type and more!

# Usage

```bash
go build
./gitlab-patches -h

Usage of ./gitlab-patches:
  -json
        Save as line-separated JSON in a file (filename format: <01-01-2024-0>.json)
  -max int
        Maximum number of commits to grab from GitLab API (max: 1000) (default 100)
  -per int
        Results to grab per page from GitLab API (max: 100) (default 100)
```

# Examples

A list of the last X patches pushed to gitlab in json format

```bash
./gitlab-patches --json
```

This can be combined with a service file or a cron to automate it.

### Output

```bash
./gitlab-patches --json

{"project_id":62541692,"project_name":"aphysica.gitlab.io","project_path":"aphysica.gitlab.io","project_namespace":"aphysica","project_web_url":"https://gitlab.com/aphysica/aphysica.gitlab.io","kind":"user","commit_id":"ea558b0ce600c0376c76c5f0261db47f1efe2195","commit_short_id":"ea558b0c","commit_title":"index","commit_message":"index","commit_created_at":"2024-10-14T00:09:37Z","commit_author_name":"aphysica","commit_patch_url":"https://gitlab.com/aphysica/aphysica.gitlab.io/-/commit/ea558b0ce600c0376c76c5f0261db47f1efe2195.patch","associated_domains":["gitlab.com"]}
{"project_id":62541798,"project_name":"shyameer - Security policy project","project_path":"shyameer-security-policy-project","project_namespace":"shyameer","project_web_url":"https://gitlab.com/shyamee11/shyameer-security-policy-project","kind":"group","commit_id":"e25c2a197d0a3c4d6ce871a8361304bdcdcbab88","commit_short_id":"e25c2a19","commit_title":"Initial commit","commit_message":"Initial commit","commit_created_at":"2024-10-14T00:12:03Z","commit_author_name":"tester bhai","commit_patch_url":"https://gitlab.com/shyamee11/shyameer-security-policy-project/-/commit/e25c2a197d0a3c4d6ce871a8361304bdcdcbab88.patch","associated_domains":null}
{"project_id":62541965,"project_name":"alphabet-soup","project_path":"alphabet-soup","project_namespace":"Colin Nguyen","project_web_url":"https://gitlab.com/colinn0803/alphabet-soup","kind":"user","commit_id":"523e8543d567390a63042f0dea4ef868f1c96e8c","commit_short_id":"523e8543","commit_title":"Minor README correction for input format.","commit_message":"Minor README correction for input format.","commit_created_at":"2019-04-29T13:47:07Z","commit_author_name":"Brett Meyers","commit_patch_url":"https://gitlab.com/colinn0803/alphabet-soup/-/commit/523e8543d567390a63042f0dea4ef868f1c96e8c.patch","associated_domains":["readme.md","eitccorp.com"]}
...
```