var app = angular.module('app', ['ngCookies', 'ui.router']);

const url = "http://ok-b.org:9992"
const tokenMinLen = 10

app.directive('a', function() {
    return {
        restrict: 'E',
        link: function(scope, elem, attrs) {
            if(attrs.ngClick || attrs.href === '' || attrs.href === '#'){
                elem.on('click', function(e){
                    e.preventDefault();
                });
            }
        }
   };
});

app.factory('Session', function($http, $cookies, $q) {
    var user = {};
    var Session = {
        setToken: function(val) {
            $cookies.token = val;
        },
        getToken: function() {
            if ($cookies.token !== undefined) {
				return $cookies.token;
            }
            return "";
        },
        getUser: function(cb) {
            var that = this;
			var scb = function(rsp) {
				user = rsp.data;
                cb(user);
            };
            if (!user || !user.Id) {
				var token = this.getToken();
                $http({
					method: 'GET',
					url: url + '/v1/user?token=' + token,
				}).then(scb);
            } else {
                cb(user);
            }
        },
        logout: function() {
            $cookies.token = "";
            user = {};
        }
    };
    return Session;
});

app.config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
    $stateProvider
        .state('index', {
            url: '/',
            templateUrl: 'partials/index.html',
            controller: 'IndexController'
        })
		.state('search', {
            url: '/',
            templateUrl: 'partials/search.html',
            controller: 'SearchController'
        })
        .state('login', {
            url: '/login',
            templateUrl: 'partials/login.html',
            controller: 'LoginController',
        })
		.state('me', {
			url: '/me',
			templateUrl: 'partials/me.html',
			controller: 'MeController',
		});
});

app.controller('IndexController', function(Session, $state, $interval, $scope, $http) {
	$scope.query = function() {
        $http.get(url + '/query', {
            "token": Session.getToken(),
        }).then(function(rsp){
        		
		}).catch(function(rsp) {
            console.log(rsp);
        });
    }
	$scope.user = {};
	$scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
	Session.getUser(function(usr) {
		console.log(usr);
        $scope.user = usr;	
	})
});

app.controller('SearchController', function(Session, $state, $interval, $scope, $http) {
});


app.controller('MeController', function(Session, $state, $interval, $scope, $http) {
	$scope.user = {};
	$scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
	Session.getUser(function(usr) {
		console.log(usr);
        $scope.user = usr;	
	})
	$scope.copyToClipboard = function() {
  		window.prompt("Copy this, press enter to close", Session.getToken());
	}
});

app.controller('HeaderController', function(Session, $state, $interval, $scope, $http) {
    $scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
	Session.getUser(function(usr) {
		console.log(usr);
        $scope.user = usr;	
	})
});

// remove this cruft once getting get params from ui.router works, ehh
function gup( name, url ) {
      if (!url) url = location.href;
      name = name.replace(/[\[]/,"\\\[").replace(/[\]]/,"\\\]");
      var regexS = "[\\?&]"+name+"=([^&#]*)";
      var regex = new RegExp( regexS );
      var results = regex.exec( url );
      return results == null ? null : results[1];
}

// Login page that exchanges the code for a token, then stores the token in a cookie
app.controller('LoginController', function(Session, $scope, $http, $rootScope, $state) {
    var code = gup("code");
	$http.post(url + '/v1/auth/github', code).then(function(rsp) {
		if (!rsp.data.Token) {
			$scope.error = rsp.data;
		} else {
			Session.setToken(rsp.data.Token);
			$state.go('search');
		} 
    }).catch(function(rsp) {
		console.log(rsp);
	});
});

