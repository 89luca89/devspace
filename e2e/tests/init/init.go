package init

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/loft-sh/devspace/pkg/devspace/config/loader/variable"

	"github.com/loft-sh/devspace/cmd"
	"github.com/loft-sh/devspace/cmd/flags"
	"github.com/loft-sh/devspace/e2e/framework"
	"github.com/loft-sh/devspace/e2e/kube"
	"github.com/loft-sh/devspace/pkg/devspace/config/versions/latest"
	"github.com/loft-sh/devspace/pkg/util/survey"
	"github.com/onsi/ginkgo"
	"gopkg.in/yaml.v2"
)

var _ = DevSpaceDescribe("init", func() {
	initialDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// create a new factory
	var (
		f          *framework.DefaultFactory
		kubeClient *kube.KubeHelper
	)

	ginkgo.BeforeEach(func() {
		f = framework.NewDefaultFactory()

		kubeClient, err = kube.NewKubeHelper()
		framework.ExpectNoError(err)
	})

	ginkgo.It("should create devspace.yml without registry details", func() {
		tempDir, err := framework.CopyToTempDir("tests/init/testdata/new")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		// set the question answer func here
		f.SetAnswerFunc(func(params *survey.QuestionOptions) (string, error) {
			if strings.Contains(params.Question, "Which registry would you want to use to push images to?") {
				return "Skip Registry", nil
			}

			return params.DefaultValue, nil
		})

		initCmd := &cmd.InitCmd{}
		err = initCmd.Run(f)
		framework.ExpectNoError(err)

		config, _, err := framework.LoadConfig(f, filepath.Join(tempDir, "devspace.yaml"))
		framework.ExpectNoError(err)

		framework.ExpectEqual(len(config.Variables()), 1+len(variable.AlwaysResolvePredefinedVars))
		framework.ExpectEqual(config.Variables()["IMAGE"], "username/app")

		ns, err := kubeClient.CreateNamespace("init")
		framework.ExpectNoError(err)
		defer framework.ExpectDeleteNamespace(kubeClient, ns)

		done := make(chan error)
		go func() {
			devCmd := &cmd.DevCmd{
				GlobalFlags: &flags.GlobalFlags{
					NoWarn:    true,
					Namespace: ns,
				},
			}
			done <- devCmd.Run(f, []string{"sh", "-c", "exit 0"})
		}()

		err = <-done
		framework.ExpectNoError(err)
	})

	ginkgo.It("should create devspace.yml without registry details and manifests deploy", func() {
		tempDir, err := framework.CopyToTempDir("tests/init/testdata/new")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		// set the question answer func here
		f.SetAnswerFunc(func(params *survey.QuestionOptions) (string, error) {
			if strings.Contains(params.Question, "Which registry would you want to use to push images to?") {
				return "Skip Registry", nil
			}

			if strings.Contains(params.Question, "How do you want to deploy this project?") {
				return cmd.ManifestsOption, nil
			}

			if strings.Contains(params.Question, "Please enter the paths to your Kubernetes manifests") {
				return "manifests/**", nil
			}

			return params.DefaultValue, nil
		})

		initCmd := &cmd.InitCmd{}
		err = initCmd.Run(f)
		framework.ExpectNoError(err)

		config, _, err := framework.LoadConfig(f, filepath.Join(tempDir, "devspace.yaml"))
		framework.ExpectNoError(err)

		framework.ExpectEqual(len(config.Variables()), 1+len(variable.AlwaysResolvePredefinedVars))
		framework.ExpectEqual(config.Variables()["IMAGE"], "username/app")

		ns, err := kubeClient.CreateNamespace("init")
		framework.ExpectNoError(err)
		defer framework.ExpectDeleteNamespace(kubeClient, ns)

		done := make(chan error)
		go func() {
			devCmd := &cmd.DevCmd{
				GlobalFlags: &flags.GlobalFlags{
					NoWarn:    true,
					Namespace: ns,
				},
			}
			done <- devCmd.Run(f, []string{"sh", "-c", "exit 0"})
		}()

		err = <-done
		framework.ExpectNoError(err)
	})

	ginkgo.It("should create devspace.yml without registry details and kustomize deploy", func() {
		tempDir, err := framework.CopyToTempDir("tests/init/testdata/new")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		// set the question answer func here
		f.SetAnswerFunc(func(params *survey.QuestionOptions) (string, error) {
			if strings.Contains(params.Question, "Which registry would you want to use to push images to?") {
				return "Skip Registry", nil
			}

			if strings.Contains(params.Question, "How do you want to deploy this project?") {
				return cmd.KustomizeOption, nil
			}

			if strings.Contains(params.Question, "Please enter path to your Kustomization folder") {
				return "./kustomization", nil
			}

			return params.DefaultValue, nil
		})

		initCmd := &cmd.InitCmd{}
		err = initCmd.Run(f)
		framework.ExpectNoError(err)

		config, _, err := framework.LoadConfig(f, filepath.Join(tempDir, "devspace.yaml"))
		framework.ExpectNoError(err)

		framework.ExpectEqual(len(config.Variables()), 1+len(variable.AlwaysResolvePredefinedVars))
		framework.ExpectEqual(config.Variables()["IMAGE"], "username/app")

		ns, err := kubeClient.CreateNamespace("init")
		framework.ExpectNoError(err)
		defer framework.ExpectDeleteNamespace(kubeClient, ns)

		done := make(chan error)
		go func() {
			devCmd := &cmd.DevCmd{
				GlobalFlags: &flags.GlobalFlags{
					NoWarn:    true,
					Namespace: ns,
				},
			}
			done <- devCmd.Run(f, []string{"sh", "-c", "exit 0"})
		}()

		err = <-done
		framework.ExpectNoError(err)
	})

	ginkgo.It("should create devspace.yml without registry details and local helm chart deploy", func() {
		tempDir, err := framework.CopyToTempDir("tests/init/testdata/new")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		// set the question answer func here
		f.SetAnswerFunc(func(params *survey.QuestionOptions) (string, error) {
			if strings.Contains(params.Question, "Which registry would you want to use to push images to?") {
				return "Skip Registry", nil
			}

			if strings.Contains(params.Question, "How do you want to deploy this project?") {
				return cmd.HelmChartOption, nil
			}

			if strings.Contains(params.Question, "Which Helm chart do you want to use?") {
				return `Use a local Helm chart (e.g. ./helm/chart/)`, nil
			}

			if strings.Contains(params.Question, "Please enter the relative path to your local Helm chart") {
				return "./chart", nil
			}

			return params.DefaultValue, nil
		})

		initCmd := &cmd.InitCmd{}
		err = initCmd.Run(f)
		framework.ExpectNoError(err)

		config, _, err := framework.LoadConfig(f, filepath.Join(tempDir, "devspace.yaml"))
		framework.ExpectNoError(err)

		framework.ExpectEqual(len(config.Variables()), 1+len(variable.AlwaysResolvePredefinedVars))
		framework.ExpectEqual(config.Variables()["IMAGE"], "username/app")

		ns, err := kubeClient.CreateNamespace("init")
		framework.ExpectNoError(err)
		defer framework.ExpectDeleteNamespace(kubeClient, ns)

		done := make(chan error)
		go func() {
			devCmd := &cmd.DevCmd{
				GlobalFlags: &flags.GlobalFlags{
					NoWarn:    true,
					Namespace: ns,
				},
			}
			done <- devCmd.Run(f, []string{"sh", "-c", "exit 0"})
		}()

		err = <-done
		framework.ExpectNoError(err)
	})

	ginkgo.It("should create devspace.yml from docker-compose.yaml", func() {
		tempDir, err := framework.CopyToTempDir("tests/init/testdata/docker-compose")
		framework.ExpectNoError(err)
		defer framework.CleanupTempDir(initialDir, tempDir)

		ns, err := kubeClient.CreateNamespace("init")
		framework.ExpectNoError(err)
		defer func() {
			err := kubeClient.DeleteNamespace(ns)
			framework.ExpectNoError(err)
		}()

		// Answer all questions with the default
		f.SetAnswerFunc(func(params *survey.QuestionOptions) (string, error) {
			return params.DefaultValue, nil
		})

		initCmd := &cmd.InitCmd{
			Reconfigure: true,
		}
		err = initCmd.Run(f)
		framework.ExpectNoError(err)

		// Created a devspace.yaml
		_, _, err = framework.LoadConfig(f, filepath.Join(tempDir, "devspace.yaml"))
		framework.ExpectNoError(err)

		// Created a .gitignore
		_, err = os.Stat(filepath.Join(tempDir, ".gitignore"))
		framework.ExpectNoError(err)

		// Created a .devspace/generated.yaml
		_, err = os.Stat(filepath.Join(tempDir, ".devspace", "generated.yaml"))
		framework.ExpectNoError(err)

		// Print the config to verify the expected deployment
		var configBuffer bytes.Buffer
		printCmd := &cmd.PrintCmd{
			GlobalFlags: &flags.GlobalFlags{
				NoWarn: true,
				Debug:  true,
			},
			Out: &configBuffer,
		}

		err = printCmd.Run(f)
		framework.ExpectNoError(err)

		generatedConfig := &latest.Config{}
		err = yaml.Unmarshal(configBuffer.Bytes(), generatedConfig)
		framework.ExpectNoError(err)

		// validate config
		framework.ExpectEqual(len(generatedConfig.Deployments), 1)
		framework.ExpectEqual(generatedConfig.Deployments[0].Name, "db")

		// ensure valid configuration by deploying the application
		deployCmd := &cmd.DeployCmd{
			GlobalFlags: &flags.GlobalFlags{
				NoWarn:    true,
				Debug:     true,
				Namespace: ns,
			},
			SkipPush: true,
			Wait:     true,
			Timeout:  120,
		}

		err = deployCmd.Run(f)
		framework.ExpectNoError(err)
	}, 120)
})
