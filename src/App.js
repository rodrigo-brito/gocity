import React, { Component } from 'react';
import FloatBox from './FloatBox';
import * as BABYLON from 'babylonjs';
import BabylonScene from './Scene';
import axios from 'axios';
import Navbar from './Nav';
import Legend from './Legend';
import Loading from './Loading';
import { feedbackEvent, getProportionalColor, searchEvent, logoBase64 } from './utils';
import swal from 'sweetalert2';
import Cookies from 'js-cookie';

const URLRegexp = new RegExp(/^(?:https:\/\/?)?(github\.com\/.*)/i);

const endpoint = Cookies.get('gocity_api') || process.env.REACT_APP_API_URL;

// TODO: isolate in the constants file
const colors = {
  PACKAGE: {
    start: { r: 255, g: 100, b: 100 },
    end: { r: 255, g: 100, b: 100 }
  },
  FILE: {
    start: { r: 255, g: 255, b: 255 },
    end: { r: 0, g: 0, b: 0 }
  },
  STRUCT: {
    start: { r: 32, g: 156, b: 238 },
    end: { r: 0, g: 0, b: 0 }
  }
};

const examples = [
  {
    branch: "master",
    name: "sirupsen/logrus",
    link: "github.com/sirupsen/logrus"
  },
  {
    branch: "master",
    name: "gin-gonic/gin",
    link: "github.com/gin-gonic/gin"
  },
  {
    branch: "master",
    name: "spf13/cobra",
    link: "github.com/spf13/cobra"
  },
  {
    branch: "master",
    name: "golang/dep",
    link: "github.com/golang/dep"
  },
  {
    branch: "master",
    name: "gohugoio/hugo",
    link: "github.com/gohugoio/hugo"
  }
];

class App extends Component {
  canvas = null;
  scene = null;
  engine = null;
  camera = null;
  light = null;

  constructor(props) {
    super(props);
    this.state = {
      feedbackFormActive: false,
      loading: false,
      repository:
        this.props.match.params.repository || "github.com/rodrigo-brito/gocity",
      branch: this.props.match.params.branch || "master"
    };

    this.addBlock = this.addBlock.bind(this);
    this.onInputChange = this.onInputChange.bind(this);
    this.onClick = this.onClick.bind(this);
    this.showTooltip = this.showTooltip.bind(this);
    this.hideTooltip = this.hideTooltip.bind(this);
    this.plot = this.plot.bind(this);
    this.process = this.process.bind(this);
    this.reset = this.reset.bind(this);
    this.initScene = this.initScene.bind(this);
    this.onMouseMove = this.onMouseMove.bind(this);
    this.updateCamera = this.updateCamera.bind(this);
    this.onSceneMount = this.onSceneMount.bind(this);
    this.onFeedBackFormClose = this.onFeedBackFormClose.bind(this);
    this.openFeedBackForm = this.openFeedBackForm.bind(this);
    this.openModal = this.openModal.bind(this);
    this.closeModal = this.closeModal.bind(this);
    this.getBadgeValue = this.getBadgeValue.bind(this);
  }

  componentDidMount() {
    if (this.state.repository) {
      this.process(this.state.repository, "", this.state.branch);
    }
  }

  onMouseMove(e) {
    this.mouse_x = e.pageX;
    this.mouse_y = e.pageY;
  }

  showTooltip(info) {
    setTimeout(() => {
      this.setState({
        infoVisible: true,
        infoData: info,
        infoPosition: { x: this.mouse_x, y: this.mouse_y }
      });
    }, 100);
  }

  hideTooltip() {
    this.setState({
      infoVisible: false
    });
  }

  reset() {
    this.scene.dispose();
    this.scene = new BABYLON.Scene(this.engine);
    this.initScene();
  }

  addBlock = data => {
    const bar = BABYLON.MeshBuilder.CreateBox(
      data.label,
      { width: data.width, depth: data.depth, height: data.height },
      this.scene
    );
    bar.receiveShadows = false;

    if (data.parent) {
      bar.parent = data.parent;

      var bounds = data.parent.getBoundingInfo();
      bar.position.y = bounds.maximum.y + data.height / 2.0;
    }
    bar.position.x = data.x || 0;
    bar.position.z = data.y || 0;

    bar.info = data.info;

    bar.actionManager = new BABYLON.ActionManager(this.scene);
    bar.actionManager.registerAction(
      new BABYLON.ExecuteCodeAction(
        BABYLON.ActionManager.OnPointerOverTrigger,
        () => {
          this.showTooltip(bar.info);
        }
      )
    );

    bar.actionManager.registerAction(
      new BABYLON.ExecuteCodeAction(
        BABYLON.ActionManager.OnPointerOutTrigger,
        this.hideTooltip
      )
    );

    // Material
    bar.material = new BABYLON.StandardMaterial(data.label + "mat", this.scene);
    bar.material.diffuseColor = data.color;

    bar.freezeWorldMatrix();

    return bar;
  };

  plot(children, parent) {
    if (!children) {
      return;
    }

    children.forEach(data => {
      var color = getProportionalColor(
        colors[data.type].start,
        colors[data.type].end,
        Math.min(100, data.numberOfLines / 2000.0)
      );

      var mesh = this.addBlock({
        x: data.position.x,
        y: data.position.y,
        width: data.width,
        depth: data.depth,
        height: data.numberOfMethods,
        color: new BABYLON.Color3(color.r / 255, color.g / 255, color.b / 255),
        parent: parent,
        info: {
          name: data.name,
          url: data.url,
          type: data.type,
          NOM: data.numberOfMethods,
          NOL: data.numberOfLines,
          NOA: data.numberOfAttributes
        }
      });

      if (parent) {
        mesh.parent = parent;
      }

      if (data.children && data.children.length > 0) {
        this.plot(data.children, mesh);
      }
    });
  }

  updateCamera(width, height) {
    if (width > 1000) {
      this.camera.useAutoRotationBehavior = false;
    } else {
      this.camera.useAutoRotationBehavior = true;
    }
    width = Math.min(width, 1000);
    height = Math.min(height, 1000);
    this.camera.setPosition(
      new BABYLON.Vector3(width / 2, width, (width + height) / 2)
    );
  }

  initScene() {
    this.scene.clearColor = new BABYLON.Color3(0.7, 0.7, 0.7);
    // This creates and positions a free camera (non-mesh)
    this.camera = new BABYLON.ArcRotateCamera(
      "camera",
      0,
      0,
      10,
      BABYLON.Vector3.Zero(),
      this.scene
    );

    // This targets the camera to scene origin
    this.camera.setTarget(BABYLON.Vector3.Zero());

    // This attaches the camera to the canvas
    this.camera.attachControl(this.canvas, true);

    this.camera.setPosition(new BABYLON.Vector3(500, 400, -100));
    this.camera.useAutoRotationBehavior = true;

    // This creates a light, aiming 0,1,0 - to the sky (non-mesh)
    var light = new BABYLON.HemisphericLight(
      "global_light",
      new BABYLON.Vector3(0, 1, 0),
      this.scene
    );

    light.intensity = 0.8;
  }

  onSceneMount(e) {
    this.scene = e.scene;
    this.canvas = e.canvas;
    this.engine = e.engine;

    this.initScene();

    this.engine.runRenderLoop(() => {
      if (this.scene) {
        this.scene.render();
      }
    });
  }

  handleKeyPress = event => {
    if (event.key === "Enter") {
      this.onClick();
    }
  };

  onInputChange(e) {
    if(e.target.id === "repository") {this.setState({ repository: e.target.value })}
    if(e.target.id === "branch") {this.setState({ branch: e.target.value })}
  }

  process(repository, json, branch) {
    if (!BABYLON.Engine.isSupported()) {
      return;
    }

    const match = URLRegexp.exec(repository);
    if (!match) {
      swal("Invalid URL", "Please inform a valid Github URL.", "error");
      return;
    }
    if (match !== this.props.match.params.repository || branch !== this.props.match.params.branch) {
      this.props.history.push(`/${match[1]}/#/${branch}`);
    }

    this.setState({
      repository: match[1],
      loading: true
    });

    let request = null;
    if (json) {
      request = axios.get(json);
    } else {
      request = axios.get(endpoint, {
        params: {
          q: match[1],
          b: branch
        }
      });
    }

    request
      .then(response => {
        this.setState({ loading: false });
        this.reset();

        if (response.data.children && response.data.children.length === 0) {
          swal("Invalid project", "Only Go projects are allowed.", "error");
        }

        this.plot(response.data.children);
        this.updateCamera(response.data.width, response.data.depth);
      })
      .catch(e => {
        this.setState({ loading: false });
        swal(
          "Error during plot",
          "Something went wrong during the plot. Try again later",
          "error"
        );
        console.error(e);
      });

    // this.scene.freezeActiveMeshes();
    this.scene.autoClear = false; // Color buffer
    this.scene.autoClearDepthAndStencil = false; // Depth and stencil, obviously
    this.scene.blockfreeActiveMeshesAndRenderingGroups = true;
    this.scene.blockfreeActiveMeshesAndRenderingGroups = false;
  }

  onClick() {
    searchEvent(this.state.repository);
    this.process(this.state.repository, "", this.state.branch);
  }

  onFeedBackFormClose() {
    this.setState({ feedbackFormActive: false });
  }

  openFeedBackForm() {
    this.setState({ feedbackFormActive: true });
    feedbackEvent();
  }

  openModal() {
    this.setState({ modalActive: true });
  }

  closeModal() {
    this.setState({ modalActive: false });
  }

  getBadgeValue(template) {
    const repo = this.state.repository;
    const baseUrl = `https://img.shields.io/static/v1?label=gocity&color=blue&style=for-the-badge&message=${repo}&logo=${logoBase64()}`;
    const templates = {
      md: `![](${baseUrl})`,
      html: `<img src="${baseUrl}" alt="checkout my repo on gocity"/>`
    };
    return templates[template];
  }

  render() {
    return (
      <main onMouseMove={this.onMouseMove}>
        <a
          href="https://github.com/rodrigo-brito/gocity"
          className="github-corner is-hidden-tablet"
          aria-label="View source on GitHub"
        >
          <svg
            width="80"
            height="80"
            viewBox="0 0 250 250"
            style={{ fill: "#151513", color: "#fff" }}
            aria-hidden="true"
          >
            <path d="M0,0 L115,115 L130,115 L142,142 L250,250 L250,0 Z" />
            <path
              d="M128.3,109.0 C113.8,99.7 119.0,89.6 119.0,89.6 C122.0,82.7 120.5,78.6 120.5,78.6 C119.2,72.0 123.4,76.3 123.4,76.3 C127.3,80.9 125.5,87.3 125.5,87.3 C122.9,97.6 130.6,101.9 134.4,103.2"
              fill="currentColor"
              style={{ transformOrigin: "130px 106px" }}
              className="octo-arm"
            />
            <path
              d="M115.0,115.0 C114.9,115.1 118.7,116.5 119.8,115.4 L133.7,101.6 C136.9,99.2 139.9,98.4 142.2,98.6 C133.8,88.0 127.5,74.4 143.8,58.0 C148.5,53.4 154.0,51.2 159.7,51.0 C160.3,49.4 163.2,43.6 171.4,40.1 C171.4,40.1 176.1,42.5 178.8,56.2 C183.1,58.6 187.2,61.8 190.9,65.4 C194.5,69.0 197.7,73.2 200.1,77.6 C213.8,80.2 216.3,84.9 216.3,84.9 C212.7,93.1 206.9,96.0 205.4,96.6 C205.1,102.4 203.0,107.8 198.3,112.5 C181.9,128.9 168.3,122.5 157.7,114.1 C157.9,116.9 156.7,120.9 152.7,124.9 L141.0,136.5 C139.8,137.7 141.6,141.9 141.8,141.8 Z"
              fill="currentColor"
              className="octo-body"
            />
          </svg>
        </a>
        <FloatBox position={this.state.infoPosition} info={this.state.infoData} visible={this.state.infoVisible} />
        <header className="header">
          <div className="container">
            <Navbar />
            <p>
              GoCity is an implementation of the Code City metaphor for visualizing Go source code. Visit our repository
              for <a href="https://github.com/rodrigo-brito/gocity">more details.</a>
            </p>
            <p>
              You can also add a custom badge for your go repository.{' '}
              <a onClick={this.openModal} href="#">
                click here
              </a>{' '}
              to generate one.
            </p>
            <div className="field has-addons">
              <div className="control is-expanded">
                <input
                  onKeyPress={this.handleKeyPress}
                  onChange={this.onInputChange}
                  className="input"
                  id="repository"
                  type="text"
                  placeholder="eg: github.com/golang/go"
                  value={this.state.repository}
                />
              </div>
              <div className="control">
                <input
                  onKeyPress={this.handleKeyPress}
                  onChange={this.onInputChange}
                  className="input"
                  id="branch"
                  type="text"
                  placeholder="eg: master"
                  value={this.state.branch}
                />
              </div>
              <div className="control">
                <a
                  id="search"
                  onClick={this.onClick}
                  className="button is-info"
                >
                  Plot
                </a>
              </div>
            </div>
            <div className="level">
              <small className="level-left">
                Examples:{" "}
                {examples.map(example => (
                  <a
                    className="m-l-10"
                    key={example.link}
                    onClick={() => {
                      this.process(example.link, example.json, example.branch);
                    }}
                  >
                    {example.name}
                  </a>
                ))}
              </small>
            </div>
          </div>
        </header>
        <section className="canvas">
          {this.state.loading ? (
            <Loading message="Fetching repository..." />
          ) : (
              <BabylonScene
                width={window.innerWidth}
                engineOptions={{ preserveDrawingBuffer: true, stencil: true }}
                onSceneMount={this.onSceneMount}
              />
            )}
        </section>
        <div className="footer-warning notification is-danger is-hidden-tablet is-paddingless is-marginless is-unselectable">
          GoCity is best viewed on Desktop
        </div>
        <Legend />
      </main>
    );
  }
}

export default App;
