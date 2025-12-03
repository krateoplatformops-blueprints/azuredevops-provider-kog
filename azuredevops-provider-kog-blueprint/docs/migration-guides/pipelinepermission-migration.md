# `PipelinePermission` migration example

**Starting point**: `PipelinePermission` resource managed by Azure DevOps Provider "classic".
**Ending point**: `PipelinePermission` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`PipelinePermission` on Azure DevOps) will be the same.

Note that the `PipelinePermission` resource is a non-namespaced resource in the context of Azure DevOps Provider "classic", while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "pipelinepermissions"'
```
Output:
```sh
NAME                                SHORTNAMES   APIVERSION                            NAMESPACED   KIND
pipelinepermissions                              azuredevops.ogen.krateo.io/v1alpha1   true         PipelinePermission
pipelinepermissions                              azuredevops.krateo.io/v1alpha2        false        PipelinePermission
```

The **starting point** for this migration is the following example of a `Pipeline` resource managed by the Azure DevOps Provider "classic":
```yaml
apiVersion: azuredevops.krateo.io/v1alpha2
kind: PipelinePermission
metadata:
  name: pipeline-permission-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  authorizeAll: false
  projectRef: 
    name: project-1-classic
    namespace: default
  pipelines:
    - pipelineRef:
        name: pipeline-1
        namespace: default
      authorized: true
  resource:
    type: environment
    resourceRef:
      name: enviroment-1-classic
      namespace: default
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that the `PipelinePermission` resource is referecing a `ConnectorConfig` resource, a `Environment` resource, and a `Pipeline` resource, which are managed by the Azure DevOps Provider "classic".

Note that the `Pipeline` referenced in the example above is a `Pipeline` resource managed by the Azure DevOps Provider "classic".  
However, the `PipelinePermission` resource managed by the Azure DevOps Provider KOG will work with both `Pipeline` resources managed by the Azure DevOps Provider "classic" and the Azure DevOps Provider KOG since it use the `id` of the pipeline as the reference.

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate pipelinepermissions.azuredevops.krateo.io pipeline-permission-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Pipeline` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch pipelinepermissions.azuredevops.krateo.io pipeline-permission-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get pipelinepermissions.azuredevops.krateo.io pipeline-permission-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha2
kind: PipelinePermission
metadata:
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  authorizeAll: false
  connectorConfigRef:
    name: connectorconfig-sample
    namespace: default
  pipelines:
  - authorized: true
    pipelineRef:
      name: pipeline-1
      namespace: default
  projectRef:
    name: project-1-classic
    namespace: default
  resource:
    resourceRef:
      name: enviroment-1-classic
      namespace: default
    type: environment
status:
  conditions:
+ - lastTransitionTime: "2025-07-16T10:10:59Z"
+   reason: ReconcilePaused
+   status: "False"
+   type: Synced
```

Now, you can create a new `PipelinePermission` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#pipelinepermission-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: PipelinePermission
metadata:
  name: pipeline-permission-1
  namespace: azuredevops-system                   # Replace with your namespace
  annotations:
    krateo.io/connector-verbose: "true"           # Optional: to enable verbose logging
    krateo.io/deletion-policy: orphan             # Optional: to ensure the external resource is not deleted when the resource is deleted
spec:
  configurationRef:                               # Reference to a PipelinePermissionConfiguration CR that contains the authentication information.
    name: my-pipelinepermission-config
    namespace: default
  
  organization: krateo-kog                        # Name of the Azure DevOps organization
  project: "project-1-classic"                    # ID or name of the project
  
  resourceType: "environment"                     # Type of the resource
  resourceId: "35"                                # ID of the resource

  allPipelines:
    authorized: false                             # Whether to authorize all pipelines for the resource

  pipelines:
    - id: 66                                      # ID of the pipeline to authorize for the resource  
      authorized: true
EOF
```

Note that:
- the `projectRef` field has been replaced with `project`, which can be either the `ID` or the `name` of the project in the case of the spec of the `PipelinePermission` resource.
- the `resource` field has been replaced with `resourceType` and `resourceId`, which are the type and ID of the resource to authorize for the pipelines.
- the `pipelineRef` fields have been replaced with `id`, which are the ID of the pipelines to authorize for the resource.

In order to dynamically retrieve the IDs, you can use a Helm `lookup` function.

An example of how to use the Helm `lookup` function to retrieve the project ID, environment ID and pipeline ID dynamically is shown below.
In this case the context is a Helm chart, so the `lookup` function is used to retrieve the `TeamProject`, `Environment` and `Pipeline` resources by their names and namespace, and then the project ID, environment ID, and pipeline ID are accessed from the status of those resources.

```yaml
{{- $project := lookup "azuredevops.krateo.io/v1alpha1" "TeamProject" .Release.Namespace (.Values.project.name | lower) }}
{{- if and $project $project.status $project.status.id }}

{{- $environment := lookup "azuredevops.krateo.io/v1alpha1" "Environment" .Release.Namespace (.Values.environment.name | lower) }}
{{- if and $environment $environment.status $environment.status.id }}

{{- $pipeline := lookup "azuredevops.krateo.io/v1alpha1" "Pipeline" .Release.Namespace (.Values.pipeline.name | lower) }}
{{- if and $pipeline $pipeline.status $pipeline.status.id }}

apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: PipelinePermission
spec:
  project: "{{ $project.status.id }}"           # Dynamically retrieve the project ID

  resourceType: "environment"                   # Type of the resource
  resourceId: "{{ $environment.status.id }}"    # Dynamically retrieve the environment ID

  pipelines:
    - id: "{{ $pipeline.status.id }}"           # Dynamically retrieve the pipeline ID  
...

{{- end }}
```

You can check the new `PipelinePermission` resource managed by Azure DevOps Provider KOG by running the following command:
```sh
kubectl get pipelinepermissions.azuredevops.ogen.krateo.io pipeline-permission-1 -n azuredevops-system
```
And the output should look like this:
```sh
NAME                    AGE   READY
pipeline-permission-1   61s   True
```

At this point, you can proceed to delete the old `PipelinePermission` resource managed by Azure DevOps Provider "classic" (note the different API group).

First, you can delete the old `Pipeline` resource managed by Azure DevOps Provider "classic":
```sh
kubectl delete pipelinepermissions.azuredevops.krateo.io pipeline-permission-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate pipelinepermissions.azuredevops.krateo.io pipeline-permission-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `PipelinePermission` resource from Azure DevOps Provider "classic" to Azure DevOps Provider KOG is complete.
