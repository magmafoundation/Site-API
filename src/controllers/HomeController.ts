import {Controller, Get, PathParams, Redirect, Res} from "@tsed/common";
import {Hidden, Responses, Summary} from "@tsed/swagger";
import {RedisService} from "../services/RedisService";
import {VersionService} from "../services/VersionService";

@Controller("/")
export class HomeController {

    constructor(private redisService: RedisService, private versionService: VersionService) {
        redisService.getAsync("myKey").then(r => console.log(r));
        redisService.set("someKey", {Hello: "Hello World"}, 10).then(value => console.log(value))
    }

    @Get("/")
    @Redirect(301, "https://api.magmafoundation.org/api-docs")
    @Hidden()
    private index() {
    }

    @Summary("List of Dev Versions")
    @Get("/api/resources/:name/:version/dev")
    private async getDev(@PathParams("name") name: string,
                         @PathParams("version") version: string, @Res() res: any): Promise<any> {

        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let release = await this.versionService.fetchDev(releaseName);

        return release;
    }

    @Summary("Latest Dev Version")
    @Get("/api/resources/:name/:version/dev/latest")
    private async getLatestDev(@PathParams("name") name: string,
                               @PathParams("version") version: string, @Res() res: any): Promise<any> {

        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let release = await this.versionService.fetchLatestDev(releaseName);

        return release;
    }

    @Summary("Latest Dev Version")
    @Responses(301, {
        description: "Download url of latest jar"
    })
    @Get("/api/resources/:name/:version/dev/latest/download")
    private async getLatestDevDownload(@PathParams("name") name: string,
                                       @PathParams("version") version: string, @Res() res: any): Promise<any> {

        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let release = await this.versionService.fetchLatestDev(releaseName);

        res.redirect(301, release.browser_download_url);
    }


    @Summary("Download Specified Dev version")
    @Get("/api/resources/:name/:version/dev/:tag/download")
    private async getDevTagDownload(@PathParams("name") name: string,
                                    @PathParams("version") version: string,
                                    @PathParams("tag") tag: string, @Res() res: any): Promise<any> {

        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let releases = await this.versionService.fetchDev(releaseName);

        let download = releases.find(value => value.tag_name == tag).assets.find(value => value.name.includes("server"));
        res.redirect(301, download.browser_download_url);
    }


    @Summary("List of Stable Versions")
    @Get("/api/resources/:name/:version/stable/")
    private async getStable(@PathParams("name") name: string,
                            @PathParams("version") version: string): Promise<any> {
        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let release = await this.versionService.fetchStable(releaseName);

        return release;
    }

    @Summary("Latest Stable Version")
    @Get("/api/resources/:name/:version/stable/latest")
    private async getLatestStable(@PathParams("name") name: string,
                                  @PathParams("version") version: string): Promise<any> {
        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let release = await this.versionService.fetchLatestStable(releaseName);

        return release;
    }

    @Summary("Download Latest Stable Version")
    @Responses(301, {
        description: "Download url of latest jar"
    })
    @Get("/api/resources/:name/:version/stable/latest/download")
    private async getLatestStableDownload(@PathParams("name") name: string,
                                          @PathParams("version") version: string, @Res() res: any): Promise<any> {
        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let release = await this.versionService.fetchLatestStable(releaseName);

        res.redirect(301, release.browser_download_url);
    }
    @Summary("Download Specified Stable version")
    @Get("/api/resources/:name/:version/dev/:tag/download")
    private async getStableTagDownload(@PathParams("name") name: string,
                                    @PathParams("version") version: string,
                                    @PathParams("tag") tag: string, @Res() res: any): Promise<any> {

        let releaseName;

        console.log(version);
        if (version == "1.12.2") {
            releaseName = name;
        } else if (version == "1.14.4") {
            releaseName = name + "-" + version
        }

        let releases = await this.versionService.fetchStable(releaseName);

        let download = releases.find(value => value.tag_name == tag).assets.find(value => value.name.includes("server"));
        res.redirect(301, download.browser_download_url);
    }


}
