import React from "react";
import Message from "./Message";
import Login from "./Login";
import Loading from "../Loading";
import PropTypes from 'prop-types';

const DATABASE_ALL = "all-messages";
const DATABASE_PUBLIC = "public-messages";
const DEFAULT_LIMIT = 100;

export default class FeedbackForm extends React.Component {
  provider = null;
  database = null;

  constructor(props) {
    super(props);
    this.state = {
      token: null,
      user: null,
      message: "",
      userMessages: [],
      loading: false,
      errorMessage: null,
      isPublic: true,
      sent: false
    };

    this.onSend = this.onSend.bind(this);
    this.onLogin = this.onLogin.bind(this);
    this.onChange = this.onChange.bind(this);
    this.onClose = this.onClose.bind(this);
    this.onPublicChange = this.onPublicChange.bind(this);

    this.database = window.firebase.firestore();
    this.database.settings({ timestampsInSnapshots: true });
  }

  componentDidMount() {
    this.setState({ loading: true });
    this.database
      .collection(DATABASE_PUBLIC)
      .orderBy("date", "desc")
      .limit(DEFAULT_LIMIT)
      .get()
      .then(querySnapshot => {
        const userMessages = querySnapshot.docs.map(node => {
          const message = node.data();
          message.date = message.date.toDate();
          return message;
        });
        this.setState({ userMessages, loading: false });
      })
      .catch(e => {
        console.error(e);
        this.setState({ loading: false });
      });
  }

  onSend(e) {
    e.preventDefault();
    e.stopPropagation();

    const data = {
      name: this.state.user.displayName,
      email: this.state.user.email,
      photo: this.state.user.photoURL,
      message: this.state.message,
      date: new Date(),
      project: this.props.projectURL
    };

    this.database
      .collection(DATABASE_ALL)
      .add(data)
      .then(docRef => {
        if (!this.state.isPublic) {
          this.setState({ sent: true, errorMessage: null });
          return;
        }

        this.database
          .collection(DATABASE_PUBLIC)
          .doc(docRef.id)
          .set({
            name: data.name,
            photo: data.photo,
            message: data.message,
            date: data.date,
            project: data.project
          })
          .then(() => {
            const userMessages = this.state.userMessages;
            userMessages.unshift(data);
            this.setState({ sent: true, errorMessage: null, userMessages });
          });
      })
      .catch(function(error) {
        this.setState({ errorMessage: "Error on send feedback!" });
        console.error("Error adding document: ", error);
      });
  }

  onLogin(provider, token, user, error) {
    if (error) {
      console.error(error);
    }

    this.setState({ user, token });
    this.provider = provider;
  }

  onChange(e) {
    this.setState({ message: e.target.value });
  }

  onPublicChange() {
    this.setState(prev => ({ isPublic: !prev.isPublic }));
  }

  onClose() {
    this.setState({
      message: "",
      errorMessage: null,
      isPublic: true,
      sent: false
    });
    this.props.onClose();
  }

  render() {
    return (
      <div className={`modal ${this.props.active ? "is-active" : ""}`}>
        <div className="modal-background" />
        <div className="modal-card">
          <header className="modal-card-head">
            <p className="modal-card-title">User Feedback</p>
            <button
              type="button"
              className="delete"
              aria-label="close"
              onClick={this.onClose}
            />
          </header>
          <section className="modal-card-body">
            {this.state.sent && (
              <div className="notification is-success">
                Message sent! Thank you for the feedback!
              </div>
            )}
            {this.state.errorMessage && (
              <div className="notification is-error">
                Erro on send message. Try again
              </div>
            )}
            {!this.state.user && (
              <Login
                onLogin={this.onLogin}
                onClose={this.onClose}
                projectURL={this.props.projectURL}
              />
            )}
            {this.state.user &&
              !this.state.sent && (
                <form className="m-b-10 is-clearfix">
                  <label htmlFor="input">
                    Hi <b>{this.state.user.displayName}</b>, leave your feedback
                    about <b>{this.props.projectURL}</b>:
                  </label>
                  <div className="field">
                    <div className="control">
                      <textarea
                        className="textarea"
                        value={this.state.message}
                        onChange={this.onChange}
                        placeholder="e.g. Hello world"
                      />
                    </div>
                  </div>
                  <label className="checkbox">
                    <input
                      type="checkbox"
                      checked={this.state.isPublic}
                      onChange={this.onPublicChange}
                    />{" "}
                    Public feedback
                  </label>
                  <button
                    className="button is-success is-pulled-right"
                    onClick={this.onSend}
                  >
                    Send Feedback
                  </button>
                </form>
              )}
            <h5 className="title is-5">Recent comments</h5>
            {this.state.loading && (
              <Loading dark message="Loading feedback..." />
            )}
            {this.state.userMessages.map((message, index) => {
              return (
                <Message
                  key={index}
                  onOpen={this.props.onOpen}
                  onClose={this.props.onClose}
                  name={message.name}
                  photo={message.photo}
                  content={message.message}
                  project={message.project}
                  date={message.date}
                />
              );
            })}
          </section>
          <footer className="modal-card-foot">
            <button className="button" onClick={this.onClose}>
              Close
            </button>
          </footer>
        </div>
      </div>
    );
  }
}

FeedbackForm.propTypes = {
  active: PropTypes.bool,
  projectURL: PropTypes.string,
  onOpen: PropTypes.func,
  onClose: PropTypes.func
}