# TO Enhance:

1. To trigger argocd app refresh before syncing ✅
2. To Trigger Sync again untill there's no diff. To sync latest changes. ✅
3. Optimise `SyncWithLabels` func, currently it has to loop through all the apps to see which has the matching label to sync. 
The current argocd app get has no filter option based on labels, to see how to optimise to avoid looking through all apps to sync with labels approach. 🚨
> Some Ref:
- https://pkg.go.dev/github.com/argoproj/argo-cd/v2/pkg/apiclient/application#ApplicationQuery.Refresh
- https://argo-cd.readthedocs.io/en/stable/user-guide/commands/argocd_app_get/