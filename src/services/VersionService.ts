import {Injectable} from "@tsed/di";
import {RedisService} from "./RedisService";
import {Release} from "../entities/Release";

import * as rm from "typed-rest-client";
import {Asset} from "../entities/Asset";

@Injectable()
export class VersionService {

  rest: rm.RestClient;


  constructor(private redisService: RedisService) {
    this.rest = new rm.RestClient("");
  }

  async fetchLatestDev(name: string): Promise<Asset> {
    let res: rm.IRestResponse<Release[]> = await this.rest.get<Release[]>(`https://api.github.com/repos/magmafoundation/${name}/releases`);
    let result: Release[] = res.result;

    return result[0].assets.find(value => {
      if (value.name.includes("server")) {
        return value;
      }
    });
  }

  async fetchDev(name: string): Promise<Release[]> {
    let res: rm.IRestResponse<Release[]> = await this.rest.get<Release[]>(`https://api.github.com/repos/magmafoundation/${name}/releases`);
    let result: Release[] = res.result;
    return result;
  }

  async fetchLatestStable(name: string): Promise<Asset> {
    let res: rm.IRestResponse<Release> = await this.rest.get<Release>(`https://api.github.com/repos/magmafoundation/${name}/releases/latest`);
    let result: Release = res.result;

    return result.assets.find(value => {
      if (value.name.includes("server")) {
        return value;
      }
    });
  }

  async fetchStable(name: string): Promise<Release> {
    let res: rm.IRestResponse<Release> = await this.rest.get<Release>(`https://api.github.com/repos/magmafoundation/${name}/releases/latest`);
    let result: Release = res.result;

    return result;
  }


}
