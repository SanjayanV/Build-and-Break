import fs from "fs";
var names = ["Sanjayan V", "Admin"];
document.getElementById("greet").textContent = "Hello " + names[0];
var file = document.getElementById("fileupload");
spath = "./uploads"+""
fs.copyFile(file,"./uploads")