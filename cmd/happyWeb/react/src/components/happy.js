function extend(destination, source) {
	destination.prototype = Object.create(source.prototype);
	return source.prototype;
}

var ytPlayer;
var scPlayer;

function MyHappyPlayer(playList) {
	this.init(playList)
	this.isInit = -2;
	var self = this;
	var waitForYoutube = function(){
		console.log("wait for youtyube");
		if (typeof YT === "undefined" || typeof YT.Player === "undefined"){
			setTimeout(waitForYoutube, 200);
		} else {
			console.log("onyoutube ready");
			ytPlayer = new YT.Player("youtube-player", {
				playerVars:{
					cc_load_policy: '0',
					iv_load_policy: '3',
				},
				events: {
					'onStateChange': function(e){
						if (e.data === 0) {
							console.log("youtube end")
							self.onSongEnded();
						}
					},
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
			setTimeout(waitForSoundCloud, 200);
		} else {
			console.log("sc ready");
			scPlayer = SC.Widget(document.getElementById('soundcloud-player'));
			scPlayer.bind("finish", function(){
				console.log("sc end")
				self.onSongEnded();
			});
			scPlayer.bind("ready", function(){
				console.log("readyasdhasd");
				scPlayer.setVolume(0);
				scPlayer.play();
				scPlayer.unbind("ready");
			});
		}
	}();
}

MyHappyPlayer.prototype = Object.create(HappyPlayer.prototype);

MyHappyPlayer.prototype.playSong = function(event) {
	if (this.isInit < 0) {
		return;
	}
	this.__proto__.__proto__.playSong.call(this, event)

	console.log("child run")
	var song = this.playList[this.currentIndex].song
	if (song.provider == "youtube") {
		scPlayer && scPlayer.pause && scPlayer.pause();
		ytPlayer && ytPlayer.loadVideoById && ytPlayer.loadVideoById(song.link);
	} else if (song.provider == "soundcloud") {
		ytPlayer && ytPlayer.pauseVideo && ytPlayer.stopVideo();
		scPlayer && scPlayer.load(song.link, {
			visual: false,
			show_artwork: false,
			auto_play: true,
			callback: function(){
				console.log("soundcloud player ready");
				scPlayer.play();
			}
		})
	}
}

MyHappyPlayer.prototype.resumeSong = function() {
	var self = this;
	if (self.playList[self.currentIndex].song.provider == "youtube") {
		ytPlayer && ytPlayer.playVideo();
	} else if (self.playList[self.currentIndex].song.provider  == "soundcloud") {
		scPlayer && scPlayer.play();
	}
	this.__proto__.__proto__.resumeSong.call(this)
}

MyHappyPlayer.prototype.pauseSong = function() {
	var self = this;
	if (self.playList[self.currentIndex].song.provider == "youtube") {
		ytPlayer && ytPlayer.pauseVideo();
	} else if (self.playList[self.currentIndex].song.provider  == "soundcloud") {
		scPlayer && scPlayer.pause();
	}
	this.__proto__.__proto__.pauseSong.call(this)
}

MyHappyPlayer.prototype.seekToBegin = function() {
	var self = this;
	if (self.playList[self.currentIndex].song.provider == "youtube") {
		ytPlayer && ytPlayer.seekTo(0, true);
	} else if (self.playList[self.currentIndex].song.provider  == "soundcloud") {
		scPlayer && scPlayer.seekTo(0) && scPlayer.play();
	}
}

MyHappyPlayer.prototype.onSongEnded = function() {
	console.log("onsong ended", this)
	if (this.isInit < 0){
		this.isInit++;
		return;
	} 

	this.__proto__.__proto__.onSongEnded.call(this)
}

MyHappyPlayer.prototype.updateUI = function() {
	console.log("update ui", this)
	if (this.playerState == "playing") {
		$("#button-play").hide()
		$("#button-pause").show()
	} else {
		$("#button-play").show()
		$("#button-pause").hide()
	}

	if(this.shuffle) {
		$("#button-shuffle").parent().addClass("active");
	} else {
		$("#button-shuffle").parent().removeClass("active");
	}

	if(this.repeatOne) {
		$("#button-repeat").parent().addClass("active");
	} else {
		$("#button-repeat").parent().removeClass("active");
	}

	var song = this.playList[this.currentIndex].song;

	$("#song-name").text(song.name);
	if (song.thumbnail) {
		$('#track-background').css('background','url("' + song.thumbnail +'") center center/cover')
		$('.art').attr('src', song.thumbnail);
	} else {
		$('#track-background').css('background','url() center center/cover')
		$('.art').attr('src', '');
	}
}

export default MyHappyPlayer
