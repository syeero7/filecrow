import { registerFile, streamFile } from "./api";
import { getElement } from "./utils";

const uploadForm = getElement("[data-upload-form]");
const pendingFiles = new Map<string, File>();

uploadForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const form = e.currentTarget as HTMLFormElement;
  const formData = new FormData(form);
  const files = formData.getAll("files") as File[];

  files.forEach(async (file) => {
    const { name, size } = file;
    if (!name && !size) return;
    const { id } = await registerFile({ type: "register",  name, size })
    pendingFiles.set(id, file);
  });

  form.reset();
});

export function streamPendingFile(fileId: string) {
  const file = pendingFiles.get(fileId);
  if (!file) return console.error(`file with id ${fileId} not found`);

  streamFile(fileId, file).then(() => pendingFiles.delete(fileId));
}
