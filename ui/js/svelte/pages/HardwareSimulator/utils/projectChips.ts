export type Dimensions = {
  width: number;
  height: number;
};

export type Position = {
  x: number;
  y: number;
};

export function calculateFileContextMenuPosition(
  containerElement: HTMLDivElement,
  event: MouseEvent,
  fileContextMenuDimensions: Dimensions,
): Position {
  const containerDimensions: Dimensions = {
    width: containerElement.clientWidth,
    height: containerElement.clientHeight,
  };

  const absoluteClickPosition: Position = {
    x: event.clientX,
    y: event.clientY,
  };

  // calculate position relative to container
  const relativeClickPosition: Position = {
    x: absoluteClickPosition.x - containerElement.getBoundingClientRect().left,
    y: absoluteClickPosition.y - containerElement.getBoundingClientRect().top,
  };

  // calculate context menu position
  let fileContextMenuX: number = relativeClickPosition.x;
  let fileContextMenuY: number = relativeClickPosition.y;

  if (
    containerDimensions.height - fileContextMenuY <
    fileContextMenuDimensions.height
  ) {
    fileContextMenuY = fileContextMenuY - fileContextMenuDimensions.height;
  }
  if (
    containerDimensions.width - fileContextMenuX <
    fileContextMenuDimensions.width
  ) {
    fileContextMenuX = fileContextMenuX - fileContextMenuDimensions.width;
  }

  return { x: fileContextMenuX, y: fileContextMenuY };
}
