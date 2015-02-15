'use strict';

var slackApp = angular.module("slackApp", [
				"ngRoute",
				"slackAppControllers",
				"slackAppServices"
				]);

slackApp.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/api/channels', {
        templateUrl: 'partials/channel-list.html',
        controller: 'ChannelListCtrl'
      }).
      when('/api/channels/:channelId', {
        templateUrl: 'partials/channel-details.html',
        controller: 'ChannelDetailsCtrl'
      }).
      otherwise({
        redirectTo: '/api/channels'
      });
  }]);
