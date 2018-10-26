export function getProportionalColor(start, end, percent) {
  percent = 1 - percent;
  const red = (start.r - end.r) * percent + end.r;
  const green = (start.g - end.g) * percent + end.g;
  const blue = (start.b - end.b) * percent + end.b;

  return { r: red, g: green, b: blue };
}

export function feedbackEvent(action = "click") {
  window.ga('send', 'event', {
    eventCategory: 'feedback',
    eventAction: action
  });
}

export function searchEvent(projectName) {
  window.ga('send', 'event', 'search', projectName);
}

export function openGithubEvent(url) {
  window.ga('send', 'event', 'github', url);
}

export function loginEvent() {
  window.ga('send', 'event', 'login');
}