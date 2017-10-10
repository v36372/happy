var app = function(){
	$(".back, .header a").click(function(){
		$(".player").toggleClass("playlist");
	});

	$("#hamburger-menu").click(function(){
		$(".player").toggleClass("playlist");
	});

	var index = 0;
	var songs = [
		{
			provider: "youtube",
			id: "B7bqAsxee4I"
		},
		{
			provider: "youtube",
			id: "B7bqAsxee4I"
		},
		{
			provider: "youtube",
			id: "B7bqAsxee4I"
		},
		{
			provider: "youtube",
			id: "7dp4GLm7sgo"
		},
	];
	var ytPlayer;
	var scPlayer;
	var songEnded = false;

	var onStateChange = function(e){
		if (e.data === 0) {
			onSongEnded()
		}
	}

	var playNextSong = function() {
		if (index < songs.length) {
			index = index+1;
			playSong(index);
		}
	}

	var onSongEnded = function() {
		console.log("on song ended");
		if (!songEnded) {
			songEnded = true;
			playNextSong();
		}
	}

	var playSong = function(index) {
		console.log(index);
		if (songs[index].provider == "youtube") {
			songEnded = false;
			ytPlayer.loadVideoById(songs[index].id);
		} else if (songs[index].provider == "soundcloud") {
			scPlayer.load(songs[index].id, {
				visual: false,
				show_artwork: false,
				auto_play: true,
				callback: function(){
					console.log("soundcloud player ready");
					scPlayer.play();
					songEnded = false;
				}
			})
		}
	}

	var waitForYoutube = function(){
		console.log("wait for youtyube");
		if (typeof YT === "undefined" || typeof YT.Player === "undefined"){
			setTimeout(waitForYoutube, 250);
		} else {
			console.log("onyoutube ready");
			ytPlayer = new YT.Player("youtube-player", {
				events: {
					'onStateChange': onStateChange,
					'onReady': function() {
						console.log("readgy");
					},
				}
			})
		}
	}();

	var waitForSoundCloud = function(){
		if (typeof SC === "undefined"){
			setTimeout(waitForSoundCloud, 250);
		} else {
			console.log("sc ready");
			scPlayer = SC.Widget("soundcloud-player", {})
			scPlayer.bind("finish", onSongEnded);
		}
	}();

	$("#button-play").click(function(){
		console.log("called")
		playSong(index);
	});
};

function onYouTubeIframeAPIReady() {
	app();
	app = function(){};
}
