
describe("ChannelListCtrl", function() {
	beforeEach(module("slackApp"));
	it("should create channels model with 4 channels", inject(function($controller) {
		var scope = {},
			ctrl = $controller("ChannelListCtrl", {$scope: scope});
		expect(scope.channels.length).toBe(4)
	}));
});
