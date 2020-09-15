## REST API postman collection

- [Collection environment](postman/collection-env)

```json
{
  "info": {
    "_postman_id": "1588ecb0-a996-4286-a3d2-3faf1e3554a1",
    "name": "Krane API Collection",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Authentication",
      "item": [
        {
          "name": "Authenticate",
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n\t\"request_id\": \"310edecd-777f-4ce6-b2f0-99f306c1beff\",\n\t\"token\" :\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InNlc3Npb25faWQiOiIxZTUxNTgwZC00ZjVhLTQyYWMtYjgwOC1lYjBmNTI3MmFmYjUifSwiZXhwIjoxNjMxMDM3MzkyLCJpc3MiOiJrcmFuZSJ9.__ZP0YbJDljdqcsuJbbtN6oZd7B8LUjxubZudOIgP6A\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{scheme}}://{{host}}/auth",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["auth"]
            }
          },
          "response": []
        },
        {
          "name": "Login",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{scheme}}://{{host}}/login",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["login"]
            }
          },
          "response": []
        }
      ],
      "protocolProfileBehavior": {}
    },
    {
      "name": "Activity",
      "item": [
        {
          "name": "Get recent history",
          "request": {
            "auth": {
              "type": "bearer",
              "bearer": [
                {
                  "key": "token",
                  "value": "{{authKey}}",
                  "type": "string"
                }
              ]
            },
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{host}}/history",
              "host": ["{{host}}"],
              "path": ["history"]
            }
          },
          "response": []
        }
      ],
      "protocolProfileBehavior": {}
    },
    {
      "name": "Deployments",
      "item": [
        {
          "name": "Create deployment",
          "request": {
            "auth": {
              "type": "bearer",
              "bearer": [
                {
                  "key": "token",
                  "value": "{{token}}",
                  "type": "string"
                }
              ]
            },
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n    \"name\": \"{{name}}\",\n    \"image\": \"{{image}}\",\n    \"alias\": [\"biensupernice.com\", \"krane.sh\"],\n    \"env\": {\n        \"NODE_ENV\": \"dev\"\n    },\n    \"secrets\": {\n        \"TOKEN\": \"@token\"\n    },\n    \"volumes\": {\n        \"/var/sock\": \"/docker\"\n    }\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{scheme}}://{{host}}/deployments",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["deployments"]
            }
          },
          "response": []
        },
        {
          "name": "Get deployment",
          "request": {
            "auth": {
              "type": "bearer",
              "bearer": [
                {
                  "key": "token",
                  "value": "{{token}}",
                  "type": "string"
                }
              ]
            },
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{scheme}}://{{host}}/deployments/{{name}}",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["deployments", "{{name}}"]
            }
          },
          "response": []
        },
        {
          "name": "Get all deployments",
          "request": {
            "auth": {
              "type": "bearer",
              "bearer": [
                {
                  "key": "token",
                  "value": "{{token}}",
                  "type": "string"
                }
              ]
            },
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{scheme}}://{{host}}/deployments",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["deployments"]
            }
          },
          "response": []
        },
        {
          "name": "Delete deployment",
          "request": {
            "auth": {
              "type": "bearer",
              "bearer": [
                {
                  "key": "token",
                  "value": "{{token}}",
                  "type": "string"
                }
              ]
            },
            "method": "DELETE",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": ""
            },
            "url": {
              "raw": "{{scheme}}://{{host}}/deployments/{{name}}",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["deployments", "{{name}}"]
            }
          },
          "response": []
        }
      ],
      "protocolProfileBehavior": {}
    },
    {
      "name": "Sessions",
      "item": [
        {
          "name": "Get all sessions",
          "request": {
            "auth": {
              "type": "bearer",
              "bearer": [
                {
                  "key": "token",
                  "value": "{{token}}",
                  "type": "string"
                }
              ]
            },
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{scheme}}://{{host}}/sessions",
              "protocol": "{{scheme}}",
              "host": ["{{host}}"],
              "path": ["sessions"]
            }
          },
          "response": []
        }
      ],
      "protocolProfileBehavior": {}
    }
  ],
  "variable": [
    {
      "id": "baseUrl",
      "key": "baseUrl",
      "value": "http://petstore.swagger.io/v1",
      "type": "string"
    }
  ],
  "protocolProfileBehavior": {}
}
```
