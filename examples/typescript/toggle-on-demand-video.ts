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

export const disableOnDemandVideo = (uuid: OnDemandVideo['uuid']) => {
     return fetch(`${API_HOST}/on-demand/${uuid}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${API_KEY}`
        }
     }).then(res => {
        return res.ok;
    });
}

export const enableOnDemandVideo = (uuid: OnDemandVideo['uuid']) => {
     return fetch(`${API_HOST}/on-demand/${uuid}`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${API_KEY}`
        }
     }).then(res => {
        return res.ok;
    });
}