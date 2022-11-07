package config

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Downloader struct {
	Url string `json:"url" yaml:"url"`
}

func (d *Downloader) Get() ([]byte, error) {
	var ruleb []byte
	u, err := url.Parse(d.Url)
	if err != nil {
		return nil, fmt.Errorf("download error:%s", err.Error())
	}
	switch u.Scheme {
	case "http", "https":
		resp, err := http.Get(d.Url)
		if err != nil {
			return nil, fmt.Errorf("download error:%s", err.Error())
		}
		defer resp.Body.Close()
		ruleb, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("download error:%s", err.Error())
		}
	case "s3":
		query := u.Query()
		sess := session.Must(session.NewSession(&aws.Config{
			Region: aws.String(query.Get("region")),
		}))
		service := s3.New(sess)

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
		defer cancel()
		bucket := u.Host
		key := u.Path
		out, err := service.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			return nil, fmt.Errorf("download error:%s", err.Error())
		}
		defer out.Body.Close()
		ruleb, err = io.ReadAll(out.Body)
		if err != nil {
			return nil, fmt.Errorf("download error:%s", err.Error())
		}
	}
	return ruleb, nil
}
