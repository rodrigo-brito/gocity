export function getProportionalColor(start, end, percent) {
  percent = 1 - percent;
  const red = (start.r - end.r) * percent + end.r;
  const green = (start.g - end.g) * percent + end.g;
  const blue = (start.b - end.b) * percent + end.b;

  return { r: red, g: green, b: blue };
}
