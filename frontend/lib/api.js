import useSWR from "swr";

// FIXME: Use separate URLs for development and production.
export const api_URL = "/api/v1";

export const fetcher = (...args) => fetch(...args).then((r) => r.json());

export const useApi = (path) => useSWR(`${api_URL}/${path}`, fetcher);
