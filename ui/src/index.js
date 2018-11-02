import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { HashRouter, Route } from "react-router-dom";

const Root = () => (
  <HashRouter basename={process.env.PUBLIC_URL}>
    <Route exact path="/:repository*" component={App} />
  </HashRouter>
);

ReactDOM.render(<Root />, document.getElementById("root"));
