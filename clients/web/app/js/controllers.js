
var slackApp = angular.module("slackApp", []);

slackApp.controller("ChannelListCtrl", function($scope) {
		$scope.currentTeam = {
			"teamId": "team1",
			"teamName": "Backbone"
		};
		$scope.channels = [
			{"name": "Channel 1"},
			{"name": "Channel 2"},
			{"name": "Channel 3"},
			{"name": "Channel 4"},
		];
});
