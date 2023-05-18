import useSWR from "swr";

// FIXME: Use separate URLs for development and production.
export const apiURL = process.env.NEXT_PUBLIC_API_URL;

export const fetcher = (url: string) => fetch(url).then((r) => r.json());

export const useApi = <Response>(path: string) =>
  useSWR<Response>(`${apiURL}/${path}`, fetcher);

// getAPIResponse accepts the same parameters as fetch, except
// that the first param is a path, not an absolute URL
// it handles authorization automatically.
export const getAPIResponse = async (path: string, init?: RequestInit) => {
  const url = `${apiURL}/${path}`;
  init = init || {};
  const token = localStorage.getItem("token");
  if (token) {
    init.headers = {
      ...init.headers,
      Authorization: `Bearer ${token}`,
    };
  }

  const response = await fetch(url, init);
  if (response.status === 401) {
    throw new Error("Unauthorized");
  }
  return response;
};
