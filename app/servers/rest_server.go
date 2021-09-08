package servers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/masagatech/nav-vts/app/models"
	"github.com/masagatech/nav-vts/app/router"
	"github.com/valyala/fasttemplate"
)

type RESTServer struct {
	App *models.App
}

func (r *RESTServer) Start(config *models.Config) {
	r.App.Fiber = fiber.New()

	r.App.Fiber.Get("/", func(c *fiber.Ctx) error {

		tmpl := "<h4>Welcome to {{name}} !</h4>"
		rndr := fasttemplate.New(tmpl, "{{", "}}")
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		// Compile template and write to fiber context
		_, err := rndr.Execute(c, map[string]interface{}{
			"name": "NAV VTS",
		})
		// Pass error to global error handler, up to you
		return err
	})
	r.Inititalize()
	fmt.Println("Rest server starting @ Port : " + strconv.Itoa(config.Servers.Rest_server.Port))
	err := r.App.Fiber.Listen(":" + strconv.Itoa(config.Servers.Rest_server.Port))

	if err != nil {
		fmt.Println(err)
		panic("Faild to start rest server")
	}

}

func (rest *RESTServer) Inititalize() {
	route := router.Router{
		App: rest.App,
	}
	//route.InitializeMiddleware()
	route.InitializeRouters()
}
