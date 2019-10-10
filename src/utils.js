export function getProportionalColor(start, end, percent) {
  percent = 1 - percent;
  const red = (start.r - end.r) * percent + end.r;
  const green = (start.g - end.g) * percent + end.g;
  const blue = (start.b - end.b) * percent + end.b;

  return { r: red, g: green, b: blue };
}

export function feedbackEvent(action = 'click') {
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

export function logoBase64() {
  return 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAMAAAAoLQ9TAAACB1BMVEUAAAAACAgALjQAMzTGxsYAsN0AutgAxsno6OgAuu0AwugA2dUAy+IA19oAw+cAz9/R5vXp8fcA4d0AzuoA2uIAzOwA2OUAw+/m8fv0+PwAy+oA3uEA398AwvAAzerN6Pjm8voA0ugA498AuPgA0ukAw/IAwfMAyu4Ay+7x+Pz6/P0AwPQA4+AA1+cA4uEA4OIAx/EA0uvo9Pz1+v0AwvMA2OgA5eEAxfMAxvMAyfIAyvEAy/AAzPAAze8Azu4Azu8Az+4A0O0A0O4A0e0A0uwA0u0A1OsA1eoA1ukA2OgA2ecA2egA3eYA4OQA5OEA5OIA5eEE3OYIz94NzdUPxNYP3OYSysoa0u0cocAenrgg1+op1Ogr1O0w4uQy2OMy5eI54eRB1O1D0edF2elQ4uVU1O5V0/FY5uNZWVlf2eNk5+NlZWVm2utp4edq5eRr3elv0vRv1fNy2PF82PN/2/KE5OeG3vGK4+iN5+eRkZGSkpKmpqan3OCskHit6emwsLC36um4nIG75O676+nD6OzG3d3G7PfJ7PnO2tjR6+3T09PT2tfX7vzX8PrY2NjZ8Pvf8/zh7u3h8vzk8/3k9P3k9fzl7u7l8/3l9P3l9fzm7e/m9P3m9fzn5+fp7u7p7u/q6urr6+vr9/3s7Ozs7u/u7+/v7+/x+f70+v70+/71+v71+/7///+DQEkoAAAAN3RSTlMAAAAAAQMDAwMPDxUdHSIiKio4QkJDQ0RNTU9kZ2hobW11epWlur29vb+/wMvS0tbZ2d/f4vf+WMgihAAAAPtJREFUGNMFwd9KwmAYB+Dft+/Tzf1xo20O18wTlQiCDoJupIvoTrqKjrqWIIpOk8AinFQaWTbb3HTv2/NIQCfAclraFgAgEO9XmWnGeMt/DPU9V7AGdQUAMajRmMwVqAwTvM9Ef4+m6wIKG+fIln6mn0jNus2hwe+G90+e5xt3496BDc3p14jCcvXLnoc6tmU7+dOHeA2M9cCdPOBTdE4FxYU4lmNWz8SPqiTmtLNdei8kiLmUO9ugeuScCb07ZcpTyU2XRoemWyRNO+XFSqJQ53EkjLxdDHs3M9JQpdcaMwDGVbqFAGBEwYX7FVwuPzaAAAAou8VltgOAf4T3Zu4fIzA/AAAAAElFTkSuQmCC';
}
