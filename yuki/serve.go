package yuki

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AugustineAurelius/fuufu/yuki/converter"
	skiplist "github.com/AugustineAurelius/fuufu/yuki/skip_list"
)

func Run(host string, port int) {
	list := skiplist.New()

	r := http.NewServeMux()

	registerRead(r, list)
	registerAdd(r, list)

	s := &http.Server{
		Handler: r,
		Addr:    host + ":" + strconv.Itoa(port),
	}

	panic(s.ListenAndServe())
}

type put struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

func registerAdd(r *http.ServeMux, list *skiplist.SkipList) {
	r.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			return
		}

		var put put
		if err := json.NewDecoder(r.Body).Decode(&put); err != nil {
			w.Write(converter.StringToBytes(err.Error()))
			return
		}

		hashedKey := make([]byte, 32)
		hashKey(converter.StringToBytes(put.Key), hashedKey)

		flatedValue, err := flateEncode(put.Value)
		if err != nil {
			w.Write(converter.StringToBytes(err.Error()))
			return
		}

		list.Put(hashedKey, flatedValue)

		w.WriteHeader(201)
	})
}

type get struct {
	Value json.RawMessage `json:"value"`
}

func registerRead(r *http.ServeMux, list *skiplist.SkipList) {
	r.HandleFunc("GET /{key}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(405)
			return
		}

		key := r.PathValue("key")
		if key == "" {
			w.Write(converter.StringToBytes("empty key"))
			return
		}

		hashedKey := make([]byte, 32)
		hashKey(converter.StringToBytes(key), hashedKey)

		flattedValue, ok := list.Get(hashedKey)
		if !ok {
			w.WriteHeader(404)
			w.Write(converter.StringToBytes("not found"))
			return
		}

		val, err := flateDecode(flattedValue)
		if err != nil {
			w.Write(converter.StringToBytes(err.Error()))
			return
		}

		g := get{
			Value: val,
		}
		data, err := json.Marshal(&g)
		if err != nil {
			w.Write(converter.StringToBytes(err.Error()))
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	})

}
