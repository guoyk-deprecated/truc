package main

import (
	"encoding/json"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/truc/ext/extmgo"
	"github.com/yankeguo/truc/ext/extos"
	"github.com/yankeguo/truc/ext/extzerolog"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Article struct {
	Lang       string    `json:"lang"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Vendor     string    `json:"vendor"`
	URL        string    `json:"url"`
	OriginalID string    `json:"original_id"`
	Date       time.Time `json:"date"`
}

var (
	optVerbose   bool
	optWorkspace = "/workspace"
)

var (
	coll *mgo.Collection
)

func init() {
	extos.EnvBool(&optVerbose, "VERBOSE")
	extos.EnvStr(&optWorkspace, "WORKSPACE")
}

func main() {
	var err error
	defer extos.Exit(&err)

	extzerolog.SetupPlainZerolog(optVerbose, false)

	var sess *mgo.Session
	if sess, err = extmgo.DialLinkedMongo(); err != nil {
		return
	}
	defer sess.Clone()

	coll = sess.DB("main").C("corpus")

	var dir *os.File
	if dir, err = os.Open(optWorkspace); err != nil {
		return
	}

	var names []string
	if names, err = dir.Readdirnames(-1); err != nil {
		return
	}
	sort.Sort(sort.StringSlice(names))
	for _, name := range names {
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		if err = handle(name); err != nil {
			return
		}
	}
}

func handle(name string) (err error) {
	log.Info().Str("name", name).Msg("file processing")
	var buf []byte
	if buf, err = ioutil.ReadFile(filepath.Join(optWorkspace, name)); err != nil {
		return
	}
	var a Article
	if err = json.Unmarshal(buf, &a); err != nil {
		return
	}
	if err = coll.Insert(bson.M{
		"lang":        a.Lang,
		"title":       a.Title,
		"content":     a.Content,
		"vendor":      a.Vendor,
		"url":         a.URL,
		"original_id": a.OriginalID,
		"date":        a.Date,
	}); err != nil {
		return
	}
	return
}
