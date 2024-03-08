const API_HOST = "https://my-api-gateway-host.com";
const API_KEY = "My api key that I will never put on the client side";

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

type PutOnDemandVideoContentRequest = {
    title: string;
    description: string;
    thumbnail: string;
    tags: string[];
}

export const putOnDemandVideoContent = (uuid: OnDemandVideo['uuid'], req: PutOnDemandVideoContentRequest) => {
    return fetch(`${API_HOST}/on-demand/${uuid}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${API_KEY}`
        },
        body: JSON.stringify(req)
    }).then(res => {
        return res.json() as Promise<OnDemandVideoContent>
    });
}