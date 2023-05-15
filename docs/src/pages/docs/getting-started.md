---
title: Getting Started
description: Getting started with Unweave
---

{% callout %}
We'll be updating this page regularly, however, if you need help more urgently or would like to get 
involved, reach out to us on [Discord](https://discord.gg/ydyVHbFjPt), [Twitter](https://twitter.com/unweaveio), or via
[email](mailto:info@unweave.io)
{% /callout %}

You're ready to train your model. There's only one piece missing, the GPU to train on. With Unweave, you're 2 mins away to start training your model. 

## Install the cli
```bash
$ brew tap unweave/unweave
$ brew install unweave
```

## Login
Let's first login: 

```bash
$ unweave login
```

## Create a new project

{% callout %}
In the next version of our cli, you'll be able to create projects from your command line directly!
{% /callout %}

Navigate to the Unweave dashboard by: 

```bash
$ unweave open
```
You'll see a big blue button to create your project. Go ahead and create it!

## Link your project

Now that you have a project created, we'll go ahead and link it to our directory: 

```bash
$ unweave link your-username/your-project-name
```

## Launch VSCode 

Last step is to launch VSCode in our brand new GPU powered machine!

```bash
unweave code --new --type rtx_5000 --image pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel
```
You're done 🚀. You now have VSCode with all the files in your repo, and you can get started. 

Note: This will launch your exact same repository but in a GPU powered machine. Here, we are using the `rtx_5000` gpu type, and the `pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel` docker image as our base. 

## Cleaning up

Once you're done using the machine, remember to power it down: 
```bash
$ unweave terminate
```

## More options

Our CLI is pretty powerful, explore more options with: 
```
$ unweave --help
```

---

## Next Steps

- [Learn more about providers](./providers)
- [Create you first project](./projects)
- [Launch and SSH into a GPU VM](./sessions)
