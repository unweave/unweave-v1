---
title: Getting Started
description: Getting started with Unweave
---

{% callout %}
We'll be updating this page regularly, however, if you need help more urgently or would like to get 
involved, reach out to us on [Discord](https://discord.gg/ydyVHbFjPt), [Twitter](https://twitter.com/unweaveio), or via
[email](mailto:info@unweave.io)
{% /callout %}


## Install the cli
To begin, you need to install the Unweave Command Line Interface (CLI) by following these steps:

```bash
$ brew tap unweave/unweave
$ brew install unweave
```

## Login
After installing the CLI, you need to log in to your Unweave account. Use the following command:

```bash
$ unweave login
```

## Create a new project

{% callout %}
In the next version of our cli, you'll be able to create projects from your command line directly!
{% /callout %}

To create a new project, navigate to the Unweave dashboard by running the command:


```bash
$ unweave open
```

Click the big blue button to create your project.

## Link your project

Once you have created your project, you need to link it to your local directory. Use the following command:

```bash
$ unweave link your-username/your-project-name
```

## Launch VSCode 

The next step is to launch VSCode in your GPU-powered machine. Run the following command:


```bash
$ unweave code --new --type rtx_5000 --image pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel
```
This command will launch VSCode with all the files in your repository, but on a GPU-powered machine. The example above uses the rtx_5000 GPU type and the pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel Docker image as the base.

Congratulations! You're now ready to start working with your GPU-accelerated environment in VSCode.

## Listing Sessions

After you have finished your work and powered down the GPU-powered machine, you may want to check the status of any running sessions. You can do this using the unweave ls command. Simply run the following command:

```bash
$ unweave ls
```

This command will display a list of all your running sessions, along with relevant information such as session ID, machine type. Remember to always clean up and terminate your sessions when you are finished to avoid unnecessary charges.

## Cleaning up

Once you are done using the GPU-powered machine, remember to power it down to avoid unnecessary costs. Run the following command:

```bash
$ unweave terminate
```

## More options

The Unweave CLI offers various powerful features. To explore more options, use the following command:

```
$ unweave --help
```

---

## Next Steps

- [Learn more about providers](./providers)
- [Create you first project](./projects)
- [Launch and SSH into a GPU VM](./sessions)
