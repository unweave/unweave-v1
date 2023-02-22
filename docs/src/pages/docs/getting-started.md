---
title: Getting Started
description: Getting started with Unweave
---

{% callout %}
We'll be updating this page regularly, however, if you need help more urgently or would like to get 
involved, reach out to us on [Discord](https://discord.gg/ydyVHbFjPt), [Twitter](https://twitter.com/unweaveio), or via
[email](mailto:info@unweave.io)
{% /callout %}

## Dashboard

The Unweave dashboard gives you an overview of your account and projects. You should first  [create
an account and login](https://app.unweave.io) before moving on to the CLI.

## Installation

### Homebrew (Mac)

```bash
brew tap unweave/unweave
brew install unweave
```

### Linux

Currently, you'll need to extract the package from the [latest release](https://github.com/unweave/cli/releases) and install it with the package manager for your platform.
Once you've downloaded the `.apk/.deb/.rpm/.tar.gz` file, depending on your package manager, you can run one of the following commands:

```bash
# APK
sudo apk add --allow-untrusted <...>.apk
# DEB
sudo dpkg -i <...>.deb
# RPM
sudo rpm -i <...>.rpm
```

### Login to the CLI


```bash
unweave login
```




---

## Next Steps

- [Learn more about providers](./providers)
- [Create you first project](./projects)
- [Launch and SSH into a GPU VM](./sessions)
