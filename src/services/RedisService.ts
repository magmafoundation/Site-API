import {Injectable} from "@tsed/di";
import {createClient, RedisClient} from "redis";

const redis = require("redis");
const {promisify} = require("util");

@Injectable()
export class RedisService {

    client: RedisClient;

    private createClient() {

        if (this.client == undefined) {
            let env = process.env as any;
            let port: number = env.REDIS_PORT;
            let host: string = env.REDIS_HOST;
            this.client = createClient(port, host);
        }

    }


    public getAsync(key: any): Promise<any> {
        this.createClient();
        const getAsyncPromisified = promisify(this.client.get).bind(this.client); // now getAsyncPromisified is a promisified version of client.get:
        return getAsyncPromisified(key);
    }

    private isPrimitive(test: any): boolean {
        return (test !== Object(test));
    };

    public set(key: any, val: any, expire?: number): Promise<any> {
        this.createClient();
        const setAsyncPromisified = promisify(this.client.set).bind(this.client); // now setAsyncPromisified is a promisified version of client.get:

        let result: Promise<any> = undefined;
        if(!this.isPrimitive(val) ) {
            result = setAsyncPromisified(key, JSON.stringify(val));
        } else {
            result = setAsyncPromisified(key, val);
        }

        if(expire != null) {
            result.then(() => {
                this.client.expire(key, expire!);
            })
        }



        return result;
    }

}
