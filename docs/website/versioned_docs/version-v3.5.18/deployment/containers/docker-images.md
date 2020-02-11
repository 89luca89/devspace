---
title: Deploy existing images
id: version-v3.5.18-docker-images
original_id: docker-images
---

# TODO @Fabian

DevSpace CLI lets you easily define Kubernetes deployments for any existing Docker image.

### Add deployments for existing images
If you want to use a Docker image from Docker Hub or any other registry, you can add a custom component to your deployments using this command:
```bash
devspace add deployment [deployment-name] --image="my-registry.tld/my-username/image"
```
Example using Docker Hub: `devspace add deployment database --image="mysql"`

> If you are using a private Docker registry, make sure to [login to this registry](../../image-building/registries/authentication).

After adding a new deployment, you need to manually redeploy in order to start the newly added component together with the remainder of your previouly existing deployments.
```bash
devspace deploy
```
