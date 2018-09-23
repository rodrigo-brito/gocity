$(function() {
    var mouse_x = 0, mouse_y = 0;
    var canvas = document.getElementById("renderCanvas");

    var createScene = function(engine) {
        var scene = new BABYLON.Scene(engine);
        var light = new BABYLON.DirectionalLight(
            "global_light",
            new BABYLON.Vector3(0, -0.5, 1.0),
            scene
        );
        var camera = new BABYLON.ArcRotateCamera(
            "camera",
            0,
            0,
            10,
            BABYLON.Vector3.Zero(),
            scene
        );
        camera.setPosition(new BABYLON.Vector3(500, 400, -100));
        light.position = new BABYLON.Vector3(0, 25, -50);

        camera.attachControl(canvas, true);

        // Data
        var playgroundSize = 1000;

        // Ground
        var ground = BABYLON.Mesh.CreateGround(
            "ground",
            playgroundSize,
            playgroundSize,
            1,
            scene,
            false
        );
        var groundMaterial = new BABYLON.StandardMaterial("ground", scene);
        groundMaterial.diffuseColor = new BABYLON.Color3(0.5, 0.5, 0.5);
        groundMaterial.specularColor = new BABYLON.Color3(0, 0, 0);
        ground.material = groundMaterial;
        ground.receiveShadows = true;
        ground.position.y = -0.1;
        ground.isVisible = false;

        var shadowGenerator = new BABYLON.ShadowGenerator(1024, light);
        shadowGenerator.usePoissonSampling = true;

        function addBlock(data) {
            var bar = BABYLON.MeshBuilder.CreateBox(data.label, {width: data.width, depth: data.width, height: data.height}, scene);
            if (data.parent) {
                bar.parent = data.parent;

                var bounds = data.parent.getBoundingInfo();
                bar.position.y = bounds.maximum.y + (data.height / 2.0);
            }
            bar.position.x = data.x || 0;
            bar.position.z = data.y || 0;

            bar.info = data.info;

            bar.actionManager = new BABYLON.ActionManager(scene);
            bar.actionManager.registerAction(new BABYLON.ExecuteCodeAction(
                {trigger: BABYLON.ActionManager.OnPointerOverTrigger},
                function () {
                    showTooltip(bar.info);
                    console.log("Hover ", bar.info);
                }
            ));

            bar.actionManager.registerAction(new BABYLON.ExecuteCodeAction(
                {trigger: BABYLON.ActionManager.OnPointerOutTrigger},
                function () {
                    hideTooltip();
                }
            ));

    //     // Animate a bit
    //     var animation = new BABYLON.Animation(
    //       "anim",
    //       "scaling",
    //       30,
    //       BABYLON.Animation.ANIMATIONTYPE_VECTOR3
    //     );
    //     animation.setKeys([
    //       { frame: 0, value: new BABYLON.Vector3(data.width, 0, data.width) },
    //       {
    //         frame: 100,
    //         value: new BABYLON.Vector3(data.width, data.height * scale, data.width)
    //       }
    //     ]);
    //     bar.animations.push(animation);

    //     animation = new BABYLON.Animation(
    //       "anim2",
    //       "position.y",
    //       30,
    //       BABYLON.Animation.ANIMATIONTYPE_FLOAT
    //     );
    //     animation.setKeys([
    //       { frame: 0, value: 0 },
    //       { frame: 100, value: data.height * scale / 2 }
    //     ]);
    //     bar.animations.push(animation);
    //     scene.beginAnimation(bar, 0, 100, false, 2.0);

            // Material
            bar.material = new BABYLON.StandardMaterial(data.label + "mat", scene);
            bar.material.diffuseColor = data.color;
            bar.material.emissiveColor = data.color.scale(0.3);
            bar.material.specularColor = new BABYLON.Color3(0, 0, 0);

            // Shadows
            shadowGenerator.getShadowMap().renderList.push(bar);

            return bar;
        }

        var colors = {
            "PACKAGE": new BABYLON.Color3(1, 0, 0),
            "FILE": new BABYLON.Color3(1, 1, 1),
            "STRUCT": new BABYLON.Color3(0, 0, 1)
        };

        function plot(children, parent){
            if (!children) {
                return
            }

            children.map((data) => {
                console.log("ploting", data.name)
                var mesh = addBlock({
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
                        NOA: data.numberOfAttributes,
                    }
                });
                if (parent) {
                    mesh.parent = parent;
                }
                if (data.children && data.children.length > 0){
                    plot(data.children, mesh)
                }
            })
        }

        $("#search").click(function (e) {
            var value = $("#repository").val();
            console.log(value);
            $.get("/api", {
                "q": value
            }).done(function (data) {
                if (data) {
                    plot(data.children);
                }
            }).fail(function (err) {
                console.error(err);
            })
        });


        // Limit camera
        camera.lowerAlphaLimit = Math.PI;
        camera.upperAlphaLimit = 2 * Math.PI;
        camera.lowerBetaLimit = 0.1;
        camera.upperBetaLimit = Math.PI / 2 * 0.99;
        camera.lowerRadiusLimit = 5;
        camera.upperRadiusLimit = 150;

        return scene;
    };

    var engine = new BABYLON.Engine(canvas, true, {
        preserveDrawingBuffer: true,
        stencil: true
    });
    var scene = createScene();

    engine.runRenderLoop(function() {
        if (scene) {
            scene.render();
        }
    });

    // Resize
    window.addEventListener("resize", function() {
        engine.resize();
    });

    function showTooltip(data) {
        var card = $(".float-card");
        card.css({"top": mouse_y, "left": mouse_x}).show();
        card.find(".lines").text(data.NOL);
        card.find(".methods").text(data.NOM);
        card.find(".attributes").text(data.NOA);
        card.find(".name").text(data.name + " [" + data.type + "]");
    }

    function hideTooltip() {
        $(".float-card").hide();
    }

    let handleMousemove = (event) => {
        mouse_x = event.x;
        mouse_y = event.y;
    };

    let throttle = (func, delay) => {
        let prev = Date.now() - delay;
        return (...args) => {
            let current = Date.now();
            if (current - prev >= delay) {
                prev = current;
                func.apply(null, args);
            }
        }
    };

    document.addEventListener('mousemove', throttle(handleMousemove, 100));
});
