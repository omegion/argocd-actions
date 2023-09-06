## Installation

You can use `go` to build ArgoCD Actions locally with:

```shell
go install github.com/cheelim1/argocd-actions@latest
```

Or, you can use the usual commands to install or upgrade:

On OS X

```shell
sudo curl -fL https://github.com/cheelim1/argocd-actions/releases/download/{{.Env.VERSION}}/argocd-actions-darwin-amd64 -o /usr/local/bin/argocd-actions \
&& sudo chmod +x /usr/local/bin/argocd-actions
```

On Linux

```shell
sudo curl -fL https://github.com/cheelim1/argocd-actions/releases/download/{{.Env.VERSION}}/argocd-actions-linux-amd64 -o /usr/local/bin/argocd-actions \
&& sudo chmod +x /usr/local/bin/argocd-actions
```

On Windows (Powershell)

```powershell
Invoke-WebRequest -Uri https://github.com/cheelim1/argocd-actions/releases/download/{{.Env.VERSION}}/argocd-actions-windows-amd64 -OutFile $home\AppData\Local\Microsoft\WindowsApps\argocd-actions.exe
```

Otherwise, download one of the releases from the [release page](https://github.com/cheelim1/argocd-actions/releases/) 
directly.

See the installation [docs](https://github.com/cheelim1/argocd-actions/releases) for more install options and instructions.

## Changelog
