
PROBLEM WITH COMPARE DURING FINDBY and GET

[14:24:10.049] DEBUG+3: External resource not up-to-date {
  "apiVersion": "azuredevops.ogen.krateo.io/v1alpha1",
  "kind": "VariableGroup",
  "name": "my-vargroup-basic",
  "namespace": "default",
  "op": "Observe",
  "reason": "ComparisonResult: IsEqual=false, Reason=values differ, FirstValue=map[API_KEY:map[isReadOnly:false isSecret:true value:secret-key-placeholder] API_URL:map[isReadOnly:false isSecret:false value:https://api.example.com] DATABASE_NAME:map[isReadOnly:false isSecret:false value:production_db]], SecondValue=map[API_KEY:map[isSecret:true value:\u003cnil\u003e] API_URL:map[value:https://api.example.com] DATABASE_NAME:map[value:production_db]]"
}


this happens also in findby operation:
current response:
{
  "count": 1,
  "value": [
    {
      "createdBy": {
        "displayName": "Leonardo Vicentini",
        "id": "32915def-23aa-609f-b91d-9a8d94384aa2",
        "uniqueName": "leonardo.vicentini@krateo.io"
      },
      "createdOn": "2026-02-03T15:16:05.2833333Z",
      "id": 13,
      "isShared": false,
      "modifiedBy": {
        "displayName": "Leonardo Vicentini",
        "id": "32915def-23aa-609f-b91d-9a8d94384aa2",
        "uniqueName": "leonardo.vicentini@krateo.io"
      },
      "modifiedOn": "2026-02-03T15:16:05.2833333Z",
      "name": "MyBasicVariableGroup",
      "type": "Vsts",
      "variableGroupProjectReferences": null,
      "variables": {
        "API_KEY": {
          "isSecret": true,
          "value": null
        },
        "API_URL": {
          "value": "https://api.example.com"
        },
        "DATABASE_NAME": {
          "isReadOnly": true,
          "value": "production_db"
        }
      }
    }
  ]
}

the problem is


        "API_KEY": {
          "isSecret": true,
          "value": null
        },

---

PROBLEM WITH COMPARE DURING FINDBY

GET /krateo-kog/project-1-classic/_apis/distributedtask/variablegroups?api-version=7.2-preview.2 HTTP/1.1
Host: dev.azure.com
User-Agent: Go-http-client/1.1
Authorization: Basic Qkt4eXNtQnJ3dlNEQjVvSDMyU1ZVT2lTVFNrRTY5T0V4M1p1VUZiRzN6Tnl4UU9KS0NEeUpRUUo5OUJGQUNBQUFBQTdrYmNHQUFBU0FaRE8xaEhIOkJLeHlzbUJyd3ZTREI1b0gzMlNWVU9pU1RTa0U2OU9FeDNadVVGYkczek55eFFPSktDRHlKUVFKOTlCRkFDQUFBQUE3a2JjR0FBQVNBWkRPMWhISA==
Accept-Encoding: gzip


HTTP/2.0 200 OK
Access-Control-Expose-Headers: Request-Context
Activityid: b210d3fc-4630-4e62-afc5-e14c8fe9ffd7
Cache-Control: no-cache, no-store, must-revalidate
Content-Type: application/json; charset=utf-8; api-version=7.2-preview.2
Date: Tue, 03 Feb 2026 15:03:36 GMT
Expires: -1
P3p: CP="CAO DSP COR ADMa DEV CONo TELo CUR PSA PSD TAI IVDo OUR SAMi BUS DEM NAV STA UNI COM INT PHY ONL FIN PUR LOC CNT"
Pragma: no-cache
Request-Context: appId=cid-v1:e6292eea-fb85-4107-bc0a-339fd28d3647
Set-Cookie: VstsSession=%7B%22PersistentSessionId%22%3A%2208d49d49-be89-4d99-abc7-10b03e0cb001%22%2C%22PendingAuthenticationSessionId%22%3A%2200000000-0000-0000-0000-000000000000%22%2C%22CurrentAuthenticationSessionId%22%3A%2200000000-0000-0000-0000-000000000000%22%2C%22SignInState%22%3A%7B%7D%7D;SameSite=None; domain=.dev.azure.com; expires=Wed, 03-Feb-2027 15:03:36 GMT; path=/; secure; HttpOnly
Strict-Transport-Security: max-age=31536000; includeSubDomains
Vary: Accept-Encoding
X-Cache: CONFIG_NOCACHE
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Msedge-Ref: Ref A: 732BC316E4594BCA89A36225599BCA4D Ref B: GVA201070811060 Ref C: 2026-02-03T15:03:36Z
X-Tfs-Processid: 6a20605c-c952-4dc5-9aa7-e6ad14492c76
X-Tfs-Session: b210d3fc-4630-4e62-afc5-e14c8fe9ffd7
X-Vss-E2eid: b210d3fc-4630-4e62-afc5-e14c8fe9ffd7
X-Vss-Senderdeploymentid: a6c8fbe9-7425-4392-883d-c602e2e7f7eb
X-Vss-Userdata: 32915def-23aa-609f-b91d-9a8d94384aa2:leonardo.vicentini@krateo.io

{
  "count": 1,
  "value": [
    {
      "createdBy": {
        "displayName": "Leonardo Vicentini",
        "id": "32915def-23aa-609f-b91d-9a8d94384aa2",
        "uniqueName": "leonardo.vicentini@krateo.io"
      },
      "createdOn": "2026-02-03T14:44:17.3866667Z",
      "id": 12,
      "isShared": false,
      "modifiedBy": {
        "displayName": "Leonardo Vicentini",
        "id": "32915def-23aa-609f-b91d-9a8d94384aa2",
        "uniqueName": "leonardo.vicentini@krateo.io"
      },
      "modifiedOn": "2026-02-03T14:44:17.3866667Z",
      "name": "MyBasicVariableGroup",
      "type": "Vsts",
      "variableGroupProjectReferences": null,
      "variables": {
        "API_URL": {
          "value": "https://api.example.com"
        },
        "DATABASE_NAME": {
          "isReadOnly": true,
          "value": "production_db"
        }
      }
    }
  ]
}

2026/02/03 15:03:36 isCRUpdated - starting comparison between mg spec and rm
2026/02/03 15:03:36 isCRUpdated - comparing mg spec with rm
2026/02/03 15:03:36 mg spec fields:
2026/02/03 15:03:36 mg spec field: name = MyBasicVariableGroup
2026/02/03 15:03:36 mg spec field: organization = krateo-kog
2026/02/03 15:03:36 mg spec field: project = project-1-classic
2026/02/03 15:03:36 mg spec field: type = Vsts
2026/02/03 15:03:36 mg spec field: variableGroupProjectReferences = [map[name:MyBasicVariableGroup projectReference:map[name:project-1-classic]]]
2026/02/03 15:03:36 mg spec field: variables = map[API_URL:map[isReadOnly:false isSecret:false value:https://api.example.com] DATABASE_NAME:map[isReadOnly:true isSecret:false value:production_db]]
2026/02/03 15:03:36 mg spec field: configurationRef = map[name:my-variablegroup-config namespace:default]
2026/02/03 15:03:36 mg spec field: description = A basic variable group created via Krateo KOG
2026/02/03 15:03:36 rm fields:
2026/02/03 15:03:36 rm field: variableGroupProjectReferences = <nil>
2026/02/03 15:03:36 rm field: variables = map[API_URL:map[value:https://api.example.com] DATABASE_NAME:map[isReadOnly:true value:production_db]]
2026/02/03 15:03:36 rm field: createdBy = map[displayName:Leonardo Vicentini id:32915def-23aa-609f-b91d-9a8d94384aa2 uniqueName:leonardo.vicentini@krateo.io]
2026/02/03 15:03:36 rm field: isShared = false
2026/02/03 15:03:36 rm field: name = MyBasicVariableGroup
2026/02/03 15:03:36 rm field: createdOn = 2026-02-03T14:44:17.3866667Z
2026/02/03 15:03:36 rm field: id = 12
2026/02/03 15:03:36 rm field: modifiedBy = map[displayName:Leonardo Vicentini id:32915def-23aa-609f-b91d-9a8d94384aa2 uniqueName:leonardo.vicentini@krateo.io]
2026/02/03 15:03:36 rm field: modifiedOn = 2026-02-03T14:44:17.3866667Z
2026/02/03 15:03:36 rm field: type = Vsts
2026/02/03 15:03:36 isCRUpdated - performing comparison, before calling CompareExisting
[15:03:36.486] DEBUG+3: External resource not up-to-date {
  "apiVersion": "azuredevops.ogen.krateo.io/v1alpha1",
  "kind": "VariableGroup",
  "name": "my-vargroup-basic",
  "namespace": "default",
  "op": "Observe",
  "reason": "ComparisonResult: IsEqual=false, Reason=values differ (one is nil), FirstValue=[map[name:MyBasicVariableGroup projectReference:map[name:project-1-classic]]], SecondValue=\u003cnil\u003e"
}

the problem is this      "variableGroupProjectReferences": null,
while it should be returned and correctly compared

