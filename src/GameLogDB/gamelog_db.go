package GameLogDB

import (
	"database/sql"
	//	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/novalagung/golpal"
)

type MySqlDb struct {
	dbInstance *sql.DB
	log_fields []string
}

var mysqldb *MySqlDb

func init() {
	mysqldb = new(MySqlDb)
	mysqldb.log_fields = []string{
		"id", "user_id", "channel", "action", "time", "diamond",
		"coin", "attr", "achievepoint", "hp", "exp", "vip_point",
		"honor", "prop", "cloth", "other", "player"}
}

func GetGameLogDB() *MySqlDb {
	return mysqldb
}

// message = message.strip()
//        if message[:4] == 'INFO':
//            # "INFO 2015-11-30 10:49:15,677  10.224.32.84  /getmailaward uid: 93  mail_award
//            # {'other': {'mailid': 11111301}, u'prop': {u'3020006': 1, u'3020011': 1}}
//            //"INFO   2015-12-30 11:25:54,170       127.0.0.1 /takedailyachievement   uid:  1905  takedailyachievement    {'other': {'daily_acv_id': 5007029}}

func log_save(records []string) {
	action := strings.TrimSpace(records[5])
	if strings.EqualFold(action, "cust_event") {
		return
	}

	data, error := golpal.New().Execute(strings.TrimSpace(records[6]))
	if error != nil {
		log.Println("log_save error:", error)
	}
	day_time := strings.Split(strings.TrimSpace(records[1]), ",")[0]
	user_id := strings.Split(strings.TrimSpace(records[4]), ":")[1]
	if len(records) > 8 {
		// channel := strings.TrimSpace(records[7])
	} else {
		// channel :=
	}

	if len(user_id) <= 0 {
		user_id = "0"
	}
	date := strings.Replace(day_time, "-", "", -1)

	log.Println("data:", data)
	log.Println("day_time:", day_time)
	log.Println("user_id:", user_id)
	log.Println("date:", date)

}

func (self *MySqlDb) connectDB() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gamelog_t2?charset=utf8")
	checkError(err)
	self.dbInstance = db
}

func (self *MySqlDb) QueryDB(prepare string, exec ...string) {

	if self.dbInstance == nil {
		self.connectDB()
	}
	if self.dbInstance == nil {
		return
	}

	defer func() {
		self.dbInstance.Close()
		self.dbInstance = nil
	}()

	if exec == nil {
		_, err := self.dbInstance.Query(prepare)
		checkError(err)
	} else {
		_, err := self.dbInstance.Prepare(prepare)
		checkError(err)
		_, err = self.dbInstance.Exec("123")
		checkError(err)
	}

}

func checkError(err error) {
	if err != nil {
		log.Println("err:", err)
		panic(err)
	}
}
