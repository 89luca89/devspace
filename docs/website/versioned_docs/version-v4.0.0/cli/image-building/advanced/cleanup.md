---
title: Cleanup Images
id: version-v4.0.0-cleanup
original_id: cleanup
---

When using Docker for image building, disk space on your local computer can get sparse after a lot of Docker builds. DevSpace provides a convenient command to clean up all images that were built with your local Docker daemon using DevSpace. This command does not remove any pushed images remotely and just clears local images and space.

In order to cleanup all created images locally, simply run the following command in your project folder:
```bash
devspace cleanup images
```

In addition it also makes sense to prune your Docker environment to free additional space with the following command:

```bash
docker system prune
```

This command will remove:
- all stopped containers
- all networks not used by at least one container
- all dangling images
- all build cache

These commands should free up a lot of space for new image builds to come.
