package main

import (
	"bufio"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/truc/cmdutil"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	batchSize = 1000
)

var (
	optVerbose   bool
	optDryrun    = true
	optMongoHost string
	optMongoPort string
	optMongoUser string
	optMongoPass string
	optWorkspace = "/workspace"
)

var (
	coll *mgo.Collection
	bulk = make([]interface{}, 0, batchSize)

	separator = regexp.MustCompile("-{3,}|[|\\s,@:]")
)

func sanitize(strs []string) (out []string) {
	out = make([]string, 0, len(strs))
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if len(str) > 0 {
			out = append(out, str)
		}
	}
	return
}

func tokenize(str string) []string {
	return sanitize(separator.Split(str, -1))
}

func init() {
	cmdutil.EnvBool(&optVerbose, "VERBOSE")
	cmdutil.EnvBool(&optDryrun, "DRYRUN")
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

	coll = sess.DB("main").C("selib")

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
		if !strings.HasPrefix(name, "part-") {
			continue
		}
		if err = handle(name); err != nil {
			return
		}
		if optDryrun {
			break
		}
	}
}

func handle(name string) (err error) {
	log.Info().Str("name", name).Msg("file processing")
	var file *os.File
	if file, err = os.Open(filepath.Join(optWorkspace, name)); err != nil {
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	var line int
	for {
		line++

		var str string
		if str, err = r.ReadString('\n'); err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		str = strings.TrimSpace(str)
		if len(str) == 0 || len(str) >= 1024 {
			continue
		}
		tokens := tokenize(str)
		if len(tokens) == 0 {
			log.Info().Str("file", name).Int("line", line).Strs("tokens", tokens).Str("content", str).Msg("line no token")
		} else if len(tokens) == 1 {
			log.Debug().Str("file", name).Int("line", line).Strs("tokens", tokens).Str("content", str).Msg("line 1 token")
		}
		if err = appendBulk(bson.M{"tokens": tokens, "content": str, "source": "easedown"}); err != nil {
			return
		}
		if optDryrun && line > 10 {
			break
		}
	}
	if err = finishBulk(); err != nil {
		return
	}
	log.Info().Str("file", name).Int("lines", line).Msg("file processed")
	return
}

func appendBulk(doc bson.M) (err error) {
	bulk = append(bulk, doc)
	if len(bulk) >= batchSize {
		if err = coll.Insert(bulk...); err != nil {
			return
		}
		bulk = bulk[0:0]
	}
	return
}

func finishBulk() (err error) {
	if len(bulk) > 0 {
		if err = coll.Insert(bulk...); err != nil {
			return
		}
		bulk = bulk[0:0]
	}
	return
}
