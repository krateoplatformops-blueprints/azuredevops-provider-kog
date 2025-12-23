
This is a Helm chart for deploying the Azure DevOps Provider KOG Project Blueprint.


how to restore a deleted project in Azure DevOps.


if project was already managed by a CR in kubernets
check if the CR still exists and the status of the CR contains the ID




if the project was deleted 
the cr never existed or is not anymore present in the cluster 

patch the ProjectCOnfiguration resource
stateFIlter to deleted

apply the CR representing the project currently in a deleted state
which we want to restore

the findby will use the state filter to find the deleted project
