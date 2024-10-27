var menu_tog = false;

function show_menu() {
  var menu_width = document.getElementById("menu").offsetWidth.toString();
  console.log(menu_width);
  if (menu_tog === false) {
    document.getElementById("menu").style.left = "0px";
    document.getElementById("slider").style.left = menu_width + "px";
    document.getElementById("main").style.left = menu_width + "px";
  } else {
    document.getElementById("menu").style.left = "-" + menu_width + "px";
    document.getElementById("slider").style.left = "0px";
    document.getElementById("main").style.left = "0px";
  }
  menu_tog = !menu_tog;
}
