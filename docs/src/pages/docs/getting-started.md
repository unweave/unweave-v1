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
brew tap unweave/unweave
brew install unweave
```

## Login
Next, log in to your Unweave account from the terminal using:

```bash
unweave login
```

## Create a new project

{% callout %}
In the next version of our CLI, you'll be able to create projects from your command line directly!
{% /callout %}

To create a new project, navigate to the Unweave dashboard by running the command:

```bash
unweave open
```

Click the big blue button to create your project.

## Link your project

Once you've created your project, you need to link it to your local directory. Linking a 
project associates a directory on your computer to the Unweave project in your account. Once 
linked, any commands you run in this directory, or any of its subdirectories, will be run
in the context of the project.

```bash
unweave link <unweave-username>/<project-name>
```

## Launch VSCode 

Unweave ships with a super handy command to launch VS Code in the context of your project.

```bash
unweave code --new --type rtx_5000 --image pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel
```

This  will launch VS Code, setup SSH access, and sync your local directory onto the VM at 
the `/home/unweave` path.  The example above uses a VM with the RTX 5000 GPU type and 
the `pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel` Docker image as the base.

Congratulations! You're now ready to start working with your GPU-accelerated environment in VS Code.

## Listing Sessions

You can check the status of any running sessions using the `unweave ls` command: 

```bash
unweave ls
```

This command will display a list of all your running sessions, along with relevant 
information such as session ID, machine type. Remember to always clean up and terminate 
your sessions when you are finished to avoid unnecessary charges.

## Cleaning up

Once you are done using the GPU-powered machine, remember to power it down to avoid 
unnecessary costs. Run the following command:

```bash
unweave terminate
```

## More options

The Unweave CLI offers various powerful features. To explore more options, use the following command:

```
unweave --help
```

---

## Next Steps

- [Learn more about providers](./providers)
- [Create you first project](./projects)
- [Launch and SSH into a GPU VM](./sessions)
