# `Users` migration example

## Scenario

Note that the name of the resource in Azure DevOps Provider KOG is `User`, while in Azure DevOps Provider (classic) it is `Users`.
Thus, the migration involves also a change in the kind of the resource.

**Starting point**: `Users` resource managed by Azure DevOps Provider (classic).
**Ending point**: `User` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`User` on Azure DevOps) will be the same.

Note that the `Users` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while `User` is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following commands):
```sh
kubectl api-resources | awk 'NR==1{print;next} $1=="users" && $0~/azuredevops/ {print; c++; if(c==2) exit}'
```
Output:
```sh
NAME                                 SHORTNAMES         APIVERSION                            NAMESPACED   KIND
users                                                   azuredevops.krateo.io/v1alpha1        false        Users
users                                                   azuredevops.ogen.krateo.io/v1alpha1   true         User
```

## Migration steps

### Step 1: Pause and set deletion policy on the old `Users` resource

The **starting point** for this migration is the following example of a `Users` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Users
metadata:
  name: user-test
  annotations:
    krateo.io/connector-verbose: "false"
spec:
  deletionPolicy: Delete
  organization: krateo-kog
  user:
    name: matteo.gastaldello@krateo.io
  groupRefs:
    - namespace: default
      name: group-sample
    - namespace: default
      name: group-test
  teamRefs:
    - namespace: default
      name: team-test
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that:
- the `Users` resource is referecing a `ConnectorConfig` resource and `Groups` and `Teams` resource, which are managed by the Azure DevOps Provider "classic".
- the field `groupRefs` is not supported in the `User` schema of Azure DevOps Provider KOG since it maps directly the Azure DevOps REST API.
- the field `teamRefs` is not supported in the `User` schema of Azure DevOps Provider KOG since it maps directly the Azure DevOps REST API.
- Azure DevOps Provider "classic" creates **memberships** in a implicit way inside the reconciliation loop logic of the `Users` resource, while Azure DevOps Provider KOG requires explicit management of memberships through dedicated resources. Specific resources for managing [memberships](https://learn.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-7.1) will likely be introduced in future versions of the Azure DevOps Provider KOG. If this functionality is critical for your use case, it may be advisable to delay the migration until such resources are available.

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate users.azuredevops.krateo.io user-test "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Users` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch users.azuredevops.krateo.io user-test \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get users.azuredevops.krateo.io user-test -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Users
metadata:
  name: user-test
  annotations:
    krateo.io/connector-verbose: "false"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  organization: krateo-kog
  user:
    name: matteo.gastaldello@krateo.io
  groupRefs:
    - namespace: default
      name: group-sample
    - namespace: default
      name: group-test
  teamRefs:
    - namespace: default
      name: team-test
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

### Step 2: Create the new `User` resource

Now, you can create a new `User` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#user-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: User
metadata:
  name: user-test
  namespace: default
  annotations:
    krateo.io/connector-verbose: "true"
    krateo.io/deletion-policy: "orphan"
spec:
  # required: reference to a Configuration CR in the cluster
  configurationRef:
    name: my-user-config
    namespace: azuredevops-system
  organization: krateo-kog

  # Use either mailAddress, principalName, or origin + originId to identify the user.

  mailAddress: "matteo.gastaldello@krateo.io"

  #principalName: "" # "This should be the principal name or upn of the user in the source AD or AAD provider.

  #origin: "" # "This should be the name of the origin provider. Example: github.com"
  #originId: "" # "This should be the object id or sid of the user from the source AD or AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Team Services will communicate with the source provider to fill all other fields on creation."
EOF
```

Note that: 
- `user.name` in the old `Users` resource is mapped to either `mailAddress`, `principalName`, or `origin` + `originId` in the new `User` resource depending on how the user is identified. In this example, we are using `mailAddress`.
- with the respect to the old `Users` resource (Azure DevOps Provider "classic"), the new `User` resource does not support using the name corresponding to the `displayName` field of API responses (e.g., "D421a46a2-30c7-4271-97ce-0b61e436596d Build Service (krateo-kog)"). This is due to the fact that `displayName` is not available in the API request for creating a user, so it is not part of the spec fields of the CR. The Azure DevOps Provider "classic" has custom reconcile logic to handle this case, while Azure DevOps Provider KOG strictly maps the API schema.

Note that you need to have already created a `UserConfiguration` resource that contains the authentication and configuration information for the `User` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can wait for the new `User` resource to be ready by running the following command:
```sh
kubectl wait users.azuredevops.ogen.krateo.io/user-test --for condition=Ready=True --namespace azuredevops-system --timeout=300s
```
```sh
users.azuredevops.ogen.krateo.io/user-test condition met
```

#### Step 3: Delete the old `Users` resource

At this point, you can proceed to delete the old `Users` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `Users` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete users.azuredevops.krateo.io user-test
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate users.azuredevops.krateo.io user-test --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `Users` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
