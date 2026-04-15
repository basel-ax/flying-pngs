const DEFAULTS = { speed: 4, count: 500, width: 1280, height: 720 };
const STORAGE_KEYS = {
  speed: "flying_speed",
  count: "flying_count",
  width: "flying_width",
  height: "flying_height",
  autoload: "flying_autoload",
  randomizeMode: "flying_randomize_mode",
  randomizeN: "flying_randomize_n",
  whiteMode: "flying_white_mode",
};

let speed = DEFAULTS.speed;
let baseSpeed = DEFAULTS.speed; // tracks user-set speed, unaffected by randomize slowdown
let windows = [];
let windowsNum = DEFAULTS.count;
let imgs = [];
let paused = false;
let started = false;
let whiteMode = false;
let config = { width: DEFAULTS.width, height: DEFAULTS.height };

// Bright tint colors for dark background
const TINT_COLORS = [
  [255, 255, 255],   // white
  [0, 174, 255],     // blue
  [255, 69, 58],     // red
  [50, 215, 75],     // green
  [255, 214, 10],    // yellow
  [191, 90, 242],    // purple
  [100, 210, 255],   // light blue
  [255, 159, 10],    // orange
  [172, 142, 104],   // tan
  [0, 199, 190],     // teal
];

// Dark tint colors for white background
const DARK_TINT_COLORS = [
  [0, 0, 0],         // black
  [0, 80, 180],      // dark blue
  [180, 30, 20],     // dark red
  [20, 130, 40],     // dark green
  [160, 120, 0],     // dark yellow
  [100, 40, 160],    // dark purple
  [0, 100, 160],     // dark cyan
  [180, 90, 0],      // dark orange
  [100, 70, 50],     // dark brown
  [0, 120, 110],     // dark teal
];

// Randomize mode state
let randomizeMode = false;
let randomizeN = 10; // max seconds for random interval
let isSlowed = false; // whether we're currently in the slowed phase
let nextRandomizeTime = 0; // millis() timestamp when next toggle fires

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
    imgs[i] = loadImage(
      "png/" + ASSET_FILES[i],
      function () {},
      function () {
        console.warn("[flying-pngs] Failed to load asset: " + ASSET_FILES[i]);
      },
    );
  }
}

class Window {
  constructor() {
    this.x = random(-width, width);
    this.y = random(-height, height);
    this.z = random(width);
    this.pz = this.z;
    this.img = random(imgs);
    this.colorIndex = floor(random(TINT_COLORS.length));
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

    const palette = whiteMode ? DARK_TINT_COLORS : TINT_COLORS;
    const c = palette[this.colorIndex];
    tint(c[0], c[1], c[2]);
    image(this.img, sx, sy, r, r);
    noTint();

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
      if (!Number.isNaN(savedSpeed)) {
        speed = savedSpeed;
        baseSpeed = savedSpeed;
      }
      const savedCount = parseInt(localStorage.getItem(STORAGE_KEYS.count), 10);
      if (!Number.isNaN(savedCount)) windowsNum = savedCount;
      const savedW = parseInt(localStorage.getItem(STORAGE_KEYS.width), 10);
      if (!Number.isNaN(savedW)) config.width = savedW;
      const savedH = parseInt(localStorage.getItem(STORAGE_KEYS.height), 10);
      if (!Number.isNaN(savedH)) config.height = savedH;

      // Randomize mode persistence
      const savedRandMode = localStorage.getItem(STORAGE_KEYS.randomizeMode);
      if (savedRandMode !== null) randomizeMode = savedRandMode === "true";
      const savedRandN = parseFloat(
        localStorage.getItem(STORAGE_KEYS.randomizeN),
      );
      if (!Number.isNaN(savedRandN) && savedRandN > 0) randomizeN = savedRandN;

      // White mode persistence
      const savedWhiteMode = localStorage.getItem(STORAGE_KEYS.whiteMode);
      if (savedWhiteMode !== null) whiteMode = savedWhiteMode === "true";
    }
  } catch (e) {
    // ignore storage errors
  }

  createCanvas(config.width, config.height);

  // Convert black silhouette PNGs to white so tint() can colorize them
  for (let i = 0; i < imgs.length; i++) {
    imgs[i].loadPixels();
    for (let j = 0; j < imgs[i].pixels.length; j += 4) {
      imgs[i].pixels[j] = 255;     // R
      imgs[i].pixels[j + 1] = 255; // G
      imgs[i].pixels[j + 2] = 255; // B
      // keep alpha as-is
    }
    imgs[i].updatePixels();
  }

  // ── Start animation wrapper ────────────────────────────────────────────────

  window.AGENT_startAnimation = function () {
    if (started) return;
    started = true;
    windows = [];
    for (let i = 0; i < windowsNum; i++) {
      windows[i] = new Window();
    }
    nextRandomizeTime = millis() + Math.random() * randomizeN * 1000;
  };

  // ── Agent setter wrappers ──────────────────────────────────────────────────

  window.AGENT_setSpeed = function (n) {
    if (typeof n === "number" && isFinite(n)) {
      baseSpeed = n;
      // Only apply directly to speed when not currently slowed
      speed = isSlowed && randomizeMode ? n / 5 : n;
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

  /**
   * Enable or disable randomize mode.
   * @param {boolean} enabled
   */
  window.AGENT_setRandomizeMode = function (enabled) {
    randomizeMode = !!enabled;
    localStorage.setItem(STORAGE_KEYS.randomizeMode, String(randomizeMode));

    if (randomizeMode) {
      // Start fresh: not slowed, schedule first event
      isSlowed = false;
      speed = baseSpeed;
      nextRandomizeTime = millis() + Math.random() * randomizeN * 1000;
    } else {
      // Restore full speed when mode is turned off
      isSlowed = false;
      speed = baseSpeed;
    }
  };

  /**
   * Set the maximum interval (in seconds) for the randomize mode timer.
   * @param {number} n  positive number of seconds
   */
  window.AGENT_setRandomizeN = function (n) {
    if (typeof n === "number" && isFinite(n) && n > 0) {
      randomizeN = n;
      localStorage.setItem(STORAGE_KEYS.randomizeN, String(n));
      // Re-schedule so the new N takes effect immediately for the next tick
      nextRandomizeTime = millis() + Math.random() * randomizeN * 1000;
    }
  };

  window.AGENT_setWhiteMode = function (enabled) {
    whiteMode = !!enabled;
    localStorage.setItem(STORAGE_KEYS.whiteMode, String(whiteMode));
  };
}

function draw() {
  background(whiteMode ? 255 : 0);
  translate(width / 2, height / 2);

  if (!started) {
    push();
    fill(whiteMode ? 0 : 255, whiteMode ? 0 : 255, whiteMode ? 0 : 255, 120);
    textSize(18);
    textAlign(CENTER, CENTER);
    text("Set parameters and press Start", 0, 0);
    pop();
    return;
  }

  // ── Randomize mode tick ───────────────────────────────────────────────────
  if (randomizeMode && !paused) {
    const now = millis();
    if (now >= nextRandomizeTime) {
      // Toggle slow/normal
      isSlowed = !isSlowed;
      speed = isSlowed ? baseSpeed / 5 : baseSpeed;
      // Schedule next toggle at a new random interval within [0, N] seconds
      nextRandomizeTime = now + Math.random() * randomizeN * 1000;
    }
  }

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
    fill(whiteMode ? 0 : 255);
    textSize(14);
    textAlign(RIGHT, BOTTOM);
    text("PAUSED (P)", width / 2 - 8, height / 2 - 8);
    pop();
  }

  // ── Randomize mode slow indicator ────────────────────────────────────────
  if (randomizeMode && isSlowed && !paused) {
    push();
    fill(255, 200, 0);
    textSize(13);
    textAlign(LEFT, BOTTOM);
    text("SLOW ×1/5", -width / 2 + 8, height / 2 - 8);
    pop();
  }
}

function keyPressed() {
  if (key === "P" || key === "p") {
    paused = !paused;
  }
}
