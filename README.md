# mesh-app
Application for testing Service Mesh. This application logs every request and responses with incoming request details: url, headers and body.

## Published Docker Image
https://hub.docker.com/layers/iglin/mesh-app/v2.0/images/sha256-e493c240a9f8f6bb6a0f4d1651f03639cc4d869020ca06f1d8c0778b00ba0659?context=repo

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