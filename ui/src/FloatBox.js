import React from "react";

export default ({ position, info, visible }) => {
  if (!visible) {
    return null;
  }

  return (
    <div className="float-card" style={{ left: position.x, top: position.y }}>
      <div className="box">
        <h4 className="name">
          {" "}
          {info.name} [{info.type}]
        </h4>
        <div className="content">
          <b>Lines: </b>
          <span className="lines">{info.NOL}</span>
          <br />
          <b>Methods: </b>
          <span className="methods">{info.NOM}</span>
          <br />
          <b>Attributes: </b>
          <span className="attributes">{info.NOA}</span>
          <br />
        </div>
      </div>
    </div>
  );
};
