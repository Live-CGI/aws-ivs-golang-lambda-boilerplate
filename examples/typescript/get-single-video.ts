const API_HOST = "https://my-api-gateway-host.com";

type OnDemandVideo = {
    active: boolean;
    content: OnDemandVideoContent;
    duration_seconds: number;
    location: string;
    start_time: number;
    uuid: string;
}

type OnDemandVideoContent = {
    description: string;
    tags: string[];
    thumbnail: string;
    title: string;
}

export const getOnDemandVideo = (uuid: OnDemandVideo['uuid']): Promise<OnDemandVideo> => {
    return fetch(`${API_HOST}/on-demand/${uuid}`).then(res => {
        return res.json() as Promise<OnDemandVideo>;
    });
}