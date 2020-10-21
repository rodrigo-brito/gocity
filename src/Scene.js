import React from "react";
import * as BABYLON from "babylonjs";
import swal from "sweetalert2";
import PropTypes from 'prop-types';

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
    } else {
      swal(
        'Browser not supported',
        'Your browser don\'t support WebGL, Plaease update your browser.',
        'error'
      );
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
      return null;
    }

    let { width, height } = this.props;

    let opts = {};

    if (width !== undefined && height !== undefined) {
      opts.width = width;
      opts.height = height;
    }

    return <canvas {...opts} ref={this.onCanvasLoaded} />;
  }
}

Scene.propTypes = {
  width: PropTypes.number,
  height: PropTypes.number,
  onSceneMount: PropTypes.oneOf([PropTypes.func, PropTypes.undefined])
}