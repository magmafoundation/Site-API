import {GlobalAcceptMimesMiddleware, ServerLoader, ServerSettings} from "@tsed/common";
import "@tsed/swagger";
import * as bodyParser from "body-parser";
import * as compress from "compression";
import * as cookieParser from "cookie-parser";
import * as methodOverride from "method-override";
import * as dotENV from "dotenv";
import * as cors from "cors";

const rootDir = __dirname;


dotENV.config();

@ServerSettings({
    rootDir,
    mount: {
        "/": `${rootDir}/controllers/**/*.ts`
    },
    componentsScan: [
        `${rootDir}/services/**/*.ts`
    ],
    acceptMimes: ["application/json"],
    swagger: [
        {
            path: "/api-docs"
        }
    ]
})
export class Server extends ServerLoader {
    /**
     * This method let you configure the express middleware required by your application to works.
     * @returns {Server}
     */
    public $beforeRoutesInit(): void | Promise<any> {
        this
            .use(cors())
            .use(GlobalAcceptMimesMiddleware)
            .use(cookieParser())
            .use(compress({}))
            .use(methodOverride())
            .use(bodyParser.json())
            .use(bodyParser.urlencoded({
                extended: true
            }));
    }

}
