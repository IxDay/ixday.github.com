#!/usr/bin/env python
# -*- coding: utf-8 -*- #
from __future__ import unicode_literals

AUTHOR = u'IxDay'
SITENAME = u'Not today...'
FOOTER_TEXT = u'Done with love... and beer'


PATH = 'content'

TIMEZONE = 'Europe/Paris'

DEFAULT_LANG = u'en'

# Feed generation is usually not desired when developing
FEED_ALL_ATOM = None
CATEGORY_FEED_ATOM = None
TRANSLATION_FEED_ATOM = None

# Blogroll
LINKS = (('Github', 'https://github.com/IxDay/'),)

DEFAULT_PAGINATION = 10

# Uncomment following line if you want document-relative URLs when developing
#RELATIVE_URLS = True

EXTRA_PATH_METADATA = {
    'extras/favicon.ico': {'path': 'favicon.ico'},
}

PLUGIN_PATHS = ['plugins']
PLUGINS = ['image_process']

IMAGE_PROCESS = {
    'article-image': {
        'type': 'image',
        'ops': ["scale_in 700 700 True"],
    },
}

STATIC_PATHS = [
    'images',
    'extras'
]

THEME = "pelican-chunk"

GOOGLE_ANALYTICS = 'UA-38228870-1'

FAVICON_URL = '/favicon.ico'

SINGLE_AUTHOR = True
MINT = False

DISQUS_SITENAME = 'IxDay'

# WRITE_SELECTED = []
