---
title: "Command - devspace use context"
sidebar_label: context
---


Tells DevSpace which kube context to use

## Synopsis


```
devspace use context [flags]
```

```
#######################################################
############### devspace use context ##################
#######################################################
Switch the current kube context

Example:
devspace use context my-context
#######################################################
```
## Options

```
  -h, --help   help for context
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

* [devspace use](../../cli/commands/devspace_use)	 - Use specific config
