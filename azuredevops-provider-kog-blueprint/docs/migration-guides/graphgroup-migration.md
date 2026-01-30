# `GraphGroup` migration example

## Scenario

Note that the name of the resource in Azure DevOps Provider KOG is `GraphGroup`, while in Azure DevOps Provider (classic) it is `Groups`.
The rationale behind this change is to be more aligned with the original naming used by Azure DevOps REST API endpoints.
Thus, the migration involves also a change in the kind of the resource.

**Starting point**: `Groups` resource managed by Azure DevOps Provider (classic).
**Ending point**: `GraphGroup` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`GraphGroup` on Azure DevOps) will be the same.

Note that the `Groups` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while `GraphGroup` is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following commands):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "groups"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
groups                                            azuredevops.krateo.io/v1alpha1        false        Groups
```
```sh
kubectl api-resources | awk 'NR==1 || $1 == "graphgroups"'
```
Output:
```sh
NAME                                 SHORTNAMES   APIVERSION                            NAMESPACED   KIND
graphgroups                                       azuredevops.ogen.krateo.io/v1alpha1   true         GraphGroup
```

## Migration steps

### Step 1: Pause and set deletion policy on the old `Groups` resource

The **starting point** for this migration is the following example of a `Groups` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Groups
metadata:
  name: group-test
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  deletionPolicy: Delete
  membership:
    # If projectRef is not specified, the group will be created in the organization if is set
    projectRef: 
      namespace: default
      name: pipeline-proj
  groupName: group-test
  groupRefs:
    - namespace: default
      name: group-sample
  description: Group created from YAML
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that:
- the `Groups` resource is referecing a `ConnectorConfig` resource and a `Project` resource, which are both managed by the Azure DevOps Provider "classic".
- the field `groupRefs` is not supported in the `GraphGroup` schema of Azure DevOps Provider KOG since it maps directly the Azure DevOps REST API.
- Azure DevOps Provider "classic" creates **memberships** in a implicit way inside the reconciliation loop logic of the `Groups` resource, while Azure DevOps Provider KOG requires explicit management of memberships through dedicated resources. Specific resources for managing [memberships](https://learn.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-7.1) will likely be introduced in future versions of the Azure DevOps Provider KOG. If this functionality is critical for your use case, it may be advisable to delay the migration until such resources are available.

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate groups.azuredevops.krateo.io group-test "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Groups` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch groups.azuredevops.krateo.io group-test \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get groups.azuredevops.krateo.io group-test -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Groups
metadata:
  name: group-test
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  membership:
    # If projectRef is not specified, the group will be created in the organization if is set
    # organization: matteogastaldello
    projectRef: 
      namespace: default
      name: pipeline-proj
  groupName: group-test
  groupRefs:
    - namespace: default
      name: group-sample
  description: Group created from YAML
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

### Step 2: Create the new `GraphGroup` resource

Now, you can create a new `GraphGroup` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#graphgroup-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: GraphGroup
metadata:
  name: group-test
  namespace: azuredevops-system           # Replace with your namespace
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  # Reference to the configuration CR
  configurationRef:
    name: my-graphgroup-config
    namespace: default

  # Azure DevOps organization name
  organization: krateo-kog

  # VSTS group creation fields
  displayName: group-test
  description: Group created from YAML

  # Scope descriptor referencing the project
  # This creates the group at the project level instead of organization level
  # curl -X GET "https://vssps.dev.azure.com/krateo-kog/_apis/graph/descriptors/99837031-4e4e-4753-9a47-73fcc4cba766?api-version=7.2-preview.1" \
  scopeDescriptor: scp.OWI1OGY0OTQtY2Y1OC00Mjg2LWEyYTItZmUxNTNmNzZhNzVl                     
EOF
```

Note that: 
- `groupName` in the old `Groups` resource is mapped to `displayName` in the new `GraphGroup` resource which is aligned with the Azure DevOps REST API which uses `displayName` to define the name of a VSTS group.
- `scopeDescriptor` is used to define the scope where the group will be created. In this case, it is set to a project scope descriptor, so the group will be created at the project level. You can retrieve the scope descriptor of a project by using the Azure DevOps REST API. Currently there is not a resource in Azure DevOps Provider KOG representing descriptors, so you need to retrieve it manually.
- the field `groupRefs` is not supported in the `GraphGroup` schema of Azure DevOps Provider KOG since it maps directly the Azure DevOps REST API.
- Azure DevOps Provider "classic" creates **memberships** in a implicit way inside the reconciliation loop logic of the `Groups` resource, while Azure DevOps Provider KOG requires explicit management of memberships through dedicated resources.
For this reason, you will need to manually create `Membership` resources, see the section about [Memberships](../../../README.md#membership) in the main README for more details.

Note that you need to have already created a `GraphGroupConfiguration` resource that contains the authentication and configuration information for the `GraphGroup` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can wait for the new `GraphGroup` resource to be ready by running the following command:
```sh
kubectl wait graphgroups.azuredevops.ogen.krateo.io/group-test --for condition=Ready=True --namespace azuredevops-system --timeout=300s
```
```sh
graphgroup.azuredevops.ogen.krateo.io/group-test condition met
```

#### Step 3: Delete the old `Groups` resource

At this point, you can proceed to delete the old `Groups` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `Groups` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete groups.azuredevops.krateo.io group-test
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate groups.azuredevops.krateo.io group-test --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `Groups` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
