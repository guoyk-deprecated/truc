package main

import (
	"encoding/json"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/truc/cmdutil"
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
	optMongoHost string
	optMongoPort string
	optMongoUser string
	optMongoPass string
	optWorkspace = "/workspace"
)

var (
	coll *mgo.Collection
)

func init() {
	cmdutil.EnvBool(&optVerbose, "VERBOSE")
	cmdutil.EnvStr(&optMongoHost, "MONGO_PORT_27017_TCP_ADDR")
	cmdutil.EnvStr(&optMongoPort, "MONGO_PORT_27017_TCP_PORT")
	cmdutil.EnvStr(&optMongoUser, "MONGO_ENV_MONGO_INITDB_ROOT_USERNAME")
	cmdutil.EnvStr(&optMongoPass, "MONGO_ENV_MONGO_INITDB_ROOT_PASSWORD")
	cmdutil.EnvStr(&optWorkspace, "WORKSPACE")
}

func main() {
	var err error
	defer cmdutil.Exit(&err)

	cmdutil.SetupPlainZerolog(optVerbose, false)

	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/main?authSource=admin", optMongoUser, optMongoPass, optMongoHost, optMongoPort)
	log.Info().Str("mongo_url", url).Msg("connect mongodb")

	var sess *mgo.Session
	if sess, err = mgo.Dial(url); err != nil {
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
