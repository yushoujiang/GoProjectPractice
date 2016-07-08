package GameLogDB

import (
	"MyUtility"
	"database/sql"
	"encoding/json"
	// "fmt"
	"log"
	// "reflect"
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
		log.Println("err needToJson:", needToJson)
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
			data["cloth"] = cloth[cloth_len-1]
			cloth = cloth[:cloth_len-1]
		}
	}

	//道具
	prop := make([]string, 0)
	prop_len := 0
	tempProp := MyUtility.GetMapWithDefault(data, "prop", make(map[string]interface{}, 0)).(map[string]interface{})
	maptempProp_len := len(tempProp)
	if maptempProp_len > 0 {
		for k, v := range tempProp {
			// log.Println("k:", k)
			// log.Println("v:", v)
			if len(strings.TrimSpace(k)) > 0 {
				prop = append(prop, strconv.Itoa(MyUtility.GetSum(MyUtility.GetString2Int(k)*1000, v)))
			}
		}

		prop_len = len(prop)
		if prop_len > 0 {
			data["prop"] = prop[prop_len-1]
			prop = prop[:prop_len-1]
		}
	}

	log.Println("before filter data:", data)
	//删除非log表字段数据
	unfields := make([]string, 0)
	for k, _ := range data {
		if MyUtility.ContainInTarget(k, mysqldb.log_fields) {
			log.Println(k, "is in log_fields")
		} else {
			unfields = append(unfields, k)
		}
	}
	for _, v := range unfields {
		log.Println("delete ", v, " in data")
		delete(data, v)
	}

	//sql语句拼接
	pre_values := []string{strings.TrimSpace(user_id), strings.TrimSpace(action), strings.TrimSpace(day_time), strings.TrimSpace(channel)}
	pre_columns := []string{"user_id", "action", "time", "channel"}

	colums := strings.Join(append(pre_columns, MyUtility.GetAllMapKey(data)...), ",")
	values := strings.Join(append(pre_values, MyUtility.GetAllMapValue(data)...), ",")

	if len(colums) > 0 {
		prepare := "INSERT INTO log_=? (=?)VALUE (=?)"
		finaExe := make([]string, 0)
		// finaExe = append(finaExe, date)
		finaExe = append(finaExe, "dialy")
		finaExe = append(finaExe, ",")
		finaExe = append(finaExe, colums)
		finaExe = append(finaExe, ",")
		finaExe = append(finaExe, values)
		mysqldb.QueryDB(prepare, finaExe...)
	}

	log.Println("after filter data:", data)
	log.Println("day_time:", day_time)
	log.Println("date:", date)
	log.Println("user_id:", user_id)
	log.Println("channel:", channel)
	log.Println("cloth:", cloth)
	log.Println("prop:", prop)
	log.Println("colums:", colums)
	log.Println("values:", values)
	log.Println("----------------------------------------------------")

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

		log.Println("prepare:", prepare)
		str := strings.Join(exec, ",")
		log.Println("exec:", exec)

		stmt, err := self.dbInstance.Prepare(prepare)
		checkError(err)
		res, err := stmt.Exec(str)
		checkError(err)
		id, err := res.LastInsertId()
		log.Println("QueryDB LastInsertId:", id)
		checkError(err)
	}

}

func checkError(err error) {
	if err != nil {
		log.Println("err:", err)
		panic(err)
	}
}
