import React from "react";
import legend from "./img/legend.png"

const Legend = () => {
  return (
    <div className="legend is-hidden-mobile">
      <img src={legend} alt=""/>
    </div>
  );
};

export default Legend;