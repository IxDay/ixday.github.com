# coding: utf8
from __future__ import unicode_literals

import livereload
import fabric.colors as colors
import fabric.contrib.project as project
import fabric.contrib.files as files
import fabric.api as fabric
import jinja2
import path
import pelican.utils as utils

import collections
import datetime
import logging
import logging.handlers
import mimetypes
import sys

# Local path configuration (can be absolute or relative to fabfile)
fabric.env.deploy_path = path.path('output')
fabric.env.content_path = path.path('content')
fabric.env.jinja = jinja2.Environment(
    loader=jinja2.PackageLoader('fabfile', 'templates')
)

MESSAGE_FORMAT = '%(levelname)s %(message)s'
LEVELS = {
    'WARNING': colors.yellow('WARN', bold=True),
    'INFO': colors.blue('INFO', bold=True),
    'DEBUG': colors.green('DEBUG', bold=True),
    'CRITICAL': colors.magenta('CRIT', bold=True),
    'ERROR': colors.red('ERROR', bold=True),
}


class FabricFormatter(logging.Formatter):
    def format(self, record):
        record.levelname = LEVELS.get(record.levelname) + ':'
        return super(FabricFormatter, self).format(record)


class Server(livereload.Server):
    def _setup_logging(self):
        super(Server, self)._setup_logging()
        server_handler = logging.getLogger('livereload').handlers[0]
        server_handler.setFormatter(FabricFormatter(MESSAGE_FORMAT))


def application(env, start_response):
    chunk_size = 64 * 1024

    def not_found(sr):
        sr('404 NOT FOUND', [('Content-Type', 'text/plain')])
        return ['Not Found']

    def serve_file(f, sr):
        mime = mimetypes.guess_type(f)

        sr('200 OK', [('Content-Type', mime[0])])
        return f.chunks(chunk_size, 'rb')

    def list_dir(d, sr):
        sr('200 OK', [('Content-Type', 'text/html')])

        context = {
            'directory': d.relpath(fabric.env.deploy_path),
            'links': [ f.relpath(d) for f in d.listdir() ]
        }

        return (fabric.env.jinja.get_template('list_dir.tplt')
                .stream(**context))


    path = fabric.env.deploy_path + env.get('PATH_INFO')
    path_index = path + 'index.html'

    if not path.exists():
        return not_found(start_response)
    if path.isfile():
        return serve_file(path, start_response)
    if path_index.exists():
        return serve_file(path_index, start_response)

    return list_dir(path, start_response)


@fabric.task
def clean():
    if fabric.env.deploy_path.isdir():
        fabric.local('rm -rf {deploy_path}'.format(**fabric.env))
        fabric.local('mkdir {deploy_path}'.format(**fabric.env))

@fabric.task
def build():
    fabric.local('pelican -s pelicanconf.py')

@fabric.task
def rebuild():
    clean()
    build()

@fabric.task
def regenerate():
    fabric.local('pelican -r -s pelicanconf.py')

@fabric.task
def serve(*args):
    port = args[0] if len(args) > 0 else 8000

    if not isinstance(port, int) or port < 1024 or port > 65535:
        print(colors.red('Port must be an integer between 1024 and 65535...'))
        return

    build()
    server = Server(application)
    server.watch(fabric.env.content_path, build)
    server.serve(port=port, debug=True)

@fabric.task
def reserve():
    build()
    serve()

@fabric.task
def new_post(*args):
    title = args[0] if len(args) > 0 else fabric.prompt('New post title?')
    title = unicode(title, 'utf8')

    date = datetime.date.today().isoformat()
    filename = '.'.join([date, utils.slugify(title), 'md'])
    filename = fabric.env.content_path / filename
    print(' '.join([LEVELS['INFO'], 'Create new post:', filename]))

    (fabric.env.jinja.get_template('new_post.tplt')
     .stream(title=title)
     .dump(filename, 'utf8'))

@fabric.task
def publish():
    build()
    fabric.local('ghp-import {0}'.format(fabric.env.deploy_path))
    fabric.local('git push origin -f gh-pages:master')

