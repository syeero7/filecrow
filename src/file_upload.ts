import { registerFile, streamFile } from "./api";
import { getElement } from "./utils";

const uploadForm = getElement("[data-upload-form]");

uploadForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const form = e.currentTarget as HTMLFormElement;
  const formData = new FormData(form);
  const files = formData.getAll("files") as File[];

  files.forEach(async (file) => {
    const { name, size } = file;
    if (!name && !size) return;
    const id = crypto.randomUUID();

    await registerFile({ type: "register", id, name, size });
    streamFile(id, file);
  });

  form.reset();
});
