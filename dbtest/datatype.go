package dbtest

type Node struct {
	IdcId    int64  `db:"idcId"`
	Zone     string `db:"zone"`
	ZoneId   int64  `db:"zoneId"`
	Region   string `db:"region"`
	RegionId int64  `db:"regionId"`
}
