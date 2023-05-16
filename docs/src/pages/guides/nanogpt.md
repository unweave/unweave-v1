---
title: NanoGPT on Unweave
description: Training a NanoGPT model on the unweave platform 
---
## 1.Introduction to NanoGPT and Unweave

[NanoGPT](https://github.com/karpathy/nanoGPT) is a powerful repository for training and fine-tuning medium-sized GPTs. Unweave is a platform that provides a seamless environment for machine learning tasks. In this guide, we will walk you through the process of training NanoGPT models using Unweave. 

By following this guide, you will learn how to set up your environment, prepare the dataset, configure the fine-tuning process, initiate training on Unweave, save the fine-tuned model, generate samples, and manage your resources effectively. Let's get started and explore the capabilities of NanoGPT on the Unweave platform!

##  2.Setting Up the Environment

Before we begin training NanoGPT models on the Unweave platform, it's important to set up the environment properly. Follow these steps to ensure a smooth setup process:

- Clone the NanoGPT repository from GitHub by running the following command in your desired directory:

```shell
git clone https://github.com/karpathy/nanoGPT
```

- Change your current directory to the cloned NanoGPT repository:

```shell
cd nanoGPT
```

- Install the Unweave Command Line Interface (CLI) by running the following commands:

```shell
brew tap unweave/unweave
brew install unweave
```

- Log in to your Unweave account using the command:

```shell
unweave login
```

- Create a new project on Unweave by navigating to the Unweave dashboard:

```shell
unweave open
```

- Click the "Create Project" button to set up your project.
- Link your local directory to the Unweave project using the command:

```shell
unweave link your-username/nanoGPT
```

Replace `your-username` and `your-project-name` with your actual username.

Once you have completed these setup steps, you'll be ready to start training NanoGPT models on the Unweave platform. Let's move on to preparing the dataset.

##  3.Preparing the Dataset

To prepare the dataset for training NanoGPT, follow these steps:

- Create a new machine with GPU support in the Unweave platform:

```bash
unweave code --new --type rtx_5000 --image pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel
```

- Ensure your NanoGPT repository is on the machine (e.g., `ls`).
- Navigate to the NanoGPT repository:

```bash
cd NanoGPT
```

- Install the required [dependencies](https://github.com/karpathy/nanoGPT#install) mentioned in the NanoGPT repository:

```bash
python -m pip install --upgrade pip
pip install torch numpy transformers datasets tiktoken wandb tqdm
```

- Prepare the dataset:

```bash
python data/shakespeare_char/prepare.py
```

Make sure that `train.bin` and `val.bin` are created in the `data/shakespeare` directory.

By following these steps, you will set up the necessary environment, ensure your repository is available, install dependencies, and prepare the dataset for training NanoGPT.

##  4.Fine-tuning and Inference

With the dataset prepared, we can now proceed with fine-tuning and generating text using NanoGPT:

- Start the fine-tuning process:

```bash
python train.py config/finetune_shakespeare.py
```

- To monitor your GPU usage during training, you can use the following command:

```bash
nvidia-smi
```

- Once the training is complete, you can generate Shakespearean text using the fine-tuned model:

```bash
sample.py --out_dir=out-shakespeare
```

Congratulations! You have successfully trained NanoGPT on the Unweave platform and generated text using the fine-tuned model.

##  5.Cleaning Up and Terminating Sessions

After you have finished working with your trained model, it's important to clean up and terminate any running sessions. Follow these steps:

- Zip your fine-tuned model:

```bash
zip out-shakespeare
```

- Move the model directory to your local machine by dragging the zip file to your desktop or desired location.

- Check the list of active Unweave sessions:

```bash
unweave ls
```

- Terminate any running sessions:

```bash
unweave terminate
```

By following these steps, you can organize your project files and free up resources on the Unweave platform.

---

In this tutorial, we covered the process of training NanoGPT models using Unweave. We set up the environment, prepared the dataset, performed fine-tuning, and generated text. Finally, we learned how to clean up and terminate sessions on the Unweave platform.

Now you can leverage the power of NanoGPT and Unweave for your own machine learning tasks. Happy training!