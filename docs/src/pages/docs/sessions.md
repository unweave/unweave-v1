---
title: Sessions
pageTitle: Sessions - Unweave
description: Unweave Sessions
---

A Session in Unweave is the lowest level resource. It is equivalent to a Virtual Machine (Node) that
Unweave spins up and monitors for you. Each session runs on one of the providers you have configured
in your Unweave account.

{% callout %}
Sessions are charged by the minute. You'll get a notification when your credit falls below $1. If
the credit falls below $0, the session will automatically be terminated.
{% /callout %}

To create a new session from the CLI, run:

```bash
unweave session create --project <project-id> --provider <provider> --type <node-type>
```

This will prompt you for an SSH public key to use for the session. You can either choose to use 
your existing `id_rsa.pub` key or create a new `.pem` key. Check the [SSH](#ssh-keys) section 
for more details.

Or, if you have a project linked to your local directory, you can skip the flags. Unweave will
automatically use the defaults configured in the `config.toml` file.

```bash
unweave session create
```

## SSH Keys

You need to provide an SSH public key when you create a Session on Unweave. You can add 
an existing key to your Unweave account by running:

```bash
# Name is optional. 
unweave ssh-key add <path-to-public-key> [name]
```

You can also generate a new SSH key pair by running:

```bash
unweave ssh-key generate [name]
```

When you create a new session, you must provide the key either by name (if you've added it to
your account) or by the path to the public key on your local machine. Unweave will 
automatically provision the key on each provider before spinning up the Session.

## Node Types

Node types are the different types of VMs you can spin up. You can list the available node types
for a provider using the `ls-node-types` command:

```bash
unweave ls-node-types --provider <provider>
```

Currently, only GPU nodes are supported. CPU nodes are coming soon.

---

## Next Steps

- [SSH into a Session](./ssh)
