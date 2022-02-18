package database

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"image"
	"io"
	"os"
	"testing"

	"github.com/fumiama/imago"
	"github.com/stretchr/testify/assert"
)

const resultJson = `{"PidP":"53538084_p0","UID":63652,"Width":650,"Height":906,"Title":"秋月照月本表紙","Author":"中乃空","R18":false,"Tags":"[\"艦隊これくしょん\",\"舰队collection\",\"C89\",\"秋月\",\"Akizuki\",\"照月\",\"Teruzuki\",\"艦ぱい\",\"shipgirl breasts\",\"尻神様\",\"尻神样\",\"即夜戦\",\"即将夜战\",\"秋月型\",\"Akizuki-class\",\"ねじ込みたい尻\",\"这屁股让人想肛\"]","Ext":"png","DatePath":"2015/11/14/00/28/41","DHash":"嗉聚裌蠼嬀","Md5":"34f4ed2a6500dc8c6822f0d4333639ba"}`

func TestPictureMarshal(t *testing.T) {
	f, err := os.Open("test.png")
	if err != nil {
		t.Fatal(err)
	}
	img, typ, err := image.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	dh, err := imago.GetDHashStr(img)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Seek(0, io.SeekStart)
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	m := md5.Sum(data)
	ms := hex.EncodeToString(m[:])
	data, err = json.Marshal(&Picture{
		PidP:     "53538084_p0",
		UID:      63652,
		Width:    650,
		Height:   906,
		Title:    "秋月照月本表紙",
		Author:   "中乃空",
		R18:      false,
		Tags:     `["艦隊これくしょん","舰队collection","C89","秋月","Akizuki","照月","Teruzuki","艦ぱい","shipgirl breasts","尻神様","尻神样","即夜戦","即将夜战","秋月型","Akizuki-class","ねじ込みたい尻","这屁股让人想肛"]`,
		Ext:      typ,
		DatePath: "2015/11/14/00/28/41",
		DHash:    dh,
		Md5:      ms,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resultJson, imago.BytesToString(data))
}
