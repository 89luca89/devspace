package v3

import (
	devspacecontext "github.com/loft-sh/devspace/pkg/devspace/context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/loft-sh/devspace/pkg/devspace/config/versions/latest"
	"github.com/loft-sh/devspace/pkg/devspace/helm/generic"
	"github.com/loft-sh/devspace/pkg/devspace/helm/types"
	"github.com/loft-sh/devspace/pkg/util/command"
	"github.com/loft-sh/devspace/pkg/util/downloader/commands"
	"github.com/loft-sh/devspace/pkg/util/git"
	"github.com/loft-sh/devspace/pkg/util/log"
)

type client struct {
	exec        command.Exec
	genericHelm generic.Client
}

// NewClient creates a new helm v3 Client
func NewClient(log log.Logger) (types.Client, error) {
	c := &client{
		exec: command.NewStreamCommand,
	}

	c.genericHelm = generic.NewGenericClient(c, log)
	return c, nil
}

func (c *client) Command() commands.Command {
	return commands.NewHelmV3Command()
}

// InstallChart installs the given chart via helm v3
func (c *client) InstallChart(ctx *devspacecontext.Context, releaseName string, releaseNamespace string, values map[string]interface{}, helmConfig *latest.HelmConfig) (*types.Release, error) {
	valuesFile, err := c.genericHelm.WriteValues(values)
	if err != nil {
		return nil, err
	}
	defer os.Remove(valuesFile)

	if releaseNamespace == "" {
		releaseNamespace = ctx.KubeClient.Namespace()
	}

	args := []string{
		"upgrade",
		releaseName,
		"--namespace",
		releaseNamespace,
		"--values",
		valuesFile,
		"--install",
	}

	// Chart settings
	if helmConfig.Chart.Git != nil {
		chartName, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}

		defer os.RemoveAll(chartName)
		repo, err := git.NewGitCLIRepository(chartName)
		if err != nil {
			return nil, err
		}
		err = repo.Clone(git.CloneOptions{
			URL:    helmConfig.Chart.Git.URL,
			Branch: helmConfig.Chart.Git.Branch,
			Tag:    helmConfig.Chart.Git.Tag,
			Args:   helmConfig.Chart.Git.CloneArgs,
			Commit: helmConfig.Chart.Git.Revision,
		})
		if err != nil {
			return nil, err
		}
		if helmConfig.Chart.Git.SubPath != "" {
			chartName = filepath.Join(chartName, helmConfig.Chart.Git.SubPath)
		}
		args = append(args, chartName)
	} else {
		chartName, chartRepo := generic.ChartNameAndRepo(helmConfig)
		args = append(args, chartName)
		if chartRepo != "" {
			args = append(args, "--repo", chartRepo)
			args = append(args, "--repository-config=''")
		}
		if helmConfig.Chart.Version != "" {
			args = append(args, "--version", helmConfig.Chart.Version)
		}
		if helmConfig.Chart.Username != "" {
			args = append(args, "--username", helmConfig.Chart.Username)
		}
		if helmConfig.Chart.Password != "" {
			args = append(args, "--password", helmConfig.Chart.Password)
		}
	}

	// Upgrade options
	if helmConfig.Atomic {
		args = append(args, "--atomic")
	}
	if helmConfig.CleanupOnFail {
		args = append(args, "--cleanup-on-fail")
	}
	if helmConfig.Wait {
		args = append(args, "--wait")
	}
	if helmConfig.Timeout != "" {
		args = append(args, "--timeout", helmConfig.Timeout)
	}
	if helmConfig.Force {
		args = append(args, "--force")
	}
	if helmConfig.DisableHooks {
		args = append(args, "--no-hooks")
	}
	args = append(args, helmConfig.UpgradeArgs...)
	output, err := c.genericHelm.Exec(ctx, args, helmConfig)

	if helmConfig.DisplayOutput {
		_, _ = ctx.Log.Writer(logrus.InfoLevel).Write(output)
	}

	if err != nil {
		return nil, err
	}

	releases, err := c.ListReleases(ctx, helmConfig)
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		if r.Name == releaseName && r.Namespace == releaseNamespace {
			return r, nil
		}
	}

	return nil, nil
}

func (c *client) Template(ctx *devspacecontext.Context, releaseName, releaseNamespace string, values map[string]interface{}, helmConfig *latest.HelmConfig) (string, error) {
	cleanup, chartDir, err := c.genericHelm.FetchChart(ctx, helmConfig)
	if err != nil {
		return "", err
	} else if cleanup {
		defer os.RemoveAll(filepath.Dir(chartDir))
	}

	if releaseNamespace == "" {
		releaseNamespace = ctx.KubeClient.Namespace()
	}

	valuesFile, err := c.genericHelm.WriteValues(values)
	if err != nil {
		return "", err
	}
	defer os.Remove(valuesFile)

	args := []string{
		"template",
		releaseName,
		chartDir,
		"--namespace",
		releaseNamespace,
		"--values",
		valuesFile,
	}
	args = append(args, helmConfig.TemplateArgs...)
	result, err := c.genericHelm.Exec(ctx, args, helmConfig)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (c *client) DeleteRelease(ctx *devspacecontext.Context, releaseName string, releaseNamespace string, helmConfig *latest.HelmConfig) error {
	if releaseNamespace == "" {
		releaseNamespace = ctx.KubeClient.Namespace()
	}

	args := []string{
		"delete",
		releaseName,
		"--namespace",
		releaseNamespace,
	}
	_, err := c.genericHelm.Exec(ctx, args, helmConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) ListReleases(ctx *devspacecontext.Context, helmConfig *latest.HelmConfig) ([]*types.Release, error) {
	args := []string{
		"list",
		"--namespace",
		ctx.KubeClient.Namespace(),
		"--max",
		strconv.Itoa(0),
		"--output",
		"json",
	}
	out, err := c.genericHelm.Exec(ctx, args, helmConfig)
	if err != nil {
		return nil, err
	}

	releases := []*types.Release{}
	err = yaml.Unmarshal(out, &releases)
	if err != nil {
		return nil, err
	}

	return releases, nil
}
