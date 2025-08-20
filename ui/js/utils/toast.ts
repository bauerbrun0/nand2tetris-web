const FADE_DURATION = 250;

const CLOSE_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>`;

const INFO_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-info"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>`;
const SUCCESS_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-check-circle"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path><polyline points="22 4 12 14.01 9 11.01"></polyline></svg>`;
const WARNING_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-alert-triangle"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path><line x1="12" y1="9" x2="12" y2="13"></line><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>`;
const ERROR_SVG = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-alert-circle"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>`;

let toastContainer: HTMLElement | null;

export type ToastVariant = "simple" | "info" | "success" | "warning" | "error";

export type ToastOptions = {
  message: string;
  variant: ToastVariant;
  duration: number;
};

export function showToast(options: ToastOptions) {
  if (!toastContainer) {
    toastContainer = document.createElement("div");
    toastContainer.classList.add("toast-container");
    document.body.appendChild(toastContainer);
  }

  const toastElement = document.createElement("div");
  toastElement.classList.add("toast-element");

  // icon
  const iconSvgElement = getIconSvgElement(options.variant);

  const toastElementIconContainer = document.createElement("div");
  toastElementIconContainer.classList.add("icon-container");
  toastElementIconContainer.classList.add(options.variant);

  if (iconSvgElement) {
    toastElementIconContainer.appendChild(iconSvgElement);
  }

  // text
  const toastElementTextContainer = document.createElement("div");
  toastElementTextContainer.classList.add("text-container");
  toastElementTextContainer.innerText = options.message;

  // close button
  const toastElementCloseButton = document.createElement("button");
  toastElementCloseButton.classList.add("close-button");

  const closeSvgElement = stringToSvgElement(CLOSE_SVG);
  toastElementCloseButton.appendChild(closeSvgElement);

  // append all to toast element
  if (options.variant !== "simple") {
    toastElement.appendChild(toastElementIconContainer);
  }
  toastElement.appendChild(toastElementTextContainer);
  toastElement.appendChild(toastElementCloseButton);

  toastContainer.prepend(toastElement);

  setTimeout(() => toastElement.classList.add("open"), 10);
  const hideTimer = setTimeout(
    () => toastElement.classList.remove("open"),
    options.duration,
  );
  const removeTimer = setTimeout(
    () => toastContainer?.removeChild(toastElement),
    options.duration + FADE_DURATION,
  );

  toastElementCloseButton.onclick = (_e) => {
    clearTimeout(hideTimer);
    clearTimeout(removeTimer);
    toastElement.classList.remove("open");
    setTimeout(() => toastContainer?.removeChild(toastElement), FADE_DURATION);
  };
}

function getIconSvgElement(variant: ToastVariant): HTMLElement | null {
  let str: string;
  switch (variant) {
    case "error":
      str = ERROR_SVG;
      break;
    case "info":
      str = INFO_SVG;
      break;
    case "warning":
      str = WARNING_SVG;
      break;
    case "success":
      str = SUCCESS_SVG;
      break;
    default:
      return null;
  }

  return stringToSvgElement(str);
}

function stringToSvgElement(str: string): HTMLElement {
  const parser = new DOMParser();
  const svgDoc = parser.parseFromString(str, "image/svg+xml");
  return svgDoc.documentElement;
}
