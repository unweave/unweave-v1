---
title: FAQ
description: Frequently Asked Questions
---

{% callout %}
If you need addition support, feel free to reach out to us on 
[Discord](https://discord.gg/ydyVHbFjPt), [Twitter](https://twitter.com/unweaveio), or via
[email](mailto:info@unweave.io)
{% /callout %}

**Q. What cloud providers do you support?**

_A._ Currently, we support [Lambda Labs](https://lambdalabs.com), [GCP - private beta](htts://cloud.google.com), and
   the Unweave hosted platform.

**Q. How much do you mark up the cost of the cloud provider? How do you make money?**

_A._ Zero! We don't mark up the cost of the cloud provider. The margins are too thin anyway. We make 
   money through our hosted SaaS service that will eventually include a team/pro pricing tier.

**Q. What is the average start-up time for a new Session?**

_A._ The start-up time depends on the provider. For the hosted version, it is currently ~20s.

**Q. How can I configure the RAM and CPUs for a Session on the Unweave Provider?**

_A._ Currently, each node comes with 4vCPUs and 64GB of RAM. We are working on making this
   configurable in the near future.

**Q. How can I configure the Python version, dependencies, etc. ?**

_A._ Currently, the Lambda Labs provider runs the [Lambda stack](https://lambdalabs.com/lambda-stack-deep-learning-software). The Unweave provider
    runs on `Ubuntu 20.04 with Python 3.10 inside a Conda environment`. You can install any
    dependencies you need using `conda`. Custom Dockerfiles are in private beta. Get in touch
    if you'd like early access.
    

**Q. Is my Session data persisted?**

_A._ No. Not yet. We're working on adding persistence but for the moment, if you terminate 
    a Session, all data will be erased. **Make sure to persist any files you care about to
    S3 or other storage backends.**