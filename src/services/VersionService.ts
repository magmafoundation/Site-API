import {Injectable} from "@tsed/di";
import {RedisService} from "./RedisService";
import {Release} from "../entities/Release";

import {Asset} from "../entities/Asset";

import axios from 'axios';

@Injectable()
export class VersionService {


  token: string;

  myAxios: any;

  constructor(private redisService: RedisService) {
    this.token = Buffer.from(`${process.env.GH_USER}:${process.env.GH_TOKEN}`, 'utf8').toString('base64')
    this.myAxios = axios.create({
      headers: {
        'Authorization': `Basic ${this.token}`
      }
    });
  }

  async fetchLatestDev(name: string): Promise<Asset> {

    let response = await this.myAxios.get(`https://api.github.com/repos/magmafoundation/${name}/releases`);

    let result: Release[] = response.data as Release[];

    return result[0].assets.find(value => {
      if (value.name.includes("server")) {
        return value;
      }
    });
  }

  async fetchDev(name: string): Promise<Release[]> {

    let response = await this.myAxios.get(`https://api.github.com/repos/magmafoundation/${name}/releases`);

    let result: Release[] = response.data as Release[];
    return result;
  }

  async fetchLatestStable(name: string): Promise<Asset> {

    let response = await this.myAxios.get(`https://api.github.com/repos/magmafoundation/${name}/releases/latest`);

    let result: Release = response.data as Release;

    return result.assets.find(value => {
      if (value.name.includes("server")) {
        return value;
      }
    });
  }

  async fetchStable(name: string): Promise<Release> {
    let response = await this.myAxios.get(`https://api.github.com/repos/magmafoundation/${name}/releases/latest`);

    let result: Release = response.data as Release;

    return result;
  }


}
