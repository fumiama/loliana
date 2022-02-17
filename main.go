package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	sqlite "github.com/FloatTech/sqlite"
	para "github.com/fumiama/go-hide-param"
	"github.com/fumiama/imago"
	"github.com/sirupsen/logrus"

	"github.com/fumiama/loliana/database"
	"github.com/fumiama/loliana/lolicon"
)

type ItemList = []lolicon.Item

var (
	storage     imago.StorageInstance
	mu          sync.RWMutex
	wg          sync.WaitGroup
	pidpreg     = regexp.MustCompile(`\d+_p\d+`)
	datepathreg = regexp.MustCompile(`\d{4}/\d{2}/d{2}/d{2}/d{2}/d{2}`)
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: <simple-storage-apiurl> <simple-storage-key> <input.json> <output.db>")
		return
	}

	apiurl := os.Args[1]
	key := os.Args[2]
	input := os.Args[3]
	output := os.Args[4]

	storage = imago.NewRemoteStorage(apiurl, key)
	if storage == nil {
		panic("wrong remote para")
	}
	para.Hide(2)

	err := storage.ScanImgs("img")
	if err != nil {
		panic(err)
	}

	inpf, err := os.Open(input)
	if err != nil {
		panic(err)
	}

	db := &sqlite.Sqlite{DBPath: output}
	err = db.Create("picture", &database.Picture{})
	if err != nil {
		panic(err)
	}

	var items ItemList
	err = json.NewDecoder(inpf).Decode(&items)
	if err != nil {
		panic(err)
	}

	totalcnt := len(items)
	n := runtime.NumCPU()
	cntperthread := totalcnt / n
	i := 0
	for n--; i < n; i++ {
		wg.Add(1)
		logrus.Println("scanning from", i*cntperthread, "to", (i+1)*cntperthread)
		go scan(items[i*cntperthread:(i+1)*cntperthread], db)
		time.Sleep(time.Millisecond * 100)
	}
	wg.Add(1)
	logrus.Println("scanning from", i*cntperthread, "to", len(items))
	scan(items[i*cntperthread:], db)
	wg.Wait()
	logrus.Println("all done")
}

func scan(items ItemList, db *sqlite.Sqlite) {
	client := &http.Client{}
	for _, item := range items {
		pidp := pidpreg.FindString(item.Original)
		mu.RLock()
		if db.CanFind("picture", "where pidp="+pidp) {
			mu.RUnlock()
			continue
		}
		mu.RUnlock()
		request, _ := http.NewRequest("GET", strings.ReplaceAll(item.Original, "i.pixiv.cat", "i.pximg.net"), nil)
		request.Header.Set("Host", "i.pximg.net")
		request.Header.Set("Referer", "https://www.pixiv.net/")
		request.Header.Set("Accept", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0) Gecko/20100101 Firefox/6.0")
		resp, err := client.Do(request)
		if err != nil {
			logrus.Errorln("get img", pidp, "resp error:", err)
			continue
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logrus.Errorln("get img", pidp, "body error:", err)
			continue
		}
		dp := datepathreg.FindString(item.Original)
		stat, dh := storage.SaveImgBytes(data, "img", true, 0)
		if dh == "" {
			logrus.Errorln("save img", pidp, "bytes error:", stat)
			continue
		}
		m := md5.Sum(data)
		ms := hex.EncodeToString(m[:])
		tags, err := json.Marshal(item.Tags)
		if err != nil {
			logrus.Errorln("parse img", pidp, "tags error:", err)
			continue
		}
		mu.Lock()
		err = db.Insert("picture", &database.Picture{
			PidP:     pidp,
			UID:      item.UID,
			Width:    item.Width,
			Height:   item.Height,
			Title:    item.Title,
			Author:   item.Author,
			R18:      item.R18,
			Tags:     imago.BytesToString(tags),
			Ext:      item.Ext,
			DatePath: dp,
		})
		for _, tag := range item.Tags {
			err = db.Create(tag, &database.Tag{})
			if err != nil {
				logrus.Errorln("create tag", tag, "error:", err)
				continue
			}
			err = db.Insert(tag, &database.Tag{
				PidP: pidp,
				UID:  item.UID,
			})
			if err != nil {
				logrus.Errorln("insert tag", tag, "error:", err)
				continue
			}
		}
		mu.Unlock()
		if err != nil {
			logrus.Errorln("insert img", pidp, "error:", err)
			continue
		}
		logrus.Println("succ", pidp, dp, dh, ms)
		logrus.Debugln("successfully insert analyzed img", pidp, "datepath", dp, "dhash", dh, "md5", ms, "tags", tags)
	}
	wg.Done()
}
