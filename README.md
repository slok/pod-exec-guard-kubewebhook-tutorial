# pod-exec-guard-kubewebhook-tutorial

## Introduction

This is a tutorial that shows how to develop a [Kubernetes admission webhooks][k8s-webhooks].

To explain this, the tutorial is split in 5 videos. We will create a webhook from scratch that will try to
recreate the webhook that [this post][wh-post] describes.

## The problem to solve

When a user makes an `exec` operation on a pod, we mark that pod and set a TTL
when that TTL expires, the pod will be deleted.

The tutorial is based on [kubewebhook] to develop the webhook and uses [kube-janitor] to delete the pods after
a specific TTL.

## Disclaimer

- The webhook it is not production ready.
- Its just made as a tutorial step by step.
- It would need more structure, tests, docs, metrics...

## Content

- [Video 0][video0]
  - Introduction and context.
- [Video 1][video1]
  - Create app structure
  - Create webhook without domain logic using [kubewebhook].
  - Running application.
- [Video 2][video2]
  - Set up dev environment (dev cluster with [kind], certs with [mkcert], tunnels with [ngrok], webhook registration).
  - End-2-end manual testing to check webhook is being called.
- [Video 3][video3]
  - Implement domain logic of webhook (Marking pod as drifted)
  - End-2-end manual testing to check webhook is marking pods.
- [Video 4][video4]
  - Deploy [Kube-janitor].
  - Implement domain logic of webhook (Marking pod with expiration time).
  - End-2-end manual testing to check webhook sets expiration and kube-janitor deletes.

[wh-post]: https://medium.com/box-tech-blog/using-k8s-admission-controllers-to-detect-container-drift-at-runtime-cc0f6c67c583
[k8s-webhooks]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers
[kubewebhook]: https://github.com/slok/kubewebhook
[kind]: https://kind.sigs.k8s.io/
[kube-janitor]: https://codeberg.org/hjacobs/kube-janitor
[mkcert]: https://github.com/FiloSottile/mkcert
[ngrok]: https://ngrok.com/
[video0]: https://drive.google.com/file/d/1svMFVFESCUqHxKG41SWnJoVchBQ2cxhc/view?usp=sharing
[video1]: https://drive.google.com/file/d/151nr5QrPRNE3r6xqZJxKzNGzm6u4crZo/view?usp=sharing
[video2]: https://drive.google.com/file/d/1FBWOvpEZMBqGMiuuo4c5Snj_EYVP0seC/view?usp=sharing
[video3]: https://drive.google.com/file/d/1AGDnYJcjaq4uRjplDAdhwmZulJEpcjfz/view?usp=sharing
[video4]: https://drive.google.com/file/d/1VJPI5xwvBnifk0fikKLksvHpY9XfSLAi/view?usp=sharing
