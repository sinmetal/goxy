package backend

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	var resp *http.Response
	var err error
	log.Infof(ctx, "Mehtod = %s", r.Method)
	if r.Method == "GET" {
		url := fmt.Sprintf("http://104.197.88.95/%s", r.URL.Path)
		req, err := http.NewRequest("GET", url, r.Body)
		if err != nil {
			log.Errorf(ctx, "http.NewRequest error. err = %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for k, v := range r.Header {
			log.Infof(ctx, "%s:%s", k, v)
			for _, value := range v {
				req.Header.Add(k, value)
			}
		}
		log.Infof(ctx, "URL.Path = %s", r.URL.Path)
		log.Infof(ctx, "URL.RawQuery = %s", r.URL.RawQuery)
		resp, err = client.Do(req)
	}
	if err != nil {
		log.Errorf(ctx, "urlfetch error. err = %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode == http.StatusOK {
		w.Header().Set("Cache-Control", "public, max-age=60")
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Errorf(ctx, "response copy error. err = %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
