
var slackAppServices = angular.module("slackAppServices", ["ngResource"]);
slackAppServices.factory("Channel", ["$resource", 
				function($resource) {
						return $resource("/api/:channelId/?format=json", {}, {
								query: {method: "GET", params: {channelId: "channels"}, isArray: true}
						});
				}]);

