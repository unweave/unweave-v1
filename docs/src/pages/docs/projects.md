---
title: Projects
description: Unweave Projects
---

Resources in Unweave are scoped to a Project. You can think of a Project at the same level as a
GitHub repository. You can create a new project from the Dashboard.

[]()

## Linking a Project

To avoid having to pass the `--project` flag to every command from the CLI, you can link a project to a local
directory. This will create a `unweave` subdirectory with a `config.toml` and `.env` file. 

```bash
unweave link <username>/<project-name> [path]
```

Any command you run from the linked directory or any subdirectory will automatically use the
project ID from the config file.

## SSH Keys

SSH keys are used to setup access to the VMs created with Unweave. These keys are shared across 
projects. You can add an SSH key from your local machine either through the CLI or the Dashboard.

To add it from the CLI, you need to provide the path to the public key. You can also optionally
provide a name for the key.

```bash
unweave ssh-keys add <public-key-path> [name]
```


---

## Next Steps

- [Create a new Session](./sessions)
