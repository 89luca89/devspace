---
title: Connect remote debuggers
id: version-v3.5.18-remote-debuggers
original_id: remote-debuggers
---

DevSpace CLI lets you easily [start applications in development mode](../../workflow-basics/development) and connect remote debuggers for your application using the following steps:
1. Configure DevSpace CLI to [use a development Dockerfile](../../development/overrides#configuring-a-different-dockerfile-during-devspace-dev) that:
   - ships with the appropriate tools for debugging your application
   - starts your application together with the debugger, e.g. setting the `ENTRYPOINT` of your Dockerfile to `node --inspect=0.0.0.0:9229 index.js` would start the Node.js remote debugger on port `9229`
2. Define port-forwarding for the port of your remote debugger (e.g. `9229`) within the `dev.ports` section of your `devspace.yaml`
3. Connect your IDE to the remote debugger (see the docs of your IDE for help)
4. Set breakpoints and debug your application directly inside Kubernetes
