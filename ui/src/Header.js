import React from "react";
export default () => {
  return (
    <div className="columns is-mobile">
      <div className="column">
        <h1 className="title">GoCity</h1>
        <span>Source code visualization</span>
      </div>
      <div className="column is-2">
        <div className="share m-t-10">
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
