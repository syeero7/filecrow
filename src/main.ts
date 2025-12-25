import "./style.css";
import "./file_upload.ts";
import "./ui.ts";
import { connectWebsocket, isWebSocketMsg } from "./websocket";
import {
  calculateProgress,
  toReadableSize,
  type FileProgress,
  type ReadableSize,
  type ReadableSpeed,
} from "./utils.ts";
import { refreshScreen } from "./ui.ts";
import { streamPendingFile } from "./file_upload.ts";

type Transfer = {
  name: string;
  size: ReadableSize;
  progress: FileProgress;
  status: "pending" | "failed" | "done" | ReadableSpeed;
};

const filesMap = new Map<string, Transfer>();
export function getFilesMap(): ReadonlyMap<string, Transfer> {
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
      streamPendingFile(msg.id);
      break;
    }

    case "done":
    case "failed": {
      const old = filesMap.get(msg.id)!;
      filesMap.set(msg.id, { ...old, status: msg.type });

      setTimeout(() => {
        filesMap.delete(msg.id);
        refreshScreen();
      }, 2000);
      break;
    }
  }

  refreshScreen();
});
