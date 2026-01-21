const DEFAULTS = { speed: 4, count: 500, width: 1280, height: 720 };
const STORAGE_KEYS = {
  speed: "flying_speed",
  count: "flying_count",
  width: "flying_width",
  height: "flying_height",
  autoload: "flying_autoload",
};

let speed = DEFAULTS.speed;
let windows = [];
let windowsNum = DEFAULTS.count;
let imgs = [];
let paused = false;
let config = { width: DEFAULTS.width, height: DEFAULTS.height };

// No sound usage: removed sound and error.mp3 dependency

const ASSET_FILES = [
  "anonymous-mask.png",
  "da-vinci.png",
  "grand-master-key.png",
  "hacker.png",
  "hanukkah-glass.png",
  "jain-symbol.png",
  "kaaba.png",
  "lilith-symbol.png",
  "lucifer-sigil.png",
  "magic-map.png",
  "mezuzah.png",
  "praying-beads.png",
  "seed-of-life.png",
  "spirituality.png",
  "star-of-david.png",
];

function preload() {
  for (let i = 0; i < ASSET_FILES.length; i++) {
    imgs[i] = loadImage("png/" + ASSET_FILES[i]);
  }
}

class Window {
  constructor() {
    this.x = random(-width, width);
    this.y = random(-height, height);
    this.z = random(width);
    this.pz = this.z;
    this.img = random(imgs);
  }

  update() {
    this.z = this.z - speed;

    if (this.z < 1) {
      this.z = width / 2;
      this.x = random(-width, width);
      this.y = random(-height, height);
      this.pz = this.z;
    }
  }

  show() {
    let sx = map(this.x / this.z, 0, 1, 0, width / 2);
    let sy = map(this.y / this.z, 0, 1, 0, height / 2);

    let r = map(this.z, 0, width / 2, 26, 4);

    image(this.img, sx, sy, r, r);
    // tint(this.color);

    this.pz = this.z;
  }
}

function setup() {
  // Load persisted values (safe guard)
  try {
    const autoloadRaw = localStorage.getItem(STORAGE_KEYS.autoload);
    const shouldAutoload = autoloadRaw === null ? true : autoloadRaw === "true";
    if (shouldAutoload) {
      const savedSpeed = parseInt(localStorage.getItem(STORAGE_KEYS.speed), 10);
      if (!Number.isNaN(savedSpeed)) speed = savedSpeed;
      const savedCount = parseInt(localStorage.getItem(STORAGE_KEYS.count), 10);
      if (!Number.isNaN(savedCount)) windowsNum = savedCount;
      const savedW = parseInt(localStorage.getItem(STORAGE_KEYS.width), 10);
      if (!Number.isNaN(savedW)) config.width = savedW;
      const savedH = parseInt(localStorage.getItem(STORAGE_KEYS.height), 10);
      if (!Number.isNaN(savedH)) config.height = savedH;
    }
  } catch (e) {
    // ignore storage errors
  }

  createCanvas(config.width, config.height);

  for (let i = 0; i < windowsNum; i++) {
    windows[i] = new Window();
  }

  // Expose agent setter wrappers
  window.AGENT_setSpeed = function (n) {
    if (typeof n === "number" && isFinite(n)) {
      speed = n;
      localStorage.setItem(STORAGE_KEYS.speed, String(n));
    }
  };
  window.AGENT_setCount = function (n) {
    if (typeof n === "number" && isFinite(n)) {
      windowsNum = n;
      localStorage.setItem(STORAGE_KEYS.count, String(n));
      if (windows.length < windowsNum) {
        for (let i = windows.length; i < windowsNum; i++)
          windows[i] = new Window();
      } else if (windows.length > windowsNum) {
        windows.length = windowsNum;
      }
    }
  };
  window.AGENT_setCanvasSize = function (w, h) {
    if (Number.isFinite(w) && Number.isFinite(h) && w > 0 && h > 0) {
      config.width = w;
      config.height = h;
      localStorage.setItem(STORAGE_KEYS.width, String(w));
      localStorage.setItem(STORAGE_KEYS.height, String(h));
      resizeCanvas(w, h);
    }
  };
}

function draw() {
  background(0);
  translate(width / 2, height / 2);
  if (!paused) {
    for (let i = 0; i < windows.length; i++) {
      windows[i].update();
      windows[i].show();
    }
  } else {
    // Pause rendering of motion; still show current frame
    for (let i = 0; i < windows.length; i++) {
      windows[i].show();
    }
    // Pause indicator
    push();
    fill(255);
    textSize(14);
    textAlign(RIGHT, BOTTOM);
    text("PAUSED (P)", width / 2 - 8, height / 2 - 8);
    pop();
  }
}

function keyPressed() {
  if (key === "P" || key === "p") {
    paused = !paused;
  }
}
