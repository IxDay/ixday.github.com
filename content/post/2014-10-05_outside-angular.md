---
title:      "Outside Angular"
date:       2014-10-05
categories: ["Tuto"]
tags:       ["javascript", "angular", "web"]
url:        "post/outside-angular"
---

There are some cases when angular is accessible but we just want to access a
specific service without bootstraping an entire application.

For example, in some tests I can load some fixtures with the $http service,
or use $compile for a simple template.

It is pretty simple to do that, but it is not clearly explained in the angular
documentation. So here is an example:

```javascript
// The module ng must be loaded
angular.injector(['ng'])

// Then we just have to load the services needed
.invoke(function ($compile, $rootScope) {

  // We create a simple template and a scope
  var tplt = '<div>{{foo}}</div>';
  var scope = $rootScope.$new();

  // Populate the scope
  scope.foo = 'bar';

  // Compile the template
  var elt = $compile(tplt)(scope);

  // We are outside of angular (no application is running), so we have to run
  // a $digest cycle
  scope.$apply();

  // And it will displays "bar" MAGIC!
  console.log(elt.html());
});
```

The only issue with this solution is that you do not have access to the config
block, so if you want to change the interpolate symbol it will not be possible,
unless...

This also works with your own module

```javascript
// Define your own module
angular.module('foo', [])
.service('hello', function () {
  this.sayHello = function () {
    console.log('hello from the angular world');
  };
});

// Then load it (module ng is mandatory)
angular.injector(['ng', 'foo']).invoke(function (hello) {
  // And it will display "hello from the angular world" MAGIC!
  hello.sayHello();
});
```

So, now it will be possible to set up a custom symbol for the interpolate
service.

```javascript
angular.module('foo', [])
.config(function($interpolateProvider) {
  // We change the symbol of the interpolate provider
  $interpolateProvider.startSymbol('{%');
  $interpolateProvider.endSymbol('%}');
})
.service('buildTemplate', function ($compile, $rootScope) {
  this.getTemplate = function () {
    // Here we use the new symbols
    var tplt = '<div>{%foo%}</div>';
    var scope = $rootScope.$new();

    scope.foo = 'bar';
    var elt = $compile(tplt)(scope);
    scope.$apply();

    return elt.html();
  }
});

angular.injector(['ng', 'foo']).invoke(function (buildTemplate) {
  // And it still displays "bar" MAGIC!
  console.log(buildTemplate.getTemplate());
});

```

Here it seems not really useful but sometimes... you know ;)

(you can find a full Plunker [here](http://plnkr.co/edit/giAIRSwdTnOPAhj3TRs7))
