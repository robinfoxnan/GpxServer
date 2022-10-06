package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var chGpxData chan *GpxData
var chanList chan *GpxDataArray

func startHttpServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/v1/regist", registUserHandlers)     // 注册
	http.HandleFunc("/v1/addfriends", registUserHandlers) //

	http.HandleFunc("/v1/gpx", gpxHandlers)         // 单个点上报
	http.HandleFunc("/v1/gpxlist", gpxListHandlers) // 一组点上报

	http.HandleFunc("/v1/track", getTrackHandlers)     // 获取轨迹
	http.HandleFunc("/v1/position", getLastPtHandlers) // 获取最后的位置

	http.HandleFunc("/v1/friends", getLastPtHandlers) // 获取好友列表

	port := config.Server.Port
	logger.Info("server start:", zap.Int("port", port))
	// os.Getenv("PORT")
	//log.Printf("Defaulting to port %s", port)

	//log.Printf("Listening on port %s", port)
	//log.Printf("Open http://localhost:%s in the browser", port)
	fmt.Printf("Open http://localhost:%d in the browser", port)
	addr := fmt.Sprintf("%s:%d", config.Server.Host, port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

}

// 默认的解析函数
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := fmt.Fprint(w, "Hello, Welcome to Bird2Fish gpx World!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// 解析gpx数据并写入管道
// 优先解析body的json数据，其次是解析URL中的数据
func gpxHandlers(w http.ResponseWriter, r *http.Request) {

	//fmt.Println(r.Method)
	//fmt.Println("req from: " + r.RemoteAddr)
	//fmt.Println("headers:")
	//if len(r.Header) > 0 {
	//	for k, v := range r.Header {
	//		fmt.Printf("%s=%s\n", k, v[0])
	//	}
	//}
	tm1 := time.Now().UnixMicro()

	var gpx *GpxData
	var err error
	// post method, should parse body
	if strings.EqualFold("Post", r.Method) {
		// 获取请求报文的内容长度
		len := r.ContentLength
		if len > 0 {
			body := make([]byte, len)
			r.Body.Read(body)
			//fmt.Println("body=" + string(body))

			gpx, err = gpx.FromJsonString(string(body))
			if err != nil {
				w.WriteHeader(500)
				fmt.Fprintln(w, "parse json meet error: %s", err)
				return
			}
		}
	}

	if gpx == nil {
		gpx = &GpxData{}
	}

	str := r.FormValue("phone")
	if len(str) > 0 {
		gpx.Phone = str
	}

	str = r.FormValue("lat")
	if len(str) > 0 {
		gpx.Lat, _ = strconv.ParseFloat(str, 64)
	}
	str = r.FormValue("lon")
	if len(str) > 0 {
		gpx.Lon, _ = strconv.ParseFloat(str, 64)
	}
	str = r.FormValue("ele")
	if len(str) > 0 {
		gpx.Ele, _ = strconv.ParseFloat(str, 64)
	}

	str = r.FormValue("speed")
	if len(str) > 0 {
		gpx.Speed, _ = strconv.ParseFloat(str, 64)
	}

	str = r.FormValue("tm")
	if len(str) > 0 {
		gpx.Tm, _ = strconv.ParseInt(str, 10, 64)
	}

	// 解析后无错误的情况下，写入管道，等待持久化
	if chGpxData != nil {
		chGpxData <- gpx
	} else {
		fmt.Println("error: channel is nil\n")
	}

	//fmt.Println(gpx.ToJsonString())

	w.WriteHeader(200)
	fmt.Fprintln(w, `{"state": "OK" }`)

	tm2 := time.Now().UnixMicro()
	delta := tm2 - tm1
	fmt.Printf("cost: %d ms\n", delta/1000.0)
}

// 当长久时间未上线后，会发送一组数据，可以使用此批量接收接口
func gpxListHandlers(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold("Post", r.Method) {
		w.WriteHeader(500)
		fmt.Fprintln(w, "gpxList should use Post method")
		return
	}

	var err error
	var list *GpxDataArray = NewGpxDataList()
	// 获取请求报文的内容长度
	len := r.ContentLength
	if len > 0 {
		body := make([]byte, len)
		r.Body.Read(body)
		//fmt.Println("body=" + string(body))

		list, err = list.FromJsonString(string(body))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "gpxList has err %s", err)
			return
		}
		// post list to channel
		chanList <- list
	}

	w.WriteHeader(200)
	fmt.Fprintln(w, `{"state": "OK" }`)
}

// 注册函数
func registUserHandlers(w http.ResponseWriter, r *http.Request) {

}

// 查询某人最后位置
func getLastPtHandlers(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold("Post", r.Method) {
		w.WriteHeader(500)
		fmt.Fprintln(w, "positon should use Post method")
		return
	}

	var err error
	param := QueryParam{}
	// 获取请求报文的内容长度
	len := r.ContentLength
	if len > 0 {
		body := make([]byte, len)
		r.Body.Read(body)
		//fmt.Println("body=" + string(body))

		err = param.FromJsonString(string(body))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "param has err: %s", err)
			return
		}
	}

	// 直接查询
	gpx, err := redisCli.FindLastGpx(param.Friend)
	if gpx == nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "param has err: %s", err)
		return
	}

	gpxJson, err := gpx.ToJsonString()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "param has err: %s", err)
		return
	}

	w.WriteHeader(200)
	data := `{"state": "OK", "pt":` + gpxJson + `}`
	fmt.Fprintln(w, data)
}

// 查询某人一段时间轨迹
func getTrackHandlers(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold("Post", r.Method) {
		w.WriteHeader(500)
		fmt.Fprintln(w, "positon should use Post method")
		return
	}

	var err error
	param := QueryParam{}
	// 获取请求报文的内容长度
	len := r.ContentLength
	if len > 0 {
		body := make([]byte, len)
		r.Body.Read(body)
		//fmt.Println("body=" + string(body))

		err = param.FromJsonString(string(body))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintln(w, "param has err: %s", err)
			return
		}
	}

	// 直接查询
	gpxList, err := redisCli.FindGpxTrack(&param)
	if gpxList == nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "find track err: %s", err)
		return
	}

	gpxJson, err := gpxList.ToJsonString()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, ": %s", err)
		return
	}

	w.WriteHeader(200)
	data := `{"state": "OK", "ptList":` + gpxJson + `}`
	fmt.Fprintln(w, data)
}
