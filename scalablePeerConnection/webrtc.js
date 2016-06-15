var AllConnection = require('./allconnection.js');
var io = require('socket.io-client');

function WebRTC(server){
	var self = this;
	var user;
	var peer;
	var peerList;
	this.latencyList = [];
	this.peerNo = 0;
	this.connectionBuilt = 0;
	//this.latencyList = {};
	this.latencyListSize = 0;
	this.allConnection = new AllConnection();;
	this.socket = io(server);

	// when a datachannel setup ready
	self.socket.on("dataChannelStatus", function(dataChannelStatusData){
		if (dataChannelStatusData.status === "success"){
			self.connectionBuilt++;
			if (self.connectionBuilt === self.peerNo){
				self.sendTimeStamp();
			}
		}
	});

	// when user and a peer finish transfering their time stamp
	self.socket.on("timeStamp", function(timeStampData){
		var timeStamp = {};
		timeStamp.peer = timeStampData.peer;
		timeStamp.latency = timeStampData.receiveTime - timeStampData.sendTime;
		self.latencyList.push(timeStamp);
							
 /* self.latencyList[timeStampData.peer] = {};
  * self.latencyList[timeStampData.peer].peer = timeStampData.peer; 
	*	self.latencyList[timeStampData.peer].latency = timeStampData.receiveTime - timeStampData.sendTime;*/
		self.latencyListSize++ ; 
		if (self.latencyListSize === self.peerNo){
			for (var a in self.latencyList){
				console.log(a);
				console.log("Latency: " + self.latencyList[a].latency);
			}
			
			self.socket.emit("newUser", {
				type: "newUser",
				user: self.user,
				latency: self.latencyList
			});
		}
		

	});

	//responde to different socket received from server

	self.socket.on("feedback", function(feedback) {
		document.getElementById("feedback").value = feedback;
	});

	//receive a sdp offer
	self.socket.on("SDPOffer", function(sdpOffer) {
		self.allConnection.onOffer(sdpOffer, function(){
			if (self.peer){
				self.allConnection.initConnection(self.peer);
			}
		});
	});

	//receive a sdp answer
	self.socket.on("SDPAnswer", function(sdpAnswer) {
		self.allConnection.onAnswer(sdpAnswer);
	});

	//receive an ice candidate
	self.socket.on("candidate", function(iceCandidate) {
		console.log("receive an ice candidate");
		self.allConnection.onCandidate(iceCandidate);
	});

	// when a user in the room disconnnected
	self.socket.on("disconnectedUser", function(disConnectedUserName) {
		console.log("user " + disConnectedUserName + " is disconnected");
		self.onUserDisconnect(disConnectedUserName);
		self.socket.emit("message", {
			type: "message",
			action: "leave",
			user: self.user,
			content: ""
		});
	});

	// initialize 1 way peer connection or start host's camera
	self.socket.on("initConnection", function(peer){
		if (self.user === peer){
			console.log("init camera");
			self.allConnection.initCamera(function(){
				/* setup camera before build connection	
				 * self.onHostSetup();
				 */
			});
		}else {		
			self.allConnection.initConnection(peer);
			self.peer = peer;
		}
	});

	// delete peer connection when peer left
	self.socket.on("deleteConnection", function(peer){
		self.allConnection.deleteConnection(peer);
		self.peer = null;
	});

	self.socket.on("message", function(messageData){
		console.log("received message");
		self.onMessage(messageData);
	});
}


//find more details of following api in readme
WebRTC.prototype.login = function(userName, successCallback, failCallback) {
	var self = this;
	this.socket.emit("login", userName);
	this.socket.on("login", function(loginResponse){
		if (loginResponse.status === "success") {
			self.user = loginResponse.userName;
			self.allConnection.init(loginResponse.userName, self.socket, loginResponse.config);
			successCallback();
		} else if (loginResponse.status === "fail") {
			failCallback();
		}
	});
}

WebRTC.prototype.createRoom = function(roomId, successCallback, failCallback){
	var self = this;
	this.socket.emit("createRoom", roomId);
	this.socket.on("createRoom", function(createRoomResponse){
		if (createRoomResponse.status === "success") {
			successCallback();
		} else if (createRoomResponse.status === "fail") {
			failCallback();
		}
	});
}

WebRTC.prototype.joinRoom = function(roomId, successCallback, failCallback) {
	var self = this;
	this.socket.emit("joinRoom", roomId);
	this.socket.on("joinRoom", function(joinRoomResponse){
		if (joinRoomResponse.status === "success") {
			self.peerList = joinRoomResponse.userList;
			self.socket.emit("message", {
				type: "message",
				action: "join",
				user: self.user,
				content: ""
			});

			for (var peer in self.peerList){
				if (peer){
					self.peerNo++;
				}
			}

			console.log(self.peerNo);
			for (var peer in self.peerList){
				self.allConnection.initConnection(peer);
			}
			console.log("finish");
			successCallback();
		} else if (joinRoomResponse.status === "fail") {
			failCallback();
		}
	});
}

WebRTC.prototype.onUserDisconnect = function(userDisconnected){
}

/*WebRTC.prototype.sendChatMessage = function(chatMessage){
	var self = this;
	self.socket.emit("message", {
		type: "message",
		action: "chat",
		user: self.user,
		content: chatMessage
	})
}*/

WebRTC.prototype.sendChatMessage = function(chatMessage){
	var self = this;
	console.log("peer is " + self.peer);
	console.log(self.peer);
	console.log(	self.allConnection.connection[self.peer]);
	console.log(	self.allConnection.connection[self.peer].dataChannel);
	self.allConnection.connection[self.peer].dataChannel.send(chatMessage);
}

WebRTC.prototype.onMessage = function(messageData){
}

WebRTC.prototype.addVideo = function(){
	var self = this;
	console.log(self.allConnection.connection[self.peer]);
	console.log(self.allConnection.stream);
	this.allConnection.connection[self.peer].addVideo(self.allConnection.stream);
}

WebRTC.prototype.setIceServer = function(iceServers){
	this.allConnection.setIceServer(iceServers);
	console.log(iceServers);
}

WebRTC.prototype.sendTimeStamp = function(){
	for (var peer in this.peerList){
		var time = Date.now();
		var timeStamp = {
				type: "timeStamp",
				sendTime: time
		}
		timeStamp = JSON.stringify(timeStamp);
		this.allConnection.connection[peer].dataChannel.send(timeStamp);
	}
}

/*
WebRTC.prototype.onHostSetup = function(){
}
 */

module.exports = WebRTC;