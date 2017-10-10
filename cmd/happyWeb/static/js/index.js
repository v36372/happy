var app = function(){
	$(".back, .header a").click(function(){
		$(".player").toggleClass("playlist");
	});

	$("#hamburger-menu").click(function(){
		$(".player").toggleClass("playlist");
	});
	$("#button-pause").hide();

	var index = 0;
	var playlistLength = document.getElementsByClassName("track").length;
	var currentPlaying;
	var ytPlayer;
	var scPlayer;
	var songEnded = true;
	var playerState = "stop";

	var getSong = function(index) {
		if (index >= playlistLength) {
			return;
		}

		var el = document.getElementById("track-" + index);
		if (el) return {
			provider: el.getAttribute("provider"),
			link: el.getAttribute("link"),
			name: el.getAttribute("name"),
		}
	}

	var onStateChange = function(e){
		if (e.data === 0) {
			onSongEnded()
		}
	}

	var playNextSong = function() {
		if (index + 1 < playlistLength) {
			index = index+1;
			playSong(index, false);
			$("#button-play").hide()
			$("#button-pause").show()
		}
	}

	var onSongEnded = function() {
		console.log("on song ended");
		if (!songEnded) {
			songEnded = true;
			$("#button-pause").hide()
			$("#button-play").show()
			playNextSong();
		}
	}

	var playSong = function(index, resume) {
		console.log(index);
		var song = getSong(index);
		if (!song) {
			return;
		}

		if (song.provider == "youtube") {
			if (resume) {
				ytPlayer && ytPlayer.playVideo();
				return;
			}
			songEnded = false;
			ytPlayer.loadVideoById(song.link);
		} else if (song.provider == "soundcloud") {
			if (resume) {
				scPlayer && scPlayer.play();
				return;
			}
			scPlayer.load(song.link, {
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
		currentPlaying = song;
		$("#song-name").text(song.name);
		$("#song-artist").text(song.name);
	}

	var pauseSong = function(index) {
		if (currentPlaying.provider == "youtube") {
			ytPlayer && ytPlayer.pauseVideo() && (playerState = "pause");
		} else if (currentPlaying.provider == "soundcloud") {
			scPlayer && scPlayer.pause() && (playerState = "pause");
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
						ytPlayer.setPlaybackQuality("small");
						ytPlayer.playVideo();
					},
				}
			});
		}
	}();

	var waitForSoundCloud = function(){
		if (typeof SC === "undefined"){
			setTimeout(waitForSoundCloud, 250);
		} else {
			console.log("sc ready");
			scPlayer = SC.Widget(document.getElementById('soundcloud-player'));
			scPlayer.bind("finish", onSongEnded);
			scPlayer.bind("ready", function(){
				console.log("readyasdhasd");
				scPlayer.setVolume(0);
				scPlayer.play();
				scPlayer.unbind("ready");
			});
		}
	}();

	$("#button-pause").click(function(){
		pauseSong(index);
		$("#button-pause").hide()
		$("#button-play").show()
	});

	$("#button-play").click(function(){
		console.log("called")
		if (playerState == "stop") {
			playSong(index, false);
		} else playSong(index, true);
		$("#button-play").hide()
		$("#button-pause").show()
	});
};

function onYouTubeIframeAPIReady() {
	app();
	app = function(){};
}
