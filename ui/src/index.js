import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { BrowserRouter as Router, Route } from "react-router-dom";

const Root = () => (
  <Router basename={process.env.PUBLIC_URL}>
    <Route exact path="/:repository*" component={App} />
  </Router>
);

ReactDOM.render(<Root />, document.getElementById("root"));
