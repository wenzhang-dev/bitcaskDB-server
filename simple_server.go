package main

import (
	"encoding/base64"
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	httpprof "net/http/pprof"

	"github.com/gorilla/mux"

	"github.com/wenzhang-dev/bitcaskDB"
)

var (
	port = flag.String("port", "8090", "http server port")
	dir  = flag.String("dir", "./bitcaskDB", "working directory")

	logFile       = flag.String("log_file", "", "log file name")
	logDir        = flag.String("log_dir", "", "log directory")
	logLevel      = flag.Int("log_level", 0, "log level")
	logSize       = flag.Int("log_size", 0, "log max size")
	logMaxBackups = flag.Int("log_max_backups", 0, "max log backups")

	walSize      = flag.Int("wal_size", 1024*1024*1024, "max wal file size")
	manifestSize = flag.Int("manifest_size", 1024*1024, "max manifest file size")

	indexCap             = flag.Int("index_cap", 100000, "index capacity")
	indexLimit           = flag.Int("index_limited", 80000, "max elements")
	indexSampleKeys      = flag.Int("index_sample_keys", 3, "sample key counts")
	indexEvictionPoolCap = flag.Int("index_eviction_pool_cap", 64, "eviction pool capacity")

	diskLimit = flag.Int("disk_limited", 0, "disk usage limit")

	nsSize   = flag.Int("ns_size", 0, "namespace length")
	etagSize = flag.Int("etag_size", 0, "etag length")

	compactionInternal = flag.Int("compaction_interval", 0, "compaction interval")
	checkDiskInterval  = flag.Int("check_disk_interval", 0, "check disk usage interval")

	pickRatio = flag.Float64("pick_ratio", 0, "picker ratio")

	disableCompaction = flag.Bool("disable_compaction", false, "disable compaction")

	recordBufferSize = flag.Int("record_buf_size", 0, "record buffer size")
)

type Server struct {
	db *bitcask.DBImpl
}

func NewServer(db *bitcask.DBImpl) *Server {
	return &Server{db}
}

func (s *Server) ValidNs(ns []byte) bool {
	opts := bitcask.GetOptions()
	return int(opts.NsSize) == len(ns)
}

func (s *Server) getNsAndKey(r *http.Request) (ns, key []byte) {
	return []byte(r.URL.Query().Get("ns")), []byte(r.URL.Query().Get("key"))
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	ns, key := s.getNsAndKey(r)
	if !s.ValidNs(ns) || len(key) == 0 {
		http.Error(w, "get error", http.StatusBadRequest)
		return
	}

	s.getHandler([]byte(ns), []byte(key), w, r)
}

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	ns, key := s.getNsAndKey(r)
	if !s.ValidNs(ns) || len(key) == 0 {
		http.Error(w, "write error", http.StatusBadRequest)
		return
	}

	s.postHandler(ns, key, w, r)
}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ns, key := s.getNsAndKey(r)
	if !s.ValidNs(ns) || len(key) == 0 {
		http.Error(w, "delete error", http.StatusBadRequest)
		return
	}

	s.deleteHandler(ns, key, w, r)
}

func (s *Server) getBase64NsAndKey(r *http.Request) (ns, key []byte) {
	q1 := r.URL.Query().Get("ns")
	q2 := r.URL.Query().Get("key")
	if q2 == "" {
		log.Printf("q1: %v, q2: %v\n", q1, q2)
		return nil, nil
	}

	var err error
	ns, err = base64.URLEncoding.DecodeString(q1)
	if err != nil {
		return nil, nil
	}

	key, err = base64.URLEncoding.DecodeString(q2)
	if err != nil {
		return nil, nil
	}

	return ns, key
}

func (s *Server) GetBytesHandler(w http.ResponseWriter, r *http.Request) {
	ns, key := s.getBase64NsAndKey(r)
	if !s.ValidNs(ns) || len(key) == 0 {
		http.Error(w, "get error", http.StatusBadRequest)
		return
	}

	s.getHandler(ns, key, w, r)
}

func (s *Server) PostBytesHandler(w http.ResponseWriter, r *http.Request) {
	ns, key := s.getBase64NsAndKey(r)
	if !s.ValidNs(ns) || len(key) == 0 {
		http.Error(w, "write error", http.StatusBadRequest)
		return
	}

	s.postHandler(ns, key, w, r)
}

func (s *Server) DeleteBytesHandler(w http.ResponseWriter, r *http.Request) {
	ns, key := s.getBase64NsAndKey(r)
	if !s.ValidNs(ns) || len(key) == 0 {
		http.Error(w, "delete error", http.StatusBadRequest)
		return
	}

	s.deleteHandler(ns, key, w, r)
}

func (s *Server) getHandler(ns, key []byte, w http.ResponseWriter, _ *http.Request) {
	// TODO: handle meta data
	val, _, err := s.db.Get(ns, key, nil)
	if err != nil {
		http.Error(w, "no found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(val)
}

func (s *Server) postHandler(ns, key []byte, w http.ResponseWriter, r *http.Request) {
	val, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "write error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// TODO: handle meta data
	meta := bitcask.NewMeta(nil)
	if err = s.db.Put(ns, key, val, meta, nil); err != nil {
		http.Error(w, "write error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deleteHandler(ns, key []byte, w http.ResponseWriter, _ *http.Request) {
	if err := s.db.Delete(ns, key, nil); err != nil {
		http.Error(w, "delete error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) CheckHealth(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ensureDir(path string) error {
	if path == "" {
		path = "./"
	}

	return os.MkdirAll(path, 0o755)
}

func main() {
	flag.Parse()

	opts := &bitcask.Options{
		Dir:                       *dir,
		LogFile:                   *logFile,
		LogDir:                    *logDir,
		LogLevel:                  int8(*logLevel),
		LogMaxSize:                uint64(*logSize),
		LogMaxBackups:             uint64(*logMaxBackups),
		WalMaxSize:                uint64(*walSize),
		ManifestMaxSize:           uint64(*manifestSize),
		IndexCapacity:             uint64(*indexCap),
		IndexLimited:              uint64(*indexLimit),
		IndexEvictionPoolCapacity: uint64(*indexEvictionPoolCap),
		IndexSampleKeys:           uint64(*indexSampleKeys),
		DiskUsageLimited:          uint64(*diskLimit),
		NsSize:                    uint64(*nsSize),
		EtagSize:                  uint64(*etagSize),
		CompactionTriggerInterval: uint64(*compactionInternal),
		CheckDiskUsageInterval:    uint64(*checkDiskInterval),
		CompactionPickerRatio:     *pickRatio,
		DisableCompaction:         *disableCompaction,
		RecordBufferSize:          uint64(*recordBufferSize),
	}

	if err := ensureDir(opts.Dir); err != nil {
		log.Fatalf("ensure dir failed: %v", err)
	}

	db, err := bitcask.NewDB(opts)
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	defer db.Close()

	s := NewServer(db)
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/healthz", s.CheckHealth).Methods("GET")

	r.HandleFunc("/api/v1/kv", s.GetHandler).Methods("GET")
	r.HandleFunc("/api/v1/kv", s.PostHandler).Methods("POST", "PUT")
	r.HandleFunc("/api/v1/kv", s.DeleteHandler).Methods("DELETE")

	r.HandleFunc("/api/v1/binary/kv", s.GetBytesHandler).Methods("GET")
	r.HandleFunc("/api/v1/binary/kv", s.PostBytesHandler).Methods("POST", "PUT")
	r.HandleFunc("/api/v1/binary/kv", s.DeleteBytesHandler).Methods("DELETE")

	r.PathPrefix("/debug/pprof/").HandlerFunc(httpprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", httpprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", httpprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", httpprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", httpprof.Trace)

	if err = http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatalf("http failed: %v", err)
	}
}
