# backend-pictures

Beep backend proxying [Minio](https://min.io) to act as a fileserver. Need a running instance of `minio`.

**To run this service securely is to run it behind traefik forwarding auth to `backend-auth`**

## Quickstart

```
go build && ./pictures
```

## Environment Variables

Supply environment variables by either exporting them or editing `.env`.

| ENV | Definition | Default |
| --- | ---------- | ------- |
| LISTEN | Host and port number to listen on | :80 |
| MINIO_ENDPOINT | Host and port of minio | minio:9000 |
| MINIO_ID | Client id to use with minio | MINIO_ID |
| MINIO_KEY | Client key to use with minio | MINIO_KEY |
| MINIO_BUCKET_NAME | Name of bucket to store files in | beep |
| MINIO_LOCATION | Minio bucket region | us-east-1 |

## API

All requests need to be passed through `traefik` calling `backend-auth` as Forward Authentication. Otherwise, populate `X-User-Claim` with:

```json
{
  "userid": "<userid>",
  "clientid": "<clientid>"
}
```

### Upload FIle

```
POST /upload
```

Upload a file to be stored. Requires a request of `Content-Type` `multipart/form-data`.

#### Body

| Name | Description |
| ---- | ----------- |
| file | File to be uploaded |

#### Success (200 OK)

Name of the file as stored.

#### Errors

| Code | Description |
| ---- | ----------- |
| 400 | Error parsing file out of request |
| 500 | Error storing file in minio |

---

### Get File

```
GET /picture/:filename
```

Retrieve a picture by filename.

#### Params

| Name | Description |
| ---- | ----------- |
| filename | Name of the file to be retrieved |

#### Success (200 OK)

Image file.

#### Errors 

| Code | Description |
| ---- | ----------- |
| 500 | Error retrieving file from minio |
