package database

type Picture struct {
	PidP     string `db:"pidp"` // PidP is primary key like 123_p0
	UID      int    `db:"uid"`  // UID is the author's uid
	Width    int    `db:"width"`
	Height   int    `db:"height"`
	Title    string `db:"title"`
	Author   string `db:"author"`
	R18      bool   `db:"r18"`
	Tags     string `db:"tags"`     // Tags is a json array of tags
	Ext      string `db:"ext"`      // Ext is like "jpg"
	DatePath string `db:"datepath"` // DatePath is https://i.pximg.net/img-original/img/<<<<<2019/08/24/08/00/02>>>>>/76428535_p0.png
	DHash    string `db:"dhash"`    // DHash is utf-8 base16384 8-bytes image dhash
	Md5      string `db:"md5"`      // Md5 is hex-encoded image digest
}
