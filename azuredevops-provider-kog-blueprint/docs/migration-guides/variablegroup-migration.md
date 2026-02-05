# `VariableGroup` migration example

## Scenario

**Starting point**: `VariableGroups` resource managed by Azure DevOps Provider (classic).
**Ending point**: `VariableGroup` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`VariableGroup` on Azure DevOps) will be the same.

Note that the `VariableGroups` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while the `VariableGroup` resource is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || /variablegroup/'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
variablegroups                                    azuredevops.krateo.io/v1alpha1        false        VariableGroups
variablegroups                                    azuredevops.ogen.krateo.io/v1alpha1   true         VariableGroup
```

## Migration steps

### Step 1: Pause and set deletion policy on the old `VariableGroups` resource

The **starting point** for this migration is the following example of a `VariableGroups` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: VariableGroups
metadata:
  name: vg-test-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  deletionPolicy: Orphan
  name: vg-test
  description: "Variable group for testing"
  variables:
    var1:
      isSecret: false
      value: "value12"
    var2:
      isSecret: true
      value: "value2"
  variableGroupProjectReferences:
    - name: vg-project-test-1
      description: "Project 1"
      projectRef:
        name: pipeline-proj
        namespace: default
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that:
- the `VariableGroups` resource is referencing a `ConnectorConfig` resource and a `Project` resource (via `projectRef`), which are both managed by the Azure DevOps Provider "classic".
- the `deletionPolicy` field is already set to `Orphan` in this example. If your resource does not have it, you will need to set it in the next step.

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following command:
```sh
kubectl annotate variablegroups.azuredevops.krateo.io vg-test-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `VariableGroups` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource (skip this step if it is already set).
You can do this with the following command:
```sh
kubectl patch variablegroups.azuredevops.krateo.io vg-test-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get variablegroups.azuredevops.krateo.io vg-test-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: VariableGroups
metadata:
  name: vg-test-1
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
  deletionPolicy: Orphan
  name: vg-test
  description: "Variable group for testing"
  variables:
    var1:
      isSecret: false
      value: "value12"
    var2:
      isSecret: true
      value: "value2"
  variableGroupProjectReferences:
    - name: vg-project-test-1
      description: "Project 1"
      projectRef:
        name: pipeline-proj
        namespace: default
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

### Step 2: Create the new `VariableGroup` resource

Now, you can create a new `VariableGroup` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#variablegroup-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: VariableGroup
metadata:
  name: vg-test-1
  namespace: azuredevops-system                   # Replace with your namespace
  annotations:
    krateo.io/connector-verbose: "true"
    krateo.io/deletion-policy: "orphan"           # Optional: to ensure the external resource is not deleted when the CR is deleted
spec:
  # required: reference to a Configuration CR in the cluster
  configurationRef:
    name: my-variablegroup-config
    namespace: default

  organization: krateo-kog                        # Name of the Azure DevOps organization
  project: "99837031-4e4e-4753-9a47-73fcc4cba766" # Project ID (project name is not supported in this field)

  name: vg-test
  description: "Variable group for testing"
  type: "Vsts"

  variables:
    var1:
      value: "value12"
      isSecret: false
      isReadOnly: false
    var2:
      value: "value2"
      isSecret: true
      isReadOnly: false

  # Project references: which project can use this variable group.
  # Array shape but only one project reference is supported.
  variableGroupProjectReferences:
    - name: "vg-test"                             # Must match spec.name
      description: "Variable group for testing"
      projectReference:
        id: "99837031-4e4e-4753-9a47-73fcc4cba766"  # Project ID
EOF
```

Note that:
- the `kind` has changed from `VariableGroups` (plural, classic) to `VariableGroup` (singular, KOG).
- `connectorConfigRef` has been replaced with `configurationRef`, pointing to a `VariableGroupConfiguration` CR that holds the authentication information and API-version settings.
- `deletionPolicy` in `spec` has moved to the `krateo.io/deletion-policy` annotation in `metadata`.
- `organization` and `project` are now explicit fields. `project` must be the **project ID**, not a name or a Kubernetes resource reference.
- inside `variableGroupProjectReferences`, the Kubernetes-style `projectRef` (name + namespace) has been replaced with `projectReference` containing the Azure DevOps project `id` directly. The `name` field inside the array entry must match `spec.name`.

In order to dynamically retrieve the project ID, you can use a Helm `lookup` function.

Note that you need to have already created a `VariableGroupConfiguration` resource that contains the authentication and configuration information for the `VariableGroup` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can wait for the new `VariableGroup` resource to be ready by running the following command:
```sh
kubectl wait variablegroups.azuredevops.ogen.krateo.io/vg-test-1 --for condition=Ready=True --namespace azuredevops-system --timeout=300s
```
```sh
variablegroup.azuredevops.ogen.krateo.io/vg-test-1 condition met
```

### Step 3: Delete the old `VariableGroups` resource

At this point, you can proceed to delete the old `VariableGroups` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `VariableGroups` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete variablegroups.azuredevops.krateo.io vg-test-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate variablegroups.azuredevops.krateo.io vg-test-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `VariableGroup` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
