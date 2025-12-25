export function connectWebsocket(onMessage: (msg: string) => void) {
  const ws = new WebSocket(`ws://${location.hostname}:8080/ws`);

  ws.addEventListener("open", () => {
    console.info("websocket connected");
  });

  ws.addEventListener("close", (e) => {
    console.error(`websocket disconnected code: ${e.code} reason: ${e.reason}`);

    if (e.code !== 1001) {
      console.info("reconnecting websocket in 1s");
      setTimeout(() => connectWebsocket(onMessage), 1000);
    }
  });

  ws.addEventListener("message", (e) => {
    if (typeof e.data !== "string") {
      console.error(`unexpected message type: ${typeof e.data}`);
      return;
    }

    onMessage(e.data);
  });
}

export type RegisterFileMsg = {
  type: "register";
  id: string;
  name: string;
  size: number;
};

type TransferProgressMsg = {
  type: "progress";
  id: string;
  current: number;
  total: number;
};

type FileTranferStateMsg = {
  type: "ready" | "done";
  id: string;
};

type WebSocketMsg = RegisterFileMsg | TransferProgressMsg | FileTranferStateMsg;

export function isWebSocketMsg(
  json: Record<string, unknown>,
): json is WebSocketMsg {
  return "id" in json && "type" in json;
}
