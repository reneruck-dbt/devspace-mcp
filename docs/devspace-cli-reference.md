# DevSpace CLI Reference

This document contains the help output for all DevSpace CLI commands relevant to the MCP server implementation.

## Table of Contents

- [Main Command](#main-command)
- [List Commands](#list-commands)
- [Development Commands](#development-commands)
- [Build & Deploy Commands](#build--deploy-commands)
- [Debugging Commands](#debugging-commands)
- [Utility Commands](#utility-commands)

---

## Main Command

### devspace --help

```
DevSpace accelerates developing, deploying and debugging applications with Docker and Kubernetes. Get started by running the init command in one of your projects:

		devspace init
		# Develop an existing application
		devspace dev
		DEVSPACE_CONFIG=other-config.yaml devspace dev

Usage:
  devspace [command]

Available Commands:
  add          Adds something to devspace.yaml
  analyze      Analyzes a kubernetes namespace and checks for potential problems
  attach       Attaches to a container
  build        Builds all defined images and pushes them
  cleanup      Cleans up resources
  completion   Outputs shell completion for the given shell (bash or zsh)
  deploy       Deploys the project
  dev          Starts the development mode
  enter        Open a shell to a container
  init         Initializes DevSpace in the current folder
  list         Lists configuration
  logs         Prints the logs of a pod and attaches to it
  open         Opens the space in the browser
  print        Prints displays the configuration
  purge        Deletes deployed resources
  remove       Removes devspace configuration
  render       Builds all defined images and shows the yamls that would be deployed
  reset        Resets an cluster token
  restart      Restarts containers where the sync restart helper is injected
  run          Executes a predefined command
  run-pipeline Starts a DevSpace pipeline
  set          Sets global configuration changes
  sync         Starts a bi-directional sync between the target container and the local path
  ui           Opens the localhost UI in the browser
  update       Updates the current config
  upgrade      Upgrades the DevSpace CLI to the newest version
  use          Uses specific config
  version      Prints version of devspace

Flags:
      --debug                        Prints the stack trace if an error occurs
      --disable-profile-activation   If true will ignore all profile activations
  -h, --help                         help for devspace
      --inactivity-timeout int       Minutes the current user is inactive (no mouse or keyboard interaction) until DevSpace will exit automatically. 0 to disable. Only supported on windows and mac operating systems
      --kube-context string          The kubernetes context to use
      --kubeconfig string            The kubeconfig path to use
  -n, --namespace string             The kubernetes namespace to use
      --no-colors                    Do not show color highlighting in log output. This avoids invisible output with different terminal background colors
      --no-warn                      If true does not show any warning when deploying into a different namespace or kube-context than before
      --override-name string         If specified will override the DevSpace project name provided in the devspace.yaml
  -p, --profile strings              The DevSpace profiles to apply. Multiple profiles are applied in the order they are specified
      --silent                       Run in silent mode and prevents any devspace log output except panics & fatals
  -s, --switch-context               Switches and uses the last kube context and namespace that was used to deploy the DevSpace project
      --var strings                  Variables to override during execution (e.g. --var=MYVAR=MYVALUE)
  -v, --version                      version for devspace
```

### Global Flags (Available on all commands)

```
      --debug                        Prints the stack trace if an error occurs
      --disable-profile-activation   If true will ignore all profile activations
      --inactivity-timeout int       Minutes the current user is inactive until DevSpace will exit automatically
      --kube-context string          The kubernetes context to use
      --kubeconfig string            The kubeconfig path to use
  -n, --namespace string             The kubernetes namespace to use
      --no-colors                    Do not show color highlighting in log output
      --no-warn                      If true does not show any warning when deploying into a different namespace or kube-context
      --override-name string         If specified will override the DevSpace project name provided in the devspace.yaml
  -p, --profile strings              The DevSpace profiles to apply. Multiple profiles are applied in the order they are specified
      --silent                       Run in silent mode and prevents any devspace log output except panics & fatals
  -s, --switch-context               Switches and uses the last kube context and namespace that was used to deploy the DevSpace project
      --var strings                  Variables to override during execution (e.g. --var=MYVAR=MYVALUE)
```

---

## List Commands

### devspace list --help

```
Usage:
  devspace list [command]

Available Commands:
  commands    Lists all custom DevSpace commands
  contexts    Lists all kube contexts
  deployments Lists and shows the status of all deployments
  namespaces  Lists all namespaces in the current context
  plugins     Lists all installed devspace plugins
  ports       Lists port forwarding configurations
  profiles    Lists all DevSpace profiles
  sync        Lists sync configuration
  vars        Lists the vars in the active config
```

### devspace list commands --help

```
Lists all DevSpace custom commands defined in the devspace.yaml

Usage:
  devspace list commands [flags]

Flags:
  -h, --help   help for commands
```

### devspace list deployments --help

```
Lists the status of all deployments

Usage:
  devspace list deployments [flags]

Flags:
  -h, --help   help for deployments
```

### devspace list ports --help

```
Lists the port forwarding configurations

Usage:
  devspace list ports [flags]

Flags:
  -h, --help            help for ports
  -o, --output string   The output format of the command. Can be either empty or json
```

### devspace list sync --help

```
Lists the sync configuration

Usage:
  devspace list sync [flags]

Flags:
  -h, --help   help for sync
```

---

## Development Commands

### devspace dev --help

```
Starts your project in development mode

Usage:
  devspace dev [flags]

Flags:
      --build-sequential            Builds the images one after another instead of in parallel
      --dependency strings          Deploys only the specified named dependencies
  -b, --force-build                 Forces to build every image
  -d, --force-deploy                Forces to deploy every deployment
      --force-purge                 Forces to purge every deployment even though it might be in use by another DevSpace project
  -h, --help                        help for dev
      --max-concurrent-builds int   The maximum number of image builds built in parallel (0 for infinite)
      --pipeline string             The pipeline to execute (default "dev")
      --render                      If true will render manifests and print them instead of actually deploying them
      --sequential-dependencies     If set set true dependencies will run sequentially
      --show-ui                     Shows the ui server
      --skip-build                  Skips building of images
      --skip-dependency strings     Skips the following dependencies for deployment
      --skip-deploy                 If enabled will skip deploying
      --skip-push                   Skips image pushing, useful for minikube deployment
      --skip-push-local-kube        Skips image pushing, if a local kubernetes environment is detected (default true)
  -t, --tag strings                 Use the given tag for all built images
```

### devspace sync --help

```
Starts a bi-directional(default) sync between the target container path
and local path:

devspace sync --path=.:/app # localPath is current dir and remotePath is /app
devspace sync --path=.:/app --image-selector nginx:latest
devspace sync --path=.:/app --exclude=node_modules,test
devspace sync --path=.:/app --pod=my-pod --container=my-container

Usage:
  devspace sync [flags]

Flags:
  -c, --container string           Container name within pod where to sync to
      --download-on-initial-sync   DEPRECATED: Downloads all locally non existing remote files in the beginning (default true)
      --download-only              If set DevSpace will only download files
  -e, --exclude strings            Exclude directory from sync
  -h, --help                       help for sync
      --image-selector string      The image to search a pod for (e.g. nginx, nginx:latest, ${runtime.images.app}, nginx:${runtime.images.app.tag})
      --initial-sync string        The initial sync strategy to use (mirrorLocal, mirrorRemote, preferLocal, preferRemote, preferNewest, keepAll)
  -l, --label-selector string      Comma separated key=value selector list (e.g. release=test)
      --no-watch                   Synchronizes local and remote and then stops
      --path string                Path to use (Default is current directory). Example: ./local-path:/remote-path or local-path:.
      --pick                       Select a pod (default true)
      --pod string                 Pod to sync to
      --polling                    If polling should be used to detect file changes in the container
      --upload-only                If set DevSpace will only upload files
      --wait                       Wait for the pod(s) to start if they are not running (default true)
```

### devspace restart --help

```
Restarts containers where the sync restart helper is injected:

devspace restart
devspace restart -n my-namespace

Usage:
  devspace restart [flags]

Flags:
  -c, --container string        Container name within pod to restart
  -h, --help                    help for restart
  -l, --label-selector string   Comma separated key=value selector list (e.g. release=test)
      --name string             The sync path name to restart
      --pick                    Select a pod (default true)
      --pod string              Pod to restart
```

---

## Build & Deploy Commands

### devspace build --help

```
Builds all defined images and pushes them

Usage:
  devspace build [flags]

Flags:
      --build-sequential            Builds the images one after another instead of in parallel
      --dependency strings          Deploys only the specified named dependencies
  -b, --force-build                 Forces to build every image (default true)
  -d, --force-deploy                Forces to deploy every deployment
      --force-purge                 Forces to purge every deployment even though it might be in use by another DevSpace project
  -h, --help                        help for build
      --max-concurrent-builds int   The maximum number of image builds built in parallel (0 for infinite)
      --pipeline string             The pipeline to execute (default "build")
      --render                      If true will render manifests and print them instead of actually deploying them
      --sequential-dependencies     If set set true dependencies will run sequentially
      --show-ui                     Shows the ui server
      --skip-build                  Skips building of images
      --skip-dependency strings     Skips the following dependencies for deployment
      --skip-deploy                 If enabled will skip deploying
      --skip-push                   Skips image pushing, useful for minikube deployment
      --skip-push-local-kube        Skips image pushing, if a local kubernetes environment is detected (default true)
  -t, --tag strings                 Use the given tag for all built images
```

### devspace deploy --help

```
Deploys the current project to a Space or namespace:

devspace deploy
devspace deploy -n some-namespace
devspace deploy --kube-context=deploy-context

Usage:
  devspace deploy [flags]

Flags:
      --build-sequential            Builds the images one after another instead of in parallel
      --dependency strings          Deploys only the specified named dependencies
  -b, --force-build                 Forces to build every image
  -d, --force-deploy                Forces to deploy every deployment
      --force-purge                 Forces to purge every deployment even though it might be in use by another DevSpace project
  -h, --help                        help for deploy
      --max-concurrent-builds int   The maximum number of image builds built in parallel (0 for infinite)
      --pipeline string             The pipeline to execute (default "deploy")
      --render                      If true will render manifests and print them instead of actually deploying them
      --sequential-dependencies     If set set true dependencies will run sequentially
      --show-ui                     Shows the ui server
      --skip-build                  Skips building of images
      --skip-dependency strings     Skips the following dependencies for deployment
      --skip-deploy                 If enabled will skip deploying
      --skip-push                   Skips image pushing, useful for minikube deployment
      --skip-push-local-kube        Skips image pushing, if a local kubernetes environment is detected (default true)
  -t, --tag strings                 Use the given tag for all built images
```

### devspace purge --help

```
Deletes the deployed kubernetes resources:

devspace purge

Usage:
  devspace purge [flags]

Flags:
      --build-sequential            Builds the images one after another instead of in parallel
      --dependency strings          Deploys only the specified named dependencies
  -b, --force-build                 Forces to build every image
  -d, --force-deploy                Forces to deploy every deployment
      --force-purge                 Forces to purge every deployment even though it might be in use by another DevSpace project
  -h, --help                        help for purge
      --max-concurrent-builds int   The maximum number of image builds built in parallel (0 for infinite)
      --pipeline string             The pipeline to execute (default "purge")
      --render                      If true will render manifests and print them instead of actually deploying them
      --sequential-dependencies     If set set true dependencies will run sequentially
      --show-ui                     Shows the ui server
      --skip-build                  Skips building of images
      --skip-dependency strings     Skips the following dependencies for deployment
      --skip-deploy                 If enabled will skip deploying
      --skip-push                   Skips image pushing, useful for minikube deployment
      --skip-push-local-kube        Skips image pushing, if a local kubernetes environment is detected (default true)
  -t, --tag strings                 Use the given tag for all built images
```

### devspace run --help

```
Executes a predefined command from the devspace.yaml

Examples:
devspace run mycommand --myarg 123
devspace run mycommand2 1 2 3
devspace --dependency my-dependency run any-command --any-command-flag

Usage:
  devspace run [flags]

Flags:
      --dependency string   Run a command from a specific dependency
  -h, --help                help for run
```

### devspace run-pipeline --help

```
Execute a pipeline:
devspace run-pipeline my-pipeline
devspace run-pipeline dev

Usage:
  devspace run-pipeline [flags]

Flags:
      --build-sequential            Builds the images one after another instead of in parallel
      --dependency strings          Deploys only the specified named dependencies
  -b, --force-build                 Forces to build every image
  -d, --force-deploy                Forces to deploy every deployment
      --force-purge                 Forces to purge every deployment even though it might be in use by another DevSpace project
  -h, --help                        help for run-pipeline
      --max-concurrent-builds int   The maximum number of image builds built in parallel (0 for infinite)
      --pipeline string             The pipeline to execute
      --render                      If true will render manifests and print them instead of actually deploying them
      --sequential-dependencies     If set set true dependencies will run sequentially
      --show-ui                     Shows the ui server
      --skip-build                  Skips building of images
      --skip-dependency strings     Skips the following dependencies for deployment
      --skip-deploy                 If enabled will skip deploying
      --skip-push                   Skips image pushing, useful for minikube deployment
      --skip-push-local-kube        Skips image pushing, if a local kubernetes environment is detected (default true)
  -t, --tag strings                 Use the given tag for all built images
```

---

## Debugging Commands

### devspace analyze --help

```
Analyze checks a namespaces events, replicasets, services
and pods for potential problems

Example:
devspace analyze
devspace analyze --namespace=mynamespace

Usage:
  devspace analyze [flags]

Flags:
  -h, --help                  help for analyze
      --ignore-pod-restarts   If true, analyze will ignore the restart events of running pods
      --patient               If true, analyze will ignore failing pods and events until every deployment, statefulset, replicaset and pods are ready or the timeout is reached
      --timeout int           Timeout until analyze should stop waiting (default 120)
      --wait                  Wait for pods to get ready if they are just starting (default true)
```

### devspace logs --help

```
Prints the last log of a pod container and attachs to it

Example:
devspace logs
devspace logs --namespace=mynamespace

Usage:
  devspace logs [flags]

Flags:
  -c, --container string        Container name within pod where to execute command
  -f, --follow                  Attach to logs afterwards
  -h, --help                    help for logs
      --image-selector string   The image to search a pod for (e.g. nginx, nginx:latest, ${runtime.images.app}, nginx:${runtime.images.app.tag})
  -l, --label-selector string   Comma separated key=value selector list (e.g. release=test)
      --lines int               Max amount of lines to print from the last log (default 200)
      --pick                    Select a pod (default true)
      --pod string              Pod to print the logs of
      --wait                    Wait for the pod(s) to start if they are not running
```

### devspace enter --help

```
Execute a command or start a new terminal in your devspace:

devspace enter
devspace enter --pick # Select pod to enter
devspace enter bash
devspace enter -c my-container
devspace enter bash -n my-namespace
devspace enter bash -l release=test
devspace enter bash --image-selector nginx:latest
devspace enter bash --image-selector "${runtime.images.app.image}:${runtime.images.app.tag}"

Usage:
  devspace enter [flags]

Flags:
  -c, --container string        Container name within pod where to execute command
  -h, --help                    help for enter
      --image-selector string   The image to search a pod for (e.g. nginx, nginx:latest, ${runtime.images.app}, nginx:${runtime.images.app.tag})
  -l, --label-selector string   Comma separated key=value selector list (e.g. release=test)
      --pick                    Select a pod / container if multiple are found (default true)
      --pod string              Pod to open a shell to
      --reconnect               Will reconnect the terminal if an unexpected return code is encountered
      --screen                  Use a screen session to connect
      --screen-session string   The screen session to create or connect to (default "enter")
      --tty                     If to use a tty to start the command (default true)
      --wait                    Wait for the pod(s) to start if they are not running
      --workdir string          The working directory where to open the terminal or execute the command
```

### devspace attach --help

```
Attaches to a running container

devspace attach
devspace attach --pick # Select pod to enter
devspace attach -c my-container
devspace attach -n my-namespace

Usage:
  devspace attach [flags]

Flags:
  -c, --container string        Container name within pod where to execute command
  -h, --help                    help for attach
      --image-selector string   The image to search a pod for (e.g. nginx, nginx:latest, ${runtime.images.app}, nginx:${runtime.images.app.tag})
  -l, --label-selector string   Comma separated key=value selector list (e.g. release=test)
      --pick                    Select a pod (default true)
      --pod string              Pod to open a shell to
```

---

## Utility Commands

### devspace print --help

```
Prints the configuration for the current or given
profile after all patching and variable substitution

Usage:
  devspace print [flags]

Flags:
      --dependency string   The dependency to print the config from. Use dot to access nested dependencies (e.g. dep1.dep2)
  -h, --help                help for print
      --skip-info           When enabled, only prints the configuration without additional information
```

### devspace version --help

```
Prints version of devspace

Usage:
  devspace version [flags]

Flags:
  -h, --help   help for version
```

---

## Notes for MCP Implementation

### Commands Suitable for MCP (Non-Interactive)

| Command | Notes |
|---------|-------|
| `devspace version` | Simple, no arguments needed |
| `devspace list *` | All list subcommands work well |
| `devspace print` | Returns YAML config |
| `devspace analyze` | Returns analysis results |
| `devspace logs` | Use `--lines` to limit output, avoid `--follow` |
| `devspace build` | Long-running, needs extended timeout |
| `devspace deploy` | Long-running, needs extended timeout |
| `devspace purge` | Deletes resources |
| `devspace run` | Runs predefined commands |
| `devspace run-pipeline` | Runs custom pipelines |
| `devspace enter <cmd>` | With `--tty=false` for non-interactive exec |

### Commands NOT Suitable for MCP (Interactive)

| Command | Reason |
|---------|--------|
| `devspace dev` | Runs continuously, requires terminal |
| `devspace attach` | Requires interactive terminal |
| `devspace sync` | Runs continuously (unless `--no-watch`) |
| `devspace enter` (shell) | Requires interactive terminal |

### Key Flags for MCP Usage

- `--no-colors` - Always use to get clean output
- `--silent` - Use when only result matters
- `--no-warn` - Suppress non-essential warnings
- `-o json` - Use when available for structured output (e.g., `list ports`)
