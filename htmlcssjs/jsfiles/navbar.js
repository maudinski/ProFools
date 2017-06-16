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
