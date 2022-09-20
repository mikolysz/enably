import useSWR from "swr";

// FIXME: Use separate URLs for development and production.
export const api_URL = "/api/v1";

export const fetcher = (url: string) => fetch(url).then((r) => r.json());

export const useApi = <Response>(path: string) =>
  useSWR<Response>(`${api_URL}/${path}`, fetcher);
