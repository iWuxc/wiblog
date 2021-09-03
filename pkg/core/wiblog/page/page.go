package page

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"text/template"
	"wiblog/pkg/conf"
	"wiblog/tools"
)

var htmlTemplate *template.Template

func init() {
	htmlTemplate = template.New("wiblog").Funcs(tools.TplFuncMap)
	root := filepath.Join(conf.WorkDir, "website")
	files := tools.ReadDirFiles(root, func(name string) bool {
		if name == ".DS_Store" {
			return true
		}
		return false
	})
	_, err := htmlTemplate.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
}

// RegisterRoutes register routes
func RegisterRoutes(e *gin.Engine) {
	e.NoRoute(handleNotFound)

	//login page
	e.GET("/admin/login", handleLoginPage)


}

// RegisterRoutesAuthz register admin
func RegisterRoutesAuthz(group gin.IRoutes) {
	// console
	group.GET("/profile", handleAdminProfile)
	//write
	group.GET("write-post", handleAdminWritePost)

	//manage
	group.GET("manage-series", handleAdminSeries)
	group.GET("/add-serie", handleAdminSerie)
}
