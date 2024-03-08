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

type ListOnDemandVideosResponse = {
    data: OnDemandVideo[],
    page: number
}

type ListOnDemandVideosParams = {
    page?: number;
    pageSize?: number;
}

export const listOnDemandVideos = (params?: ListOnDemandVideosParams): Promise<ListOnDemandVideosResponse> => {
    let urlParams = "";

    if (params) {
        urlParams = new URLSearchParams({
            page: String(params.page),
            pageSize: String(params.pageSize)
        }).toString()
    }

    return fetch(`${API_HOST}/on-demand?${urlParams}`).then(res => {
        return res.json() as Promise<ListOnDemandVideosResponse>;
    });
}