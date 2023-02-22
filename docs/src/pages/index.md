---
title: About
pageTitle: About - Unweave
description: ML without the FML.
---

Unweave is an open source platform for creating and managing development environments for
machine learning projects.

{% callout %}
This page is an overview of the motivation behind Unweave. If you're eager to jump straight in, 
you can skip to the [Getting Started](./docs/getting-started) page.
{% /callout %}

## Why?

Machine learning has seen incredible, almost explosive, growth over the last few years. Unfortunately, 
tooling for ML hasn't caught up yet, and doesn't get the same amount of love as web or mobile dev, for instance.
Although there are a lot of great open-source tools for ML, they often require endless
YAML-file tweaking to setup and configure. The problem is that ML is a unique flavor of software development. 
It's more akin to running a science experiment than the usual feature centric software workflow. Current
ML tools don't quite fit under the `git-commit-push-deploy` umbrella.

Unweave aims to change that. Our goal is to bring the same incredible developer experience you get 
when working on a web or mobile app, to the world of machine learning.

An essential asterisk in the ML workflow is its unusual resource hungriness. While you can get up-and-running
with a React app on you grandad's tablet, for ML, you'll most likely need a beefy GPU. That's a spanner
in the works, because it just so happens that GPUs are in short supply.

Therefore, an essential first step in _**Unweaving**_ 🙃 the cloud VM/IAM/VPC spaghetti, is to make access 
to compute _much, much,_ easier. Unweave does this with [Sessions](./docs/sessions.md) by abstracting 
away the boilerplate around spinning up GPU VMs and lands you directly in an SSH terminal on an instance of your choice.

We do this by building from the ground up in a cloud agnostic way. Unweave has foundational support for
[Providers](./docs/providers.md), a simple interface that allows you to run Unweave on any platform 
that implements it. Currently, we support [Lambda Labs](https://lambdalabs.com), and [Google Cloud](https://cloud.google.com) (in private
beta), in addition to the hosted platform.

Lastly, since it sucks to keep reinventing the wheel, [Unweave is open source](https://github.com/unweave/unweave).
Our core philosophy is to build on top of or integrate with the best in class open-source tools.


## Architecture

A essential goal of Unweave is to track all the metadata required to make ML workflows fully reproducible.
It stores this metadata in a Postgres database. 

