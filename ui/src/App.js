import React, { Component } from "react";
import FloatBox from "./FloatBox";
import * as BABYLON from "babylonjs";
import BabylonScene from "./Scene";
import "./App.css";

class App extends Component {
  onSceneMount = e => {
    const { canvas, scene, engine } = e;

    // This creates and positions a free camera (non-mesh)
    var camera = new BABYLON.FreeCamera(
      "camera1",
      new BABYLON.Vector3(0, 5, -10),
      scene
    );

    // This targets the camera to scene origin
    camera.setTarget(BABYLON.Vector3.Zero());

    // This attaches the camera to the canvas
    camera.attachControl(canvas, true);

    // This creates a light, aiming 0,1,0 - to the sky (non-mesh)
    var light = new BABYLON.HemisphericLight(
      "light1",
      new BABYLON.Vector3(0, 1, 0),
      scene
    );

    // Default intensity is 1. Let's dim the light a small amount
    light.intensity = 0.7;

    // Our built-in 'sphere' shape. Params: name, subdivs, size, scene
    var sphere = BABYLON.Mesh.CreateSphere("sphere1", 16, 2, scene);

    // Move the sphere upward 1/2 its height
    sphere.position.y = 1;

    // Our built-in 'ground' shape. Params: name, width, depth, subdivs, scene
    // TODO var ground = BABYLON.Mesh.CreateGround("ground1", 6, 6, 2, scene);

    engine.runRenderLoop(() => {
      if (scene) {
        scene.render();
      }
    });
  };
  render() {
    return (
      <main>
        <FloatBox />
        <header className="header">
          <div className="container">
            <h1 className="title">GoCity</h1>
            <span className="subtitle">Source code visualization</span>
            <div className="field has-addons">
              <div className="control">
                <input
                  className="input"
                  id="repository"
                  type="text"
                  placeholder="eg: github.com/golang/go"
                  value="github.com/rodrigo-brito/go-async-benchmark"
                />
              </div>
              <div className="control">
                <a id="search" className="button is-info">
                  Plot
                </a>
              </div>
            </div>
          </div>
        </header>
        <section className="canvas">
          <BabylonScene onSceneMount={this.onSceneMount} />
        </section>
      </main>
    );
  }
}

export default App;
