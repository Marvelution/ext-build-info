# ext-build-info

## About this plugin
The ext-build-info plugin collects issues from git comment messages and branch names and validates that they are known by the tracker. 

## Installation with JFrog CLI
Installing the latest version:

`$ jf plugin install ext-build-info`

Installing a specific version:

`$ jf plugin install ext-build-info@version`

Uninstalling a plugin

`$ jf plugin uninstall ext-build-info`

## Usage
### Commands
* collect-issues
  - Arguments:
    - build name - The name of the build.
    - build number - The number of the build.
    - path to .git - Path to a directory containing the .git directory. If not specified, the .git directory is assumed to be in the 
      current directory or in one of the parent directories.
  - Flags
    - --server-id - [Optional] Server ID configured using the config command.
    - --project - [Optional] Project where the pipeline belongs to.
    - --tracker - [Optional] Tracker to use to collect related issue from.
    - --tracker-url - [Optional] Tracker base url to use to collect related issue from.
    - --tracker-username - [Optional] Tracker username to use to collect related issue from.
    - --tracker-token - [Optional] Tracker token to use to collect related issue from.
    - --regexp - [Optional] A regular expression used for matching the git commit messages.
    - --key-group-index - [Default: 1] The capturing group index in the regular expression used for retrieving the issue key.
    - --git-log-limit - [Default: 100] The maximum number of git commit messages to process.
    - --aggregate - [Default: false] Set to true, if you wish all builds to include issues from previous builds.
    - --aggregation-status - [Optional] If aggregate is set to true, this property indicates how far in time should the issues be 
      aggregated. In the above example, issues will be aggregated from previous builds, until a build with a RELEASE status is found. 
      Build statuses are set when a build is promoted using the jf rt build-promote command. 

  - Example:
    ```
    $ jf ext-build-info collect-issues --tracker=Jira MyBuild 1

    [Info] Reading the git branch, revision and remote URL and adding them to the build-info.
    [Info] Collecting build issues from VCS...  
    [Info] Searching Jira using request: {"jql": "issue IN (EX-1, EX2)", "fields": ["key", "summary"], "startAt": 0, maxResult: 100, "validateQuery": "warn"}
    [Info] Found Jira issue: EX-1
    [Info] Collected 1 issue details for MyBuild/1.
    ```
* clean-slate
  - Arguments:
    - build name - The name of the build.
    - build number - The number of the build.
  - Flags
    - --project - [Optional] Project where the pipeline belongs to.

  - Example:
    ```
    $ jf ext-build-info clean-slate MyBuild 1

    [Info] Clearing all existing build-info to start from a clean slate.
    [Info] Removing build-info directory: /path/to/build/info
    ```

* send-build-info
  - Arguments
    - build name - The name of the build.
    - build number - The number of the build.
  - Flags
    - --server-id - [Optional] Server ID configured using the config command, this needs to an Artifactory integration that uses an 
      Access Token.
    - --project - [Optional] Project where the pipeline belongs to.
    - --jira-id - [Optional] Jira ID to use to collect related issue from.
    - --jira-url - [Optional] Jira Url base url to use to collect related issue from.
    - --jira-client-id - [Optional] The OAuth clientId generated by Jira.
    - --jira-secret - [Optional] The OAuth secret generated by Jira.
    - --dry-run - [Optional] Enable to only log what would be send to Jira.
    - --include-pre-post-runs - [Optional] Enable to include pipeline preRun and postRun steps.
    - --fail-on-reject - [Optional] Enable to error out if any builds are rejected by Jira.
  - Example:
    ```
    $ jf ext-build-info send-build-info --server-id ArtifactoryAT --jira-id JiraOAuth MyBuild 1

    [Info] 12:02:45 [Info] Build MyBuild #1 was accepted by Jira
    ``` 

* send-deployment-info
  - Arguments
    - build name - The name of the build.
    - build number - The number of the build.
  - Flags
    - --server-id - [Optional] Server ID configured using the config command, this needs to an Artifactory integration that uses an
      Access Token.
    - --project - [Optional] Project where the pipeline belongs to.
    - --jira-id - [Optional] Jira ID to use to collect related issue from.
    - --jira-url - [Optional] Jira Url base url to use to collect related issue from.
    - --jira-client-id - [Optional] The OAuth clientId generated by Jira.
    - --jira-secret - [Optional] The OAuth secret generated by Jira.
    - --dry-run - [Optional] Enable to only log what would be send to Jira.
    - --include-pre-post-runs - [Optional] Enable to include pipeline preRun and postRun steps.
    - --fail-on-reject - [Optional] Enable to error out if any builds are rejected by Jira.
    - --environment - [Optional] The environment that the deployment targeted, default to environment variable named `environmentName`
  - Example:
    ```
    $ jf ext-build-info send-deployment-info --server-id ArtifactoryAT --jira-id JiraOAuth MyBuild 1

    [Info] 12:02:45 [Info] Deployment 1234567890 #1 was accepted by Jira
    ```

* notify-slack
  - Arguments
    - build name - The name of the build.
    - build number - The number of the build.
  - Flags
    - --server-id - [Optional] Server ID configured using the config command, this needs to an Artifactory integration that uses an
      Access Token.
    - --project - [Optional] Project where the pipeline belongs to.
    - --slack - The Slack integration name to send the message to.
    - --include-pre-post-runs - [Optional] Enable to include pipeline preRun and postRun steps.

### Environment variables
The plugin can lookup integration variables like url, username and token using the JFrog Pipelines integration environment variables.  

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
