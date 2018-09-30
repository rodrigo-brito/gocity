import React, { Component } from "react";
import FloatBox from "./FloatBox";
import * as BABYLON from "babylonjs";
import BabylonScene from "./Scene";
import "./App.css";
import axios from "axios";
import Header from "./Header";
import Legend from "./Legend";

const playgroundSize = 1000;

const URLRegexp = new RegExp(/^(?:https:\/\/?)?(github\.com\/.*)/i);

const endpoint = "/api";

// TODO: isolate in the constants file
const colors = {
  PACKAGE: new BABYLON.Color3(0.5, 0.5, 0.5),
  FILE: new BABYLON.Color3(1, 1, 1),
  STRUCT: new BABYLON.Color3(32 / 255, 156 / 255, 238 / 255)
};

const examples = [
  {
    name: "rodrigo-brito/go-async-benchmark",
    link: "github.com/rodrigo-brito/go-async-benchmark"
  },
  {
    name: "rodrigo-brito/gocity",
    link: "github.com/rodrigo-brito/gocity"
  },
  {
    name: "sirupsen/logrus",
    link: "github.com/sirupsen/logrus"
  }
];

class App extends Component {
  canvas = null;
  scene = null;
  engine = null;
  camera = null;

  constructor(props) {
    super(props);
    this.state = {
      repository: "github.com/rodrigo-brito/go-async-benchmark"
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
  }

  onMouseMove(e) {
    this.mouse_x = e.pageX;
    this.mouse_y = e.pageY;
  }

  showTooltip(info) {
    console.log("INFO: ", info);
    this.setState({
      infoVisible: true,
      infoData: info,
      infoPosition: { x: this.mouse_x, y: this.mouse_y }
    });
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
      { width: data.width, depth: data.width, height: data.height },
      this.scene
    );
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
          console.log("Hover ", bar.info);
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
    bar.material.emissiveColor = data.color.scale(0.3);
    bar.material.specularColor = new BABYLON.Color3(0, 0, 0);

    // Shadows
    // shadowGenerator.getShadowMap().renderList.push(bar);

    return bar;
  };

  plot(children, parent) {
    if (!children) {
      return;
    }

    children.forEach(data => {
      var mesh = this.addBlock({
        x: data.position.x,
        y: data.position.y,
        width: data.size,
        height: data.numberOfMethods,
        label: "teste",
        color: colors[data.type],
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

  updateCamera(size) {
    this.camera.setPosition(new BABYLON.Vector3(size, size, size));
  }

  initScene() {
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
      "light1",
      new BABYLON.Vector3(0, 1, 0),
      this.scene
    );

    // Default intensity is 1. Let's dim the light a small amount
    light.intensity = 0.7;

    var ground = BABYLON.Mesh.CreateGround(
      "ground",
      playgroundSize,
      playgroundSize,
      1,
      this.scene,
      false
    );

    var groundMaterial = new BABYLON.StandardMaterial("ground", this.scene);

    groundMaterial.diffuseColor = new BABYLON.Color3(0.5, 0.5, 0.5);
    groundMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
    ground.material = groundMaterial;
    ground.receiveShadows = true;
    ground.position.y = -0.1;
    ground.isVisible = false;
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

  onInputChange(e) {
    this.setState({ repository: e.target.value });
  }

  process(repository) {
    const match = URLRegexp.exec(repository);
    if (!match) {
      alert("url inválida");
      return;
    }

    this.setState({ repository: match[1] });

    axios
      .get(endpoint, {
        params: {
          q: match[1]
        }
      })
      .then(response => {
        console.log(response);
        this.reset();
        this.plot(response.data.children);
        this.updateCamera(response.data.size);
      })
      .catch(e => {
        alert("Erro ao processar projeto");
        console.error(e);
      });
  }

  onClick() {
    this.process(this.state.repository);
  }

  render() {
    return (
      <main onMouseMove={this.onMouseMove}>
        <FloatBox
          position={this.state.infoPosition}
          info={this.state.infoData}
          visible={this.state.infoVisible}
        />
        <header className="header">
          <div className="container">
            <Header />
            Examples:{" "}
            <span>
              {examples.map(example => (
                <a
                  className="m-l-10"
                  key={example.link}
                  onClick={() => {
                    this.process(example.link);
                  }}
                >
                  {example.name}
                </a>
              ))}
            </span>
            <div className="field has-addons">
              <div className="control is-expanded">
                <input
                  onChange={this.onInputChange}
                  className="input"
                  id="repository"
                  type="text"
                  placeholder="eg: github.com/golang/go"
                  value={this.state.repository}
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
          </div>
        </header>
        <section className="canvas">
          <BabylonScene
            engineOptions={{ preserveDrawingBuffer: true, stencil: true }}
            onSceneMount={this.onSceneMount}
          />
        </section>
        <Legend />
      </main>
    );
  }
}

export default App;
