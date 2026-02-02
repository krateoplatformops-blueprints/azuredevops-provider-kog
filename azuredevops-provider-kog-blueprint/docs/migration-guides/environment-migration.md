# `Environment` migration example

## Scenario

**Starting point**: `Environment` resource managed by Azure DevOps Provider (classic).
**Ending point**: `Environment` resource managed by Azure DevOps Provider KOG.
Note: the external resource (`Environment` on Azure DevOps) will be the same.

Note that the `Environment` resource is a non-namespaced resource in the context of Azure DevOps Provider (classic), while it is a namespaced resource in the context of Azure DevOps Provider KOG (you can check this by running the following command):
```sh
kubectl api-resources | awk 'NR==1 || $1 == "environments"'
```
Output:
```sh
NAME                                 SHORTNAMES         APIVERSION                            NAMESPACED   KIND
environments                                            azuredevops.krateo.io/v1alpha1        false        Environment
environments                                            azuredevops.ogen.krateo.io/v1alpha1   true         Environment
```

## Migration steps

### Step 1: Pause and set deletion policy on the old `Environment` resource

The **starting point** for this migration is the following example of a `Environment` resource managed by the Azure DevOps Provider (classic):
```yaml
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Environment
metadata:
  name: environment-sample-1
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  name: test-enviroment-1
  description: test description
  deletionPolicy: Delete
  projectRef:
    namespace: default
    name: project-1
  connectorConfigRef:
    namespace: default
    name: connectorconfig-sample
```

Note that:
- the `Environment` resource is referecing a `ConnectorConfig` resource and a `Project` resource, which are both managed by the Azure DevOps Provider "classic".

To ensure that the old version of the resource is not reconciled while you are migrating to the new version, you should set the `krateo.io/paused: true` annotation.
You can do this by running the following commands:
```sh
kubectl annotate environments.azuredevops.krateo.io environment-sample-1 "krateo.io/paused=true"
```

In addition, in order to keep the external resource (on Azure DevOps) after the deletion of the old `Environment` resource on Kubernetes, you need to set `deletionPolicy: Orphan` in the `spec` of the resource.
You can do this with the following command:
```sh
kubectl patch environments.azuredevops.krateo.io environment-sample-1 \
  --type=merge \
  -p '{"spec":{"deletionPolicy": "Orphan"}}'
```

You can check the changes in the resource with the following command:
```sh
kubectl get environments.azuredevops.krateo.io environment-sample-1 -o yaml
```

And the output should look like this:
```diff
apiVersion: azuredevops.krateo.io/v1alpha1
kind: Environment
metadata:
  name: environment-sample-1
  annotations:
    krateo.io/connector-verbose: "true"
+   krateo.io/paused: "true"
spec:
+ deletionPolicy: Orphan
  name: test-enviroment-1
  description: test description
  projectRef:
    namespace: default
    name: project-1
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

### Step 2: Create the new `Environment` resource

Now, you can create a new `Environment` resource using the Azure DevOps Provider KOG following the new schema.
You can find the schema in the specific section of the [README](../../../README.md#environment-schema) file of this chart.
You can apply the following example:
```sh
kubectl apply -f - <<EOF
apiVersion: azuredevops.ogen.krateo.io/v1alpha1
kind: Environment
metadata:
  name: environment-1
  namespace: azuredevops-system # Replace with the desired namespace
  annotations:
    krateo.io/connector-verbose: "true"
spec:
  # required: reference to a Configuration CR in the cluster
  configurationRef:
    name: my-environment-config
    namespace: default
  organization: krateo-kog
  projectId: project-1               # Replace with the target Project ID
  name: test-enviroment-1
  description: "test description"
                
EOF
```

Note that: 
- the `projectRef` field has been replaced with `projectId`, which can be either the `ID` or the `name` of the project in the case of the spec of the `Environment` resource.

In order to dynamically retrieve the IDs, you can use a Helm `lookup` function.

Note that you need to have already created a `EnvironmentConfiguration` resource that contains the authentication and configuration information for the `Project` resource.
See the main [README](../../../README.md#configuration) for more details about configuration resources.

You can wait for the new `Environment` resource to be ready by running the following command:
```sh
kubectl wait environments.azuredevops.ogen.krateo.io/environment-1 --for condition=Ready=True --namespace azuredevops-system --timeout=300s
```
```sh
environment.azuredevops.ogen.krateo.io/environment-1 condition met
```

#### Step 3: Delete the old `Environment` resource

At this point, you can proceed to delete the old `Environment` resource managed by Azure DevOps Provider (classic) (note the different API group).

First, you can delete the old `Environment` resource managed by Azure DevOps Provider (classic):
```sh
kubectl delete environments.azuredevops.krateo.io environment-sample-1
```

You need also to change the `krateo.io/paused` annotation to `false` to allow the resource to be deleted.
Either you `CTRL+C` the previous command (that is hanging) and run the following command, or you can run it in a separate terminal:
```sh
kubectl annotate environments.azuredevops.krateo.io environment-sample-1 --overwrite "krateo.io/paused=false"
```

At this point, the migration of the `Environment` resource from Azure DevOps Provider (classic) to Azure DevOps Provider KOG is complete.
