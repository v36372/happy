var app = function(){
	$(".button-back").click(function(){
		$(".player").toggleClass("playlist");
	});

	$("#hamburger-menu").click(function(){
		$(".player").toggleClass("playlist");
	});
	$("#button-pause").hide();

	var index = 0;
	var playlist = document.getElementsByClassName("track").length;
	var currentPlaying;
	var ytPlayer;
	var scPlayer;
	var songEnded = true;
	var playerState = "stop";
	var isInit = -2;

	var getSong = function(i) {
		if (i >= playlist) {
			return;
		}

		var el = document.getElementById(i)
		if (el) return {
			provider: el.getAttribute("provider"),
			link: el.getAttribute("link"),
			name: el.getAttribute("name"),
			thumbnail: el.getAttribute("thumbnail"),
		}
	}

	var onStateChange = function(e){
		if (e.data === 0) {
			onSongEnded()
		}
	}

	var playNextSong = function() {
		if (index + 1 < playlist) {
			index = index+1;
			playSong(index, false);
			$("#button-play").hide()
			$("#button-pause").show()
		}
	}

	var playLastSong = function() {
		if (index - 1 >= 0 && playlist > 0) {
			index = index -1;
			playSong(index, false);
			$("#button-play").hide()
			$("#button-pause").show()
		}
	}

	var onSongEnded = function(e) {
		console.log("on song ended");
		if (isInit < 0) isInit++;
		console.l
		if (!songEnded) {
			songEnded = true;
			playerState = "stop";
			$("#button-pause").hide()
			$("#button-play").show()
			playNextSong();
		}
	}

	var playSong = function(i, resume) {
		if (isInit < 0) {
			return;
		}
		var song = getSong(i);
		console.log(song)
		if (!song) {
			return;
		}

		$("#button-play").hide()
		$("#button-pause").show()
		playerState = "playing"
		if (song.provider == "youtube") {
			if (resume) {
				ytPlayer && ytPlayer.playVideo();
				return;
			}
			scPlayer && scPlayer.pause && scPlayer.pause();
			songEnded = false;
			ytPlayer && ytPlayer.loadVideoById && ytPlayer.loadVideoById(song.link);
		} else if (song.provider == "soundcloud") {
			if (resume) {
				scPlayer && scPlayer.play();
				return;
			}
			ytPlayer && ytPlayer.pauseVideo && ytPlayer.pauseVideo();
			scPlayer && scPlayer.load(song.link, {
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

		// change player info
		currentPlaying = song;
		index = Number(i);
		$("#song-name").text(song.name);
		if (song.thumbnail) {
			$('#track-background').css('background','url("' + song.thumbnail +'") center center/cover')
			$('.art').attr('src', song.thumbnail);
		} else {
			$('#track-background').css('background','url() center center/cover')
			$('.art').attr('src', '');
		}
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
			setTimeout(waitForYoutube, 50);
		} else {
			console.log("onyoutube ready");
			ytPlayer = new YT.Player("youtube-player", {
				playerVars:{
					cc_load_policy: '0',
					iv_load_policy: '3',
				},
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
			setTimeout(waitForSoundCloud, 50);
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

	$("#button-next").click(function(){
		playNextSong();
	});

	$("#button-prev").click(function(){
		playLastSong();
	});

	$("#button-pause").click(function(){
		pauseSong(index);
		$("#button-pause").hide()
		$("#button-play").show()
	});

	$("#button-add").click(function(){
		
	});

	$("#button-play").click(function(){
		console.log("called")
		if (playerState == "stop") {
			playSong(index, false);
		} else playSong(index, true);
	});

	$(".track").click(function(event){
		var el = event.currentTarget;
		playSong(el.id, false)
	})
};

function onYouTubeIframeAPIReady() {
	app();
	app = function(){};
}
