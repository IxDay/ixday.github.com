+++
title = "Test promises with jasmine"
date = "2014-10-30"
categories = ["Snippet"]
tags = ["jasmine", "tests"]
+++

Jasmine is a good testing framework, which I really like, it is really powerfull
and has just the amount of features to perform a huge variety of tests.

At some point I had to tests promises, and more generally testing that some
part of a function is not called (you will have to adapt the snippet but the
idea is here).

It is pretty simple, but not well known (I checked some stackoverflow threads
before finding this).

```javascript
/*
  We will define the promise helper to be accessible through the this keyword in
  jasmine
 */
beforeEach(function () {

  /*
   The this keyword is accessible in every jasmine tests, and you can populate
   it, this has been set up in jasmine 2.0 in order to avoid the tricky variable
   scoping of javascript.
   */
  this.promiseHelper = function (promise) {

    var successCallback = jasmine.createSpy('successCallback');
    var failCallback = jasmine.createSpy('failCallback');

    /*
     Promises are chainable so we can make some tests before sending the promise
     to this function, here we just insert the two spies in the
     */
    promise.then(successCallback, failCallback)['finally'](
      function () {
        expect(successCallback).toHaveBeenCalled();
        expect(failCallback).not.toHaveBeenCalled();
      });
  }
});

```

Do not forget to test your code, it will save you a lot of time and debugging in
the future if tests are correctly done ;)


