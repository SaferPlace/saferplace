package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"safer.place/internal/language"
	"safer.place/internal/web"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
}

type SearchRequest struct {
	Lang  string `form:"lang"`
	Input string `form:"q"`
}

type Meta struct {
	Languages  []language.Info
	NormalFont string
	FancyFont  string
}

type Response struct {
	Meta       Meta
	Lang       string
	InputValue string
	language.Language
}

type DetailsResponse struct {
	Response

	Address                string
	CoordX                 float64
	CoordY                 float64
	RoundedScore           int
	TrueScore              float64
	DistanceToUniversities map[string]int
}

var templateFuncs = template.FuncMap{
	"html": func(value string) template.HTML {
		return template.HTML(value)
	},
	"times": func(n int) []int {
		return make([]int, n)
	},
	"subtract": func(a, b int) int {
		return a - b
	},
}

func run() error {
	r := gin.Default()

	langs, err := language.Languages()
	if err != nil {
		return fmt.Errorf("unable to load languages: %w", err)
	}

	meta := Meta{}

	langInfo := make([]language.Info, 0, len(langs))
	codeToInfo := make(map[string]language.Info, len(langs))
	for info := range langs {
		langInfo = append(langInfo, info)
		codeToInfo[info.Code] = info
	}
	meta.Languages = langInfo

	prepResp := func(c *gin.Context) Response {
		var req SearchRequest
		if err := c.ShouldBind(&req); err != nil {
			log.Printf("cannot bind: %v", err)
			// What do we do here?
		}

		// default lang to english
		if _, ok := codeToInfo[req.Lang]; !ok {
			req.Lang = "en"
		}

		return Response{
			Meta:       meta,
			Lang:       req.Lang,
			InputValue: req.Input,
			Language:   langs[codeToInfo[req.Lang]],
		}
	}

	r.SetFuncMap(templateFuncs)

	r.SetHTMLTemplate(template.Must(
		template.New("").
			Funcs(templateFuncs).
			ParseFS(web.Templates, "**.html"),
	))

	r.GET("/about", func(c *gin.Context) {
		res := prepResp(c)
		c.HTML(http.StatusOK, "about.html", res)
	})

	r.GET("/details", func(c *gin.Context) {
		// TODO: Add query parameters
		c.Redirect(http.StatusPermanentRedirect, "/search")
	})
	r.GET("/search", func(c *gin.Context) {
		res := prepResp(c)
		c.HTML(http.StatusOK, "details.html", DetailsResponse{
			Response:     res,
			Address:      "TEST ADDRESS",
			CoordX:       53.42737,
			CoordY:       -6.24611,
			RoundedScore: 4,
			TrueScore:    4.20,
		})
	})

	r.GET("/", func(c *gin.Context) {
		res := prepResp(c)
		c.HTML(http.StatusOK, "index.html", res)
	})

	r.Run()
	return nil
}

type Config struct {
	Font            string
	FancyFont       string
	DefaultLanguage string
}
