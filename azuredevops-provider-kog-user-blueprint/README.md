
This is a Helm chart for deploying the Azure DevOps Provider KOG User Blueprint.



differences wrt classic controller


cannot search for 
      "displayName": "D421a46a2-30c7-4271-97ce-0b61e436596d Build Service (krateo-kog)",

      since it is not available in the api request
      so it is not part of the spec fields of the CR

RDC has a stricter reconcile loop and custom reconcile logic cannot be used anymore

