from __future__ import unicode_literals

import livereload
import fabric.colors as colors
import fabric.contrib.project as project
import fabric.contrib.files as files
import fabric.api as fabric
import jinja2
import pelican.utils as utils

import datetime
import os
import os.path
import sys

# Local path configuration (can be absolute or relative to fabfile)
fabric.env.deploy_path = 'output'
fabric.env.content_path = 'content'
fabric.env.jinja = jinja2.Environment(
    loader=jinja2.PackageLoader('fabfile', 'templates')
)

DEPLOY_PATH = env.deploy_path


def clean():
    if os.path.isdir(DEPLOY_PATH):
        local('rm -rf {deploy_path}'.format(**env))
        local('mkdir {deploy_path}'.format(**env))

def build():
    local('pelican -s pelicanconf.py')

def rebuild():
    clean()
    build()

def regenerate():
    local('pelican -r -s pelicanconf.py')

def serve(*args):
    port = args[0] if len(args) > 0 else 8000

    if not isinstance(port, int) or port < 1024 or port > 65535:
        print(colors.red('Port must be an integer between 1024 and 65535...'))
        return

    server = livereload.Server()
    server.watch(fabric.env.content_path, build)
    server.serve(port=port, root=fabric.env.deploy_path)

def reserve():
    build()
    serve()

def preview():
    local('pelican -s publishconf.py')

def new_post(*args):
    title = args[0] if len(args) > 0 else fabric.prompt('New post title?')
    title = unicode(title, 'utf8')

    date = datetime.date.today().isoformat()
    filename = '.'.join([date, utils.slugify(title), 'md'])
    filename = os.path.join(fabric.env.content_path, filename)
    print(' '.join([colors.green('[create new post]'), filename]))

    (fabric.env.jinja.get_template('new_post.tplt')
     .stream(title=title)
     .dump(filename, 'utf8'))


