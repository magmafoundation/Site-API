import {Asset} from "./Asset";

export class Release {

   name: string;
   tag_name: string;
   id: number;
   url: string;
   assets: Asset[];
}
