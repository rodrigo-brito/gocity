import React from "react";

export default ({message, dark}) => {
  return (
    <div className={`loading m-t-50 m-b-50 has-text-white ${dark ? "lds-dual-ring-dark" : ""}`}>
      <div className="lds-dual-ring" />
      {message && <p className="has-text-white">{message}</p>}
    </div>
  )
}
