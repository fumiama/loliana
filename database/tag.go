package database

type Tag struct {
	PidP string `db:"pidp"` // PidP is primary key like 123_p0
	UID  int    `db:"uid"`  // UID is the author's uid
}
