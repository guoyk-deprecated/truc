package main

import (
	"bufio"
	"flag"
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
	optVerbose    bool
	optHost       string
	optDB         string
	optCollection string
	optDir        string
)

var (
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

func main() {
	var err error
	defer cmdutil.Exit(&err)

	flag.BoolVar(&optVerbose, "verbose", false, "verbose mode")
	flag.StringVar(&optHost, "host", "localhost:27017", "mongodb host")
	flag.StringVar(&optDB, "db", "main", "mongodb url")
	flag.StringVar(&optCollection, "collection", "library", "mongodb collection")
	flag.StringVar(&optDir, "dir", ".", "data directory")
	flag.Parse()

	cmdutil.SetupPlainZerolog(optVerbose, true)

	var sess *mgo.Session
	if sess, err = mgo.Dial(optHost); err != nil {
		return
	}

	var colle *mgo.Collection
	colle = sess.DB(optDB).C(optCollection)

	var dir *os.File
	if dir, err = os.Open(optDir); err != nil {
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
		if err = handle(name, colle); err != nil {
			return
		}
		break
	}
}

func handle(name string, colle *mgo.Collection) (err error) {
	log.Info().Str("name", name).Msg("file processing")
	var file *os.File
	if file, err = os.Open(filepath.Join(optDir, name)); err != nil {
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	var line int
	docs := make([]interface{}, 0, batchSize)
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
		if len(str) == 0 {
			continue
		}
		tokens := tokenize(str)
		if len(tokens) == 0 {
			log.Info().Str("file", name).Int("line", line).Strs("tokens", tokens).Str("content", str).Msg("line no token")
		} else if len(tokens) == 1 {
			log.Debug().Str("file", name).Int("line", line).Strs("tokens", tokens).Str("content", str).Msg("line 1 token")
		}
		docs = append(docs, bson.M{"tokens": tokens, "content": str, "source": "easedown"})
		if len(docs) >= batchSize {
			if err = colle.Insert(docs...); err != nil {
				return
			}
			docs = docs[0:0]
		}
	}
	if len(docs) > 0 {
		if err = colle.Insert(docs...); err != nil {
			return
		}
	}
	log.Info().Str("file", name).Int("line", line).Msg("file processed")
	return
}
