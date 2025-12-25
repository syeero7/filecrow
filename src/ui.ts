import { getDownloadURL } from "./api";
import { getFilesMap } from "./main.ts";
import { getElement } from "./utils";

const fileList = getElement("[data-files]");
const fileTemplate = getElement<HTMLTemplateElement>("[data-file-temp]");

function renderElements() {
  const fragment = document.createDocumentFragment();
  getFilesMap().forEach((ft, fileId) => {
    const clone = fileTemplate.content.cloneNode(true);
    if (!(clone instanceof DocumentFragment)) return;

    getElement("[data-file-size]", clone).textContent = ft.size;
    getElement("[data-file-status]", clone).textContent = ft.status;

    const nameEl = getElement<HTMLLabelElement>("[data-file-name]", clone);
    nameEl.textContent = ft.name;
    nameEl.htmlFor = fileId;

    const progressEl = getElement<HTMLProgressElement>(
      "[data-file-progress]",
      clone,
    );
    progressEl.id = fileId;
    if (ft.status !== "pending") {
      progressEl.textContent = `${ft.progress.percentage} %`;
      progressEl.value = ft.progress.percentage;
    }

    if (ft.status === "ready") {
      getElement<HTMLAnchorElement>("[data-file-url]", clone).href =
        getDownloadURL(fileId);
    }

    fragment.appendChild(clone);
  });

  fileList.appendChild(fragment);
}

function removeElements() {
  while (fileList.firstChild) {
    fileList.firstChild.remove();
  }
}

export function refreshScreen() {
  requestAnimationFrame(() => {
    removeElements();
    renderElements();
  });
}
