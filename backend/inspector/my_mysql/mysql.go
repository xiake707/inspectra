package my_mysql

import "time"
import "inspector/checks"
import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type Host_Data struct {
	Host     string
	User     string
	Port     int
	Password string
	KeyFile  string
}

func SetLastTime(checkTime time.Time) error {
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
	if _, err = db.Exec("truncate table check_lastest_time"); err != nil {
		fmt.Println("清除旧数据失败")
		return err
	}
	inserthost := "insert into check_lastest_time (check_time_id) values (?)"
	_, err = db.Exec(inserthost, checkTime)
	if err != nil {
		fmt.Println("插入时间失败，检查数据失败!")
		return err
	}
	return nil
}

func SetHostData(host *struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	KeyFile  string `json:"key_file"`
}) error {
	const host_db = "root:redhat@tcp(mysql:3306)/go_inspector"
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
	inserthost := "insert into host_config (Host, User_Port, Password, KeyFile) values (?, ?, ?, ?, ?)"
	_, err = db.Exec(inserthost, host.Host, host.User, host.Port, host.Password, host.KeyFile)
	if err != nil {
		fmt.Printf("插入主机%s检查数据失败!", host.Host)
		return err
	}
	return nil
}

func GetHost() ([]Host_Data, error) {
	const host_db = "root:redhat@tcp(mysql:3306)/go_inspector"
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
		"select * from host_config",
	)
	if err != nil { //这样写的话rows成了if块中的元素，块外看不到。
		fmt.Println("查询数据失败...")
		return nil, err
	}
	defer rows.Close()
	var hosts []Host_Data
	for rows.Next() {
		var host, user, password, keyfile string
		var port int

		if err = rows.Scan(&host, &user, &port, &password, &keyfile); err != nil {
			fmt.Println("数据写入失败..")
			return nil, err
		}

		hosts = append(hosts, Host_Data{
			Host:     host,
			User:     user,
			Port:     port,
			Password: password,
			KeyFile:  keyfile,
		})
	}
	return hosts, nil
}

func SetData(checkResults *[]checks.HostResult) error {
	dsn := "root:redhat@tcp(mysql:3306)/go_inspector?parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("初始化连接池失败")
		return err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println("与数据库建立连接失败")
		return err
	}
	inserthost := "insert into host_check_tasks (host_id, check_time_id, status, error) values (?, ?, ?, ?)"
	insertitem := "insert into host_check_items (host_id, check_time_id, item_name, item_value, status, detail) values (?, ?, ?, ?, ?, ?)"
	for _, checkResult := range *checkResults {
		_, err := db.Exec(inserthost, checkResult.Host, checkResult.CheckTime, checkResult.Status, checkResult.Error)
		if err != nil {
			fmt.Printf("插入主机%s检查数据失败!", checkResult.Host)
			return err
		}
		for _, item := range checkResult.Items {
			_, err = db.Exec(insertitem, checkResult.Host, checkResult.CheckTime, item.Name, item.Value, item.Status, item.Detail)
			if err != nil {
				fmt.Printf("插入主机%s检查指标失败!", item.Name)
				return err
			}
		}
	}

	return nil
}
