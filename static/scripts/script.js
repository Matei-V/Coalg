var menu_tog = false;
function show_menu() {
  var menu_width = document.getElementById("menu").offsetWidth;
  console.log(menu_width);
  if (menu_tog === false) {
    document.getElementById("menu").style.left = "0px";
    document.getElementById("slider").style.left = menu_width.toString() + "px";
    document.getElementById("content").style.left = (menu_width + 30).toString() + "px";
  } else {
    document.getElementById("menu").style.left = "-" + menu_width.toString() + "px";
    document.getElementById("slider").style.left = "0px";
    document.getElementById("content").style.left = "30px";
  }
  menu_tog = !menu_tog;
}

function slider_pos() {
  var menu_width = document.getElementById("menu").offsetWidth;
  document.getElementById("slider").style.left = "220px";
  console.log("aAAAAA");
}


