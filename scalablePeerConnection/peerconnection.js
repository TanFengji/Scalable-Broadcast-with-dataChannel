
function PeerConnection(local, peer, socket, localVideo){
	var p2pConnection;
	var indicator;
	var dataChannel;
	this.user = local;
	this.remote = peer;
	this.socket = socket;
	this.localVideo = localVideo;
	this.configuration = {
			"iceServers": [{ 
			"credential": "aed9a3ac-31f9-11e6-9ed0-6cf1b3d414d3",
			"url": "turn:turn02.uswest.xirsys.com:80?transport=udp",
			"username": "aed9a316-31f9-11e6-8b25-3ad12cf92d7c"
			}]
	};
	console.log("local video is " + localVideo);
}

//Visitor setup the p2p connection with a peer
PeerConnection.prototype.visitorSetupPeerConnection = function(peer,/* streamCallback,*/ cb) {
	var self = this;
	// Setup stream listening
	console.log("listen to stream");
	this.p2pConnection.onaddstream = function (e) {
		self.localVideo.src = window.URL.createObjectURL(e.stream);
		streamCallback(e.stream);
		console.log("received a stream");
		console.log(e.stream);
	};



//	Setup ice handling
	console.log("start ice handling");
	this.p2pConnection.onicecandidate = function (event) {
		if (event.candidate) {
			console.log(event.candidate);
			self.socket.emit("candidate", {
				type: "candidate",
				local: self.user,
				remote: peer,
				candidate: event.candidate
			});
		}
	};
	cb();
}

//Host setup the p2p connection with a peer
PeerConnection.prototype.hostSetupPeerConnection = function(peer, stream, cb) {
	var self = this;
	// Add stream
	//this.p2pConnection.addStream(stream);

	// Setup ice handling
	this.p2pConnection.onicecandidate = function (event) {
		if (event.candidate) {
			console.log("send an ice candidate");
			console.log(event.candidate);
			self.socket.emit("candidate", {
				type: "candidate",
				local: self.user,
				remote: peer,
				candidate: event.candidate
			});
		}
	};
	cb();
}

PeerConnection.prototype.addVideo = function(stream) {
	// Add stream
	var self = this;
	this.p2pConnection.addStream(stream);
	this.makeOffer( function(sdpOffer){
		console.log(sdpOffer);
		sdpOffer = JSON.stringify(sdpOffer);
		self.dataChannel.send(sdpOffer);
	});
}

PeerConnection.prototype.onAddVideo = function(sdpOffer) {
	// Add stream
	var self = this;
	this.p2pConnection.onaddstream = function (e) {
		self.localVideo.src = window.URL.createObjectURL(e.stream);
		console.log("received a stream");
		console.log(e.stream);
	};
	this.receiveOffer(sdpOffer, function(sdpAnswer){
		sdpAnswer = JSON.stringify(sdpAnswer);
		self.dataChannel.send(sdpAnswer);
	});
}

//initialise p2pconnection at the start of a peer connection 
PeerConnection.prototype.startConnection = function(cb){
	this.p2pConnection = new RTCPeerConnection(this.configuration);
	cb();
}

PeerConnection.prototype.openDataChannel = function(cb){
	var self = this;
	var dataChannelOptions = {
			ordered: true,
			reliable: true,
			negotiated: true,
			id: "myChannel"
	};

	self.dataChannel = this.p2pConnection.createDataChannel("label", dataChannelOptions);
	console.log("new Data channel");

	self.dataChannel.onerror = function (error) {
		console.log("Data Channel Error:", error);
	};

	self.dataChannel.onmessage = function (msg) {
		console.log("Got Data Channel Message:");
		console.log(msg.data);

		if (isJson(msg.data)){
			message = JSON.parse(msg.data);
			console.log(message);
			if (message.type === "offer"){
				self.onAddVideo(message);
			}
			else if (message.type === "answer"){
				self.receiveAnswer(message);
			}
		} else {
			message = msg.data + "<br />"
			document.getElementById("info").innerHTML += message;
		}
	};

	self.dataChannel.onopen = function () {
		self.dataChannel.send("connected.");
	};

	self.dataChannel.onclose = function () {
		console.log("The Data Channel is Closed");
	};
	cb();
}


//make an sdp offer
PeerConnection.prototype.makeOffer = function(cb)	{
	var self = this;
	this.p2pConnection.createOffer(function (sdpOffer) {
		sdpOffer.sdp = sdpOffer.sdp.replace(/a=sendrecv/g,"a=sendonly");
		self.p2pConnection.setLocalDescription(sdpOffer);
		cb(sdpOffer);
	}, function(error){
		console.log(error);
	});
}

//receive an sdp offer and create an sdp answer
PeerConnection.prototype.receiveOffer = function(sdpOffer, cb){
	var self = this;
	this.p2pConnection.setRemoteDescription(sdpOffer, function(){
		self.p2pConnection.createAnswer(function (answer) {
			answer.sdp = answer.sdp.replace(/a=sendrecv/g,"a=recvonly");
			self.p2pConnection.setLocalDescription(answer);
			console.log(self.p2pConnection.localDescription);
			console.log(self.p2pConnection.remoteDescription);
			cb(answer);
		},function(error){
			console.log(error);
		});
	}, function(){});
}

//receive an spd answer
PeerConnection.prototype.receiveAnswer = function(sdpAnswer){
	this.p2pConnection.setRemoteDescription(sdpAnswer,function(){}, function(){});
	console.log(this.p2pConnection.localDescription);
	console.log(this.p2pConnection.remoteDescription);
}

//add ice candidate when receive one
PeerConnection.prototype.addCandidate = function(iceCandidate) {
	this.p2pConnection.addIceCandidate(new RTCIceCandidate(iceCandidate.candidate));
}

function isJson(str) {
	try {
		JSON.parse(str);
	} catch (e) {
		return false;
	}
	return true;
}

module.exports = PeerConnection;