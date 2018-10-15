import React from "react";
import Message from "./form/Message";
import Login from "./form/Login";
import Loading from "./Loading";

const DATABASE_NAME = "public-feedback";

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
        loading: false
      };

      this.onSend = this.onSend.bind(this);
      this.onLogin = this.onLogin.bind(this);
      this.onChange = this.onChange.bind(this);

      this.database = window.firebase.firestore();
    }

    componentDidMount(){
      this.setState({loading: true});
      this.database.collection(DATABASE_NAME).get().then((querySnapshot) => {
        console.log(querySnapshot);
        const userMessages = querySnapshot.docs.map((node) => node.data());
        console.log(userMessages);
        this.setState({userMessages, loading: false});
      }).catch((e) => {
        console.error(e);
        this.setState({loading: false});
      });
    }

    onSend(e) {
      e.preventDefault();
      e.stopPropagation();

      const data = {
        name: this.state.user.displayName,
        email: this.state.user.email,
        photo : this.state.user.photoURL,
        message: this.state.message,
        date: new Date()
      };

      this.database.collection(DATABASE_NAME).add(data).then((docRef) => {
        console.log("Document written with ID: ", docRef.id);
        alert("Feedback sent!");
        const userMessages = this.state.userMessages;
        userMessages.unshift(data);
        this.setState({userMessages});
      })
      .catch(function(error) {
        console.error("Error adding document: ", error);
        alert("Fail to send feedback!")
      });
    }

    onLogin(provider, token, user, error) {
      if (error) {
        console.error(error);
      }

      this.setState({user, token});
      this.provider = provider;
    }

    onChange(e) {
      this.setState({message: e.target.value});
    }

    render() {
      return (
        <div className={`modal ${this.props.active ? "is-active" : ""}`}>
          <div className="modal-background" />
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">User Feedback</p>
              <button type="button" className="delete" aria-label="close" onClick={this.props.onClose} />
            </header>
            <section className="modal-card-body">
              {!this.state.user && <Login onLogin={this.onLogin}/>}
              {this.state.user && <form className="m-b-10 is-clearfix">
                  <label htmlFor="input">Hi <b>{this.state.user.displayName}</b>, leave your feedback:</label>
                <div className="field">
                  <div className="control">
                    <textarea className="textarea"
                      value={this.state.message}
                      onChange={this.onChange}
                      placeholder="e.g. Hello world" />
                  </div>
                </div>
                <label className="checkbox">
                  <input type="checkbox" />
                  {" "}Make feedback public
                </label>
                <button className="button is-success is-pulled-right" onClick={this.onSend}>Send Feedback</button>
              </form>}
              <h5 className="title is-5">Recent feedback</h5>
              {this.state.loading && <Loading dark message="Loading feedback..." />}
              {this.state.userMessages.map((message, index) => {
                return <Message key={index}
                  name={message.name}
                  photo={message.photo}
                  content={message.message}
                  date={message.date} />
              })}
            </section>
            <footer className="modal-card-foot">
              <button className="button" onClick={this.props.onClose}>Close</button>
            </footer>
          </div>
        </div>
      )
    }
}