# GitHub Exporter

GitHub Exporter collects all .md files from the specified path in the git-repository of a service.

## Authentication

Github-Exporter requires authentication token in order to access private repositories hosted on GitHub. 

Personal Access Token - repo is sufficient.
Guide: [GitHub API Docs](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)

Paste a newly generated token to field `authToken` in `helm/github-exporter/values.yaml`. This token can be reused by all members of the team.

## Configuration

### Target Deployment

Target deployment has to provide the following annotations in metadata:
```yaml
    madock8s.exporter/githubExporter.baseurl: https://api.github.com/repos/MaibornWolff/madock8s/contents/<path>
```

Annotaions Breakdown:
- __baseurl__: The base URL of the project, accessible via GitHub API. Use <path> if the target document is not in the repository's root.

### Delete Mode

Configure behaviour for md-files when target deployment is deleted.
Simply update helm/values.yaml/env.deletionStrategy to one of the following values:
- IGNORE - keep existing md-files without modification;
- UPDATE - add "Deployments of the service were deleted on timestamp" to existing md-file (default);
- DELETE - remove the respective md-file from storage.
