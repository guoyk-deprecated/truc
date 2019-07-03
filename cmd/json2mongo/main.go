package main

import (
	"encoding/json"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog/log"
	"go.guoyk.net/ext/extmgo"
	"go.guoyk.net/ext/extos"
	"go.guoyk.net/ext/extzerolog"
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
	bulk := extmgo.NewBulk(coll, 1024)

	if err = extos.ReaddirFiles(optWorkspace, extos.ReaddirFilesOptions{
		BeforeFile: func(name string) bool {
			return strings.HasSuffix(name, ".json")
		},
		Handle: func(buf []byte, name string) (err error) {
			log.Info().Str("file", name).Msg("file entered")
			var a Article
			if err = json.Unmarshal(buf, &a); err != nil {
				return
			}
			if err = bulk.Append(bson.M{
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
		},
	}); err != nil {
		return
	}

	if err = bulk.Finish(); err != nil {
		return
	}

}
