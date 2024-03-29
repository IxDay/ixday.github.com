---
title:      "Getting rid of gulp bunch of dependencies"
date:       2016-04-10
categories: ["Snippet"]
tags:       ["dev", "javascript"]
url:        "post/getting_rid_of_gulp_bunch_of_dependencies"
---

Recently Nodejs environment broke due the removal from npm
of a small library (11 SLOC): `leftpad`. As it hit the world
and broke a bunch of projects and CIs, I asked myself if my projects
contains so much dependencies that if one break, everything collapse.

## The problem

For developing my frontend I use a tool which I really like:
[Gulp](http://gulpjs.com/). The issue there, is that for working
with multiple building process involved a lot of glue and third
party libraries. Here is an example: https://github.com/gulpjs/gulp/blob/master/docs/recipes/browserify-uglify-sourcemap.md
Just for using [browserify](http://browserify.org/)
(which is another great tool). Each of those libraries involved other
dependencies and so on until we download the whole internet to perform the
most simple tasks.

## The solution ?

According to this, I started looking at those libraries in order to see
if I can implement them with a reduced number of dependencies and
lines of code. The answer is: yes, and moreover it is quite simple and
helps me learned some new things.

Here is the goal of the exercise:

  - rewrite `vinyl-source-stream`, `merge-stream` and `vinyl-buffer`
  - fix browserify in the gulp environment, according to
    [this issue](http://stackoverflow.com/questions/30077567/browserify-errors-ending-gulp-watch-task)

## The implementation

To perform the implementation I will only keep the
[through2 module](https://github.com/rvagg/through2) for better creating
streams (don't worry as it is a requirements of all the gulp plugins it will
not add a new dependency)
and [gulp util](https://github.com/gulpjs/gulp-util) which will allow me
to have some helpers to deal with gulp (as through2 it will not create new
dependencies as it is a requirements for basically all the gulp plugins).

I will create a simple `utils.js` file aside my `gulpfile.js` to store
those implementations.

Here is the requirementents of the file:

```js
var through = require('through2');
var gutil = require('gulp-util');
```

### vinyl-source-stream

According to its presentation this module only provide a convenient wrapper
in order to link legacy nodejs streams and gulp implementation of streams
vinyl.

Here is the code I use:

```js
module.exports.source = function(filename) {

  // this is basically a stream, I will use the javascript scoping for
  // convenience
  var ins = through();

  // here is the piping to the previous stream, we need an object stream as we
  // will push a vinyl file in it.
  // The use will be the following:
  // someLegacyStream.pipe(utils.source()).pipe(whateverInGulpWorld);
  return through.obj(function(chunk, _, cb) {

    // Checks if we have initialized the stream with a vinyl file
    // basically, this just happen once at startup of streaming.
    if (!this._ins) {
      this._ins = ins; // this is just a convenient way of keeping a state

      // Create a vinyl file and pass "ins" stream as a content,
      // and an optional filename
      this.push(new gutil.File({contents: ins, path: filename}));
    }

    // push the chunk into our new stream in order to unify output
    ins.push(chunk);

    // as this is asynchronous just notify the system that we handled
    // the chunk
    cb();

  }, function() {
    // close the stream by pushing "null"
    ins.push(null);
    this.push(null);
  });
};
```

Okay, so, we replace a whole module by ~14 SLOC, additionally I get back some
understanding on how nodejs streams work, and on gulp plugins implementation.

This is a good start! Let's continue this way

### merge-streams

This module simply merge multiple streams in a single one, this allow the
developper to compose with multiple inputs.

```js

// the function treats arguments as a list of streams
module.exports.merge = function(/* streams... */) {

  var sources = []; // keep a track of streams merged
  var output  = gutil.noop(); // this will be the output,
                              // a simple stream which does nothing

  // this function will be called when unpiping and when a source stream ends
  function remove(source) {
    // remove the stream from sources array
    sources = sources.filter(function(it) { return it !== source; });

    // if it is the last stream opened and if the output is not yet closed,
    // we close it
    if (!sources.length && output.readable) output.end();
  }

  // when a stream is unpiped we remove it
  output.on('unpipe', remove);

  // for each stream (arguments is not a regular array this is why we use
  // this syntax)
  Array.prototype.slice.call(arguments).forEach(function(source) {

    // add the stream to our array of sources
    sources.push(source);

    // bind the remove function to the end event of the source
    source.once('end', remove.bind(null, source));

    // pipe the stream to our output, and let it open (the output)
    // even when the stream ends (in order to handle the others)
    source.pipe(output, {end: false});
  });

  return output;
};
```

Here again we replace an entire module with few lines of code (~15 SLOC).

### vinyl-buffer

This final module is also part of the vinyl utilitaries. It takes the chunks
from a stream and return them as a nodejs buffer. Like the others, this
one is quite simple and only requires to know a bit of node internal
operations and libraries.

```js

module.exports.buffer = function() {

  // like the others we will create a stream object
  return through.obj(function(file, _, cb) {

    var that = this; // keep an internal reference of this across js scoping
    var bufs = []; // array of buffer we will populate

    // if it is already a buffer or it contains nothing, just push and finish
    if (file.isNull() || file.isBuffer()) {
      that.push(file);
      return cb();
    }

    // otherwise we take the content of the stream and we pipe it
    file.contents.pipe(through.obj(
      function(data, _, cb) {
        // create a new buffer with data if it is not and push it to our array
        bufs.push(Buffer.isBuffer(data) ? data : new Buffer(data));
        cb();
      },
      function() {
        // when we have retrieved all the chunks, create a copy of file
        file = file.clone();
        // and replace the content with only one huge buffer
        file.contents = Buffer.concat(bufs);

        // push it and that's it
        that.push(file);
        cb();
      }
    ));
  });
};
```

We are still under 20 lines (we start to have a pattern here... just trolling).
This part is done !

### browserify error handling with gulp

As I mentionned it quickly at the beginning I also want to fix an issue I
have with browserify and gulp.

I finally found the solution on stack overflow and I am surprised that not a
lot of people ran in this problem before.

I also use the `gulp.src` syntax to retrieve files instead of loading the
globbing module which does exactly the same things
(one dependency removed, Yay \o/).

Here is how I implemented this:

```js
// we need the browserify module
var browserify = require('browserify');

// we create a utils.browserify kind of plugin here, which can take options
module.exports.browserify = function(options) {

  // initialize browserify with options
  var b = browserify(options || {debug: true});

  // here we create a stream which will be pluggable through a pipe
  // in order to avoid spending resources we will use it like this
  // gulp.src('**/*.js', {read: false}).pipe(utils.browserify()).pipe(whatever)
  // note the "read: false" which will avoid reading file content and only
  // provide vinyl file object (not opened) to the stream
  var s = through.obj(

    // here we only retrieve the files provided by the stream
    function(file, _, cb) {
      b.add(file.path);
      cb();
    },

    // and at flush, we bundle the result through browserify
    function() {
      b.bundle()

      // on error we provide a helpful message to the user through gutil.log
      .on('error', function(err) {
        var message = err.annotated || err.toString();
        gutil.log(new gutil.PluginError('browserify', message).toString());

        // and here is the magic not provided before, we notice gulp that
        // the stream as ended, and it can continue to watch our files states
        s.emit('end');
      })

      // when data is provided by the stream we simply push it into the
      // returned one
      .on('data', function(chunk) { s.push(chunk); })

      // do not forget to pass the end event also
      .on('end', function() { s.emit('end'); });
    }
  );

  return s;
};
```

This one was quite more complicated but it finally works and I am happy to stay
in the "gulp world" and do not provide any tricky thing.

## Conclusion

This experience was really interesting, reducing code dependencies,
gaining power on underlying nodejs concepts, fixing some bugs also.
This was not a waste of time and I really enjoyed it.

I also will not make any comment on node ecosystem because I think everything
has already been told. But this exercise proved that some libraries are really
not complicated and can be reimplemented in order to avoid bad surprise in the
future.
