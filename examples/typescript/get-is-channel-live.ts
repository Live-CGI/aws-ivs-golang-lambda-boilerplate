const API_HOST = "https://my-api-gateway-host.com";

type CreateChannelResponse = {
    arn: string,
    rtmpAddress: string,
    streamKey: string,
    uuid: string
}

type GetIsChannelLiveResponse = {
    state: 'live' | 'offline',
    timestamp: number;
}

export const getIsChannelLive = (channel: CreateChannelResponse['uuid']): Promise<GetIsChannelLiveResponse> => {
    return fetch(`${API_HOST}/channels/${channel}/live`).then(res => {
        return res.json() as Promise<GetIsChannelLiveResponse>;
    });
}