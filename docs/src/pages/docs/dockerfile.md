---
title: Dockerfile (beta)
pageTitle: Dockerfile - Unweave
description: Custom Environment with Dockerfile
---

{% callout %}
This feature is currently is private beta. If  you'd like early access, reach out to us on
[Discord](https://discord.gg/ydyVHbFjPt), [Twitter](https://twitter.com/unweaveio), or via
[email](mailto:info@unweave.io)
{% /callout %}

By default, Sessions on Unweave run on a standard Ubuntu 20.04 image customized for the machine
learning stack. If you need to customize the environment further, you can use a Dockerfile to 
specify the base image to use for the Session.

You can enable custom images by passing the `--x-dockerfile` experimental flag when creating a new
Session. Unweave will detect a `Dockerfile` at the root of your project and use it 
to build the Session image.

```bash
unweave new --ex-dockerfile
```

Images are automatically built and uploaded to the Unweave registry and tagged with the `sessionID`.
