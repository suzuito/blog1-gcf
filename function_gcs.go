package blog1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/rs/zerolog/log"
	"github.com/suzuito/blog1-go/inject"
	"github.com/suzuito/blog1-go/usecase"
	"golang.org/x/xerrors"
)

type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

func BlogUpdateArticle(ctx context.Context, ev GCSEvent) error {
	cdeps, closeFunc, err := inject.NewContextDepends(ctx, env)
	if err != nil {
		return xerrors.Errorf("Cannot inject.NewContextDepends : %w", err)
	}
	defer closeFunc()
	if ev.Bucket != env.GCPBucketArticle {
		return xerrors.Errorf("Invalid backet name exp:%s != real:%s", env.GCPBucketArticle, ev.Bucket)
	}
	u := usecase.NewImpl(env, cdeps.DB, cdeps.Storage, gdeps.MDConverter)
	log.Info().Str("file", fmt.Sprintf("%s/%s", ev.Bucket, ev.Name)).Msgf("Update")
	if err := u.UpdateArticle(ctx, ev.Name); err != nil {
		return xerrors.Errorf("Cannot u.SyncArticles : %w", err)
	}
	return nil
}

var GCSFunctions = map[string]func(context.Context, GCSEvent) error{
	"BlogUpdateArticle": BlogUpdateArticle,
}

func RegisterLocalRunner(ctx context.Context) error {
	for p := range GCSFunctions {
		fn := GCSFunctions[p]
		p = fmt.Sprintf("/bg-gcs/%s", p)
		if err := funcframework.RegisterHTTPFunctionContext(ctx, p, func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "err: %+v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			ev := GCSEvent{}
			if err := json.Unmarshal(body, &ev); err != nil {
				fmt.Fprintf(w, "err: %+v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if err := fn(ctx, ev); err != nil {
				fmt.Fprintf(w, "err: %+v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "err: %+v\n", err)
		}
		fmt.Println(p)
	}
	return nil
}
