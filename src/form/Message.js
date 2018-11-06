import React from "react";

export default ({ name, photo, content, date, project, onOpen, onClose }) => {
  return (
    <div className="card">
      <div className="card-content">
        <div className="media">
          <div className="media-left">
            <figure className="image is-48x48">
              <img src={photo} alt="User" />
            </figure>
          </div>
          <div className="media-content">
            <p className="title is-4">{name}</p>
            <p className="subtitle is-6">
              <time dateTime={date}>{date.toLocaleString()}</time> - seeing{" "}
              <a
                onClick={() => {
                  onOpen(project);
                  onClose();
                }}
              >
                {project}
              </a>
            </p>
          </div>
        </div>

        <div className="content">{content}</div>
      </div>
    </div>
  );
};
