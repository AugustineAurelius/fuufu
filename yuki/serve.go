package yuki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/AugustineAurelius/fuufu/yuki/converter"
	skiplist "github.com/AugustineAurelius/fuufu/yuki/skip_list"
	"github.com/AugustineAurelius/fuufu/yuki/wal"
)

func Run(host string, port int) {
	list := skiplist.New()

	wal, err := wal.OpenWAL()
	if err != nil {
		panic(err)
	}
	defer wal.Close()

	if err = wal.FillMemtable(list); err != nil {
		fmt.Println(err)
	}

	r := http.NewServeMux()

	registerRead(r, list)
	registerAdd(r, list, wal)
	registerPprof(r)

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

func registerAdd(r *http.ServeMux, list *skiplist.SkipList, wal *wal.Wal) {

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

		var buf bytes.Buffer
		if err := flateEncode(put.Value, &buf); err != nil {
			w.Write(converter.StringToBytes(err.Error()))
			return
		}

		if err := wal.Add(hashedKey, buf.Bytes()); err != nil {
			slog.Error(err.Error())
		}

		list.Put(hashedKey, buf.Bytes())

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

func registerPprof(r *http.ServeMux) {
	r.HandleFunc("GET /debug/pprof/", pprof.Index)
	r.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	r.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
}
