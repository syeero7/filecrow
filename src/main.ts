import "./style.css";
import "./file_upload.ts";
import "./ui.ts";
import { connectWebsocket, isWebSocketMsg } from "./websocket";
import {
  calculateProgress,
  toReadableSize,
  type FileProgress,
} from "./utils.ts";
import { refreshScreen } from "./ui.ts";

type Transfer = {
  name: string;
  size: string;
  status: "pending" | "ready" | "done" | string;
  progress: FileProgress;
};

const filesMap = new Map<string, Transfer>();
export function getFilesMap(): Readonly<typeof filesMap> {
  return filesMap;
}

connectWebsocket((str) => {
  const msg = JSON.parse(str);
  if (!isWebSocketMsg(msg)) return console.log("unknown message", msg);

  switch (msg.type) {
    case "register": {
      filesMap.set(msg.id, {
        name: msg.name,
        status: "pending",
        size: toReadableSize(msg.size),
        progress: { percentage: 0 },
      });
      break;
    }

    case "progress": {
      const { current, total, id } = msg;
      const old = filesMap.get(id)!;
      const { speed, newProgress } = calculateProgress(
        old.progress,
        current,
        total,
      );
      filesMap.set(id, { ...old, progress: newProgress, status: speed });
      break;
    }

    case "ready": {
      const old = filesMap.get(msg.id)!;
      filesMap.set(msg.id, { ...old, status: "ready" });
      break;
    }

    case "done": {
      const old = filesMap.get(msg.id)!;
      filesMap.set(msg.id, { ...old, status: "done" });

      setTimeout(() => {
        filesMap.delete(msg.id);
        refreshScreen();
      }, 2000);
      break;
    }
  }

  refreshScreen();
});
