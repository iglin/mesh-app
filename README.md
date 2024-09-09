# mesh-app
Application for testing Service Mesh. This application logs every request and responses with incoming request details: url, headers and body.

## Published Docker Image
https://hub.docker.com/layers/iglin/mesh-app/v2.2/images/sha256-e5473dd2a44ca7d38e64d14f32c9be4cf3a1645f1f7a7c168bae96b00ba763f0?context=repo

## Proxy requests
This service is also capable of sending requests to some other destination, to do so include the following request body in your HTTP request to the service:

```json
{
  "url": "<target URL>",
  "method": "<HTTP method>",
  "headers": {
    "<headerName>": [ "<headerValue1>", "<headerValue2>" ]
  }
}
```