package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog/log"
	"go.guoyk.net/ext"
	"go.guoyk.net/ext/extmgo"
	"go.guoyk.net/ext/extos"
	"go.guoyk.net/ext/extzerolog"
	"regexp"
)

var (
	optVerbose   bool
	optWorkspace = "/workspace"
	optSource    string

	separator = regexp.MustCompile("-{3,}|[|\\s,@:]")
)

func tokenize(str string) []string {
	return ext.SanitizeStrSlice(separator.Split(str, -1))
}

func init() {
	extos.EnvBool(&optVerbose, "VERBOSE")
	extos.EnvStr(&optWorkspace, "WORKSPACE")
	extos.EnvStr(&optSource, "SOURCE")
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

	coll := sess.DB("main").C("selib")
	bulk := extmgo.NewBulk(coll, 1024)
	var oldname string

	if err = extos.ReaddirLines(optWorkspace, extos.ReaddirLinesOptions{
		Handle: func(line0 []byte, name string, lineno int) (err error) {
			line := string(line0)
			if oldname != name {
				log.Info().Str("file", name).Msg("file entered")
				oldname = name
			}
			if len(line) == 0 || len(line) >= 1024 {
				return
			}
			tokens := tokenize(line)
			if len(tokens) == 0 {
				log.Info().Str("file", name).Int("lineno", lineno).Strs("tokens", tokens).Str("line", line).Msg("line no token")
			} else if len(tokens) == 1 {
				log.Debug().Str("file", name).Int("lineno", lineno).Strs("tokens", tokens).Str("line", line).Msg("line 1 token")
			}
			if err = bulk.Append(bson.M{"tokens": tokens, "content": line, "source": optSource}); err != nil {
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
