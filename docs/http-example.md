# HTTP Request Example

The following payload is an example of how the collector submits logs via HTTP:

```
POST /url-path HTTP/1.1
Host: URLHOST.xxx
User-Agent: Go-http-client/1.1
Content-Length: 1000
Accept: */*
Accept-Encoding: gzip, deflate
Authorization: Bearer ABC123
Content-Type: application/json

{
  "results": [
    {"id": "1"},
    {"id": "2"},
    {"id": "3"},
  ]
}
```