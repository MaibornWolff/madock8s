# GitLab Exporter

Gitlab Exporter collects all .md files in the git-repository of a service.

## Authentication

Gitlab-Exporter requires authentication token in order to access private repositories hosted on GitLab. 
Supported authentication:
- *Recommended*: Project Access Token - read_api is sufficient.
Guide: [GitLab API Docs](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html)


- Personal Access Token - read_api is sufficient.
Guide: [GitLab API Docs](https://docs.gitlab.com/ee/api/README.html#personalproject-access-tokens)

Paste a newly generated token to field `gitlabToken` in `helm/gitlab-exporter/values.yaml`. This token can be reused by all members of the team.

## Configuration

### Target Deployment

Target deployment has to provide the following annotations in metadata:
```
    madock8s.exporter/gitlabExporter.baseurl: https://git.example.com/api/v4/projects/9999/repository/
    madock8s.exporter/gitlabExporter.path: exporter/sample-metrics
    madock8s.exporter/gitlabExporter.pattern: .\\.md
    madock8s.exporter/gitlabExporter.recursive: true
    madock8s.exporter/gitlabExporter.ref: master
```

Parameters Breakdown:
- __baseurl__: The base URL of the project, accessible via GitLab API. NOTE: 9999 is Project ID, not name.
- __path__: The relative path to the root. The exporter searches for files from here.
- __pattern__ (optional): Pattern that is used to search files. The default behaviour is to search all .md files.
- __recursive__ (optional): Boolean flag, Enables search in sub-directories for the path. "true" (default): search in sub-directories; "false": search only in path.
- __ref__: Name of git-branch where required files are stored.

### Delete Mode

Configure behaviour for md-files when target deployment is deleted.
Simply update helm/values.yaml/env.deletionStrategy to one of the following values:
- IGNORE - keep existing md-files without modification;
- UPDATE - add "Deployments of the service were deleted on timestamp" to existing md-file (default);
- DELETE - remove the respective md-file from storage.
