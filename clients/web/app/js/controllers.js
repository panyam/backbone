
var slackAppControllers = angular.module("slackAppControllers", []);

slackAppControllers.controller("ChannelListCtrl", ["$scope", "Channel", function($scope, Channel) {
	$scope.channels = Channel.query();
}]);

slackAppControllers.controller("ChannelDetailsCtrl", ["$scope", "$routeParams", "Channel", function($scope, Channel) {
	$scope.Channel = Channel.get({channelId: $routeParams.channelId},
			function(channel) {
			});
}]);

slackAppControllers.controller("MessageListCtrl", ["$scope", "Channel", function($scope, Channel) {
	$scope.channels = []; // Channel.query();
}]);
