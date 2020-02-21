import {Controller, Get} from "@tsed/common";
import {Summary} from "@tsed/swagger";
import {RedisService} from "../services/RedisService";

@Controller("/")
export class HomeController {

    constructor(private redisService: RedisService) {
        redisService.getAsync("myKey").then(r => console.log(r));
        redisService.set("someKey", {Hello: "Hello World"}, 10).then(value => console.log(value))
    }

    @Summary("Index")
    @Get("/")
    index(): string {
        return "Hello, World!"
    }

}
