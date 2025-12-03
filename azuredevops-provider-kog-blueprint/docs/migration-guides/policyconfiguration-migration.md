# `PolicyConfiguration` migration example

Note that the name of the resource in Azure DevOps Provider KOG is `PolicyConfiguration`, while in Azure DevOps Provider (classic) it is `Policy`.
The rationale behind this change is to be more aligned with the original naming used by Azure DevOps REST API.
Thus, the migration involves also a change in the kind of the resource.

**Starting point**: `Policy` resource managed by Azure DevOps Provider (classic).
**Ending point**: `PolicyConfiguration` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`PolicyConfiguration` on Azure DevOps) will be the same.

Note that the `Policy` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following commands):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "policies"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
policies                                          azuredevops.krateo.io/v1alpha1        false        Policy
```
```sh
kubectl api-resources | awk 'NR==1 || $1 == "policyconfigurations"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
policyconfigurations                              azuredevops.ogen.krateo.io/v1alpha1   true         PolicyConfiguration
```

The **starting point** for this migration is the following example of a `Policy` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Policy
metadata:
  name: policy-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  connectorConfigRef:
    name: connectorconfig-sample
    namespace: default
  policyBody:
    isBlocking: false
    isDeleted: false
    isEnabled: true
    isEnterpriseManaged: false
    projectRef:
      name: pipeline-proj
      namespace: default
    settings:
      buildDefinitionId: 17
      scope:
      - matchKind: Exact
        refName: refs/heads/main
        repositoryRef:
          name: policy-repo
          namespace: default
    type:
      id: "0609b952-1397-4640-95ec-e00a01b2c241"
```

Note: the `id` field in the `type` section corresponds to the policy description: *"This policy will require a successful build has been performed before updating protected refs."*. You can find more information about the policy types in the related section of the [troubleshooting guide](../troubleshooting_guide.md#policyconfiguration)

Note that the `Policy` resource is referecing a `ConnectorConfig` resource, a `Project` resource and a `GitRepository` resource, which are by the Azure DevOps Provider "classic".

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate policies.azuredevops.krateo.io policy-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Policy` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch policies.azuredevops.krateo.io policy-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get policies.azuredevops.krateo.io policy-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Policy
metadata:
  name: policy-1
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  connectorConfigRef:
    name: connectorconfig-sample
    namespace: default
  policyBody:
    isBlocking: false
    isDeleted: false
    isEnabled: true
    isEnterpriseManaged: false
    projectRef:
      name: pipeline-proj
      namespace: default
    settings:
      buildDefinitionId: 17
      scope:
      - matchKind: Exact
        refName: refs/heads/main
        repositoryRef:
          name: policy-repo
          namespace: default
    type:
      id: "0609b952-1397-4640-95ec-e00a01b2c241"
status:
  conditions:
+  - lastTransitionTime: "2025-07-14T15:35:05Z"
+    reason: ReconcilePaused
+    status: "False"
+    type: Synced
```

Now, you can create a new `PolicyConfiguration` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#policyconfiguration-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: PolicyConfiguration
metadata:
  name: policy-1
  namespace: azuredevops-system
  annotations:
    krateo.io/connector-verbose: "true"
    krateo.io/deletion-policy: orphan             # Optional: to ensure the external resource is not deleted when the resource is deleted
spec:
  configurationRef:
    name: my-policyconfiguration-config
    namespace: default

  isBlocking: false             # required        
  isDeleted: false            
  isEnabled: true               # required
  isEnterpriseManaged: false
  
  organization: krateo-kog      # required
  project: project-1-classic    # required
  
  type:
    id: "0609b952-1397-4640-95ec-e00a01b2c241"

  settings:
    buildDefinitionId: 17
    scope:
      - repositoryId: "5605b0ba-e2fa-4aab-af0b-0888321b3a08"
        refName: "refs/heads/main"
        matchKind: Exact
EOF
```

Note that:
- the `projectRef` field has been replaced with `project`, which can be either the `ID` or the `name` of the project in the case of the spec of the `PolicyConfiguration` resource.
- the `repositoryRef` field has been replaced with `repositoryId`, which is the ID of the repository in the case of the spec of the `PolicyConfiguration` resource.

In order to dynamically retrieve the IDs, you can use a Helm `lookup` function.

An example of how to use a Helm `lookup` function to retrieve the GitRepository ID dynamically is shown below.
In this case the context is a Helm chart, so the `lookup` function is used to retrieve the Azure DevOps KOG `GitRepository` and `TeamProject` resources by their names and namespace, and then the repository ID and project ID are accessed from the status of those resources.

```yaml
{{- $project := lookup "azuredevops.krateo.io/v1alpha1" "TeamProject" .Release.Namespace (.Values.project.name | lower) }}
{{- if and $project $project.status $project.status.id }}

{{- $repository := lookup "azuredevops.ogen.krateo.io/v1alpha1" "GitRepository" .Release.Namespace (.Values.repository.name | lower) }}
{{- if and $repository $repository.status $repository.status.id }}

apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: PolicyConfiguration
spec:
...
  project: "{{ $project.status.id }}"                # Dynamically retrieve the project ID
  settings:
    scope:
      - repositoryId: "{{ $repository.status.id }}"  # Dynamically retrieve the repository ID
...

{{- end }}
```

Note that you need to have already created a `PolicyConfigurationConfiguration` resource that contains the authentication and configuration information for the `PolicyConfiguration` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can check the new `PolicyConfiguration` resource managed by Azure DevOps Provider KOG by running the following command:
```sh
kubectl get policyconfigurations.azuredevops.ogen.krateo.io repo-1 -n azuredevops-system
```
And the output should look like this:
```sh
NAME       AGE    READY
policy-1   10s    True
```

At this point, you can proceed to delete the old `Policy` resource managed by Azure DevOps Provider (classic) (note the different API group and kind).

First, you can delete the old `Policy` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete policies.azuredevops.krateo.io policy-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate policies.azuredevops.krateo.io policy-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `Policy` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
