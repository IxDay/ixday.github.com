Title: Angular $parse
Date: 2014-10-23
Category: Tuto
Tags: javascript, angular, web

Hey there, I have started to be tired of weak example of angular power,
so I will go deeper on angular services and directives and wrote some articles
about it.

The service `$parse` is the one who runs on the html to bind data with your
javascript. It provides a lot of useful features which can be really interesting
especially with directive manipulation.

So, we will illustrate with some examples:

```js
//This notation allow us to stay in pure javascript world
angular.injector(['ng'])
.invoke(function ($parse, $rootScope) {

  //new clean scope (optionnal, it will work directly on $rootScope)
  var $scope = $rootScope.$new();

  //allocate simple value
  $scope.foo = 'bar';

  //here the magic
  var foo = $parse('foo');

  /*
    The function which has been returned look for a context to interpolate the
    value
  */
  console.log(foo($scope));
  // display => bar

  /*
    Redefine foo value
   */
  foo.assign($scope, 'foo');

  /*
    $scope has changed Oo
   */
  console.log($scope.foo);
  // display => foo
});
```

Pretty cool uh? But wait there is more

```js
angular.injector(['ng'])
.invoke(function ($parse, $rootScope) {
  var $scope = $rootScope.$new();

  //allocate function
  $scope.foo = function (argFoo) {
    return argFoo;
  };

  //allocate value
  $scope.someValue = 'foo'

  var foo = $parse('foo(someValue)');

  console.log(foo($scope));
  // display => foo


});
```

It also interpolates functions?! But wait there is more

```js
angular.injector(['ng'])
.invoke(function ($parse, $rootScope) {
  var $scope = $rootScope.$new();

  //allocate function
  $scope.foo = function (argFoo) {
    return argFoo;
  };

  //allocate value
  $scope.someValue = 'foo';

  var foo = $parse('foo(someValue)');

  /*
    The function returned by the $parse service has a second argument wich will
    override the context in case of conflict
   */
  console.log(foo($scope, {someValue: 'bar'}));
  // display => bar
});
```

Now you will ask: Is there a use case for that? Of course, passing additionnal
arguments in directive is possible. For example, the ng-click directive
([doc](https://docs.angularjs.org/api/ng/directive/ngClick))

This pattern allows you to retrieve the `$event` object directly in your
function handler.

Here is a simplification with explanation of the ng-click directive


```js
angular.module('myModule', [])
.directive('myNgClick', ['$parse', function ($parse) {
  return {
    link: function (scope, elt, attr) {
      /*
       Gets the function you have passed to ng-click directive
       Parse returns a function which has a context and extra params which
       overrides the context
      */
      var handler = $parse(attr['myNgClick']);

      /*
      here you bind on click event you can look at the documentation
      https://docs.angularjs.org/api/ng/function/angular.element
      */
      elt.on('click', function (event) {

        //callback is here for the explanation
        var callback = function () {

          /*
           Here handler will do the following, it will call the dedicated
           function and fill the arguments with the elements found in the scope
           (if possible), the second argument will override the $event attribute
           in the scope (if there is some) and provide the event element of the
           click
          */
          handler(scope, {$event: event});
        };

        //$apply force angular to run a digest cycle in order to propagate the
        //changes
        scope.$apply(callback);
      });
    }
  };
}]);
```

And it will work, example => [here](http://jsbin.com/hujeluqigo/2/edit?html,js,console,output)

You can also find example of the two first parts [here](http://jsbin.com/nahivi/edit?js,console)

Hope you enjoyed the trip, see ya o/
