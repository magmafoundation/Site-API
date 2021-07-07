package app

import (
	mgithub "MagmaAPI/github"
	mredis "MagmaAPI/redis"
	"MagmaAPI/utils"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func getRepoFromVersion(version string) string {

	switch version {
	case "1.15.2":
		return "Magma-1.15.x"
	case "1.16":
		return "Magma-1.16.x"
	case "1.16.5":
		return "Magma-1.16.x"
	default:
		return "Magma"
	}

}

func Start() {
	mredis.RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	app := fiber.New()

	app.Use(cors.New())

	util := utils.VersionUtils{}
	mgithub.Setup(os.Getenv("GH_TOKEN"))

	app.Get("/api/resources/:name/:version/dev", func(ctx *fiber.Ctx) error {
		return ctx.JSON(util.GetPreReleases(getRepoFromVersion(ctx.Params("version"))))
	})

	app.Get("/api/resources/:name/:version/dev/latest", func(ctx *fiber.Ctx) error {
		return ctx.JSON(util.GetPreReleases(getRepoFromVersion(ctx.Params("version")))[0])
	})

	app.Get("/api/resources/:name/:version/stable", func(ctx *fiber.Ctx) error {
		return ctx.JSON(util.GetStableReleases(getRepoFromVersion(ctx.Params("version"))))
	})

	app.Get("/api/resources/:name/:version/stable/latest", func(ctx *fiber.Ctx) error {
		return ctx.JSON(util.GetStableReleases(getRepoFromVersion(ctx.Params("version")))[0])
	})

	app.Get("/api/resources/:name/:version/dev/:tag/download", func(ctx *fiber.Ctx) error {

		releases := util.GetPreReleases(getRepoFromVersion(ctx.Params("version")))
		for _, release := range releases {
			if strings.Contains(*release.TagName, ctx.Params("tag")) {
				for _, asset := range release.Assets {
					if strings.Contains(asset.GetBrowserDownloadURL(), "server") {
						return ctx.Redirect(asset.GetBrowserDownloadURL(), 301)
					}
				}
			}

		}
		return ctx.JSON(fiber.Map{
			"error": "No release found for tag",
		})
	})

	app.Get("/api/resources/:name/:version/dev/latest/download", func(ctx *fiber.Ctx) error {

		assets := util.GetPreReleases(getRepoFromVersion(ctx.Params(":version")))[0].Assets

		for _, asset := range assets {
			if strings.Contains(asset.GetBrowserDownloadURL(), "server") {
				return ctx.Redirect(asset.GetBrowserDownloadURL(), fiber.StatusMovedPermanently)
			}
		}

		return ctx.JSON(fiber.Map{
			"error": "No release found",
		})

	})

	app.Get("/api/resources/:name/:version/stable/latest/download", func(ctx *fiber.Ctx) error {

		assets := util.GetPreReleases(getRepoFromVersion(ctx.Params(":version")))[0].Assets

		for _, asset := range assets {
			if strings.Contains(asset.GetBrowserDownloadURL(), "server") {
				return ctx.Redirect(asset.GetBrowserDownloadURL(), fiber.StatusMovedPermanently)
			}
		}

		return ctx.JSON(fiber.Map{
			"error": "No release found",
		})

	})

	app.Get("/api/stats", func(ctx *fiber.Ctx) error {

		return ctx.JSON(utils.GetStats())

	})

	panic(app.Listen(":" + os.Getenv("APP_PORT")))
}
