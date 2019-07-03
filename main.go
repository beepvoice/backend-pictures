package main

import (
  "bytes"
  "context"
  "crypto/rand"
  "encoding/hex"
  "encoding/json"
  "io"
  "log"
  "net/http"
  "os"
  "strings"

  "github.com/joho/godotenv"
  "github.com/julienschmidt/httprouter"
  "github.com/minio/minio-go"
)

const MaxBiteSize = 1024 * 1024 * 100 // 100MB

var listen string
var minioClient *minio.Client
var bucketName string

func main() {
  // Load .env
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
  listen = os.Getenv("LISTEN")
  minioEndpoint := os.Getenv("MINIO_ENDPOINT")
  minioID := os.Getenv("MINIO_ID")
  minioKey := os.Getenv("MINIO_KEY")
  bucketName = os.Getenv("MINIO_BUCKET_NAME")
  minioLocation := os.Getenv("MINIO_LOCATION")

  // Minio client
  minioClient, err = minio.New(minioEndpoint, minioID, minioKey, false)
  if err != nil {
    log.Fatal("Error loading minio")
  }

  // Create bucket if it doesn't exist
  err = minioClient.MakeBucket(bucketName, minioLocation)
  if err != nil {
    exists, err := minioClient.BucketExists(bucketName)
    if err == nil && exists {
      log.Printf("Bucket %s already exists", bucketName)
    } else {
      log.Printf("%s", err)
      log.Fatal("Error creating bucket")
    }
  } else {
    log.Printf("Created bucket %s", bucketName)
  }

  // Routes
  router := httprouter.New()
  router.POST("/upload", AuthMiddleware(Upload))
  router.GET("/picture/:filename", GetFile)

  // Start server
  log.Printf("starting server on %s", listen)
  log.Fatal(http.ListenAndServe(listen, router))
}

// Pull Auth header
type RawClient struct {
  UserId string `json:"userid"`
  ClientId string `json:"clientid"`
}
func AuthMiddleware(next httprouter.Handle) httprouter.Handle {
  return func (w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    ua := r.Header.Get("X-User-Claim")
    if ua == "" {
      http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
  		return
    }

    var client RawClient
    err := json.Unmarshal([]byte(ua), &client)

    if err != nil {
      http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
  		return
    }

    if client.UserId == "" || ClientId == "" {
      http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
  		return
    }

    context := context.WithValue(r.Context(), "user", client)
    next(w, r.WithContext(context), p)
  }
}

func Upload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  client := r.Context().Value("user").(RawClient)
  var buf bytes.Buffer
  file, header, err := r.FormFile("file")
  if err != nil {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
  }
  defer file.Close()

  originalName := strings.Split(header.Filename, ".")

  io.Copy(&buf, file)
  reader := bytes.NewReader(buf.Bytes())

  fileName := RandomHex() + "." + originalName[1]
  options := minio.PutObjectOptions{
    UserMetadata: make(map[string] string),
  }
  options.UserMetadata["owner"] = client.UserId

  _, err = minioClient.PutObject(bucketName, fileName, reader, header.Size, options)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  }

  io.WriteString(w, fileName)
}

func RandomHex() string {
  b := make([]byte, 16)
  _, err := rand.Read(b)
  if err != nil {
    panic("unable to generate 16 bits of randomness")
  }
  return hex.EncodeToString(b)
}

func GetFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  fileName := p.ByName("filename")
  options := minio.GetObjectOptions{}

  reader, err := minioClient.GetObject(bucketName, fileName, options)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
  }

  io.Copy(w, reader)
}
