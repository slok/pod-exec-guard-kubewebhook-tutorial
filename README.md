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

- **[Watch video 1][video1]** ([Download][video1-dl])
  - Introduction and context.
- **[Watch video 2][video2]** ([Download][video2-dl])
  - Create app structure
  - Create webhook without domain logic using [kubewebhook].
  - Running application.
- **[Watch video 3][video3]** ([Download][video3-dl])
  - Set up dev environment (dev cluster with [kind], certs with [mkcert], tunnels with [ngrok], webhook registration).
  - End-2-end manual testing to check webhook is being called.
- **[Watch video 4][video4]** ([Download][video4-dl])
  - Implement domain logic of webhook (Marking pod as drifted)
  - End-2-end manual testing to check webhook is marking pods.
- **[Watch video 5][video5]** ([Download][video5-dl])
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
[video1]: https://youtu.be/ujCzvjGXO08
[video2]: https://youtu.be/3gsrYSQcgJI
[video3]: https://youtu.be/3hqQWN7oTrU
[video4]: https://youtu.be/miCVIbKZdXw
[video5]: https://youtu.be/LidzzFRat3k
[video1-dl]: https://drive.google.com/file/d/1svMFVFESCUqHxKG41SWnJoVchBQ2cxhc/view?usp=sharing
[video2-dl]: https://drive.google.com/file/d/151nr5QrPRNE3r6xqZJxKzNGzm6u4crZo/view?usp=sharing
[video3-dl]: https://drive.google.com/file/d/1FBWOvpEZMBqGMiuuo4c5Snj_EYVP0seC/view?usp=sharing
[video4-dl]: https://drive.google.com/file/d/1AGDnYJcjaq4uRjplDAdhwmZulJEpcjfz/view?usp=sharing
[video5-dl]: https://drive.google.com/file/d/1VJPI5xwvBnifk0fikKLksvHpY9XfSLAi/view?usp=sharing
