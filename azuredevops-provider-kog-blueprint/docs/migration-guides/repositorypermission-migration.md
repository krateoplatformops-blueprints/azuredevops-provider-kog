# `RepositoryPermission` migration example

## Scenario

**Starting point**: `RepositoryPermission` resource managed by Azure DevOps Provider (classic).
**Ending point**: `RepositoryPermission` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`RepositoryPermission` on Azure DevOps) will be the same.
Note that the `RepositoryPermission` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "repositorypermissions"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
repositorypermissions                             azuredevops.krateo.io/v1alpha1        false        RepositoryPermission
repositorypermissions                             azuredevops.ogen.krateo.io/v1alpha1   true         RepositoryPermission
```

## Migration steps

### Step 1: Pause and set deletion policy on the old `RepositoryPermission` resource

The **starting point** for this migration is the following example of a `RepositoryPermission` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: RepositoryPermission
metadata:
  name: repository-permission-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  repositoryRef:
    name: repo-perms-sample
    namespace: default
  permissions: 
    merge: false
    identity: 
      type: azure-group
      projectRef:
        name: teamproject
        namespace: default
      name: Contributors
    allowList:
       - genericContribute
       - forcePush
    denyList:
      - createTag
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that:
- the `RepositoryPermission` resource is referecing a `ConnectorConfig` resource and a `Project` resource, which are both managed by the Azure DevOps Provider "classic".

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate repositorypermissions.azuredevops.krateo.io repository-permission-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `RepositoryPermission` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch repositorypermissions.azuredevops.krateo.io repository-permission-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get repositorypermissions.azuredevops.krateo.io repository-permission-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: RepositoryPermission
metadata:
  name: repository-permission-1
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan  
  repositoryRef:
    name: repo-perms-sample
    namespace: default
  permissions: 
    merge: false
    identity: 
      type: azure-group
      projectRef:
        name: teamproject
        namespace: default
      name: Contributors
    allowList:
       - genericContribute
       - forcePush
    denyList:
      - createTag
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

### Step 2: Create the new `RepositoryPermission` resource

Now, you can create a new `RepositoryPermission` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#repositorypermission-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: RepositoryPermission
metadata:
  name: repository-permission-1
  namespace: azuredevops-system                   # Replace with your namespace
  annotations:
    krateo.io/connector-verbose: "true"
    krateo.io/deletion-policy: orphan             # Optional: to ensure the external resource is not deleted when the resource is deleted  
spec:
  configurationRef:
    name: my-repositorypermission-config
    namespace: default

  organization: krateo-kog                            # required
  projectId: 99837031-4e4e-4753-9a47-73fcc4cba766     # required
  repositoryId: c2f8d804-7c2f-4f3f-bd69-70c3169cfb8c  # required

  permissions:
    identity:
      type: azure-group
      name: Contributors
    allow:
      GenericContribute: true
      ForcePush: true
    deny:
      CreateTag: true
EOF
```

Note that: 
- the `projectRef` field has been replaced with `projectId`.
- the `repositoryRef` field has been replaced with `repositoryId`.

In order to dynamically retrieve the IDs, you can use a Helm `lookup` function.

An example of how to use a Helm `lookup` function to retrieve the project ID, and repository ID dynamically is shown below.
In this case the context is a Helm chart, so the `lookup` function is used to retrieve the `TeamProject` and `GitRepository` resources by their names and namespace, and then the project ID and repository ID are accessed from the status of those resources.

```yaml
{{- $project := lookup "azuredevops.krateo.io/v1alpha1" "TeamProject" .Release.Namespace (.Values.project.name | lower) }}
{{- if and $project $project.status $project.status.id }}

{{- $repository := lookup "azuredevops.ogen.krateo.io/v1alpha1" "GitRepository" .Release.Namespace (.Values.repository.name | lower) }}
{{- if and $repository $repository.status $repository.status.id }}

apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: RepositoryPermission
spec:
...
  projectId: "{{ $project.status.id }}"  # Dynamically retrieve the project ID
  repositoryId: "{{ $repository.status.id }}"  # Dynamically retrieve the repository ID
...

{{- end }}
```

Note that you need to have already created a `RepositoryPermissionConfiguration` resource that contains the authentication and configuration information for the `RepositoryPermission` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can wait for the new `RepositoryPermission` resource to be ready by running the following command:
```sh
kubectl wait repositorypermissions.azuredevops.ogen.krateo.io/repository-permission-1 --for condition=Ready=True --namespace azuredevops-system --timeout=300s
```
```sh
repositorypermission.azuredevops.ogen.krateo.io/repository-permission-1 condition met
```

### Step 3: Delete the old `RepositoryPermission` resource

At this point, you can proceed to delete the old `RepositoryPermission` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `RepositoryPermission` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete repositorypermissions.azuredevops.krateo.io repository-permission-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate repositorypermissions.azuredevops.krateo.io repository-permission-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `RepositoryPermission` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
