package lolicon

type Item struct {
	Pid int `json:"pid"` // Pid is primary key

	/*Urls struct {
		Original string `json:"original"`
		Regular  string `json:"regular"`
	} `json:"urls"`*/

	P        int      `json:"p"`
	UID      int      `json:"uid"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Title    string   `json:"title"`
	Author   string   `json:"author"`
	R18      bool     `json:"r18"`
	Tags     []string `json:"tags"`
	Ext      string   `json:"ext"`
	Original string   `json:"original"`
	Regular  string   `json:"regular"`
}
