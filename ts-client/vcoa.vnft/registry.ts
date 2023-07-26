import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCreateClass } from "./types/vcoa/vnft/tx";
import { MsgBurnNFT } from "./types/vcoa/vnft/tx";
import { MsgMintNft } from "./types/vcoa/vnft/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/vcoa.vnft.MsgCreateClass", MsgCreateClass],
    ["/vcoa.vnft.MsgBurnNFT", MsgBurnNFT],
    ["/vcoa.vnft.MsgMintNft", MsgMintNft],
    
];

export { msgTypes }