# `Pipeline` migration example

**Starting point**: `Pipeline` resource managed by Azure DevOps Provider "classic".
**Ending point**: `Pipeline` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`Pipeline` on Azure DevOps) will be the same.

Note that the `Pipeline` resource is a non-namespaced resource in the context of Azure DevOps Provider "classic", while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "pipelines"'
```
Output:
```sh
NAME                                SHORTNAMES   APIVERSION                            NAMESPACED   KIND
pipelines                                        azuredevops.ogen.krateo.io/v1alpha1   true         Pipeline
pipelines                                        azuredevops.krateo.io/v1alpha1        false        Pipeline
```

The **starting point** for this migration is the following example of a `Pipeline` resource managed by the Azure DevOps Provider "classic":
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Pipeline
metadata:
  name: pipeline-1
spec:
  name: pipeline-1
  configurationType: yaml
  definitionPath: azure-pipelines.yml
  repositoryType: azureReposGit
  repositoryRef:
    name: repo-1
    namespace: default
  projectRef:
    name: project-1-classic
    namespace: default
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that the `Pipeline` resource is referecing a `ConnectorConfig` resource, a `Project` resource, and a `GitRepository` resource, which are managed by the Azure DevOps Provider "classic".
Note that the `GitRepository` referenced in the example above is a `GitRepository` resource managed by the Azure DevOps Provider "classic". 
However, the `Pipeline` resource managed by the Azure DevOps Provider KOG will work with both `GitRepository` resources managed by the Azure DevOps Provider "classic" and the Azure DevOps Provider KOG since it use the `id` of the repository as the reference.

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate pipelines.azuredevops.krateo.io pipeline-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Pipeline` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch pipelines.azuredevops.krateo.io pipeline-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get pipelines.azuredevops.krateo.io pipeline-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Pipeline
metadata:
  annotations:
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  configurationType: yaml
  connectorConfigRef:
    name: connectorconfig-sample
    namespace: default
  definitionPath: azure-pipelines.yml
  name: pipeline-1
  projectRef:
    name: project-1-classic
    namespace: default
  repositoryRef:
    name: repo-1
    namespace: default
  repositoryType: azureReposGit
status:
  conditions:
+  - lastTransitionTime: "2025-07-16T08:33:23Z"
+    reason: ReconcilePaused
+    status: "False"
+    type: Synced
  id: "65"
  url: https://dev.azure.com/krateo-kog/82575162-69b2-4a88-8fd7-bda0d05c0284/_apis/pipelines/65?revision=1
```

Now, you can create a new `Pipeline` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../README.md#pipeline-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: Pipeline
metadata:
  name: pipeline-1
  namespace: azuredevops-system                   # Replace with your namespace
  annotations:
    krateo.io/connector-verbose: "true"           # Optional: to enable verbose logging
    krateo.io/deletion-policy: orphan             # Optional: to ensure the external resource is not deleted when the resource is deleted
spec:
  configurationRef:                               # Reference to a PipelineConfiguration CR that contains the authentication information.
    name: my-pipeline-config
    namespace: default
  
  organization: krateo-kog                        # Name of the Azure DevOps organization
  project: "project-1-classic"                    # ID or name of the project
  
  configuration:
    path: azure-pipelines.yml                      # Path to the pipeline configuration file within the repository
    repository: 
      id: "63885da2-bfea-4b13-b19d-cfdb414bccaf"   # ID of the repository where the pipeline is defined
      type: azureReposGit                          # Type of the repository, e.g., gitHub, azureReposGit, etc.
    type: yaml                                     # Type of the pipeline configuration, e.g., yaml, designer, etc.

  name: pipeline-1                                 # Name of the pipeline
EOF
```

Note that:
- the `projectRef` field has been replaced with `project`, which can be either the `ID` or the `name` of the project in the case of the spec of the `Pipeline` resource.
- the `repositoryRef` field has been replaced with `configuration.repository.id`, which is the ID of the repository where the pipeline is defined.

In order to dynamically retrieve the IDs, you can use a Helm `lookup` function.

An example of how to use the `lookup` function to retrieve the project ID and repository ID dynamically is shown below.
In this case the context is a Helm chart, so the `lookup` function is used to retrieve the `TeamProject` and `GitRepository` resources by their names and namespace, and then the project ID and repository ID are accessed from the status of those resources.

```yaml
{{- $project := lookup "azuredevops.krateo.io/v1alpha1" "TeamProject" .Release.Namespace (.Values.project.name | lower) }}
{{- if and $project $project.status $project.status.id }}

{{- $repository := lookup "azuredevops.krateo.io/v1alpha1" "GitRepository" .Release.Namespace (.Values.repository.name | lower) }}
{{- if and $repository $repository.status $repository.status.id }}

apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: Pipeline
spec:
  project: "{{ $project.status.id }}"  # Dynamically retrieve the project ID

  configuration:
    repository:
      id: "{{ $repository.status.id }}"  # Dynamically retrieve the repository ID
...

{{- end }}
```

Note that you need to have already created a `PipelineConfiguration` resource that contains the authentication and configuration information for the `Pipeline` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can check the new `Pipeline` resource managed by Azure DevOps Provider KOG by running the following command:
```sh
kubectl get pipelines.azuredevops.ogen.krateo.io pipeline-1 -n azuredevops-system
```
And the output should look like this:
```sh
NAME         AGE   READY
pipeline-1   86s   True
```

At this point, you can proceed to delete the old `Pipeline` resource managed by Azure DevOps Provide "classic" (note the different API group).

First, you can delete the old `Pipeline` resource managed by Azure DevOps Provider "classic":
```sh
kubectl delete pipelines.azuredevops.krateo.io pipeline-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate pipelines.azuredevops.krateo.io pipeline-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `Pipeline` resource from Azure DevOps Provider "classic" to Azure DevOps Provider KOG is complete.
