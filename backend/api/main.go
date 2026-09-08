package main

import "log"
import "os/exec"
import "api/mysql"
import "encoding/json"
import "fmt"
import "net/http"
import "os"
import "strings"

// 2. 统一响应结构
type ResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, response ResponseData) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://sreplatform:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func main() {
	getTrends := func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		//获取请求参数
		ip := r.URL.Query().Get("host")

		var trends mysql.Trends // 修改这里
		fmt.Println("host:",ip)
		trends, err := mysql.GetTrends(ip)
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取数据失败:%v", err)
			writeJSON(w, http.StatusInternalServerError, ResponseData{
				Code:    http.StatusInternalServerError,
				Message: "获取趋势数据失败",
				Error:   err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// 4. 设置 HTTP 状态码
		w.WriteHeader(http.StatusOK)

		// 5. 序列化并写入响应体 Body
		if err := json.NewEncoder(w).Encode(trends); err != nil {
			// 写入失败处理（通常在实际项目中记录日志）
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("成功返回数据")
	}

	getHosts := func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		var hosts []mysql.HostData
		hosts, err := mysql.GetHost()
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取数据失败:%v", err)
			writeJSON(w, http.StatusInternalServerError, ResponseData{
				Code:    http.StatusInternalServerError,
				Message: "获取主机检查数据失败",
				Error:   err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// 4. 设置 HTTP 状态码
		w.WriteHeader(http.StatusOK)

		// 5. 序列化并写入响应体 Body
		if err := json.NewEncoder(w).Encode(hosts); err != nil {
			// 写入失败处理（通常在实际项目中记录日志）
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("成功返回数据")
	}
	getHostLists := func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		var hosts []mysql.HostList
		hosts, err := mysql.GetHostList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取数据失败:%v", err)
			writeJSON(w, http.StatusInternalServerError, ResponseData{
				Code:    http.StatusInternalServerError,
				Message: "获取主机列表失败",
				Error:   err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// 4. 设置 HTTP 状态码
		w.WriteHeader(http.StatusOK)

		// 5. 序列化并写入响应体 Body
		if err := json.NewEncoder(w).Encode(hosts); err != nil {
			// 写入失败处理（通常在实际项目中记录日志）
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("成功返回数据")
	}

	/*	health := func(w http.ResponseWriter, r *http.Request){
		var builder strings.Builder
		builder.WriteString("{\n")
		builder.WriteString("    \"status\": \"ok\"\n")
		builder.WriteString("}")

		result := builder.String()
		fmt.Fprintf(w, result)
	}*/
	/*user := func(w http.ResponseWriter, r *http.Request){
		var builder strings.Builder
		builder.WriteString("{\n")
		builder.WriteString("    \"user\": \"xwang\"\n")
		builder.WriteString("}")

		result := builder.String()
		fmt.Fprintf(w, result)
	}*/
	mysqll := func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, ResponseData{
				Code:    http.StatusMethodNotAllowed,
				Message: "请求方法不被允许",
			})
			return
		}

		// 限制读取的 Body 大小（防止恶意大包把内存打爆，这里限制为 1MB）
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var req mysql.HostList
		// 3. 直接从 r.Body 中流式解码 JSON
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, ResponseData{
				Code:    http.StatusBadRequest,
				Message: "请求数据格式错误",
				Error:   err.Error(),
			})
			return
		}
		if strings.TrimSpace(req.Host) == "" || strings.TrimSpace(req.User) == "" {
			writeJSON(w, http.StatusBadRequest, ResponseData{
				Code:    http.StatusBadRequest,
				Message: "主机地址和 SSH 用户不能为空",
			})
			return
		}
		if req.Port < 1 || req.Port > 65535 {
			writeJSON(w, http.StatusBadRequest, ResponseData{
				Code:    http.StatusBadRequest,
				Message: "SSH 端口必须在 1 到 65535 之间",
			})
			return
		}

		if err := mysql.SetHostData(&req); err != nil {
			fmt.Fprintf(os.Stderr, "添加主机 %s 失败: %v\n", req.Host, err)
			writeJSON(w, http.StatusInternalServerError, ResponseData{
				Code:    http.StatusInternalServerError,
				Message: "主机信息添加失败",
				Error:   err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusCreated, ResponseData{
			Code:    http.StatusCreated,
			Message: "主机信息添加成功",
		})

	}

	execGoInspector := func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置跨域头
		setCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 3. 执行系统命令
		binPath := "./go-inspector"
		cmd := exec.Command(binPath)

		// Run 会执行命令并等待其结束；若退出码非 0 则返回 error
		if err := cmd.Run(); err != nil {
			log.Printf("go-inspector 执行失败: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError) // 执行失败返回 500
			return
		}

		// 4. 执行成功，返回 204 No Content（或 200 OK）
		w.WriteHeader(http.StatusNoContent)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hosts", getHosts)
	// 使用普通路径注册，让 POST 和浏览器的 OPTIONS 预检请求都能进入 mysqll。
	mux.HandleFunc("/post/data", mysqll)
	mux.HandleFunc("/items", getTrends)
	mux.HandleFunc("/hostlist", getHostLists)
	mux.HandleFunc("/check", execGoInspector)
	serve := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("监听端口8080......")
	if err := serve.ListenAndServe(); err != nil {
		fmt.Println("监听端口失败...")
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

}
