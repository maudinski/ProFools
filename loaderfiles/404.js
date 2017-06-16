function exhibitColor(event){
	var pathname = window.location.pathname;
	var firstTrunc = pathname.substring(1, pathname.length)
	var secondTrunc = firstTrunc.substring(firstTrunc.indexOf("/")+1, firstTrunc.length)
	var exhibit = secondTrunc.substring(0, secondTrunc.indexOf("/"))
	var navlist = document.getElementsByClassName("highlight_li_navbar");
	var found = false
	for (i = 0; i < navlist.length; i++){
		if(navlist[i].id == exhibit){
			navlist[i].style.backgroundColor = "#d81515";
			found = true
			break;
		}
	}
	if (!found) {
		navlist[0].style.backgroundColor = "#d81515";
	}
}
/*function is for the top bar to be good looking at all window widths*/
function headerDynamics(event){
    var width = window.innerWidth;
    var sign_links = document.getElementById("sign_links_div");
    var search_bar = document.getElementById('search_div');
    var phrase = document.getElementById("phrase_div");
    /*hard coded as fuck i dont give a shit*/
    if(width < 762){
        search_bar.style.left = "360px";
        search_bar.style.position = "relative";
        search_bar.style.top = "-8px";
        sign_links.style.left = "655px";
        sign_links.style.position = "relative";
        sign_links.style.top = "-30px";
        phrase.style.position = "relative";
    } else {
        search_bar.style.left = "initial";
        search_bar.style.position = "absolute";
        search_bar.style.top = "10px";
        sign_links.style.left = "initial";
        sign_links.style.position = "absolute";
        sign_links.style.top = "10px";
        phrase.style.position = "absolute";
    }
}
/*spaces out the sub li's to look nice*/
/*function setNavBar(id){
    var ul = document.getElementById(id);
    var items = ul.getElementsByTagName("li");
    var spacer = 10;
    var left_dist = 0;
    for(var i = 0; i < items.length; ++i){
        items[i].style.left = String(left_dist + spacer).concat("px");
        left_dist += items[i].offsetWidth + spacer;
    }
}*/
/*stuff that till happen on load*/
function loadStuff(event){
    /*get correct header placements*/
    headerDynamics(event);
    /*colors appropriate "tab" exhibit*/
    exhibitColor(event);
}
window.onload = loadStuff;
window.addEventListener('resize', headerDynamics);
