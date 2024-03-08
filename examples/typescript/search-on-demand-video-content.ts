const API_HOST = "https://my-api-gateway-host.com";

type OnDemandVideoContent = {
    description: string;
    tags: string[];
    thumbnail: string;
    title: string;
}

type ListOnDemandVideosParams = {
    page?: number;
    pageSize?: number;
}

type SearchOnDemandVideoContentResponse = {
    data: OnDemandVideoContent[],
    page: number
}

export const searchOnDemandVideos = (searchTerm: string, params?: ListOnDemandVideosParams): Promise<SearchOnDemandVideoContentResponse> => {
    const searchParams = new URLSearchParams({ search: searchTerm });

    if (params) {
        searchParams.set('page', String(params.page));
        searchParams.set('pageSize', String(params.pageSize));
    }

    return fetch(`${API_HOST}/on-demand/search?${searchParams.toString()}`).then(res => {
        return res.json() as Promise<SearchOnDemandVideoContentResponse>;
    });
}