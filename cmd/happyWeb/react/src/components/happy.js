function extend(destination, source) {
	destination.prototype = Object.create(source.prototype);
	return source.prototype;
}

function MyHappyPlayer(playList) {
	this.init(playList)
	this.isInit = -2;
}

var parent = extend(MyHappyPlayer, HappyPlayer)

MyHappyPlayer.prototype.playSong = function(event) {
	parent.playSong(event)

	console.log("child run")
	console.log(this)
	if (this.playList[event.data].provider == "youtube") {
		scPlayer && this.scPlayer.pause && scPlayer.pause();
		ytPlayer && ytPlayer.loadVideoById && ytPlayer.loadVideoById(song.link);
	} else if (this.playList[event.data].provider == "soundcloud") {
		ytPlayer && ytPlayer.pauseVideo && ytPlayer.pauseVideo();
		scPlayer && scPlayer.load(song.link, {
			visual: false,
			show_artwork: false,
			auto_play: true,
			callback: function(){
				console.log("soundcloud player ready");
				scPlayer.play();
			}
		})
		scPlayer.bind("finish", parent.onSongEnded);
	}
}

MyHappyPlayer.prototype.onSongEnded = function() {
	if (this.isInit < 0){
		isInit++;
		return;
	} 

	parent.onSongEnded()
}

MyHappyPlayer.prototype.updateUI = function() {
	if (this.playerState == "playing") {
		$("#button-play").hide()
		$("#button-pause").show()
	} else {
		$("#button-play").show()
		$("#button-pause").hide()
	}

	var song = this.playList[this.currentIndex];

	$("#song-name").text(song.name);
	if (song.thumbnail) {
		$('#track-background').css('background','url("' + song.thumbnail +'") center center/cover')
		$('.art').attr('src', song.thumbnail);
	} else {
		$('#track-background').css('background','url() center center/cover')
		$('.art').attr('src', '');
	}
}

$(".button-back").click(function(){
	$(".player").toggleClass("playlist");
});

$("#hamburger-menu").click(function(){
	$(".player").toggleClass("playlist");
});
$("#button-pause").hide();

var ytPlayer;
var scPlayer;

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
		scPlayer.bind("ready", function(){
			console.log("readyasdhasd");
			scPlayer.setVolume(0);
			scPlayer.play();
			scPlayer.unbind("ready");
		});
	}
}();


export default MyHappyPlayer
