# Unweave
<img align="right" src="https://unweave.io/favicon.svg" height="150px" style="margin: 2rem 1rem" alt="Unweave logo">

Unweave is [Supabase](https://supabase.com) for machine learning. It's an open source tool
for creating and managing development environments for ML projects. Read the  [docs](https://docs.unweave.io) to see how it works.
### Features

- Setup and SSH into dev environments on cloud VMs with single command
- Uniform API across cloud providers
- Automatic search for available GPU instances
- Metadata tracked in Postgres
- SSH keypair and credentials management
- Plugin based architecture for adding new cloud providers (or local environments)

### Why?

Machine learning has seen incredible growth over the last few years. Unfortunately, 
tooling for ML hasn't caught up yet. Although there are a lot of great open-source tools 
for ML, they often require endless YAML-file tweaking to setup and configure. You shouldn't
have to manually configure VPCs and IAM policies just for some extra compute crunch.

The goal of Unweave is use the best in class open-source tools and make them "just work" for
the entire ML lifecycle.


### Overview

Unweave is composed of three parts: the platform, the [CLI](https://github.com/unweave/cli),
and the Dashboard. The easiest way to get started is using the hosted platform at 
[unweave.io](https://app.unweave.io). You can also [self-host or develop locally](#self-hosting-and-local-development).

### Installation

**Homebrew (Mac):**

```bash
brew tap unweave/unweave
brew install unweave
```

**Linux:**

Currently, you'll need to extract the package from the [latest release](https://github.com/unweave/cli/releases)
and install it with the package manager for your platform.


### Getting Started

Login to the CLI:
```bash
unweave login
```

**Create a new session:**

VM instances in Unweave are called Sessions. You can start a new Session on a supported 
cloud provider by using the `--provider` flag. The default provider is Unweave.

```bash
unweave new --project <project-id> --provider lambdalabs --type gpu_1x_a10
```

You can then ssh into the Session you create as usual. You can find the host by running 
`unweave ls`.

```bash
ssh ubuntu@<host>
```

[//]: # ()
[//]: # (**One shot ssh:**)

[//]: # ()
[//]: # (You can also combine the create and ssh steps into a single command:)

[//]: # (```bash)

[//]: # (unweave ssh --create --provider lambdalabs)

[//]: # (```)

**Linking a project:**

You can avoid manually adding the provider and project flags each time by linking to an Unweave project. 
This will initialize an `unweave` folder with a `config.toml` file. You can add you default
provider and session preferences here.

```bash
unweave link <project-id>
unweave new
```


### Self-hosting and Local Development

The [docker-compose.yml](./docker-compose.yml) file contains all the services needed to run
the Unweave platform. You'll need to use the [.env.example](./.env.example) when running 
docker compose to configure the environment.

To use the CLI with the platform, you'll need to set the `UNWEAVE_ENV=dev` variable for the
CLI.

The self-hosted platform does not include any authentication or project handling. 


### Getting Help

- [Documentation](https://docs.unweave.io)
- [Discord](https://discord.gg/ydyVHbFjPt)
- [Twitter](https://twitter.com/intent/follow?screen_name=unweaveio)
- [GitHub Issues](https://github.com/unweave/unweave/issues)

### Roadmap

Our goal is to provide a "works by default" _**dev → production**_ platform for machine 
learning teams that builds on top of existing open-source tooling. All parts of Unweave 
are or will also be open-source and modular, so you can continue using it even if your 
requirements change. 

We're starting with a unified API for versions ML dev environments in VMs built on top of 
SSH and Git.  

- [ ] VMs for ML dev environments
  - [ ] Providers
    - [x] [LambdaLabs](https://lambdalabs.com)
    - [x] [Unweave](https://unweave.io)
    - [x] [GCP](https://cloud.google.com) _(private beta)_
    - [ ] [DigitalOcean](https://digitalocean.com)
    - [ ] [AWS](https://aws.amazon.com)
    - [ ] [Azure](https://azure.microsoft.com)
    - [ ] Localhost
- [ ] Custom Docker images
- [ ] Serverless executions `unweave exec python train.py`
- [ ] Blob Storage
- [ ] GitHub Integration
- [ ] Notebooks
- [ ] Deployments
  - [ ] CPU only
  - [ ] GPU

### Contributing

Contributions are welcome! However, we're still very early and the codebase and API are
likely to change a lot. 

The best way to start is developing is to follow the [local dev](#self-hosting-and-local-development)
section above.
