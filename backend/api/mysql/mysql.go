package mysql

import "strconv"
import "strings"
import "fmt"
import "database/sql"
import "time"
import _ "github.com/go-sql-driver/mysql"

type Item struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}
type Item_ struct {
	Name  string `json: "name"`
	Color string `json: "color"`
	Datas []int  `json: "datas"`
}

type HostData struct {
	Host      string    `json:"host"`
	CheckTime time.Time `json: "checktime"`
	Items     []Item    `json:"items"`
	Status    string    `json: "status"`
	Error     string    `json:"error,omitempty"`
}
type Trends struct {
	// 修改 xAxis 可控制趋势图的横轴检查日期。
	XAxis []string `json: "xAxis"`
	// 修改 series[].data 可控制每条 y 轴使用率曲线。
	Series []Item_ `json: "series"`
}
type HostList struct {
	Host     string `json:"Host"`
	User     string `json:"User"`
	Port     int    `json:"Port"`
	Password string `json:"Password"`
	KeyFile  string `json:"KeyFile"`
}

func GetHostList() ([]HostList, error) {
	const host_db = "root:redhat@tcp(mysql:3306)/go_inspector?parseTime=true&loc=Local"
	db, err := sql.Open("mysql", host_db)
	if err != nil {
		fmt.Println("初始化连接池失败...")
		return nil, err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println("连接数据库失败...")
		return nil, err
	}
	rows, err := db.Query(
		"select Host, User, Port, Password, KeyFile from host_config",
	)
	if err != nil {
		fmt.Println("查询最近一次时间点失败...")
		return nil, err
	}
	defer rows.Close()

	var hostlists []HostList
	for rows.Next() {
		var hostlist HostList
		if err = rows.Scan(&hostlist.Host, &hostlist.User, &hostlist.Port, &hostlist.Password, &hostlist.KeyFile); err != nil {
			fmt.Println("数据写入失败..")
			return nil, err
		}
		hostlists = append(hostlists, hostlist)
	}
	return hostlists, nil
}

// 查询可统计数据
func GetTrends(host string) (Trends, error) {
	var trends Trends
	const host_db = "root:redhat@tcp(mysql:3306)/go_inspector?parseTime=true&loc=Local"
	db, err := sql.Open("mysql", host_db)
	if err != nil {
		fmt.Println("初始化连接池失败...")
		return Trends{}, err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println("连接数据库失败...")
		return Trends{}, err
	}
	trendsSelect := "select item_value, check_time_id from host_check_items where host_id=? and item_name=? order by check_time_id asc"
	var itemnames = []struct {
		Name  string
		Color string
	}{
		{"CPU", "#2e78d2"},
		{"Memory", "#198b68"},
		{"Disk", "#d48329"},
	}
	timeCondition := true
	for _, itemname := range itemnames {
		item := Item_{
			Name:  itemname.Name,
			Color: itemname.Color,
		}
		rows, err := db.Query(trendsSelect, host, itemname.Name)
		defer rows.Close()
		if err != nil {
			fmt.Println("查询这台主机最近的数据失败...")
			return Trends{}, err
		}
		for rows.Next() {
			var checktime string
			var itemvalue string
			if err = rows.Scan(&itemvalue, &checktime); err != nil {
				fmt.Println("数据写入失败..")
				return Trends{}, err
			}
			cleanStr := strings.TrimSpace(strings.TrimSuffix(itemvalue, "%"))

			var value int
			if cleanStr != "" {
				var parseErr error
				value, parseErr = strconv.Atoi(cleanStr)
				if parseErr != nil {
					// 如果数据异常（如包含字母或特殊符号），可以打印警告并保留默认值 0，避免打断整个请求
					fmt.Printf("警告: 指标值解析异常 [%s]: %v\n", itemvalue, parseErr)
					value = 0
				}
			} else {
				// 数据为空时默认设为 0（根据业务也可以选择 continue 跳过该点）
				value = 0
			}
			if err != nil {
				fmt.Println("转换转换为数值失败:", err)
				return Trends{}, err
			}
			if timeCondition {
				trends.XAxis = append(trends.XAxis, checktime[5:16])
				fmt.Println(checktime[5:16])
			}
			item.Datas = append(item.Datas, value)
		}
		trends.Series = append(trends.Series, item)
		timeCondition = false
	}

	return trends, nil
}

// 将主机数据写入到数据库中
func SetHostData(host *HostList) error {
	const host_db = "root:redhat@tcp(mysql:3306)/go_inspector?parseTime=true&loc=Local"
	db, err := sql.Open("mysql", host_db)
	if err != nil {
		fmt.Println("初始化连接池失败...")
		return err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println("连接数据库失败...")
		return err
	}
	inserthost := "insert into host_config (Host, User, Port, Password, KeyFile) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(inserthost, host.Host, host.User, host.Port, host.Password, host.KeyFile)
	if err != nil {
		fmt.Printf("插入主机%s检查数据失败!", host.Host)
		return err
	}
	return nil
}

func GetHost() ([]HostData, error) {
	const host_db = "root:redhat@tcp(mysql:3306)/go_inspector?parseTime=true&loc=Local"
	db, err := sql.Open("mysql", host_db)
	if err != nil {
		fmt.Println("初始化连接池失败...")
		return nil, err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println("连接数据库失败...")
		return nil, err
	}
	rows, err := db.Query(
		"select check_time_id from check_lastest_time",
	)
	defer rows.Close()
	if err != nil { //这样写的话rows成了if块中的元素，块外看不到。
		fmt.Println("查询最近一次时间点失败...")
		return nil, err
	}

	var check_time_id time.Time
	for rows.Next() {
		if err = rows.Scan(&check_time_id); err != nil {
			fmt.Println("数据写入失败..")
			return nil, err
		}
	}

	selectsql := "select host_id, status, error from host_check_tasks where check_time_id=?"
	rows, err = db.Query(selectsql, check_time_id)
	if err != nil {
		fmt.Println("查询数据失败tasks...")
		return nil, err
	}

	var hosts []HostData
	var checkCondition string
	var host_id string
	var error_id string
	for rows.Next() {
		var host HostData

		if err = rows.Scan(&host_id, &checkCondition, &error_id); err != nil {
			fmt.Println("数据写入失败..")
			return nil, err
		}
		host.Host = host_id
		host.CheckTime = check_time_id
		host.Status = checkCondition
		host.Error = error_id
		if checkCondition == "UNCON" {
			hosts = append(hosts, host)
			continue
		}
		//处理item的值
		selectsql = "select item_name, item_value, status, detail from host_check_items where check_time_id=? and host_id=?"
		rrows, err := db.Query(selectsql, check_time_id, host_id)
		if err != nil {
			fmt.Println("查询数据失败items...")
			return nil, err
		}
		defer rrows.Close()
		for rrows.Next() {
			var item Item
			if err = rrows.Scan(&item.Name, &item.Value, &item.Status, &item.Detail); err != nil {
				fmt.Println("数据读取失败..")
				return nil, err
			}
			host.Items = append(host.Items, item)
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}
