package GameLogDB

import (
	"database/sql"
	//	"fmt"
	"MyUtility"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/novalagung/golpal"
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

	//对python传过来的数据处理下
	needToJson := strings.TrimSpace(records[6])
	needToJson = strings.Replace(needToJson, "'", "\"", -1)
	needToJson = strings.Replace(needToJson, "u\"", "", -1)
	// log.Println("needToJson=", needToJson)

	data := make(map[string]interface{})
	error := json.Unmarshal([]byte(needToJson), &data)
	if error != nil {
		log.Println("log_save error:", error)
	}
	day_time := strings.Split(strings.TrimSpace(records[1]), ",")[0]
	user_id := strings.Split(strings.TrimSpace(records[4]), ":")[1]
	var channel string = ""
	if len(records) > 8 {
		channel = strings.TrimSpace(records[7])
	} else {
		channel = MyUtility.GetMapRetString(data, "channel", "unknow123")
	}

	if len(user_id) <= 0 {
		user_id = "0"
	}
	date := strings.Replace(day_time, "-", "", -1)

	//开始数据分析啦

	//服装
	cloth := make([]string, 0)
	tempCloth := MyUtility.GetMapWithDefault(data, "cloth", make(map[string]interface{}, 0)).(map[string]interface{})
	mapcloth_len := len(tempCloth)
	if mapcloth_len > 0 {
		// log.Println("cloth_len:", cloth_len)
		for k, v := range tempCloth {
			// log.Println("k:", k)
			// log.Println("v:", v)
			num := MyUtility.GetMapRetInt(v.(map[string]interface{}), "num", 0)
			// log.Println("num:", num)
			cloth_id, err := strconv.Atoi(k)
			if err == nil {
				cloth = append(cloth, strconv.Itoa(cloth_id*1000+num))
			} else {
				log.Println("strconv.Atoi err:", err)
			}
		}

		cloth_len := len(cloth)
		if cloth_len > 0 {
			data["num"] = cloth[cloth_len-1]
			cloth = cloth[:cloth_len-1]
		}
	}

	//道具

	log.Println("data:", data)
	log.Println("day_time:", day_time)
	log.Println("date:", date)
	log.Println("user_id:", user_id)
	log.Println("channel:", channel)
	log.Println("cloth:", cloth)

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
