import React from "react";
import PropTypes from 'prop-types';

const Message = ({ name, photo, content, date, project, onOpen, onClose }) => {
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
              <button className="link-like-button"
                onClick={() => {
                  onOpen(project);
                  onClose();
                }}
              >
                {project}
              </button>
            </p>
          </div>
        </div>

        <div className="content">{content}</div>
      </div>
    </div>
  );
};

Message.propTypes = {
  name: PropTypes.string,
  photo: PropTypes.string,
  content: PropTypes.string,
  date: PropTypes.instanceOf(Date),
  project: PropTypes.string,
  onOpen: PropTypes.func,
  onClose: PropTypes.func,
}

export default Message;
