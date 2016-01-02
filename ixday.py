#! /usr/bin/env python
# coding: utf8

import click

import livereload
import jinja2
import path
import pelican.utils as utils

import datetime
import mimetypes
import subprocess

import SimpleHTTPServer
import SocketServer

PELICAN_CONF = 'pelicanconf.py'
OUTPUT = 'output'
CONTENT = 'content'
THEME = './theme'

class Config(object):
    def __init__(self):
        self.settings = path.path(PELICAN_CONF)
        self.output = path.path(OUTPUT)
        self.content = path.path(CONTENT)
        self.theme = path.path(THEME)
        self.jinja = jinja2.Environment(
            loader=jinja2.PackageLoader('ixday', 'templates')
        )


class Application(object):

    def __init__(self, config):
        self.config = config
        self.chunk_size = 64 * 1024

    def _not_found(self, sr):
        sr('404 NOT FOUND', [('Content-Type', 'text/plain')])
        return ['Not Found']

    def _serve_file(self, f, sr):
        mime = mimetypes.guess_type(f)

        sr('200 OK', [('Content-Type', mime[0])])
        return f.chunks(self.chunk_size, 'rb')

    def _list_dir(self, d, sr):
        sr('200 OK', [('Content-Type', 'text/html')])

        context = {
            'directory': d.relpath(self.config.output),
            'links': [f.relpath(d) for f in d.listdir()]
        }

        return (self.config.jinja.get_template('list_dir.tplt')
                .stream(**context))

    def __call__(self, env, start_response):
        path = self.config.output + env.get('PATH_INFO')
        path_index = path + 'index.html'

        if not path.exists():
            return self._not_found(start_response)
        if path.isfile():
            return self._serve_file(path, start_response)
        if path_index.exists():
            return self._serve_file(path_index, start_response)

        return self._list_dir(path, start_response)


pass_config = click.make_pass_decorator(Config, ensure=True)


@click.group()
def cli():
    pass


@cli.command()
@pass_config
def clean(config):
    '''Clean the "output" directory'''
    click.echo('Clean output directory')
    for dir in config.output.dirs():
        dir.rmtree()
    for file in config.output.files():
        file.remove()


@cli.command()
@pass_config
def build(config):
    '''Build the pelican static site'''
    subprocess.call(['pelican', '-s', config.settings])


@cli.command()
@click.pass_context
def rebuild(ctx):
    '''Run clean then build'''
    ctx.invoke(clean)
    ctx.invoke(build)


@cli.command()
@click.option('--no-debug', is_flag=True, help='remove the debug output')
@click.option('--no-lr', is_flag=True, help='remove the livereload support')
@click.argument('port', default=8000, required=False)
@pass_config
@click.pass_context
def serve(ctx, config, no_debug, no_lr, port):
    '''Serve the 'output' directory, default port is 8000'''
    if port < 1024 or port > 65535:
        error_msg = 'port must be an integer between 1024 and 65535'
        raise click.BadParameter(error_msg, param_hint='port')

    debug = not no_debug
    inner_build = lambda: ctx.invoke(build)

    inner_build()
    if no_lr:
        with config.output:
            Handler = SimpleHTTPServer.SimpleHTTPRequestHandler
            httpd = SocketServer.TCPServer(('', port), Handler)

            click.echo('serving at port %d' % port)
            httpd.serve_forever()
    else:
        server = livereload.Server(Application(config))
        server.watch(config.content, inner_build)
        server.watch(config.theme, inner_build)
        server.serve(port=port, debug=debug)


@cli.command()
@click.argument('title')
@pass_config
def new_post(config, title):
    '''Create a new blog entry'''
    date = datetime.date.today().isoformat()
    filename = '.'.join([date, utils.slugify(title), 'md'])
    filename = config.blog / filename
    click.echo('Create new post: %s' % click.style(filename, fg='green'))

    (config.jinja.get_template('new_post.tplt')
     .stream(title=title)
     .dump(filename, 'utf8'))


@cli.command()
@pass_config
@click.pass_context
def publish(ctx, config):
    ctx.invoke(build)
    subprocess.call(['ghp-import', config.output])
    subprocess.call(['git', 'push', 'origin', '-f', 'gh-pages:master'])


if __name__ == '__main__':
    cli()
