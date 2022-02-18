package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	sql "github.com/FloatTech/sqlite"
	"github.com/fumiama/imago"
	"github.com/fumiama/loliana/database"
	"github.com/fumiama/loliana/lolicon"
	"github.com/sirupsen/logrus"
)

var pidpreg = regexp.MustCompile(`\d+_p\d+`)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <inputjson> <outputdb>")
		return
	}
	logrus.SetLevel(logrus.ErrorLevel)

	inp := os.Args[1]
	db := sql.Sqlite{DBPath: os.Args[2]}

	inpf, err := os.Open(inp)
	if err != nil {
		panic(err)
	}

	err = db.Create("picture", &database.Picture{})
	if err != nil {
		panic(err)
	}

	var items []lolicon.Item
	err = json.NewDecoder(inpf).Decode(&items)
	if err != nil {
		panic(err)
	}
	_ = inpf.Close()

	nohupf, err := os.Open("nohup.out.txt")
	if err != nil {
		panic(err)
	}

	nohupmap := make(map[string]*database.Picture, 40000)
	s := bufio.NewScanner(nohupf)
	for s.Scan() {
		r := strings.Split(s.Text(), " ")
		pidp := r[0]
		dp := r[1]
		dh := r[2]
		ms := r[3]
		nohupmap[pidp] = &database.Picture{
			PidP:     pidp,
			DatePath: dp,
			DHash:    dh,
			Md5:      ms,
		}
	}

	for _, item := range items {
		pidp := pidpreg.FindString(item.Original)
		tags, err := json.Marshal(item.Tags)
		if err != nil {
			panic(err)
		}
		p, ok := nohupmap[pidp]
		if ok {
			p.UID = item.UID
			p.Width = item.Width
			p.Height = item.Height
			p.Title = item.Title
			p.Author = item.Author
			p.R18 = item.R18
			p.Tags = imago.BytesToString(tags)
			p.Ext = item.Ext
			err = db.Insert("picture", p)
			if err != nil {
				logrus.Errorln("insert img", pidp, "error:", err)
				continue
			}
			logrus.Println("succ", p)
		} else {
			logrus.Errorln("fail", pidp)
		}
	}
	_ = db.Close()
}
