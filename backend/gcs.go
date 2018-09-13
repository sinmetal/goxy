package backend

import (
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/gcs/", handlerGCS)
}

func handlerGCS(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	log.Infof(ctx, "Handler GCS Path : %s", r.URL.Path)
	paths := strings.Split(r.URL.Path, "/")

	log.Infof(ctx, "%+v", paths)
	bucketName := paths[2]
	objectName := strings.Join(paths[3:], "/")
	log.Infof(ctx, "Bucket:%s, Object:%s", bucketName, objectName)

	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bkt := client.Bucket(bucketName)
	reader, err := bkt.Object(objectName).NewReader(ctx)
	if err != nil {
		if err == storage.ErrBucketNotExist {
			http.Error(w, "ErrBucketNotExist", http.StatusNotFound)
			return
		}
		if err == storage.ErrObjectNotExist {
			http.Error(w, "ErrObjectNotExist", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	attrs, err := bkt.Object(objectName).Attrs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof(ctx, "Content-Type:%s", attrs.ContentType)

	w.Header().Set("Content-Type", attrs.ContentType)
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, reader)
	if err != nil {
		log.Errorf(ctx, "response copy error. err = %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}