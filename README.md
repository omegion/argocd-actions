# ArgoCD Application Actions

[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Find%20and%20Replace-blue.svg?colorA=24292e&colorB=0366d6&style=flat&longCache=true&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAA4AAAAOCAYAAAAfSC3RAAAABHNCSVQICAgIfAhkiAAAAAlwSFlzAAAM6wAADOsB5dZE0gAAABl0RVh0U29mdHdhcmUAd3d3Lmlua3NjYXBlLm9yZ5vuPBoAAAERSURBVCiRhZG/SsMxFEZPfsVJ61jbxaF0cRQRcRJ9hlYn30IHN/+9iquDCOIsblIrOjqKgy5aKoJQj4O3EEtbPwhJbr6Te28CmdSKeqzeqr0YbfVIrTBKakvtOl5dtTkK+v4HfA9PEyBFCY9AGVgCBLaBp1jPAyfAJ/AAdIEG0dNAiyP7+K1qIfMdonZic6+WJoBJvQlvuwDqcXadUuqPA1NKAlexbRTAIMvMOCjTbMwl1LtI/6KWJ5Q6rT6Ht1MA58AX8Apcqqt5r2qhrgAXQC3CZ6i1+KMd9TRu3MvA3aH/fFPnBodb6oe6HM8+lYHrGdRXW8M9bMZtPXUji69lmf5Cmamq7quNLFZXD9Rq7v0Bpc1o/tp0fisAAAAASUVORK5CYII=)](https://github.com/omegion/argocd-actions)
[![Actions Status](https://github.com/omegion/argocd-actions/workflows/Build/badge.svg)](https://github.com/omegion/argocd-actions/actions)
[![Actions Status](https://github.com/omegion/argocd-actions/workflows/Integration%20Test/badge.svg)](https://github.com/omegion/argocd-actions/actions)

This action will sync ArgoCD application.

## Usage

### Example workflow

This example replaces syncs ArgoCD application.

```yaml
name: My Workflow
on: [ push, pull_request ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Sync ArgoCD Application
        uses: omegion/argocd-actions@v1
        with:
          address: "argocd.example.com"
          token: ${{ secrets.ARGOCD_TOKEN }}
          action: sync
          appName: "my-example-app"
```

### Inputs

| Input     | Description                            |
|-----------|----------------------------------------|
| `address` | ArgoCD server address.                 |
| `token`   | ArgoCD Token.                          |
| `action`  | ArgoCD Action i.e. sync.               |
| `appName` | Application name to execute action on. |

**Optional** Labels to sync the ArgoCD app with. If provided, the action will sync the app based on these labels, instead of the app name.

## Examples

### Sync Application

You can sync ArgoCD application after building an image etc.

```yaml
name: My Workflow
on: [ push, pull_request ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Sync ArgoCD Application
        uses: omegion/argocd-actions@master
        with:
          address: "vault.example.com"
          token: ${{ secrets.ARGOCD_TOKEN }}
          action: sync
          appName: "my-example-app"
```

### Example syncing with labels
```yaml
name: My Workflow
on: [ push, pull_request ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Sync ArgoCD Application
        uses: omegion/argocd-actions@master
        with:
          address: "vault.example.com"
          token: ${{ secrets.ARGOCD_TOKEN }}
          action: sync
          labels: "env=production,team=myteam" # Replace with your label key-value pairs
```

## Publishing

To publish a new version of this Action we need to update the Docker image tag in `action.yml` and also create a new
release on GitHub.

- Work out the next tag version number.
- Update the Docker image in `action.yml`.
- Create a new release on GitHub with the same tag.
