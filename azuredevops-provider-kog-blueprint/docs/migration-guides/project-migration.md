# `Project` migration example

## Scenario

Note that the name of the resource in Azure DevOps Provider KOG is `Project`, while in Azure DevOps Provider (classic) it is `TeamProject`.
The rationale behind this change is to be more aligned with the original naming used by Azure DevOps REST API endpoints.
Thus, the migration involves also a change in the kind of the resource.

**Starting point**: `TeamProject` resource managed by Azure DevOps Provider (classic).
**Ending point**: `Project` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`Project` on Azure DevOps) will be the same.

Note that the `TeamProject` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while `Project` is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following commands):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "teamprojects"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
teamprojects                                      azuredevops.krateo.io/v1alpha1        false        TeamProject
```
```sh
kubectl api-resources | awk 'NR==1 || $1 == "projects"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
projects                                          azuredevops.ogen.krateo.io/v1alpha1   true         Project
```

## Migration steps

### Step 1: Pause and set deletion policy on the old `TeamProject` resource

The **starting point** for this migration is the following example of a `TeamProject` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: TeamProject
metadata:
  name: teamproject-sample
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  name: Test Project n.1
  organization: krateo-kog
  description: Lorem ipsum dolor sit amet, consectetur adipiscing elit olè
  visibility: private
  capabilities:
    versioncontrol:
      sourceControlType: Git
    processTemplate:
      templateTypeId: 6b724908-ef14-45cf-84f8-768b5384da45
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that:
- the `TeamProject` resource is referecing a `ConnectorConfig` resource, which is managed by the Azure DevOps Provider "classic".

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate teamprojects.azuredevops.krateo.io teamproject-sample "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `TeamProject` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch teamprojects.azuredevops.krateo.io teamproject-sample \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get teamprojects.azuredevops.krateo.io teamproject-sample -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: TeamProject
metadata:
  name: teamproject-sample
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  name: Test Project n.1
  organization: krateo-kog
  description: Lorem ipsum dolor sit amet, consectetur adipiscing elit olè
  visibility: private
  capabilities:
    versioncontrol:
      sourceControlType: Git
    processTemplate:
      templateTypeId: 6b724908-ef14-45cf-84f8-768b5384da45
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

### Step 2: Create the new `Project` resource

Now, you can create a new `Project` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#project-schema) file of this chart.
You can apply the following example:
```sh
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: Project
metadata:
  name: project-1
  namespace: default
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  # required: reference to a Configuration CR in the cluster
  configurationRef:
    name: my-project-config
    namespace: default
  organization: krateo-kog
  name: Test Project n.1
  description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit olè"
  visibility: private
  capabilities:
    versioncontrol:
      sourceControlType: "Git"
    processTemplate:
      templateTypeId: "6b724908-ef14-45cf-84f8-768b5384da45" # Scrum template                 
EOF
```

Note that: 
- `templateTypeId` is representing the Scrum process template and it is the same used in the old `TeamProject` resource. The list of available process templates can be retrieved via Azure DevOps REST API (reference: [Process Templates - Get](https://learn.microsoft.com/it-it/rest/api/azure/devops/processes/processes/list)).

Note that you need to have already created a `ProjectConfiguration` resource that contains the authentication and configuration information for the `Project` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can wait for the new `Project` resource to be ready by running the following command:
```sh
kubectl wait projects.azuredevops.ogen.krateo.io/project-1 --for condition=Ready=True --namespace default --timeout=300s
```
```sh
project.azuredevops.ogen.krateo.io/project-1 condition met
```

#### Step 3: Delete the old `TeamProject` resource

At this point, you can proceed to delete the old `TeamProject` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `TeamProject` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete teamprojects.azuredevops.krateo.io teamproject-sample
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate teamprojects.azuredevops.krateo.io teamproject-sample --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `TeamProject` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
