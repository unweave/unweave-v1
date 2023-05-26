---
title: NanoGPT on Unweave
description: Training a NanoGPT model on the unweave platform 
---
## 1.Introduction to NanoGPT and Unweave

[NanoGPT](https://github.com/karpathy/nanoGPT) is a powerful repository for training and fine-tuning 
tiny GPTs, mostly for instructional purposes. In this guide, we'll walk you through the process of 
training NanoGPT models using Unweave. 

By following this guide, you'll learn how to set up your environment, prepare the dataset, 
configure the fine-tuning process, initiate training on Unweave, save the fine-tuned model, 
generate samples, and manage your resources effectively. 

##  2.Setting Up the Environment

We'll first setup the NanoGPT repository on our local machine and link it to a new project on Unweave.

- Clone the NanoGPT repository from GitHub:

```shell
git clone https://github.com/karpathy/nanoGPT
```

- Change your current directory to the cloned NanoGPT repository:

```shell
cd nanoGPT
```

- If you haven't already, install the Unweave Command Line Interface (CLI):

```shell
brew tap unweave/unweave
brew install unweave
```

- Log in to your Unweave account on the CLI. This will redirect you to the Unweave dashboard in your browser:

```shell
unweave login
```

- Create a new project on Unweave by navigating to the Unweave dashboard:

```shell
unweave open
```

- Click the "Create Project" button to set up your project. Call it `nanoGPT`
- Link your local directory to the Unweave project using the command:

```shell
unweave link <your-username>/nanoGPT
```

Replace `your-username` with your Unweave username above.

We're now ready to start training NanoGPT models. Let's move on to preparing the dataset.

##  3.Preparing the Dataset

We'll do all the data prep and training on Unweave using the Pytorch 2.0 Docker image:

- Create a new machine with GPU support using the Unweave CLI. In this case, we don't need a huge machine, 
  so we'll use the `rtx_5000` GPU type. Since the pytorch image is quite large, this might take 1-2 minutes 
  to start-up:

```bash
unweave code --gpu-type rtx_5000 --image pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel
```

- This should land you in the `/home/unweave` directory with your code copied over. Make sure it is
 (e.g., `ls`). 

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

##  4.Fine-tuning and Inference

With the dataset prepared, we can now proceed to fine-tuning and generating text using NanoGPT:

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

Congratulations! You have successfully trained your NanoGPT to talk Shakespeare! 

##  5.Cleaning Up and Terminating Sessions

After you have finished working with your trained model, it's important to clean up and terminate 
any running sessions. Follow these steps:

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

---

In this tutorial, we covered the process of training NanoGPT models using Unweave. We set up the 
environment, prepared the dataset, performed fine-tuning, and generated text. Finally, we learned 
how to clean up and terminate sessions on the Unweave.
