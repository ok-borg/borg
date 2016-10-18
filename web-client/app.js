var app = angular.module('app', ['ngCookies', 'ui.router']);

const url = window.location.protocol.includes("https") ? "https://ok-b.org:9993" : "http://borg.crufter.com:9992"
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

app.factory('Session', function($http, $cookies, $q, $window) {
    var user = {};
    var Session = {
        setToken: function(val) {
            var now = new $window.Date();
            var exp = new $window.Date(now.getFullYear(), now.getMonth()+6, now.getDate());
            $cookies.put('token', val, {
                expires: exp
            })
        },
        getToken: function() {
            var token = $cookies.get('token');
            if (token !== undefined) {
				return token;
            }
            return "";
        },
        getUser: function(cb) {
            if (this.getToken().length == 0) {
                cb({});
                return    
            }
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

app.config(function ($locationProvider, $stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
    $locationProvider.html5Mode(true);
    $stateProvider
        .state('index', {
            url: '/',
            templateUrl: 'partials/index.html',
            controller: 'IndexController'
        })
		.state('search', {
            url: '/s/:query',
            templateUrl: 'partials/search.html',
            controller: 'SearchController',
            params: {
                query: ""
            }
        })
        .state('login', {
            url: '/login',
            templateUrl: 'partials/login.html',
            controller: 'LoginController',
        })
        .state('single', {
            url: '/t/:id/:slug',
            templateUrl: 'partials/single.html',
            controller: 'SingleController',
        })
        .state('edit', {
            url: '/edit/:id',
            templateUrl: 'partials/edit.html',
            controller: 'EditController',
        })
        .state('latest', {
            url: '/latest',
            templateUrl: 'partials/latest.html',
            controller: 'LatestController',
        })
        .state('new', {
            url: '/new',
            templateUrl: 'partials/new.html',
            controller: 'NewController',
        })
		.state('me', {
			url: '/me',
			templateUrl: 'partials/me.html',
			controller: 'MeController',
		});
});

app.controller('MainController', function($scope, $rootScope, $location, $window) {
    $rootScope.$on('$locationChangeStart', function() {
        $window.ga('send', 'pageview', { page: $location.url() });
    })
    $scope.title = "OK borg - the quickest solution to your bash woes"
    $scope.noIndex = false
    $rootScope.$on('titleChange', function(e, d) {
        $scope.title = d;
    })
    $rootScope.$on('noIndex', function(e, d) {
        $scope.noIndex = d;
    })
});

app.controller('IndexController', function(Session, $state, $interval, $scope, $http) {
    $scope.apiBaseUrl = url;
    $scope.submit = function() {
        $state.go('search', {query: $scope.query});
    }
    $scope.user = {};
    $scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
    Session.getUser(function(usr) {
        $scope.user = usr;     
    })
})

app.controller('SearchController', function(Session, $window, $state, $interval, $scope, $http) {
	var search = function(q) {
        $window.ga('send', 'event', 'search', 'frontend', q);
        $http.get(url + '/v1/query', {
            params: {
				"t": Session.getToken(),
				"q": q
			}
     	}).then(function(rsp){
			$scope.results = rsp.data;
		}).catch(function(rsp) {
            console.log(rsp);
    	});
	}
	$scope.slugify = function(text) {
		return text
        .toLowerCase()
        .replace(/[^\w ]+/g,'')
        .replace(/ +/g,'-')
	}
	$scope.body = function(bodies) {
        return bodies.join("\n")
	}
	search($state.params.query);
	$scope.$on('query-submitted', function(event, args) {
        $state.go('search', {query: args.query});
        search(args.query);
	});
});

app.controller('LatestController', function(Session, $rootScope, $state, $interval, $scope, $http) {
	$rootScope.$emit('titleChange', "Latest")
    var search = function() {
		$http.get(url + '/v1/latest', {
            params: {
				"t": Session.getToken()
			}
     	}).then(function(rsp){
			$scope.results = rsp.data;
		}).catch(function(rsp) {
            console.log(rsp);
    	});
	}
    search()
	$scope.slugify = function(text) {
		return text
        .toLowerCase()
        .replace(/[^\w ]+/g,'')
        .replace(/ +/g,'-')
	}
	$scope.body = function(bodies) {
        return bodies.join("\n")
	}
});

app.controller('SingleController', function(Session, $rootScope, $state, $interval, $scope, $http) {
	var related = function() {
        $http.get(url + '/v1/query', {
            params: {
				"t": Session.getToken(),
				"q": $scope.single.Title,
                "l": 6
			}
     	}).then(function(rsp){
			$scope.related = rsp.data;
		}).catch(function(rsp) {
            console.log(rsp);
    	});
    }
    var f = function() {
		$http.get(url + '/v1/p/' + $state.params.id).then(function(rsp){
			$scope.single = rsp.data;
            $rootScope.$emit('titleChange', rsp.data.Title)
            $rootScope.$emit('noIndex', !rsp.data.CreatedBy || rsp.data.CreatedBy.length == 0)
            related();
		}).catch(function(rsp) {
            console.log(rsp);
    	});
	}
	f()
   	$scope.slugify = function(text) {
		return text
        .toLowerCase()
        .replace(/[^\w ]+/g,'')
        .replace(/ +/g,'-')
	}
	$scope.body = function(bodies) {
        return bodies.join("\n")
	}
    $scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
});

app.controller('EditController', function(Session, $rootScope, $state, $interval, $scope, $http) {
	$rootScope.$emit('titleChange', "Edit")
    var f = function() {
		$http.get(url + '/v1/p/' + $state.params.id).then(function(rsp){
			$scope.single = rsp.data;
		}).catch(function(rsp) {
            console.log(rsp);
    	});
	}
	f()
	$scope.slugify = function(text) {
		return text
        .toLowerCase()
        .replace(/[^\w ]+/g,'')
        .replace(/ +/g,'-')
	}
	$scope.body = function(bodies) {
        return bodies.join("\n")
	}
    $scope.save = function() {
        var solutions = [];
        $("textarea").each(function(index, el) {
            var t = $(el).val().trim();
            if (t.length == 0) {
                return;
            }
            solutions.push({Body: [t]})
        })
        var p = {
            Id: $scope.single.Id,
            Title: $scope.single.Title,
            Solutions: solutions 
        }
        $http({
            url: url + '/v1/p',
            method: "PUT",
            headers: {Authorization: Session.getToken()},
            data: p
        }).then(function(rsp){
			$state.go('single', {id: $scope.single.Id, slug: $scope.slugify($scope.single.Title)})
		}).catch(function(rsp) {
            console.log(rsp);
    	});
    }
});

app.controller('NewController', function(Session, $rootScope, $state, $interval, $scope, $http) {
    $rootScope.$emit('titleChange', "Submit new")
    $scope.slugify = function(text) {
		return text
        .toLowerCase()
        .replace(/[^\w ]+/g,'')
        .replace(/ +/g,'-')
	}
    $scope.save = function() {
        var t = $scope.single.Title.trim();
        var b = $scope.single.Body.trim();
        if (t.length == 0 || b.length == 0) {
            return
        }
        var p = {
            Title: t,
            Solutions: [{"Body": [b]}] 
        }
        $http({
            url: url + '/v1/p',
            method: "POST",
            headers: {Authorization: Session.getToken()},
            data: p
        }).then(function(rsp){
			$state.go('single', {id: rsp.data.Id, slug: $scope.slugify(rsp.data.Title)})
		}).catch(function(rsp) {
            console.log(rsp);
    	});
    }
});

app.controller('MeController', function(Session, $rootScope, $state, $interval, $scope, $http) {
	$scope.user = {};
	$scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
	Session.getUser(function(usr) {
        $scope.user = usr;	
	})
	$scope.copyToClipboard = function() {
  		window.prompt("Copy this, press enter to close", Session.getToken());
	}
});

app.controller('HeaderController', function(Session, $state, $rootScope, $scope) {
    $scope.formData = {query: $state.params.query};
	$scope.isLoggedIn = Session.getToken().length >= tokenMinLen;
	Session.getUser(function(usr) {
        $scope.user = usr;	
	})
	$scope.submit = function() {
		if ($state.current.name != 'search') {
			$state.go('search', {query: $scope.formData.query});
            return;
		}
		$rootScope.$broadcast('query-submitted', {query: $scope.formData.query});
	}
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
			$state.go('index');
		} 
    }).catch(function(rsp) {
        console.log(rsp);
	});
});

