package GameLogDB

import (
	"MyUtility"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	// "reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	// "github.com/novalagung/golpal"
)

type MySqlDb struct {
	// dbInstance *sql.DB
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
	} else {
		log.Println("first deal data:", data)
	}
	day_time := strings.Split(strings.TrimSpace(records[1]), ",")[0]
	user_id := strings.Split(strings.TrimSpace(records[4]), ":")[1]
	var channel string = ""
	if len(records) > 8 {
		channel = strings.TrimSpace(records[7])
	} else {
		channel = MyUtility.GetMapRetString(data, "channel", "unknow")
	}

	if len(user_id) <= 0 {
		user_id = "0"
	}
	date := strings.Replace(day_time, "-", "", -1)

	//开始数据分析啦

	//属性
	attr_num := map[string]int{"health": 5, "virtue": 6, "vigour": 7}
	attrs := make([]string, 0)
	attrs_len := 0
	tempattr, attrErr := MyUtility.GetMapWithDefault(data, "attr", make(map[string]interface{}, 0)).(map[string]interface{})
	if attrErr == false {
		delete(data, "attr")
	} else {
		maptempattr_len := len(tempattr)
		if maptempattr_len > 0 {
			for k, v := range tempattr {
				// log.Println("k:", k)
				// log.Println("v:", v)
				if len(strings.TrimSpace(k)) > 0 {
					target, attr_finded := attr_num[k]
					if attr_finded {
						attrs = append(attrs, strconv.Itoa(MyUtility.GetSum(target*1000, v)))
					}
				}
			}

			attrs_len = len(attrs)
			if attrs_len > 0 {
				data["attr"] = attrs[attrs_len-1]
				attrs = attrs[:attrs_len-1]
			}
		}
	}

	//服装
	cloth := make([]string, 0)
	cloth_len := 0
	tempCloth, clothErr := MyUtility.GetMapWithDefault(data, "cloth", make(map[string]interface{}, 0)).(map[string]interface{})
	if clothErr == false {
		delete(data, "cloth")
	} else {
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

			cloth_len = len(cloth)
			if cloth_len > 0 {
				data["cloth"] = cloth[cloth_len-1]
				cloth = cloth[:cloth_len-1]
			}
		}
	}

	//道具
	prop := make([]string, 0)
	prop_len := 0
	tempProp, propErr := MyUtility.GetMapWithDefault(data, "prop", make(map[string]interface{}, 0)).(map[string]interface{})
	if propErr == false {
		delete(data, "prop")
	} else {
		maptempProp_len := len(tempProp)
		if maptempProp_len > 0 {
			for k, v := range tempProp {
				// log.Println("k:", k)
				// log.Println("v:", v)
				if len(strings.TrimSpace(k)) > 0 {
					prop = append(prop, strconv.Itoa(MyUtility.GetSum(MyUtility.GetString2Int(k)*1000, v)))
					prop_len += 1
				}
			}

			if prop_len > 0 {
				data["prop"] = prop[prop_len-1]
				prop = prop[:prop_len-1]
			}
		}
	}

	// log.Println("before filter data:", data)
	//删除非log表字段数据
	unfields := make([]string, 0)
	for k, _ := range data {
		if MyUtility.ContainInTarget(k, mysqldb.log_fields) {
			// log.Println(k, "is in log_fields")
		} else {
			unfields = append(unfields, k)
		}
	}
	for _, v := range unfields {
		// log.Println("delete ", v, " in data")
		delete(data, v)
	}

	//sql语句拼接
	pre_values := []string{strings.TrimSpace(user_id), strings.TrimSpace(action), strings.TrimSpace(day_time), strings.TrimSpace(channel)}
	pre_columns := []string{"user_id", "action", "time", "channel"}

	mapKeyList := MyUtility.GetAllMapKey(data)
	colums := strings.Join(append(pre_columns, mapKeyList...), ",")
	values := strings.Join(append(pre_values, MyUtility.GetAllMapValueOrderByKey(data, mapKeyList)...), "','")

	if len(colums) > 0 {
		prepare := fmt.Sprintf("INSERT INTO log_%s (%s)VALUES ('%s')", "dialy", colums, values)

		mysqldb.QueryDB(prepare)
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

	max_Len := MyUtility.GetIntMax(attrs_len, cloth_len, prop_len)

	if max_Len > 1 {
		log.Println("|||||||||||||---------------------||||||||||||")
		rest_prop := make([]string, max_Len-prop_len)
		rest_cloth := make([]string, max_Len-cloth_len)
		rest_attrs := make([]string, max_Len-attrs_len)

		prop = append(prop, rest_prop...)
		cloth = append(cloth, rest_cloth...)
		attrs = append(attrs, rest_attrs...)

		log.Println("prop:", prop)
		log.Println("cloth:", cloth)
		log.Println("attrs:", attrs)

		rest_values := make([]string, max_Len-1)
		for i := 0; i < max_Len-1; i++ {
			finalAttrs := attrs[i]
			finalProp := prop[i]
			finalCloth := cloth[i]
			if len(finalAttrs) <= 0 {
				finalAttrs = "0"
			}
			if len(finalProp) <= 0 {
				finalProp = "0"
			}
			if len(finalCloth) <= 0 {
				finalCloth = "0"
			}
			tmp := fmt.Sprintf("%s','%s','%s", finalAttrs, finalProp, finalCloth)
			log.Println("tmp[a]:", i, "=", tmp)
			rest_values[i] = tmp
		}
		// rest_values = strings.Join(rest_values, ",")
		log.Println("rest_values:", rest_values)
		rest_columns := append(pre_columns, "attr", "prop", "cloth", "other")
		values_sql := "'"

		for k, v := range rest_values {
			tmp := strings.Join(append(pre_values, v), "','") + "','" + MyUtility.GetMapRetString(data, "other", "")
			log.Println("k:", k, ",v:", v)
			log.Println("tmp[b]:", tmp)
			if k == len(rest_values)-1 {
				values_sql += tmp + "'"
			} else {
				values_sql += tmp + "'), ('"
			}
		}
		rest_sql := fmt.Sprintf("INSERT INTO log_%s (%s)VALUES (%s)", "dialy", strings.Join(rest_columns, ","), values_sql)

		mysqldb.QueryDB(rest_sql)
	}

	log.Println("----------------------------------------------------")

}

func (self *MySqlDb) connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gamelog_t2?charset=utf8")
	checkError(err)
	// self.dbInstance = db
	return db
}

func (self *MySqlDb) QueryDB(prepare string, exec ...string) {

	dbInstance := self.connectDB()

	if dbInstance == nil {
		return
	}

	defer func() {
		dbInstance.Close()
		dbInstance = nil
	}()

	log.Println("prepare:", prepare)
	res, err := dbInstance.Exec(prepare)
	checkError(err)
	id, err := res.LastInsertId()
	log.Println("QueryDB LastInsertId:", id)
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Println("err:", err)
		panic(err)
	}
}
