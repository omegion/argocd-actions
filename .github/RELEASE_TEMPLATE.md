## Installation

You can use `go` to build S3 Secrets Manager locally with:

```shell
go install github.com/omegion/argocd-actions
```

Or, you can use the usual commands to install or upgrade:

On OS X

```shell
curl -L https://github.com/omegion/argocd-actions/releases/download/{{.Env.VERSION}}/argocd-actions-darwin-amd64 >/usr/local/bin/argocd-actions 
&& \
  chmod +x /usr/local/bin/argocd-actions
```

On Linux

```shell
curl -L https://github.com/omegion/argocd-actions/releases/download/{{.Env.VERSION}}/argocd-actions-linux-amd64 >/usr/local/bin/argocd-actions 
&& \
    chmod +x /tmp/argocd-actions && \
    sudo cp /tmp/argocd-actions /usr/local/bin/argocd-actions
```

Otherwise, download one of the releases from the [release page](https://github.com/omegion/argocd-actions/releases/)
directly.

See the install [docs](https://argocd-actions.omegion.dev) for more install options and instructions.

## Changelog
