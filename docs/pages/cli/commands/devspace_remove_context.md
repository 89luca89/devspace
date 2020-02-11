---
title: "Command - devspace remove context"
sidebar_label: context
---


Removes a kubectl-context

## Synopsis


```
devspace remove context [flags]
```

```
#######################################################
############# devspace remove context #################
#######################################################
Removes a kubectl-context

Example:
devspace remove context myspace
devspace remove context --all-spaces
#######################################################
```
## Options

```
      --all-spaces        Remove all kubectl contexts belonging to Spaces
  -h, --help              help for context
      --provider string   The cloud provider to use
```

### Options inherited from parent commands

```
      --debug                 Prints the stack trace if an error occurs
      --kube-context string   The kubernetes context to use
  -n, --namespace string      The kubernetes namespace to use
      --no-warn               If true does not show any warning when deploying into a different namespace or kube-context than before
  -p, --profile string        The devspace profile to use (if there is any)
      --silent                Run in silent mode and prevents any devspace log output except panics & fatals
  -s, --switch-context        Switches and uses the last kube context and namespace that was used to deploy the DevSpace project
      --var strings           Variables to override during execution (e.g. --var=MYVAR=MYVALUE)
```

## See Also

* [devspace remove](../../cli/commands/devspace_remove)	 - Changes devspace configuration
