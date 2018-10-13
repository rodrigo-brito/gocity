import React from "react";
export default () => {
  return (
    <div className="legend">
      <h3 className="title is-5">Legend</h3>
      <p>
        <i className="legend-ico blue" /> Struct
      </p>
      <p>
        <i className="legend-ico white" /> File
      </p>
      <p>
        <i className="legend-ico red" /> Folder
      </p>
    </div>
  );
};
