---
title: NanoGPT on Unweave
description: Training a NanoGPT model on the unweave platform 
---

## 1. 🚀 Kickstart Your Journey with OpenAI Whisper and Unweave

Welcome aboard! 🎉 You're about to embark on an exciting journey with [OpenAI Whisper](https://openai.com/research/whisper/), a supercharged automatic speech recognition (ASR) system, and Unweave, your trusty sidekick for all things machine learning. Whisper is a multilingual maestro, trained on a vast dataset of diverse audio, capable of speech recognition, translation, and language identification. 🌍

In this guide, we'll show you how to transcribe a podcast episode using Whisper and Unweave. You'll learn how to set up your environment, prepare your audio file, and transcribe like a pro! So buckle up and let's dive right in! 🏊‍♂️

## 2. 🛠 Setting Up Your Workspace

Before we start transcribing, let's get your workspace ready. Follow these steps to ensure a smooth setup process:

- First, install the Unweave Command Line Interface (CLI) by running these commands:

```shell
brew tap unweave/unweave
brew install unweave
```

- Next, log in to your Unweave account:

```shell
unweave login
```

- Now, let's create a new project on Unweave:

```shell
unweave open
```

- Click the "Create Project" button and you're halfway there! 🎉
- Link your local directory to the Unweave project:

```shell
unweave link your-username/your-project-name
```

Don't forget to replace `your-username` and `your-project-name` with your actual username and project name.

- Finally, install the latest release of Whisper:

```shell
pip install -U openai-whisper
```

- Make sure `ffmpeg` is installed on your system. If not, install it using your system's package manager.

Voila! Your workspace is all set and ready to go! 🚀

## 3. 🎙 Preparing Your Audio File

Now, let's get your audio file ready for transcription:

- Create a new machine with GPU support in the Unweave platform:

```bash
unweave code --new --type rtx_5000 --image pytorch/pytorch:2.0.0-cuda11.7-cudnn8-devel
```

- Make sure your audio file is on the machine (e.g., `ls`).

And just like that, your audio file is ready for the spotlight! 🎬

## 4. 📝 Transcribing Your Podcast Episode

With your audio file prepped and ready, it's time to let Whisper work its magic:

- Start the transcription process:

```bash
whisper your_audio_file.wav --model medium
```

Remember to replace `your_audio_file.wav` with your actual audio file name.

- If your audio file is not in English, no worries! Whisper speaks multiple languages. Just specify the language:

```bash
whisper your_audio_file.wav --model medium --language Japanese
```

And there you have it! You've just transcribed a podcast episode using OpenAI Whisper on the Unweave platform. 🎉

## 5. 🧹 Cleaning Up and Wrapping Up

After all the fun, it's time to clean up:

- Check the list of active Unweave sessions:

```bash
unweave ls
```

- Terminate any running sessions:

```bash
unweave terminate
```

And that's a wrap

! 🎬 By following these steps, you've kept your project files tidy and freed up resources on the Unweave platform.

---

In this tutorial, we've journeyed through the process of transcribing a podcast episode using OpenAI Whisper and Unweave. We've set up the environment, prepared the audio file, and performed transcription. Finally, we've learned how to clean up and terminate sessions on the Unweave platform.

Now, you're all set to leverage the power of OpenAI Whisper and Unweave for your own transcription tasks. So go forth and transcribe! 🚀 Happy transcribing!