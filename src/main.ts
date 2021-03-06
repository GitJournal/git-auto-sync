import { app, BrowserWindow, Menu, Tray, nativeImage } from "electron";
import * as path from "path";

function createWindow() {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, "preload.js"),
    },
    width: 800,
  });

  // and load the index.html of the app.
  mainWindow.loadFile(path.join(__dirname, "../index.html"));

  // Open the DevTools.
  mainWindow.webContents.openDevTools();
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on("ready", () => {
  createWindow();

  app.on("activate", function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});

// In this file you can include the rest of your app"s specific main process
// code. You can also put them in separate files and require them here.

//
// Tray Code
//
var trayIconPath = path.join(__dirname, "../assets/folder-green-git-icon.png");

let tray = null;
app.whenReady().then(() => {
  var image = nativeImage.createFromPath(trayIconPath);

  tray = new Tray(image.resize({ width: 32, height: 32 }));
  console.log(`Loaded ${trayIconPath} ${tray}`);

  const contextMenu = Menu.buildFromTemplate([
    { label: "Item1", type: "normal" },
    { type: "separator" },
    { label: "Item3", type: "radio", checked: true },
    { label: "Item4", type: "radio" },
    { label: "Quit", role: "quit" },
  ]);

  tray.setToolTip("GitAutoSync");
  tray.setContextMenu(contextMenu);
});

console.log("huh?");

// Main logic
/*
  cd /Users/vishesh/notes
  gstatus=`git status --porcelain`

  if [ ${#gstatus} -ne 0 ]
  then
      git add --all
      git commit -m "$gstatus"

      git pull --rebase
      git push
  else
    git pull
  fi
*/

import simpleGit, { SimpleGit, SimpleGitOptions } from "simple-git";

const options: Partial<SimpleGitOptions> = {
  baseDir: process.cwd(),
  binary: "git",
  maxConcurrentProcesses: 1,
};

async function mainLoop() {
  const git: SimpleGit = simpleGit(options);
  var statusResult = await git.status();
  console.log(statusResult);

  console.log(statusResult.files);
}

mainLoop();

// I need some test data!
