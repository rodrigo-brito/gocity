import React from "react";

export default () => {
  return (
    <div className="loading m-t-50 m-b-50">
      <div className="lds-dual-ring"></div>
      <p className="title-white">Fetching repository...</p>
    </div>
  )
}
