package app

import (
	mgithub "MagmaAPI/github"
	mredis "MagmaAPI/redis"
	"MagmaAPI/utils"
	"context"
	json2 "encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mileusna/crontab"
	"log"
	"os"
	"strings"
	"time"
)
import swagger "github.com/arsmn/fiber-swagger/v2"



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

func FetchAndCacheStats()  {
	log.Println("Fetching stats...")
	stats := utils.GetStats()
	json, _ := json2.Marshal(stats)
	mredis.RDB.Set(context.Background(), "stats", json, 5*time.Minute)
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
		val, err := mredis.RDB.Get(context.Background(), "stats").Result()
		var stats utils.Stats

		if err == redis.Nil {
			stats = utils.GetStats()
			json, _ := json2.Marshal(stats)
			mredis.RDB.Set(context.Background(), "stats", json, 5*time.Minute)
			return ctx.JSON(stats)
		}

		_ = json2.Unmarshal([]byte(val), &stats)
		return ctx.JSON(stats)

	})
	app.Static("/docs", "./docs", fiber.Static{
		Compress: false,
		CacheDuration: time.Duration(-1),
	})

	app.Get("/api-docs/*", swagger.New(swagger.Config{ // custom
		URL: "/docs/doc.json",
		DeepLinking: false,
	}))



	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/api-docs/index.html")
	})


	FetchAndCacheStats()

	ctab := crontab.New() // create cron table

	// Fetch stats every 5 min
	ctab.MustAddJob("*/5 * * * *", FetchAndCacheStats)

	panic(app.Listen(":" + os.Getenv("APP_PORT")))
}
