package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/snakem982/pandora-box/pkg/proxy"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/utils"
)

func WebTest(r chi.Router) {
	r.Mount("/webtest", webtestRouter())
}

func webtestRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getWebTest)
	r.Post("/delete", deleteWebTest)
	r.Put("/", updateWebTest)
	r.Get("/order", saveWebTestOrder)
	r.Post("/delay", delayWebTest)
	r.Post("/ip", getWebTestIp)

	return r
}

func getWebTest(w http.ResponseWriter, r *http.Request) {
	// Get the webtest from the database
	var res []models.WebTest
	_ = cache.GetList(constant.PrefixWebTest, &res)

	// 返回默认列表
	if len(res) == 0 {
		_ = json.Unmarshal(internal.DefaultWebTest, &res)
		for _, webTest := range res {
			_ = cache.Put(webTest.Id, webTest)
		}
		render.JSON(w, r, res)
		return
	}

	var order []models.WebTest
	_ = cache.Get(constant.WebTestOrder, &order)

	// If the order is empty, return the webtest as is
	if len(order) == 0 {
		render.JSON(w, r, res)
		return
	}

	// Create a map for quick lookup of webtest by ID
	webtestMap := make(map[string]models.WebTest)
	for _, webtest := range res {
		webtestMap[webtest.Id] = webtest
	}

	// Sort res based on the order
	var sortedRes []models.WebTest
	for _, item := range order {
		if webtest, exists := webtestMap[item.Id]; exists {
			sortedRes = append(sortedRes, webtest)
		}
	}

	render.JSON(w, r, sortedRes)
}

func deleteWebTest(w http.ResponseWriter, r *http.Request) {
	webtest := &models.WebTest{}
	if err := render.DecodeJSON(r.Body, webtest); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// Delete the webtest from the database
	_ = cache.Delete(webtest.Id)

	render.NoContent(w, r)
}

func updateWebTest(w http.ResponseWriter, r *http.Request) {
	webtest := &models.WebTest{}
	if err := render.DecodeJSON(r.Body, webtest); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	if webtest.Id == "" {
		webtest.Id = fmt.Sprintf("%s%d", constant.PrefixWebTest, utils.SnowflakeId())
	}

	// Add the webtest to the database
	_ = cache.Put(webtest.Id, webtest)

	render.JSON(w, r, webtest)
}

func delayWebTest(w http.ResponseWriter, r *http.Request) {
	var list []models.WebTest
	if err := render.DecodeJSON(r.Body, &list); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	if len(list) == 0 {
		render.JSON(w, r, list)
		return
	}

	// 进行订阅请求
	pool := utils.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(list))
	for i, web := range list {
		list[i].Delay = -1
		url := web.TestUrl
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorln("Delay测试失败 URL= %s, 错误: %v", url, err)
				}
				done <- struct{}{}
			}()
			// 获取当前时间
			start := time.Now()
			code, err := utils.SendHead(url, proxy.GetProxyUrl())
			// 获取以毫秒为单位的执行时间
			elapsed := time.Since(start).Milliseconds()
			if err != nil {
				return
			}
			if code != 404 && code != 500 && code != 0 {
				list[i].Delay = int(elapsed)
			}
		}, 9*time.Second)
	}
	pool.StartAndWait()

	render.JSON(w, r, list)
	for _, test := range list {
		_ = cache.Put(test.Id, test)
	}
}

func getWebTestIp(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Url string `json:"url"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Accept":     "application/json",
	}
	body, _, err := utils.SendGet(req.Url, headers, proxy.GetProxyUrl())
	if err != nil {
		ErrorResponse(w, r, err)
		return
	}
	if body == "" {
		ErrorResponse(w, r, fmt.Errorf("body is empty"))
		return
	}

	render.PlainText(w, r, body)
}
