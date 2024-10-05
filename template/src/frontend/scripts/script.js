var menu_tog = false;

function show_menu() {
  if (menu_tog === false) {
    document.getElementById("menu").style.left = "0px";
    document.getElementById("slider").style.left = "39px";
  } else {
    document.getElementById("menu").style.left = "-40px";
    document.getElementById("slider").style.left = "0px";
  }
  menu_tog = !menu_tog;
}
