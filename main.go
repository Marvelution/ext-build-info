package main

import (
	"errors"
	"fmt"
	artifactoryUtils "github.com/jfrog/jfrog-cli-core/v2/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/marvelution/ext-build-info/commands"
	"os"
	"strconv"
)

func main() {
	plugins.PluginMain(components.App{
		Name:        "ext-build-info",
		Description: "Extended build info.",
		Version:     "v1.6.1",
		Commands: []components.Command{
			{
				Name:        "collect-issues",
				Description: "Collect issue details from git and add them to a build.",
				Aliases:     []string{"ci"},
				Flags: []components.Flag{
					components.StringFlag{
						Name:        "server-id",
						Description: "Server ID configured using the config command.",
					},
					components.StringFlag{
						Name:        "project",
						Description: "Artifactory project key.",
					},
					components.StringFlag{
						Name:        "tracker",
						Description: "Tracker to use to collect related issue from.",
					},
					components.StringFlag{
						Name:        "tracker-url",
						Description: "Tracker base url to use to collect related issue from.",
					},
					components.StringFlag{
						Name:        "tracker-username",
						Description: "Tracker username to use to collect related issue from.",
					},
					components.StringFlag{
						Name:        "tracker-token",
						Description: "Tracker token to use to collect related issue from.",
					},
					components.StringFlag{
						Name:        "regexp",
						Description: "A regular expression used for matching the git commit messages.",
					},
					components.StringFlag{
						Name:         "key-group-index",
						Description:  "The capturing group index in the regular expression used for retrieving the issue key.",
						DefaultValue: "1",
					},
					components.StringFlag{
						Name:         "git-log-limit",
						Description:  "The maximum number of git commit messages to process.",
						DefaultValue: "100",
					},
					components.BoolFlag{
						Name:         "aggregate",
						Description:  "Set to true, if you wish all builds to include issues from previous builds.",
						DefaultValue: false,
					},
					components.StringFlag{
						Name: "aggregation-status",
						Description: "If aggregate is set to true, " +
							"this property indicates how far in time should the issues be aggregated. In the above example, " +
							"issues will be aggregated from previous builds, until a build with a RELEASE status is found. " +
							"Build statuses are set when a build is promoted using the jf rt build-promote command.",
					},
				},
				Arguments: []components.Argument{
					{
						Name:        "build name",
						Description: "The name of the build.",
					},
					{
						Name:        "build number",
						Description: "The number of the build.",
					},
					{
						Name:        "path to .git",
						Description: "Path to a directory containing the .git directory. If not specified, the .git directory is assumed to be in the current directory or in one of the parent directories.",
					},
				},
				Action: func(c *components.Context) error {
					return collectIssuesCmd(c)
				},
			},
			{
				Name:        "clean-slate",
				Description: "Clean the build-info directory to start from a clean slate.",
				Aliases:     []string{"cs"},
				Flags: []components.Flag{
					components.StringFlag{
						Name:        "project",
						Description: "Artifactory project key.",
					},
				},
				Arguments: []components.Argument{
					{
						Name:        "build name",
						Description: "The name of the build.",
					},
					{
						Name:        "build number",
						Description: "The number of the build.",
					},
				},
				Action: func(c *components.Context) error {
					return cleanSlateCmd(c)
				},
			},
			{
				Name:        "send-build-info",
				Description: "Send build-info to Jira",
				Aliases:     []string{"sbi"},
				Flags: []components.Flag{
					components.StringFlag{
						Name:        "server-id",
						Description: "Server ID configured using the config command.",
					},
					components.StringFlag{
						Name:        "project",
						Description: "Artifactory project key.",
					},
					components.StringFlag{
						Name:        "jira-id",
						Description: "Jira integration name.",
					},
					components.StringFlag{
						Name:        "jira-url",
						Description: "Jira base url.",
					},
					components.StringFlag{
						Name:        "jira-client-id",
						Description: "The OAuth clientId generated by Jira.",
					},
					components.StringFlag{
						Name:        "jira-secret",
						Description: "The OAuth secret generated by Jira.",
					},
					components.BoolFlag{
						Name:         "dry-run",
						Description:  "Enable to only log what would be send to Jira.",
						DefaultValue: false,
					},
					components.BoolFlag{
						Name:         "include-pre-post-runs",
						Description:  "Enable to include pipeline preRun and postRun steps.",
						DefaultValue: false,
					},
					components.BoolFlag{
						Name:         "fail-on-reject",
						Description:  "Enable to error out if any builds are rejected by Jira.",
						DefaultValue: false,
					},
				},
				Arguments: []components.Argument{
					{
						Name:        "build name",
						Description: "The name of the build.",
					},
					{
						Name:        "build number",
						Description: "The number of the build.",
					},
				},
				Action: func(c *components.Context) error {
					return sendBuildInfoCmd(c)
				},
			},
			{
				Name:        "send-deployment-info",
				Description: "Send deployment-info to Jira",
				Aliases:     []string{"sdi"},
				Flags: []components.Flag{
					components.StringFlag{
						Name:        "server-id",
						Description: "Server ID configured using the config command.",
					},
					components.StringFlag{
						Name:        "project",
						Description: "Artifactory project key.",
					},
					components.StringFlag{
						Name:        "jira-id",
						Description: "Jira integration name.",
					},
					components.StringFlag{
						Name:        "jira-url",
						Description: "Jira base url.",
					},
					components.StringFlag{
						Name:        "jira-client-id",
						Description: "The OAuth clientId generated by Jira.",
					},
					components.StringFlag{
						Name:        "jira-secret",
						Description: "The OAuth secret generated by Jira.",
					},
					components.BoolFlag{
						Name:         "dry-run",
						Description:  "Enable to only log what would be send to Jira.",
						DefaultValue: false,
					},
					components.BoolFlag{
						Name:         "include-pre-post-runs",
						Description:  "Enable to include pipeline preRun and postRun steps.",
						DefaultValue: false,
					},
					components.BoolFlag{
						Name:         "fail-on-reject",
						Description:  "Enable to error out if any builds are rejected by Jira.",
						DefaultValue: false,
					},
					components.StringFlag{
						Name:        "environment",
						Description: "The environment that the deployment targeted.",
					},
				},
				Arguments: []components.Argument{
					{
						Name:        "build name",
						Description: "The name of the build.",
					},
					{
						Name:        "build number",
						Description: "The number of the build.",
					},
				},
				Action: func(c *components.Context) error {
					return sendDeploymentInfoCmd(c)
				},
			},
			{
				Name:        "notify-slack",
				Description: "Send build-info to Slack",
				Aliases:     []string{"ns"},
				Flags: []components.Flag{
					components.StringFlag{
						Name:        "server-id",
						Description: "Server ID configured using the config command.",
					},
					components.StringFlag{
						Name:        "project",
						Description: "Artifactory project key.",
					},
					components.StringFlag{
						Name:        "slack",
						Description: "Slack integration name.",
					},
					components.BoolFlag{
						Name:         "include-pre-post-runs",
						Description:  "Enable to include pipeline preRun and postRun steps.",
						DefaultValue: false,
					},
				},
				Arguments: []components.Argument{
					{
						Name:        "build name",
						Description: "The name of the build.",
					},
					{
						Name:        "build number",
						Description: "The number of the build.",
					},
				},
				Action: func(c *components.Context) error {
					return notifySlackCmd(c)
				},
			},
		},
	})
}

func collectIssuesCmd(c *components.Context) error {
	nargs := len(c.Arguments)
	if nargs > 3 {
		return errors.New(fmt.Sprintf("Wrong number of arguments (%d).", nargs))
	}
	buildConfiguration := CreateBuildConfiguration(c)
	if err := buildConfiguration.ValidateBuildParams(); err != nil {
		return err
	}

	issueConfiguration, err := CreateIssueConfiguration(c)
	if err != nil {
		return err
	}
	if err := issueConfiguration.ValidateIssueConfiguration(); err != nil {
		return err
	}

	collectIssueCommand := commands.NewCollectIssueCommand().SetBuildConfiguration(buildConfiguration).SetIssuesConfig(issueConfiguration)
	if nargs == 3 {
		collectIssueCommand.SetDotGitPath(c.Arguments[2])
	} else if nargs == 1 {
		collectIssueCommand.SetDotGitPath(c.Arguments[0])
	}
	return collectIssueCommand.Run()
}

func cleanSlateCmd(c *components.Context) error {
	nargs := len(c.Arguments)
	if nargs > 2 {
		return errors.New(fmt.Sprintf("Wrong number of arguments (%d).", nargs))
	}
	buildConfiguration := CreateBuildConfiguration(c)
	if err := buildConfiguration.ValidateBuildParams(); err != nil {
		return err
	}

	cleanslateCommand := commands.NewCleanSlateCommand().SetBuildConfiguration(buildConfiguration)
	return cleanslateCommand.Run()
}

func sendBuildInfoCmd(c *components.Context) error {
	nargs := len(c.Arguments)
	if nargs > 2 {
		return errors.New(fmt.Sprintf("Wrong number of arguments (%d).", nargs))
	}
	buildConfiguration := CreateBuildConfiguration(c)
	if err := buildConfiguration.ValidateBuildParams(); err != nil {
		return err
	}

	jiraConfiguration := CreateJiraConfiguration(c)
	if err := jiraConfiguration.ValidateJiraConfiguration(); err != nil {
		return err
	}

	sendBuildInfoCommand := commands.NewSendBuildInfoCommand().SetBuildConfiguration(buildConfiguration).SetJiraConfiguration(jiraConfiguration)
	return sendBuildInfoCommand.Run()
}

func sendDeploymentInfoCmd(c *components.Context) error {
	nargs := len(c.Arguments)
	if nargs > 2 {
		return errors.New(fmt.Sprintf("Wrong number of arguments (%d).", nargs))
	}
	buildConfiguration := CreateBuildConfiguration(c)
	if err := buildConfiguration.ValidateBuildParams(); err != nil {
		return err
	}

	jiraConfiguration := CreateJiraConfiguration(c)
	if err := jiraConfiguration.ValidateJiraConfiguration(); err != nil {
		return err
	}

	deploymentInfo, err := CreateDeploymentInfo(c)
	if err != nil {
		return err
	}
	sendDeploymentInfoCommand := commands.NewSendDeploymentInfoCommand().SetBuildConfiguration(buildConfiguration).SetJiraConfiguration(
		jiraConfiguration).SetDeploymentInfo(deploymentInfo)
	return sendDeploymentInfoCommand.Run()
}

func notifySlackCmd(c *components.Context) error {
	nargs := len(c.Arguments)
	if nargs > 2 {
		return errors.New(fmt.Sprintf("Wrong number of arguments (%d).", nargs))
	}
	buildConfiguration := CreateBuildConfiguration(c)
	if err := buildConfiguration.ValidateBuildParams(); err != nil {
		return err
	}

	slackConfiguration := CreateSlackConfiguration(c)
	if err := slackConfiguration.ValidateSlackConfiguration(); err != nil {
		return err
	}

	notifySlackCommand := commands.NewNotifySlackCommand().SetBuildConfiguration(buildConfiguration).SetSlackConfiguration(slackConfiguration)
	return notifySlackCommand.Run()
}

func CreateBuildConfiguration(c *components.Context) *artifactoryUtils.BuildConfiguration {
	buildConfiguration := new(artifactoryUtils.BuildConfiguration)
	buildNameArg, buildNumberArg := "", ""
	if len(c.Arguments) >= 2 {
		buildNameArg, buildNumberArg = c.Arguments[0], c.Arguments[1]
	}
	if buildNameArg == "" || buildNumberArg == "" {
		buildNameArg = ""
		buildNumberArg = ""
	}
	buildConfiguration.SetBuildName(buildNameArg).SetBuildNumber(buildNumberArg).SetProject(c.GetStringFlagValue("project"))
	return buildConfiguration
}

func CreateIssueConfiguration(c *components.Context) (*commands.IssuesConfiguration, error) {
	issueConfiguration := new(commands.IssuesConfiguration)
	issueConfiguration.SetServerID(c.GetStringFlagValue("server-id"))
	issueConfiguration.SetTracker(c.GetStringFlagValue("tracker"))
	if url := c.GetStringFlagValue("tracker-url"); url != "" {
		issueConfiguration.SetTrackerDetails(url, c.GetStringFlagValue("tracker-username"), c.GetStringFlagValue("tracker-token"))
	}
	issueConfiguration.SetRegexp(c.GetStringFlagValue("regexp"))
	if index := c.GetStringFlagValue("key-group-index"); index != "" {
		groupIndex, err := strconv.Atoi(index)
		if err != nil {
			return nil, err
		}
		issueConfiguration.SetKeyGroupIndex(groupIndex)
	}
	issueConfiguration.SetAggregate(c.GetBoolFlagValue("aggregate"))
	issueConfiguration.SetAggregationStatus(c.GetStringFlagValue("aggregation-status"))
	return issueConfiguration, nil
}

func CreateJiraConfiguration(c *components.Context) *commands.JiraConfiguration {
	jiraConfiguration := new(commands.JiraConfiguration)
	jiraConfiguration.SetServerID(c.GetStringFlagValue("server-id"))
	jiraConfiguration.SetJiraID(c.GetStringFlagValue("jira-id"))
	if url := c.GetStringFlagValue("jira-url"); url != "" {
		jiraConfiguration.SetJiraDetails(url, c.GetStringFlagValue("jira-client-id"), c.GetStringFlagValue("jira-secret"))
	}
	jiraConfiguration.SetDryRun(c.GetBoolFlagValue("dry-run"))
	jiraConfiguration.SetIncludePrePostRunSteps(c.GetBoolFlagValue("include-pre-post-runs"))
	jiraConfiguration.SetFailOnReject(c.GetBoolFlagValue("fail-on-reject"))
	return jiraConfiguration
}

func CreateDeploymentInfo(c *components.Context) (*commands.DeploymentInfo, error) {
	environment := c.GetStringFlagValue("environment")
	if environment == "" {
		environment = os.Getenv("environmentName")
	}
	if environment == "" {
		return nil, errorutils.CheckErrorf("Missing deployment environment")
	}
	return commands.NewDeploymentInfo(environment), nil
}

func CreateSlackConfiguration(c *components.Context) *commands.SlackConfiguration {
	slackConfiguration := new(commands.SlackConfiguration)
	slackConfiguration.SetServerID(c.GetStringFlagValue("server-id"))
	slackConfiguration.SetSlack(c.GetStringFlagValue("slack"))
	slackConfiguration.SetIncludePrePostRunSteps(c.GetBoolFlagValue("include-pre-post-runs"))
	slackConfiguration.SetFailOnReject(c.GetBoolFlagValue("fail-on-reject"))
	return slackConfiguration
}
