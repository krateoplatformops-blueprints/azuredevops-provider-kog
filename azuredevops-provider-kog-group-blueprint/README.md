
This is a Helm chart for deploying the Azure DevOps Provider KOG GraphGroup Blueprint.


Three Group Types:
VSTS Native Groups (GraphGroupVstsCreationContext)
Uses: displayName, description
Created within Azure DevOps

AAD Groups via Mail Address (GraphGroupMailAddressCreationContext)
Uses: mailAddress
Materializes existing AAD group by email

AAD Groups via Origin ID (GraphGroupOriginIdCreationContext)
Uses: originId
Materializes existing AAD group by object ID