# `Team` migration example

**Starting point**: `Team` resource managed by Azure DevOps Provider (classic).
**Ending point**: `Team` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`Team` on Azure DevOps) will be the same.

Note that the `Team` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "teams"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
teams                                             azuredevops.krateo.io/v1alpha1        false        Team
teams                                             azuredevops.ogen.krateo.io/v1alpha1   true         Team
```

The **starting point** for this migration is the following example of a `Team` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Team
metadata:
  name: team-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  projectRef: 
    namespace: default
    name: pipeline-proj
  name: team-1
  groupRefs:
    - namespace: default
      name: group-sample
  description: Team created from YAML
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that the `Team` resource is referecing a `ConnectorConfig` resource and a `Project` resource, which are both managed by the Azure DevOps Provider "classic".

Note also that the field `groupRefs` is not supported in the `Team` schema of Azure DevOps Provider KOG since the Azure DevOps REST API does not directly support managing groups associated with a team. 
Specific resources for managing group memberships will likely be introduced in future versions of the Azure DevOps Provider KOG.

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate teams.azuredevops.krateo.io team-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Team` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch teams.azuredevops.krateo.io team-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get teams.azuredevops.krateo.io team-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Team
metadata:
  name: team-1
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  projectRef: 
    namespace: default
    name: pipeline-proj
  name: team-1
  groupRefs:
    - namespace: default
      name: group-sample
  description: Team created from YAML
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

Now, you can create a new `Team` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#team-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: Team
metadata:
  name: team-1
  namespace: azuredevops-system
  annotations:
    krateo.io/connector-verbose: "true"
    krateo.io/deletion-policy: orphan             # Optional: to ensure the external resource is not deleted when the resource is deleted  
spec:
  # required: reference to a Configuration CR in the cluster
  configurationRef:
    name: my-team-config
    namespace: default
  organization: krateo-kog
  projectId: 99837031-4e4e-4753-9a47-73fcc4cba766 # Note: this must be a projectId, not a project name since the API returns id in the response
  description: "Team created via KOG"
  name: team-1
EOF
```

Note that: 
- the `projectRef` field has been replaced with `projectId`, which can is the ID of the project on Azure DevOps. Name cannot be used here since the API returns only the ID in the response.

In order to dynamically retrieve the IDs, you can use a Helm `lookup` function.

An example of how to use a Helm `lookup` function to retrieve the project ID dynamically is shown below.
In this case the context is a Helm chart, so the `lookup` function is used to retrieve the `TeamProject` resource by its name and namespace, and then the project ID is accessed from the status of that resource.

```yaml
{{- $project := lookup "azuredevops.krateo.io/v1alpha1" "TeamProject" .Release.Namespace (.Values.project.name | lower) }}
{{- if and $project $project.status $project.status.id }}

apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: Team
spec:
...
  projectId: "{{ $project.status.id }}"  # Dynamically retrieve the project ID
...

{{- end }}
```

Note that you need to have already created a `TeamConfiguration` resource that contains the authentication and configuration information for the `Team` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can check the new `Team` resource managed by Azure DevOps Provider KOG by running the following command:
```sh
kubectl get teams.azuredevops.ogen.krateo.io team-1 -n azuredevops-system
```
And the output should look like this:
```sh
NAME     AGE    READY
team-1   10s    True
```

At this point, you can proceed to delete the old `Team` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `Team` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete teams.azuredevops.krateo.io team-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate teams.azuredevops.krateo.io team-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `Team` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
