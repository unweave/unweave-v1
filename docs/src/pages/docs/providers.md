---
title: Providers
pageTitle: Providers - Unweave
description: Uniform API across cloud providers
---


Unweave lets you choose where your code runs when you create a new SSH session. This means you can
pick a VM across cloud providers using a unified interface.
Currently, in addition to the Unweave Hosted Platform, we support [Lambda Labs](https://lambdalabs.com),
and [Google Cloud](https://cloud.google.com) (in private beta). We're working on adding a number of
additional providers in the near future. [See here](https://github.com/unweave/unweave#roadmap) for the Roadmap.

---

## Unweave Provider

Every Unweave account comes with the Unweave provider enabled by default. Any resources you create are 
charged by the minute.


## Configure the Lambda Labs Provider

To configure Lambda Labs, you'll need to create a Lambda Labs account and generate an API key.
Once you have the key, head over to your Unweave [account settings](https://app.unweave.io/settings) and
click on `connect` next to the Lambda Labs provider.

![](./images/ll-connect.png)


To test that your provider is configured correctly, you can try listing the available VMs
on Lambda Labs from the CLI:


```bash
unweave ls-node-types --provider lambdalabs
```


## Configure the Google Cloud Provider

{% callout %}
This feature is currently is private beta. If  you'd like early access, reach out to us on
[Discord](https://discord.gg/ydyVHbFjPt), [Twitter](https://twitter.com/unweaveio), or via
[email](mailto:info@unweave.io)
{% /callout %}

## Next Steps

- [Create your first project](./projects)
