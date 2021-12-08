package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"safer.place/internal/address/roughprefix"
	"safer.place/internal/language"
	"safer.place/internal/stations"
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

	// For now we just want something, we don't care what
	addrResolver := roughprefix.New()

	meta := Meta{}

	langInfo := make([]language.Info, 0, len(langs))
	codeToInfo := make(map[string]language.Info, len(langs))
	for info := range langs {
		langInfo = append(langInfo, info)
		codeToInfo[info.Code] = info
	}
	meta.Languages = langInfo

	// station locations
	stations := stations.New()

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

	scoreForCoordinates := func(x, y float64) float64 {
		nearest := stations.Nearest(x, y, 3)

		log.Println("nearest:", nearest)

		sum := 0.0
		for _, s := range nearest {
			sum += s.ScoreAverage(5)
		}
		// TODO: Add weights etc, but for now we just cap it at 5
		return math.Min(sum, 5)
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

		address, x, y, err := addrResolver.Resolve(res.InputValue)
		if err != nil {
			log.Printf("unable to resolve: %v", err)
		}
		score := scoreForCoordinates(x, y)
		c.HTML(http.StatusOK, "details.html", DetailsResponse{
			Response:     res,
			Address:      address,
			CoordX:       x,
			CoordY:       y,
			RoundedScore: int(math.Round(score)),
			TrueScore:    score,
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
