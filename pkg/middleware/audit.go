package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/AugustineAurelius/fuufu/internal/repository/event"
	"github.com/google/uuid"
)

type Geolocator interface {
	GetLocationRaw(ip string) ([]byte, error)
}

func AuditMiddleware(geoClient Geolocator, db *event.CommandRepository, next http.Handler, excludedPath ...string) http.Handler {
	slices.Sort(excludedPath)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := slices.BinarySearch(excludedPath, r.URL.Path); !ok {
			eventRaw, err := geoClient.GetLocationRaw(strings.Split(r.RemoteAddr, ":")[0])
			if err != nil {
				fmt.Println(err)
			}

			if err = db.Create(r.Context(), &event.Event{
				ID:        uuid.New(),
				Payload:   eventRaw,
				CreatedAt: time.Now().UTC(),
			}); err != nil {
				fmt.Println(err)
			}
		}

		next.ServeHTTP(w, r)
	})
}
