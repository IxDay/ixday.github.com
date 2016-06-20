+++
title = "Listing directory with wsgi"
date = "2015-07-02"
categories = ["Project"]
tags = ["dev", "python"]
+++

Recently I ran into an issue with my blog and pelican
(the blogging engine I use). For some reasons (which I explain [here]())
I had to develop a small wsgi app which act like the
[python SimpleHTTPServer](https://docs.python.org/2/library/simplehttpserver.html).

I tried a lot of things but they never worked as I wanted them to.
So, I decided to do this by myself.

## Specifications

Here are the features I wanted:

* Define a directory as the root (/) directory from which all the files will be
served.
* If this directory contains an index.html file, serves it to the user.
* If there is no index.html file list all the subdirectories and files in a
simple page and provide links to navigate.

## Tools

* For this little project I need a simple server which is available directly
in the [stdlib](https://docs.python.org/2/library/simplehttpserver.html).
* I will use the famous [path.py lib](https://github.com/jaraco/path.py) for
handling all the interactions with the FS.
* Then [Jinja](https://github.com/jaraco/path.py) for all the rendering.

## Hello world!

First of all, making it works for the most simple use case
```python
import wsgiref.simple_server

def application(environ, start_response):
	"""Return a simple Hello World when accessing an url (anyone) on the server
	"""
	# send status_code and headers
	start_response('200 OK', [('Content-Type', 'text/html')])
	# send body
	return ['Hello World!']

if __name__ == '__main__':
	# create a server on 'localhost' and port 8080 using the application
	srv = wsgiref.simple_server.make_server('localhost', 8080, application)
	# serve...
	srv.serve_forever()
```

## Routing

Then I have to handle the request, translate the url in a possible path in the
FS, and dispatch to the different choices available.

```python
import wsgiref.simple_server
import path

ROOT = path.path('/tmp') # The root directory from which I will serve the files

def application(environ, start_response):
	"""Route to the dedicated functions"""

	# build the possible path of the file to serve or directory to list
	path = ROOT + environ.get('PATH_INFO')

	# build the possible path of the index file if the path is a directory
    path_index = path + 'index.html'


	# checks if the path exists, if not want to send a 404 not found
    if not path.exists():
        return not_found(start_response)

	# checks if the path links to a file, in this case serve the file
    if path.isfile():
        return serve_file(path, start_response)

	# now we are sure that the user asks for a directory, checks if the
	# directory has an index.html file, if so, serve it.
    if path_index.exists():
        return serve_file(path_index, start_response)

	# last case we want to list what is in the directory
    return list_dir(path, start_response)


if __name__ == '__main__':
	srv = wsgiref.simple_server.make_server('localhost', 8080, application)
	srv.serve_forever()

```

## Serve content

### Return 404

Most simple use case, return a 404 error and a small message
```python
	def not_found(start_response):
		# status code and header
        start_response('404 NOT FOUND', [('Content-Type', 'text/plain')])
		# simplest message ever
        return ['Not Found']
```

### Serve file
```python
import mimetypes

chunk_size = 64 * 1024

def serve_file(filepath, start_response):
	# retrieve mimetype for serving purpose
	mime = mimetypes.guess_type(f)

	# start response with the given mimetype
	sr('200 OK', [('Content-Type', mime[0])])

	# yield the file content through network (chunks function from path.py)
	return f.chunks(chunk_size, 'rb')
```

### List directory

Last but not least, the listing of the directory. Here I use Jinja for the
templating.

```python
import jinja2

# setup my jinja environment with the templates directory containing my
# template files
env = jinja2.Environment(
    loader=jinja2.PackageLoader('my_wsgi_listdir', 'templates')
)


def list_dir(directory, start_response):

	# start a response as an html page
	start_response('200 OK', [('Content-Type', 'text/html')])

	# retrieve informations from FS through path.py API
	context = {
		'directory': directory.relpath(ROOT),
		'links': [ f.relpath(directory) for f in d.listdir() ]
	}

	# render the template with the informations and return the stream
	return env.get_template('list_dir.tplt').stream(**context))

```

Here is the `list_dir.tplt` file


```jinja
<html>
<head>
  <title>Directory listing for {{directory}}</title>
</head>
<body>
  <h2>Directory listing for {{directory}}</h2>
  <hr>
  <ul>
  {% for link in links %}
    <li><a href="{{ link }}">{{ link }}</a></li>
  {% endfor %}
  </ul>
  <hr>
</body>
</html>
```

And that's all, simple and efficient this little wsgi application fits my
needs.
