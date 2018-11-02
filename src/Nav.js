import React from "react";
import logo from "./img/logo.png"
export default () => {
  return (
    <div className="level">
      <div className="level-left">
        <div className="level-item">
          <div className="control">
            <a href="/">
              <img className="m-r-10" width={50} src={logo} alt="GoCity: Source code visualization"/>
            </a>
          </div>
          <a href="/">
            <h1 className="title"> GoCity</h1>
          </a>
        </div>
      </div>
      <div className="level-right">
        <div className="level-item is-hidden-mobile">
          <a
            className="github-button"
            href="https://github.com/rodrigo-brito/gocity"
            data-size="large"
            data-show-count="true"
            aria-label="Star gocity on GitHub"
          >
            Star
          </a>
          <span className="m-l-10" />
          <a
            className="github-button"
            href="https://github.com/rodrigo-brito/gocity/fork"
            data-icon="octicon-repo-forked"
            data-show-count="true"
            data-size="large"
            aria-label="Fork gocity on GitHub"
          >
            Fork
          </a>
        </div>
      </div>
    </div>
  );
};
