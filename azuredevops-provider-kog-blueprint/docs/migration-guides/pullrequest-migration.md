# `PullRequest` migration example

**Starting point**: `PullRequest` resource managed by Azure DevOps Provider (classic).
**Ending point**: `PullRequest` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`PullRequest` on Azure DevOps) will be the same.
Note that the `PullRequest` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "pullrequests"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
pullrequests                                      azuredevops.krateo.io/v1alpha1        false        PullRequest
pullrequests                                      azuredevops.ogen.krateo.io/v1alpha1   true         PullRequest
```

The **starting point** for this migration is the following example of a `PullRequest` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: PullRequest
metadata:
  name: pullrequest-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  projectRef:
    name: pipeline-proj
    namespace: default
  repositoryRef:
    name: gitrepository-sample
    namespace: default
  pullRequest:
    sourceRefName: refs/heads/prova
    targetRefName: refs/heads/main
    title: A new feature-1
    description: This is a new feature
    status: active
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that the `PullRequest` resource is referecing a `ConnectorConfig` resource, a `Project` resource, and a `GitRepository` resource, which are managed by the Azure DevOps Provider "classic".

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate pullrequests.azuredevops.krateo.io pullrequest-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `PullRequest` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch pullrequests.azuredevops.krateo.io pullrequest-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get pullrequests.azuredevops.krateo.io pullrequest-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: PullRequest
metadata:
  name: pullrequest-1
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  projectRef:
    name: pipeline-proj
    namespace: default
  repositoryRef:
    name: gitrepository-sample
    namespace: default
  pullRequest:
    sourceRefName: refs/heads/prova
    targetRefName: refs/heads/main
    title: A new feature-1
    description: This is a new feature
    status: active
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
status:
  conditions:
+  - lastTransitionTime: "2025-07-14T15:35:05Z"
+    reason: ReconcilePaused
+    status: "False"
+    type: Synced
```

Now, you can create a new `PullRequest` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#pullrequest-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: PullRequest
metadata:
  name: pullrequest-1
  namespace: default
  annotations:
    krateo.io/connector-verbose: "true"
    krateo.io/deletion-policy: orphan             # Optional: to ensure the external resource is not deleted when the resource is deleted
spec:
  configurationRef:
    name: sample-pullrequest-config
    namespace: default
  completionOptions:
    deleteSourceBranch: false
    mergeCommitMessage: "Squash merge for feature X, test from Krateo"
    mergeStrategy: squash
    transitionWorkItems: true
  description: "This is a new feature"
  isDraft: false # this field cannot be updated after creation of the CR as it is not among updatable fields allowed by the external API
  mergeOptions:
    conflictAuthorshipCommits: true
    detectRenameFalsePositives: false
    disableRenames: false
  organization: krateo-kog
  project: project-1-classic
  repositoryId: 5605b0ba-e2fa-4aab-af0b-0888321b3a08
  sourceRefName: "refs/heads/prova"
  status: active
  supportsIterations: true
  targetRefName: "refs/heads/main"
  title: "A new feature-1"
EOF
```

Note that:
- the `projectRef` field has been replaced with `project`, which can be either the `ID` or the `name` of the project in the case of the spec of the `PullRequest` resource.
- the `repositoryRef` field has been replaced with `repositoryId`, which is the `ID` of the repository in Azure DevOps.

In order to dynamically retrieve IDs, you can use a Helm `lookup` function.

An example of how to use a Helm `lookup` function to retrieve the `GitRepository` ID dynamically is shown below.
In this case the context is a Helm chart, so the `lookup` function is used to retrieve the `GitRepository` resource by its name and namespace, and then the ID is accessed from the status of that resource.

```yaml
{{- $repository := lookup "azuredevops.ogen.krateo.io/v1alpha1" "GitRepository" .Release.Namespace (.Values.repository.name | lower) }}
{{- if and $repository $repository.status $repository.status.id }}

apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: PullRequest
spec:
...
  repositoryId: "{{ $repository.status.id }}"  # Dynamically retrieve the repository ID
...

{{- end }}
```

Note that you need to have already created a `PullRequestConfiguration` resource that contains the authentication and configuration information for the `PullRequest` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can check the new `PullRequest` resource managed by Azure DevOps Provider KOG by running the following command:
```sh
kubectl get pullrequests.azuredevops.ogen.krateo.io pullrequest-1 -n azuredevops-system
```
And the output should look like this:
```sh
NAME     AGE   READY
pullrequest-1   10s    True
```

At this point, you can proceed to delete the old `PullRequest` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `PullRequest` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete pullrequests.azuredevops.krateo.io pullrequest-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate pullrequests.azuredevops.krateo.io pullrequest-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `PullRequest` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
