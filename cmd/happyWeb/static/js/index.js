$("main").addClass("pre-enter").removeClass("with-hover");
asd
setTimeout(function(){
	$("main").addClass("on-enter");
}, 500);
setTimeout(function(){
	$("main").removeClass("pre-enter on-enter");
	setTimeout(function(){
		$("main").addClass("with-hover");
	}, 50);
}, 2000);

$(".flip, .front a").click(function(){
	console.log("vo ham` click roi ne`")
	$(".player").toggleClass("playlist");
});

$(".bottom a").not(".flip").click(function(){
	$(this).toggleClass("active");
});
