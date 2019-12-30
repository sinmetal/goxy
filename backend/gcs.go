package backend

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
)

func HandlerGCS(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	paths := strings.Split(r.URL.Path, "/")

	bucketName := paths[2]
	objectName := strings.Join(paths[3:], "/")
	fmt.Printf("Bucket:%s, Object:%s\n", bucketName, objectName)

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
	defer func() {
		if err := reader.Close(); err != nil {
			fmt.Printf("failed gcs reader close. err=%+v\n", err)
		}
	}()

	attrs, err := bkt.Object(objectName).Attrs(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", attrs.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=600")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, reader)
	if err != nil {
		fmt.Printf("response copy error. err = %+v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
