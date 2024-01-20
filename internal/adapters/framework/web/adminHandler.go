package web

import (
	"encoding/json"
	"github.com/Renewdxin/selfMade/internal/ports/app/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandlerAdapter struct {
	UserApp user.AdminApplicationPorts
}

func NewAdminHandlerAdapter() AdminHandlerAdapter {
	return AdminHandlerAdapter{}
}

func (adapter AdminHandlerAdapter) HomePage(c *gin.Context, title string) {
	// 传递给模板的动态数据
	data := gin.H{
		"title": title,
	}

	// 渲染HTML页面，使用已有的HTML文件 "your_existing_template.html"
	c.HTML(http.StatusOK, "your_existing_template.html", data)
}

func (adapter AdminHandlerAdapter) ShowJobApply(c *gin.Context) {
	result := adapter.UserApp.ShowJobsApply()

	// 使用循环遍历结构体数组
	for _, job := range result {
		// 将每个结构体转换为 JSON
		jsonData, err := json.Marshal(job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// 输出 JSON 数据
		c.String(http.StatusOK, string(jsonData))

	}
}

func (adapter AdminHandlerAdapter) ShowAllJobs(c *gin.Context) {

}

func (adapter AdminHandlerAdapter) ShowJobDetails(c *gin.Context) {

}

func (adapter AdminHandlerAdapter) ApproveJobs(c *gin.Context) {

}