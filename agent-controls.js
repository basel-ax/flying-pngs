(function () {
  "use strict";

  const DEFAULTS = {
    speed: 4,
    count: 500,
    width: 1280,
    height: 720,
    randomizeMode: false,
    randomizeN: 10,
  };

  const STORAGE_KEYS = {
    speed: "flying_speed",
    count: "flying_count",
    width: "flying_width",
    height: "flying_height",
    autoload: "flying_autoload",
    randomizeMode: "flying_randomize_mode",
    randomizeN: "flying_randomize_n",
  };

  function clamp(n, min, max) {
    return Math.max(min, Math.min(max, n));
  }

  function readParams() {
    const params = new URLSearchParams(window.location.search);
    const count = parseInt(params.get("count"), 10);
    const speed = parseInt(params.get("speed"), 10);
    const width = parseInt(params.get("width"), 10);
    const height = parseInt(params.get("height"), 10);
    const randomizeN = parseFloat(params.get("randomizeN"));
    const randomizeModeRaw = params.get("randomizeMode");
    return {
      count: Number.isFinite(count) ? count : null,
      speed: Number.isFinite(speed) ? speed : null,
      width: Number.isFinite(width) ? width : null,
      height: Number.isFinite(height) ? height : null,
      randomizeN:
        Number.isFinite(randomizeN) && randomizeN > 0 ? randomizeN : null,
      randomizeMode:
        randomizeModeRaw !== null ? randomizeModeRaw === "true" : null,
    };
  }

  function readStored(key, fallback) {
    const raw = localStorage.getItem(key);
    if (raw === null || raw === "") return fallback;
    const n = parseInt(raw, 10);
    return Number.isFinite(n) ? n : fallback;
  }

  function readStoredFloat(key, fallback) {
    const raw = localStorage.getItem(key);
    if (raw === null || raw === "") return fallback;
    const n = parseFloat(raw);
    return Number.isFinite(n) && n > 0 ? n : fallback;
  }

  function readStoredBool(key, fallback) {
    const raw = localStorage.getItem(key);
    if (raw === null || raw === "") return fallback;
    return raw === "true";
  }

  function persistDefaultsIfMissing() {
    if (localStorage.getItem(STORAGE_KEYS.speed) === null) {
      localStorage.setItem(STORAGE_KEYS.speed, String(DEFAULTS.speed));
    }
    if (localStorage.getItem(STORAGE_KEYS.count) === null) {
      localStorage.setItem(STORAGE_KEYS.count, String(DEFAULTS.count));
    }
    if (localStorage.getItem(STORAGE_KEYS.width) === null) {
      localStorage.setItem(STORAGE_KEYS.width, String(DEFAULTS.width));
    }
    if (localStorage.getItem(STORAGE_KEYS.height) === null) {
      localStorage.setItem(STORAGE_KEYS.height, String(DEFAULTS.height));
    }
    if (localStorage.getItem(STORAGE_KEYS.autoload) === null) {
      localStorage.setItem(STORAGE_KEYS.autoload, "true");
    }
    if (localStorage.getItem(STORAGE_KEYS.randomizeMode) === null) {
      localStorage.setItem(
        STORAGE_KEYS.randomizeMode,
        String(DEFAULTS.randomizeMode),
      );
    }
    if (localStorage.getItem(STORAGE_KEYS.randomizeN) === null) {
      localStorage.setItem(
        STORAGE_KEYS.randomizeN,
        String(DEFAULTS.randomizeN),
      );
    }
  }

  function ensurePanel() {
    if (document.getElementById("agent-panel")) return;

    const panel = document.createElement("div");
    panel.id = "agent-panel";
    panel.style.position = "fixed";
    panel.style.bottom = "12px";
    panel.style.right = "12px";
    panel.style.background = "rgba(0,0,0,0.6)";
    panel.style.color = "#fff";
    panel.style.padding = "10px";
    panel.style.borderRadius = "8px";
    panel.style.fontFamily = "Arial, sans-serif";
    panel.style.fontSize = "12px";
    panel.style.zIndex = "9999";
    panel.style.minWidth = "260px";
    panel.style.boxShadow = "0 2px 8px rgba(0,0,0,.4)";

    panel.innerHTML = `
      <div style="margin-bottom:6px; font-weight:700; text-align:left;">Agent Controls</div>
      <div style="margin-bottom:6px;">
        <label>Speed: <span id="ag_speed_val" style="margin-left:6px; width:28px; display:inline-block;">${DEFAULTS.speed}</span></label><br/>
        <input id="ag_speed" type="range" min="1" max="20" step="1" value="${DEFAULTS.speed}" style="width:180px;">
      </div>
      <div style="margin:6px 0 6px 0;">
        <label>Count: <span id="ag_count_val" style="margin-left:6px; width:40px; display:inline-block;">${DEFAULTS.count}</span></label><br/>
        <input id="ag_count" type="range" min="50" max="1000" step="10" value="${DEFAULTS.count}" style="width:180px;">
      </div>
      <div style="margin-top:6px; display:flex; gap:8px; align-items:center;">
        <label>Width: <input id="ag_width" type="number" min="320" max="1920" value="${DEFAULTS.width}" style="width:90px;"></label>
        <label>Height: <input id="ag_height" type="number" min="240" max="1080" value="${DEFAULTS.height}" style="width:90px;"></label>
      </div>
      <div style="margin-top:8px; display:flex; justify-content:space-between; align-items:center;">
        <label style="display:flex; align-items:center; gap:6px;">
          <input id="ag_autoload" type="checkbox" checked>
          Auto-load last session
        </label>
        <button id="ag_reset_defaults">Reset to defaults</button>
      </div>

      <div style="margin-top:10px; border-top:1px solid rgba(255,255,255,0.2); padding-top:8px;">
        <div style="font-weight:700; margin-bottom:6px; display:flex; align-items:center; gap:8px;">
          Randomize Mode
          <label style="display:flex; align-items:center; gap:4px; font-weight:normal; cursor:pointer;">
            <input id="ag_randomize_mode" type="checkbox">
            <span id="ag_randomize_mode_label">Off</span>
          </label>
        </div>
        <div style="display:flex; align-items:center; gap:8px;">
          <label style="display:flex; align-items:center; gap:4px;">
            Max interval&nbsp;<em>N</em>&nbsp;(s):
            <input
              id="ag_randomize_n"
              type="number"
              min="1"
              max="300"
              step="1"
              value="${DEFAULTS.randomizeN}"
              style="width:56px; margin-left:4px;"
            >
          </label>
        </div>
        <div style="margin-top:4px; font-size:11px; color:rgba(255,255,255,0.55); line-height:1.4;">
          When on: every&nbsp;<em>random(N)</em>&nbsp;s the speed<br>drops to&nbsp;×1/5 and then restores.
        </div>
      </div>
    `;

    document.body.appendChild(panel);

    document
      .getElementById("ag_reset_defaults")
      .addEventListener("click", () => {
        const w = DEFAULTS.width;
        const h = DEFAULTS.height;
        const s = DEFAULTS.speed;
        const c = DEFAULTS.count;

        document.getElementById("ag_width").value = w;
        document.getElementById("ag_height").value = h;
        document.getElementById("ag_speed").value = s;
        document.getElementById("ag_count").value = c;

        // Also reset randomize mode
        document.getElementById("ag_randomize_mode").checked =
          DEFAULTS.randomizeMode;
        document.getElementById("ag_randomize_mode_label").textContent =
          DEFAULTS.randomizeMode ? "On" : "Off";
        document.getElementById("ag_randomize_n").value = DEFAULTS.randomizeN;

        if (typeof window.AGENT_setCanvasSize === "function")
          window.AGENT_setCanvasSize(w, h);
        if (typeof window.AGENT_setSpeed === "function")
          window.AGENT_setSpeed(s);
        if (typeof window.AGENT_setCount === "function")
          window.AGENT_setCount(c);
        if (typeof window.AGENT_setRandomizeMode === "function")
          window.AGENT_setRandomizeMode(DEFAULTS.randomizeMode);
        if (typeof window.AGENT_setRandomizeN === "function")
          window.AGENT_setRandomizeN(DEFAULTS.randomizeN);

        localStorage.setItem(STORAGE_KEYS.width, String(w));
        localStorage.setItem(STORAGE_KEYS.height, String(h));
        localStorage.setItem(STORAGE_KEYS.speed, String(s));
        localStorage.setItem(STORAGE_KEYS.count, String(c));
        localStorage.setItem(
          STORAGE_KEYS.randomizeMode,
          String(DEFAULTS.randomizeMode),
        );
        localStorage.setItem(
          STORAGE_KEYS.randomizeN,
          String(DEFAULTS.randomizeN),
        );

        const speedVal = document.getElementById("ag_speed_val");
        const countVal = document.getElementById("ag_count_val");
        if (speedVal) speedVal.textContent = String(s);
        if (countVal) countVal.textContent = String(c);
      });
  }

  function wire() {
    const speedInput = document.getElementById("ag_speed");
    const countInput = document.getElementById("ag_count");
    const widthInput = document.getElementById("ag_width");
    const heightInput = document.getElementById("ag_height");
    const autoloadInput = document.getElementById("ag_autoload");
    const speedVal = document.getElementById("ag_speed_val");
    const countVal = document.getElementById("ag_count_val");
    const randomizeModeInput = document.getElementById("ag_randomize_mode");
    const randomizeModeLabel = document.getElementById(
      "ag_randomize_mode_label",
    );
    const randomizeNInput = document.getElementById("ag_randomize_n");

    function updateDisplay() {
      speedVal.textContent = speedInput.value;
      countVal.textContent = countInput.value;
    }

    speedInput.addEventListener("input", () => {
      const v = parseInt(speedInput.value, 10);
      if (typeof window.AGENT_setSpeed === "function") window.AGENT_setSpeed(v);
      localStorage.setItem(STORAGE_KEYS.speed, String(v));
      updateDisplay();
    });

    countInput.addEventListener("input", () => {
      const v = parseInt(countInput.value, 10);
      if (typeof window.AGENT_setCount === "function") window.AGENT_setCount(v);
      localStorage.setItem(STORAGE_KEYS.count, String(v));
      updateDisplay();
    });

    function applySize() {
      const w = parseInt(widthInput.value, 10) || DEFAULTS.width;
      const h = parseInt(heightInput.value, 10) || DEFAULTS.height;
      if (typeof window.AGENT_setCanvasSize === "function")
        window.AGENT_setCanvasSize(w, h);
      localStorage.setItem(STORAGE_KEYS.width, String(w));
      localStorage.setItem(STORAGE_KEYS.height, String(h));
    }

    widthInput.addEventListener("change", applySize);
    heightInput.addEventListener("change", applySize);

    autoloadInput.addEventListener("change", () => {
      localStorage.setItem(
        STORAGE_KEYS.autoload,
        autoloadInput.checked ? "true" : "false",
      );
    });

    // ── Randomize mode wiring ────────────────────────────────────────────────

    randomizeModeInput.addEventListener("change", () => {
      const enabled = randomizeModeInput.checked;
      randomizeModeLabel.textContent = enabled ? "On" : "Off";
      if (typeof window.AGENT_setRandomizeMode === "function")
        window.AGENT_setRandomizeMode(enabled);
      localStorage.setItem(STORAGE_KEYS.randomizeMode, String(enabled));
    });

    randomizeNInput.addEventListener("change", () => {
      const v = parseFloat(randomizeNInput.value);
      if (!Number.isFinite(v) || v <= 0) return;
      const clamped = Math.min(300, Math.max(1, v));
      randomizeNInput.value = clamped;
      if (typeof window.AGENT_setRandomizeN === "function")
        window.AGENT_setRandomizeN(clamped);
      localStorage.setItem(STORAGE_KEYS.randomizeN, String(clamped));
    });

    updateDisplay();
  }

  function applyStoredIfEnabled() {
    const autoload = readStoredBool(STORAGE_KEYS.autoload, true);
    const autoloadInput = document.getElementById("ag_autoload");
    if (autoloadInput) autoloadInput.checked = autoload;

    if (!autoload) return;

    const s = readStored(STORAGE_KEYS.speed, DEFAULTS.speed);
    const c = readStored(STORAGE_KEYS.count, DEFAULTS.count);
    const w = readStored(STORAGE_KEYS.width, DEFAULTS.width);
    const h = readStored(STORAGE_KEYS.height, DEFAULTS.height);
    const rMode = readStoredBool(
      STORAGE_KEYS.randomizeMode,
      DEFAULTS.randomizeMode,
    );
    const rN = readStoredFloat(STORAGE_KEYS.randomizeN, DEFAULTS.randomizeN);

    const speedInput = document.getElementById("ag_speed");
    const countInput = document.getElementById("ag_count");
    const widthInput = document.getElementById("ag_width");
    const heightInput = document.getElementById("ag_height");
    const randomizeModeInput = document.getElementById("ag_randomize_mode");
    const randomizeModeLabel = document.getElementById(
      "ag_randomize_mode_label",
    );
    const randomizeNInput = document.getElementById("ag_randomize_n");

    if (speedInput) speedInput.value = s;
    if (countInput) countInput.value = c;
    if (widthInput) widthInput.value = w;
    if (heightInput) heightInput.value = h;
    if (randomizeModeInput) {
      randomizeModeInput.checked = rMode;
      if (randomizeModeLabel)
        randomizeModeLabel.textContent = rMode ? "On" : "Off";
    }
    if (randomizeNInput) randomizeNInput.value = rN;

    if (typeof window.AGENT_setSpeed === "function") window.AGENT_setSpeed(s);
    if (typeof window.AGENT_setCount === "function") window.AGENT_setCount(c);
    if (typeof window.AGENT_setCanvasSize === "function")
      window.AGENT_setCanvasSize(w, h);
    if (typeof window.AGENT_setRandomizeMode === "function")
      window.AGENT_setRandomizeMode(rMode);
    if (typeof window.AGENT_setRandomizeN === "function")
      window.AGENT_setRandomizeN(rN);

    const speedVal = document.getElementById("ag_speed_val");
    const countVal = document.getElementById("ag_count_val");
    if (speedVal) speedVal.textContent = String(s);
    if (countVal) countVal.textContent = String(c);
  }

  function applyParams() {
    const params = readParams();
    if (params.count != null) {
      const v = clamp(params.count, 50, 1000);
      document.getElementById("ag_count").value = v;
      if (typeof window.AGENT_setCount === "function") window.AGENT_setCount(v);
      localStorage.setItem(STORAGE_KEYS.count, String(v));
    }
    if (params.speed != null) {
      const v = clamp(params.speed, 1, 20);
      document.getElementById("ag_speed").value = v;
      if (typeof window.AGENT_setSpeed === "function") window.AGENT_setSpeed(v);
      localStorage.setItem(STORAGE_KEYS.speed, String(v));
    }
    if (params.width != null || params.height != null) {
      const w =
        params.width != null
          ? clamp(params.width, 320, 1920)
          : parseInt(document.getElementById("ag_width").value, 10);
      const h =
        params.height != null
          ? clamp(params.height, 240, 1080)
          : parseInt(document.getElementById("ag_height").value, 10);
      document.getElementById("ag_width").value = w;
      document.getElementById("ag_height").value = h;
      if (typeof window.AGENT_setCanvasSize === "function")
        window.AGENT_setCanvasSize(w, h);
      localStorage.setItem(STORAGE_KEYS.width, String(w));
      localStorage.setItem(STORAGE_KEYS.height, String(h));
    }
    if (params.randomizeMode !== null) {
      const enabled = params.randomizeMode;
      const input = document.getElementById("ag_randomize_mode");
      const label = document.getElementById("ag_randomize_mode_label");
      if (input) {
        input.checked = enabled;
      }
      if (label) {
        label.textContent = enabled ? "On" : "Off";
      }
      if (typeof window.AGENT_setRandomizeMode === "function")
        window.AGENT_setRandomizeMode(enabled);
      localStorage.setItem(STORAGE_KEYS.randomizeMode, String(enabled));
    }
    if (params.randomizeN !== null) {
      const v = Math.min(300, Math.max(1, params.randomizeN));
      const input = document.getElementById("ag_randomize_n");
      if (input) input.value = v;
      if (typeof window.AGENT_setRandomizeN === "function")
        window.AGENT_setRandomizeN(v);
      localStorage.setItem(STORAGE_KEYS.randomizeN, String(v));
    }
  }

  document.addEventListener("DOMContentLoaded", () => {
    persistDefaultsIfMissing();
    ensurePanel();
    wire();
    applyStoredIfEnabled();
    applyParams();
  });
})();
