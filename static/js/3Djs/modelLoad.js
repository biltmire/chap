var scene = new THREE.Scene();
var camera = new THREE.PerspectiveCamera( 75, window.innerWidth / window.innerHeight, 0.1, 1000 );
camera.position.z = 30;
var renderer = new THREE.WebGLRenderer();
renderer.setSize( window.innerWidth, window.innerHeight );
document.body.appendChild( renderer.domElement );
var controls = new THREE.OrbitControls(camera, renderer.domElement);
controls.enableDamping = true;
controls.dampingFactor = 0.25;
controls.enableZoom = true;
var keyLight = new THREE.DirectionalLight(0xffffff, 1.0);
keyLight.position.set(-100, 0, 100);
var backLight = new THREE.DirectionalLight(0xffffff, 0.);
backLight.position.set(100, 0, -100).normalize();
var softLight = new THREE.AmbientLight( 0x404040, 0.5); // soft white light
softLight.position.set(100, 0, -100).normalize();
scene.add(softLight);
scene.add(keyLight);
scene.add(backLight);

var objLoader = new THREE.OBJLoader();
objLoader.setPath('/static/js/3Djs/model/');
objLoader.load('chess.obj', function (object) {
scene.add(object);
object.position.set(0, 0, 0);
object.rotation.set(-20.2,50.25,-37.55);

        });
var animate = function () {
    requestAnimationFrame( animate );
    controls.update();
    renderer.render(scene, camera);
};
animate();