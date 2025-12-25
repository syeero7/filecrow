import type { RegisterFileMsg } from "./websocket";

export function registerFile(body: RegisterFileMsg) {
  return fetcher("/register", "POST", body, ["json"]);
}

export function streamFile(fileId: string, body: File) {
  return fetcher(`/stream?id=${fileId}`, "POST", body, ["multipart"]);
}

export function getDownloadURL(fileId: string) {
  return `${getBaseURL()}/download?id=${fileId}`;
}

function getBaseURL() {
  return import.meta.env.DEV ? "/api" : `http://${location.host}`;
}

async function fetcher(
  path: string,
  method: "GET" | "POST",
  body?: Record<string, unknown> | File,
  headers?: ("json" | "multipart")[],
) {
  const url = `${getBaseURL()}${path}`;
  const options: RequestInit = { method };
  const tmp: Record<string, string> = {};

  try {
    if (body && headers) {
      headers.forEach((header) => {
        switch (header) {
          case "json": {
            tmp["Content-Type"] = "application/json";
            options.body = JSON.stringify(body);
            break;
          }

          case "multipart": {
            if (!(body instanceof Blob)) throw new Error("body is not a Blob");
            options.body = body;
            break;
          }
        }
      });

      options.headers = tmp;
    }

    const res = await fetch(url, options);
    if (!res.ok) throw res;
  } catch (err) {
    if (err instanceof Response) {
      console.error("failed to fetch: ", err.statusText);
      return;
    }

    console.error(err instanceof Error ? err.message : err);
  }
}
