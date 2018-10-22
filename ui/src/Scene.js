import React from "react";
import * as BABYLON from "babylonjs";

export default class Scene extends React.Component {
  scene = null;
  engine = null;
  canvas = null;

  constructor(props) {
    super(props);
    this.state = {
      hideSupportMessage: false,
      support: BABYLON.Engine.isSupported()
    };

    this.hideSupportMessage = this.hideSupportMessage.bind(this);
  }

  onResizeWindow = () => {
    if (this.engine) {
      this.engine.resize();
    }
  };

  componentDidMount() {
    if (this.state.support) {
      this.engine = new BABYLON.Engine(this.canvas, true, {
        preserveDrawingBuffer: true,
        stencil: true
      });

      this.scene = new BABYLON.Scene(this.engine);

      if (typeof this.props.onSceneMount === "function") {
        this.props.onSceneMount({
          scene: this.scene,
          engine: this.engine,
          canvas: this.canvas
        });
      } else {
        console.error("onSceneMount function not available");
      }

      // Resize the babylon engine when the window is resized
      window.addEventListener("resize", this.onResizeWindow);
    }
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.onResizeWindow);
  }

  onCanvasLoaded = c => {
    if (c !== null) {
      this.canvas = c;
    }
  };

  hideSupportMessage() {
    this.setState({ hideSupportMessage: true });
  }

  render() {
    if (!this.state.support) {
      return (
        <div className={`modal ${this.state.hideSupportMessage ? "" : "is-active"}`}>
          <div className="modal-background" />
          <div className="modal-card">
            <header className="modal-card-head">
              <p className="modal-card-title">Browser not supported</p>
              <button className="delete" aria-label="close" onClick={this.hideSupportMessage} />
            </header>
            <section className="modal-card-body">
              <h1 className="title">Your browser don't support WebGL</h1>
              <p>Plaease update your browser. See more information <a href="https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API">here.</a></p>
            </section>
            <footer className="modal-card-foot">
              <button className="button" onClick={this.hideSupportMessage}>Close</button>
            </footer>
          </div>
        </div>
      );
    }
    // 'rest' can contain additional properties that you can flow through to canvas:
    // (id, className, etc.)
    let { width, height } = this.props;

    let opts = {};

    if (width !== undefined && height !== undefined) {
      opts.width = width;
      opts.height = height;
    }

    return <canvas {...opts} ref={this.onCanvasLoaded} />;
  }
}
