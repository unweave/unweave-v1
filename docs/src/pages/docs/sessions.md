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

Or, if you have a project linked to your local directory, you can skip the flags. Unweave will
automatically use the defaults configured in the `config.toml` file.

```bash
unweave session create
```


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
