var PeerConnection = require('./peerconnection.js');
var Indicator = require('./indicator.js');

function AllConnection(){
	var local;
	var stream;
	var socket;
	var configuration;
	var localVideo;
	this.connection = {};
	this.indicator = new Indicator();
	this.ms = new MediaSource();
	console.log(this.ms);
}

//initialise the setup of AllConnection
AllConnection.prototype.init = function(user, socket, config){
	this.local = user;
	this.socket = socket;
	this.configuration = config;
}

//initialise the setup of own camera
AllConnection.prototype.initCamera = function(cb){
	var self = this;
	this.localVideo = document.getElementById("localVideo");
	this.localVideo.src = URL.createObjectURL(self.ms);
	this.localVideo.autoplay = true;

	self.ms.addEventListener('sourceopen', function(){
		// this.readyState === 'open'. Add source buffer that expects webm chunks.
		window.sourceBuffer = self.ms.addSourceBuffer('video/webm; codecs="vorbis,vp8"');

//		To Do: Problem: create 2 video when 2 users enter simultaneously
		if (self.indicator.hasUserMedia()) {
			console.log(sourceBuffer);
			navigator.getUserMedia({ video: true, audio: true }, function(stream){
				self.startRecording(stream);
			}, function (error) {
				console.log(error);
			});
			console.log(self);
			cb();
		} else {
			alert("Sorry, your browser does not support WebRTC.");
		}
	});
}

//initialise a connection with peers
AllConnection.prototype.initConnection = function(peer){	
	var self = this;
	self.localVideo = document.getElementById("localVideo");
	self.localVideo.autoplay = true;
	self.connection[peer] = new PeerConnection(self.local, peer, self.socket, self.localVideo, self.configuration);
	self.connection[peer].startConnection(function(){
		self.connection[peer].openDataChannel(function(){
			self.connection[peer].hostSetupPeerConnection(peer, self.stream, function(){
				self.connection[peer].makeOffer( function(offer){
					self.socket.emit("SDPOffer", {
						type: "SDPOffer",
						local: self.local,
						remote: peer,
						offer: offer
					});
				});
			});
		});
	});
}

//when receive an spd offer
AllConnection.prototype.onOffer = function(sdpOffer, cb){
	var self = this;
	self.localVideo = document.getElementById("localVideo");
	self.localVideo.autoplay = true;
	peer = sdpOffer.remote;
	self.connection[peer] = new PeerConnection(self.local, peer, self.socket, self.localVideo, self.configuration);
	self.connection[peer].startConnection(function(){
		self.connection[peer].openDataChannel(function(){
			self.connection[peer].visitorSetupPeerConnection(peer, /*function(stream){
				self.stream = stream;
				cb();
			}, */function(){
				self.connection[sdpOffer.remote].receiveOffer(sdpOffer.offer, function(sdpAnswer){
					self.socket.emit("SDPAnswer", {
						type: "SDPAnswer",
						local: self.local,
						remote: sdpOffer.remote,
						answer: sdpAnswer
					});
				});
			});
		});
	});
}

//when receive an spd answer
AllConnection.prototype.onAnswer = function(sdpAnswer, cb){
	this.connection[sdpAnswer.remote].receiveAnswer(sdpAnswer.answer);
}

//when receive an ice candidate
AllConnection.prototype.onCandidate = function(iceCandidate){
	this.connection[iceCandidate.remote].addCandidate(iceCandidate);
}

AllConnection.prototype.deleteConnection = function(peer){
	self.connection[peer] = null;
}

//set the ICE server 
AllConnection.prototype.setIceServer = function(iceServers){
	this.iceServers = iceServers;
}

AllConnection.prototype.setLocalStream = function(streamStatus){
	this.stream = this.connection[streamStatus.host].stream;
}

/*
	var chunks = [];
	console.log('Starting...');
	mediaRecorder = new MediaRecorder(stream);
	setTimeout(function(){
		mediaRecorder.stop();
	}, 5000);
	mediaRecorder.start();

	mediaRecorder.ondataavailable = function(e) {
		chunks.push(e.data);
		console.log(e.data);
		console.log(e);
	};

	mediaRecorder.onerror = function(e){
		log('Error: ' + e);
		console.log('Error: ', e);
	};


	mediaRecorder.onstart = function(){
		console.log('Started, state = ' + mediaRecorder.state);
	};

	mediaRecorder.onstop = function(){
		console.log('Stopped, state = ' + mediaRecorder.state);

		var blob = new Blob(chunks, {type: "video/webm"});
		chunks = [];

		var videoURL = window.URL.createObjectURL(blob);
		var downloadLink = document.getElementById("download");
		console.log(videoURL);
		self.localVideo.src = videoURL;
		downloadLink.innerHTML = 'Download video file';

		var rand = Math.floor((Math.random() * 10000000));
		var name = "video_"+rand+".webm" ;
		downloadLink.setAttribute( "href", videoURL);
		downloadLink.setAttribute( "download", name);
		downloadLink.setAttribute( "name", name);

	};

	mediaRecorder.onwarning = function(e){
		console.log('Warning: ' + e);
	};*/


AllConnection.prototype.startRecording = function(stream) {
	var self = this;
	//self.localVideo.play();
	console.log(this);
	console.log("here");

	console.log(sourceBuffer);
//	sourceBuffer.mode = "sequence";
	sourceBuffer.timestampOffset = 2.5;
	console.log(sourceBuffer);
	var chunks = [];
	console.log('Starting...');
	var mediaRecorder = new MediaRecorder(stream);
	setTimeout(function(){
		mediaRecorder.stop();
	}, 7000);

	mediaRecorder.start(500);

	mediaRecorder.ondataavailable = function (e) {
		var reader = new FileReader();
		reader.addEventListener("loadend", function () {
			var arr = new Uint8Array(reader.result);
			try{
				sourceBuffer.appendBuffer(arr);
				console.log(sourceBuffer);
				console.log("correct");
			}catch(e){
				console.log(e);
			}
		});
		reader.readAsArrayBuffer(e.data);
	};

	mediaRecorder.onerror = function(e){
		console.log('Error: ', e);
	};


	mediaRecorder.onstart = function(){
		console.log('Started, state = ' + mediaRecorder.state);
	};

	mediaRecorder.onstop = function(){
		console.log('Stopped, state = ' + mediaRecorder.state);
		console.log(self.ms.sourceBuffers);
		/*
		var blob = new Blob(chunks, {type: "video/webm"});
		chunks = [];

		var videoURL = window.URL.createObjectURL(blob);
		var downloadLink = document.getElementById("download");
		console.log(videoURL);
		self.localVideo.src = videoURL;
		downloadLink.innerHTML = 'Download video file';

		var rand = Math.floor((Math.random() * 10000000));
		var name = "video_"+rand+".webm" ;
		downloadLink.setAttribute( "href", videoURL);
		downloadLink.setAttribute( "download", name);
		downloadLink.setAttribute( "name", name);
		 */
	};

	mediaRecorder.onwarning = function(e){
		console.log('Warning: ' + e);
	};

}

module.exports = AllConnection;
